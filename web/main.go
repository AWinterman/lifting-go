package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"regexp"
	"time"
	"cloud.google.com/go/civil"
	"github.com/go-playground/form"
	"github.com/awinterman/lifting"
)

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
	log.Printf("Listening http://localhost %v", port)
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

	context, err := h.getContext()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var edit = regexp.MustCompile(`/edit/(?P<ID>\d\d*)/`) // Contains "abc"

	log.Println(path, path == "/")

	switch {
	case path == "/" || path == "":
		log.Println("matched index handler")

		h.contextHandler(w, r, context, "index.html")
		break
	case ((path == "/create/" || path == "/create") && r.Method == "POST") || (edit.MatchString(path) && r.Method == "PUT"):
		var repetition lifting.Repetition
		r.ParseForm()

		err = h.Decoder.Decode(&repetition, r.Form)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		 reps := make([]lifting.Repetition, 1)
		 reps[0] = repetition
		 err := h.Storage.Load(reps)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		 http.Redirect(w, r, "/index/", 301)


	case path == "/create/" || path == "/create":
		log.Println("matched create handler")
		h.contextHandler(w, r, context, "form.html")
		break
	case edit.MatchString(path):
		matches := edit.FindStringSubmatch(path)
		var IDString interface{} = matches[1]
		ID, ok := IDString.(int)
		if !ok {
			context.Message = "Specified ID is not valid"
			h.contextHandler(w, r, context, "something-bad.html")
			return
		}

		repetition, err := h.Storage.GetByID(ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		context.Repetition = &repetition

		h.contextHandler(w, r, context, "form.html")
		break
	default:

		var templates = template.Must(template.ParseFiles("templates/404.html", "templates/base.html"))
		err := templates.ExecuteTemplate(w, "404.html", Context{})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		break

	}

}

func (h *Handlers) getContext() (*Context, error) {
	reps, err := h.Storage.GetLast(10, 0)
	if err != nil {
		return nil, err
	}

	categories, err := h.Storage.GetUniqueCategories()
	log.Printf("categories: %v", categories)

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

func (h *Handlers) contextHandler(w http.ResponseWriter, r *http.Request, context *Context, t string) {
	templates, err := template.ParseFiles(
		fmt.Sprintf("templates/%s", t),
		"templates/base.html",
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = templates.ExecuteTemplate(w, t, context)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	return
}
