package lifting

import (
	"testing"

	"cloud.google.com/go/civil"
)

func TestGrouping(t *testing.T) {
	duration := civil.Time{Minute: 30}

	reps := []Repetition{
		Repetition{
			Exercise:    "run",
			Effort:      70,
			Volume:      2,
			Weight:      0,
			SessionDate: civil.Date{Year: 2018, Month: 12, Day: 20},
			Units:       "miles",
			Elapsed:     duration,
			Failure:     false,
			Category:    "aerobic/recovery",
		},
		Repetition{
			Exercise:    "squat",
			Effort:      70,
			Volume:      5,
			Weight:      180,
			SessionDate: civil.Date{Year: 2018, Month: 12, Day: 25},
			Units:       "lbs",
			Failure:     false,
			Category:    "strength",
		},
		Repetition{
			Exercise:    "squat",
			Effort:      70,
			Volume:      5,
			Weight:      180,
			SessionDate: civil.Date{Year: 2018, Month: 12, Day: 26},
			Units:       "lbs",
			Failure:     false,
			Category:    "strength",
		},
		Repetition{
			Exercise:    "squat",
			Effort:      70,
			Volume:      5,
			Weight:      180,
			SessionDate: civil.Date{Year: 2018, Month: 12, Day: 27},
			Units:       "lbs",
			Failure:     false,
			Category:    "strength",
		},
		Repetition{
			Exercise: "yoga",
			Category: "mobility",
		},
		Repetition{
			Exercise:    "overhead press",
			Effort:      90,
			Volume:      5,
			Weight:      95,
			SessionDate: civil.Date{Year: 2018, Month: 12, Day: 25},
			Units:       "lbs",
			Failure:     false,
			Category:    "strength",
		},
		Repetition{
			Exercise:    "overhead press",
			Effort:      90,
			Volume:      5,
			Weight:      95,
			SessionDate: civil.Date{Year: 2018, Month: 12, Day: 26},
			Units:       "lbs",
			Failure:     false,
			Category:    "strength",
		},
		Repetition{
			Exercise:    "overhead press",
			Effort:      90,
			Volume:      5,
			Weight:      95,
			SessionDate: civil.Date{Year: 2018, Month: 12, Day: 27},
			Units:       "lbs",
			Failure:     false,
			Category:    "strength",
		},
	}

	group(reps)



}
