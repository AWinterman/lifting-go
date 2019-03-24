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
)

func now() string {
	return civil.DateOf(time.Now()).String()
}

// Page defines pagination.
type Page struct {
	Offset int
	Count  int
}

func (p *Page) Next() Page {
	return Page{
		Offset: p.Offset + p.Count,
		Count:  p.Count,
	}
}

func (p *Page) Previous() Page {
	offset := p.Offset - p.Count
	if offset < 0 {
		offset = 0
	}
	return Page{
		Offset: offset,
		Count:  p.Count,
	}
}

// Context is basic context for the site.
type Context struct {
	History    []lifting.Repetition
	Categories []string
	Exercises  []string
	Units      []string
	Repetition *lifting.Repetition
	Now        string
	Message    string
	Next       Page
	Previous       Page
	Current       Page
	CanGoLater bool
	CanGoEarlier bool
}

func (h *Handlers) getContext(page Page) (*Context, error) {
	reps, err := h.Storage.GetLast(page.Count, page.Offset)
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
		Next:       page.Next(),
		Previous:       page.Previous(),
		Current: page,
		CanGoLater: page.Offset > 0,
		CanGoEarlier: len(reps) == page.Count,
	}, nil
}

var edit = regexp.MustCompile(`/edit/(?P<ID>\d\d*)(/)?`)
var copy = regexp.MustCompile(`/copy/(?P<ID>\d\d*)(/)?`)
var delete = regexp.MustCompile(`/delete/(?P<ID>\d\d*)(/)?`)

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
	units       = "Units"
)

// Handlers is all the http handlers
type Handlers struct {
	Storage lifting.Storage
	step    int
}

func (h *Handlers) handle(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "Text/HTML")
	path := r.URL.Path

	switch {
	case path == "/" || path == "":
		log.Println("matched index")
		h.index(w, r)
	case (path == "/create/" || path == "/create"):
		h.handleCreate(w, r)
	case edit.MatchString(path):
		log.Println("matched edit")
		h.handleEdit(w, r)
	case copy.MatchString(path):
		log.Println("matched copy")
		h.handleCopy(w, r)
	case delete.MatchString(path):
		h.handleDelete(w, r)
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

func (h *Handlers) handleErrors(w http.ResponseWriter, r *http.Request, err error, code int) {
	log.Printf("returning error response because %v", fmt.Errorf(err.Error()))
	http.Error(w, err.Error(), code)
	return

}

func (h *Handlers) contextHandler(w http.ResponseWriter, r *http.Request, context interface{}, t string) {
	templates, err := template.ParseFiles(
		fmt.Sprintf("templates/%s", t),
		"templates/table.html",
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



func (h *Handlers) getPage(r *http.Request) (Page, error) {
	var err error
	page := Page{Count: h.step, Offset: 0}

	query := r.URL.Query()

	countStringArray := query["count"]
	offsetStringArray := query["offset"]

	if len(countStringArray) > 0 {
		page.Count, err = parseInt(countStringArray[0])
		if err != nil {
			return page, err
		}
	}

	if len(offsetStringArray) > 0 {
		page.Offset, err = parseInt(offsetStringArray[0])
		if err != nil {
			return page, err
		}
	}

	return page, nil

}

func (h *Handlers) index(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		h.handleErrors(w, r, fmt.Errorf("Unexpected method"), http.StatusMethodNotAllowed)
		return
	}

	page, err := h.getPage(r)
	if err != nil {
		h.handleErrors(w, r, err, http.StatusInternalServerError)
		return
	}

	context, err := h.getContext(page)
	if err != nil {
		h.handleErrors(w, r, err, http.StatusInternalServerError)
		return
	}
	h.contextHandler(w, r, context, "index.html")

}

func (h *Handlers) handleDelete(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	matches := delete.FindStringSubmatch(path)

	repetition, err := h.getRep(matches[1])
	if err != nil {
		h.handleErrors(w, r, err, http.StatusInternalServerError)
		return
	}

	if r.Method == "GET" {
		h.contextHandler(w, r, repetition, "delete.html")
		return
	} else if r.Method == "POST" {
		err = h.Storage.Delete(*repetition.ID)

		if err != nil {
			h.handleErrors(w, r, err, http.StatusBadRequest)
			return
		}
		http.Redirect(w, r, "/", 301)
	}

}

func (h *Handlers) handleCopy(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	matches := copy.FindStringSubmatch(path)

	repetition, err := h.getRep(matches[1])
	if err != nil {
		h.handleErrors(w, r, err, http.StatusInternalServerError)
		return
	}

	// nil it so when it gets sent back we make a new one.
	repetition.ID = nil

	if r.Method == "GET" {
		page, err := h.getPage(r)
		if err != nil {
			h.handleErrors(w, r, err, http.StatusInternalServerError)
			return
		}
		context, err := h.getContext(page)
		if err != nil {
			h.handleErrors(w, r, err, http.StatusInternalServerError)
			return
		}

		context.Repetition = repetition

		h.contextHandler(w, r, context, "form.html")
		return
	} else if r.Method == "POST" {
		log.Println("call new creation script with copy")
		h.handleCreatePost(w, r, repetition)
	}

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

		page, err := h.getPage(r)
		if err != nil {
			h.handleErrors(w, r, err, http.StatusInternalServerError)
			return
		}
		context, err := h.getContext(page)
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

func (h *Handlers) handleCreate(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		h.handleCreatePost(w, r, nil)
	case "GET":
		h.handleCreateGet(w, r)
	default:
		h.handleErrors(w, r, fmt.Errorf("Unexpected method"), http.StatusMethodNotAllowed)
	}
}

func (h *Handlers) handleCreateGet(w http.ResponseWriter, r *http.Request) {

	page, err := h.getPage(r)
	if err != nil {
		h.handleErrors(w, r, err, http.StatusInternalServerError)
		return
	}
	context, err := h.getContext(page)
	if err != nil {
		h.handleErrors(w, r, err, http.StatusInternalServerError)
		return
	}
	h.contextHandler(w, r, context, "form.html")
}

func (h *Handlers) handleCreatePost(w http.ResponseWriter, r *http.Request, existing *lifting.Repetition) {
	var (
		err error
	)

	repetition := &lifting.Repetition{}
	if existing != nil {
		repetition = existing
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

		if len(value) > 0 {
			log.Println(key, value)
			switch key {
			case (category):
				repetition.Category = value[0]
			case (sessionDate):
				repetition.SessionDate, err = civil.ParseDate(value[0])
			case (exercise):
				repetition.Exercise = value[0]
			case (volume):
				repetition.Volume, err = strconv.ParseFloat(value[0], 64)
			case (weight):
				repetition.Weight, err = parseInt(value[0])
			case (hour):
				repetition.Elapsed.Hour, err = parseInt(value[0])
			case (minute):
				repetition.Elapsed.Minute, err = parseInt(value[0])
			case (second):
				repetition.Elapsed.Second, err = parseInt(value[0])
			case (sets):
				repetition.Sets, err = parseInt(value[0])
			case (comment):
				repetition.Comment = value[0]
			case (effort):
				repetition.Effort, err = parseInt(value[0])
			case (id):
				var ID int
				ID, err = parseInt(value[0])
				repetition.ID = &ID
			case (units):
				repetition.Units = value[0]

			default:
				log.Panicln("unknown field ", key)
			}

			if err != nil {
				h.handleErrors(w, r, err, http.StatusBadRequest)
				return
			}
		} else {
			log.Printf("skipping %s because no value", key)
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
