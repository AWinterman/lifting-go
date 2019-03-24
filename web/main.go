package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"cloud.google.com/go/civil"
	"github.com/awinterman/lifting"
	"github.com/go-playground/form"
)

const (
	category    = "Category"
	sessionDate = "SessionDate"
	exercise    = "Exercise"
	volume      = "Volume"
	weight      = "Weight"
	hour        = "DurationHour"
	minute      = "DurationMinute"
	second      = "DurationSecond"
	effort      = "Effort"
	failure     = "Failure"
	sets        = "Sets"
	comment     = "Comment"
	id          = "ID"
)

var edit = regexp.MustCompile(`/edit/(?P<ID>\d\d*)(/)?`) // Contains "abc"

// Context is basic context for the site.
type Context struct {
	History    []lifting.Repetition
	Categories []string
	Exercises  []string
	Units      []string
	Repetition *lifting.Repetition
	Now        string
	Message    string
}

func now() string {
	return civil.DateOf(time.Now()).String()
}

func main() {
	var storage, err = lifting.CreateStorage(".lift.sqlite", nil)

	if err != nil {
		panic(err)
	}

	mux := http.NewServeMux()
	fs := http.FileServer(http.Dir("static/"))

	handlers := Handlers{Storage: storage, Decoder: form.NewDecoder()}

	mux.Handle("/stylesheets/", fs)
	mux.HandleFunc("/", handlers.handle)

	port := ":9000"
	log.Printf("Listening http://localhost:%v", port)
	http.ListenAndServe(port, mux)
}

type handlerFunc func(w http.ResponseWriter, r *http.Request)

// Handlers is all the http handlers
type Handlers struct {
	Storage lifting.Storage
	Decoder *form.Decoder
}

func (h *Handlers) handle(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "Text/HTML")
	path := r.URL.Path

	switch {
	case path == "/" || path == "":
		if r.Method != "GET" {
			h.handleErrors(w, r, fmt.Errorf("Unexpected method"), http.StatusMethodNotAllowed)
		}

		context, err := h.getContext()
		if err != nil {
			h.handleErrors(w, r, err, http.StatusInternalServerError)
			return
		}
		h.contextHandler(w, r, context, "index.html")
		break
	case ((path == "/create/" || path == "/create") && r.Method == "POST") || (edit.MatchString(path) && r.Method == "PUT"):
		h.handleCreatePost(w, r, nil)
	case path == "/create/" || path == "/create" && r.Method == "GET":
		h.handleCreateGet(w, r)
	case path == "/create/" || path == "/create":
		h.handleErrors(w, r, fmt.Errorf("Unexpected method"), http.StatusMethodNotAllowed)
	case edit.MatchString(path):
		h.handleEdit(w, r)
		break
	default:
		var templates = template.Must(template.ParseFiles("templates/404.html", "templates/base.html"))
		err := templates.ExecuteTemplate(w, "404.html", Context{})
		if err != nil {
			h.handleErrors(w, r, err, http.StatusInternalServerError)
			return
		}
		break

	}
}

func (h *Handlers) getRep(id string) (*lifting.Repetition, error) {
	ID, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}
	repetition, err := h.Storage.GetByID(int(ID))
	if err != nil {
		return nil, err
	}
	return repetition, nil
}

func (h *Handlers) handleEdit(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	matches := edit.FindStringSubmatch(path)

	repetition, err := h.getRep(matches[1])
	if err != nil {
		h.handleErrors(w, r, err, http.StatusInternalServerError)
		return
	}

	if r.Method == "GET" {
		context, err := h.getContext()
		if err != nil {
			h.handleErrors(w, r, err, http.StatusInternalServerError)
			return
		}

		context.Repetition = repetition

		h.contextHandler(w, r, context, "form.html")
		return
	} else if r.Method == "POST" {
		h.handleCreatePost(w, r, repetition)
	}
}

func (h *Handlers) getContext() (*Context, error) {
	reps, err := h.Storage.GetLast(10, 0)
	if err != nil {
		return nil, err
	}

	categories, err := h.Storage.GetUniqueCategories()

	if err != nil {
		return nil, err
	}
	exercises, err := h.Storage.GetUniqueExercises()
	if err != nil {
		return nil, err
	}

	units, err := h.Storage.GetUniqueUnits()
	if err != nil {
		return nil, err
	}

	return &Context{
		History:    reps,
		Categories: categories,
		Exercises:  exercises,
		Repetition: nil,
		Now:        now(),
		Units:      units,
	}, nil
}

func (h *Handlers) handleCreateGet(w http.ResponseWriter, r *http.Request) {
	context, err := h.getContext()
	if err != nil {
		h.handleErrors(w, r, err, http.StatusInternalServerError)
		return
	}
	h.contextHandler(w, r, context, "form.html")
}

func (h *Handlers) handleCreatePost(w http.ResponseWriter, r *http.Request, repetition *lifting.Repetition) {
	var (
		err error
	)

	if repetition == nil {
		repetition = &lifting.Repetition{}

	}
	r.ParseForm()

	required := map[string]bool{}

	for key, value := range r.Form {
		if required[key] && len(value) < 1 {
			// todo: send you back to the same url with an error message
			http.Error(w, "Missing required field",
				http.StatusBadRequest)
			return
		}

		match := func(expected string) bool {
			return key == expected && len(value) > 0 && len(value[0]) > 0
		}

		switch {
		case match(category):
			repetition.Category = value[0]
		case match(sessionDate):
			repetition.SessionDate, err = civil.ParseDate(value[0])
		case match(exercise):
			repetition.Exercise = value[0]
		case match(volume):
			repetition.Volume, err = strconv.ParseFloat(value[0], 64)
		case match(weight):
			repetition.Weight, err = strconv.Atoi(value[0])
		case match(hour):
			repetition.Elapsed.Hour, err = strconv.Atoi(value[0])
		case match(minute):
			repetition.Elapsed.Minute, err = strconv.Atoi(value[0])
		case match(second):
			repetition.Elapsed.Second, err = strconv.Atoi(value[0])
		case match(sets):
			repetition.Sets, err = strconv.Atoi(value[0])
		case match(comment):
			repetition.Comment = value[0]
		case match(id):
			var ID int
			ID, err = strconv.Atoi(value[0])
			repetition.ID = &ID
		}

		if err != nil {
			h.handleErrors(w, r, err, http.StatusBadRequest)
			return
		}
	}

	reps := make([]lifting.Repetition, 1)
	reps[0] = *repetition
	err = h.Storage.Load(reps)

	if err != nil {
		h.handleErrors(w, r, err, http.StatusBadRequest)
		return
	}
	http.Redirect(w, r, "/", 301)

}

func (h *Handlers) contextHandler(w http.ResponseWriter, r *http.Request, context *Context, t string) {
	templates, err := template.ParseFiles(
		fmt.Sprintf("templates/%s", t),
		"templates/base.html",
	)

	if err != nil {
		h.handleErrors(w, r, err, http.StatusInternalServerError)
		return
	}
	err = templates.ExecuteTemplate(w, t, context)

	if err != nil {
		h.handleErrors(w, r, err, http.StatusInternalServerError)
		return
	}

	return
}

func (h *Handlers) handleErrors(w http.ResponseWriter, r *http.Request, err error, code int) {
	log.Printf("returning error response because %v", fmt.Errorf(err.Error()))
	http.Error(w, err.Error(), code)
	return

}
