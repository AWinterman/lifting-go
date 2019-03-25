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

	result := group(reps)
	twentyfifth, present := result[civil.Date{Year: 2018, Month: 12, Day: 25}]

	if !present {
		t.Fatal("no twenty fifth", result)
	}

	squat, present := twentyfifth["squat"]

	if !present {
		t.Fatal("no squat", twentyfifth)
	}

	if len(squat) != 1 {
		t.Fatal("wrong number of squats", squat)
	}

	overhead, present := twentyfifth["overhead press"]

	if !present {
		t.Fatal("no overhead", twentyfifth)
	}

	if len(overhead) != 1 {
		t.Fatal("wrong number of overheads", overhead)
	}

}
