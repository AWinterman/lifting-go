package lifting

import (
	"sort"
	"strconv"
	"time"

	"cloud.google.com/go/civil"
)

func parseInt(value string) (int, error) {
	if len(value) < 1 {
		return 0, nil
	}
	i, errs := strconv.Atoi(value)

	return i, errs
}

// Category is the category of exercises on a date
type Category struct {
	Category string
	Reps     []Repetition
}

// Group groups by date
type Group struct {
	Date       civil.Date
	Categories []Category
}

// Weekday computes the weekday for the group's Date
func (r *Group) Weekday() string {
	date := time.Date(r.Date.Year, r.Date.Month, r.Date.Day, 0, 0, 0, 0, time.UTC)
	return date.Weekday().String()
}

func mapGroup(reps []Repetition) map[civil.Date]map[string][]Repetition {
	m := make(map[civil.Date]map[string][]Repetition)

	for _, rep := range reps {
		date, present := m[rep.SessionDate]
		if !present {
			m[rep.SessionDate] = make(map[string][]Repetition)
			m[rep.SessionDate][rep.Category] = []Repetition{rep}
		} else {
			exercise, present := date[rep.Category]
			if !present {
				m[rep.SessionDate][rep.Category] = []Repetition{rep}
			} else {
				m[rep.SessionDate][rep.Category] = append(exercise, rep)
			}
		}
	}

	return m

}

// Groups is a collection of Group.
type Groups []Group

func (s Groups) Len() int {
	return len(s)
}

func (s Groups) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s Groups) Less(i, j int) bool {

	if s[i].Date.Year == s[j].Date.Year {
		if s[i].Date.Month > s[j].Date.Month {
			return true
		}

		if s[i].Date.Month == s[j].Date.Month {
			return s[i].Date.Day > s[j].Date.Day
		}
	}

	if s[i].Date.Year > s[j].Date.Year {
		return true
	}

	return false
}

func group(reps []Repetition) Groups {
	m := mapGroup(reps)

	gs := make(Groups, 0)

	for date, exercises := range m {
		g := Group{
			Date:       date,
			Categories: make([]Category, 0),
		}
		for exercise, reps := range exercises {
			g.Categories = append(g.Categories, Category{Category: exercise, Reps: reps})
		}
		gs = append(gs, g)
	}

	sort.Sort(gs)

	return gs

}
