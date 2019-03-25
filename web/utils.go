package web

import (
	"strconv"
	"cloud.google.com/go/civil"
	"github.com/awinterman/lifting"
)

func parseInt(value string) (int, error) {
	if len(value) < 1 {
		return 0, nil
	}
	i, errs := strconv.Atoi(value)

	return i, errs
}

func group(reps []lifting.Repetition) map[civil.Date]map[string][]lifting.Repetition {
	m := make(map[civil.Date]map[string][]lifting.Repetition)

	for _, rep := range reps {
		date, present := m[rep.SessionDate] 
		if !present {
			m[rep.SessionDate] = make(map[string][]lifting.Repetition)
			m[rep.SessionDate][rep.Exercise] = []lifting.Repetition{rep}
		} else {
			exercise, present := date[rep.Exercise]
			if !present {
				m[rep.SessionDate][rep.Exercise] = []lifting.Repetition{rep}
			} else {
				m[rep.SessionDate][rep.Exercise] = append(exercise, rep)
			}
		}
	}

	return m
}
