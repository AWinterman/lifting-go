package lifting

import (
	"cloud.google.com/go/civil"
	"fmt"
	"testing"
)

func TestSqlite(t *testing.T) {
	storage, err := CreateStorage("test.db", nil)
	defer storage.Drop()
	if err != nil {
		t.Error(err)
	}
	duration := civil.Time{Minute: 30}
	if err != nil {
		t.Error(err)
	}

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
			Category:       "aerobic/recovery",
		},
		Repetition{
			Exercise:    "squat",
			Effort:      70,
			Volume:      5,
			Weight:      180,
			SessionDate: civil.Date{Year: 2018, Month: 12, Day: 26},
			Units:       "lbs",
			Failure:     false,
			Category:       "strength",
		},
		Repetition{
			Exercise:    "squat",
			Effort:      70,
			Volume:      5,
			Weight:      180,
			SessionDate: civil.Date{Year: 2018, Month: 12, Day: 26},
			Units:       "lbs",
			Failure:     false,
			Category:       "strength",
		},
		Repetition{
			Exercise:    "squat",
			Effort:      70,
			Volume:      5,
			Weight:      180,
			SessionDate: civil.Date{Year: 2018, Month: 12, Day: 26},
			Units:       "lbs",
			Failure:     false,
			Category:       "strength",
		},
		Repetition{
			Exercise:    "yoga",
			Category:       "mobility",
		},
		Repetition{
			Exercise:    "overhead press",
			Effort:      90,
			Volume:      5,
			Weight:      95,
			SessionDate: civil.Date{Year: 2018, Month: 12, Day: 26},
			Units:       "lbs",
			Failure:     false,
			Category:       "strength",
		},
		Repetition{
			Exercise:    "overhead press",
			Effort:      90,
			Volume:      5,
			Weight:      95,
			SessionDate: civil.Date{Year: 2018, Month: 12, Day: 26},
			Units:       "lbs",
			Failure:     false,
			Category:       "strength",
		},
		Repetition{
			Exercise:    "overhead press",
			Effort:      90,
			Volume:      5,
			Weight:      95,
			SessionDate: civil.Date{Year: 2018, Month: 12, Day: 26},
			Units:       "lbs",
			Failure:     false,
			Category:       "strength",
		},
	}
	err = storage.Load(reps)
	if err != nil {
		t.Error(err)
	}

	last, err := storage.GetLast(1, 0)

	if err != nil {
		t.Fatal(err)
		return
	}

	expected := reps[len(reps)-1]
	if len(last) == 0 {
		t.Error("nothing returned by get last")
		return
	}
	last[0].ID = 0
	if expected != last[0] {
		t.Fatal("mimsatch",
			fmt.Sprintf("expected %#v", expected),
			fmt.Sprintf("found %#v", last[0]),
		)
	}

	Categorys, err := storage.GetUniqueCategories()
	if Categorys[0] != "aerobic/recovery" {
		t.Fatal("mimsatch",
			fmt.Sprintf("expected %#v", "aerobic/recovery"),
			fmt.Sprintf("found %#v", Categorys[0]),
		)
	}
	if Categorys[1] != "strength" {
		t.Fatal("mimsatch",
			fmt.Sprintf("expected %#v", "strength"),
			fmt.Sprintf("found %#v", Categorys[1]),
		)
	}

	lastStrength, err := storage.GetByCategory("aerobic", 1, 0)

	if err != nil {
		t.Fatal(err)
		return
	}

	if len(lastStrength) == 0 {
		t.Error("nothing returned by get getByCategory")
		return
	}
	lastStrength[0].ID = 0
	expected = reps[0]
	expected.Sets = 1
	if expected != lastStrength[0] {
		t.Fatal("mimsatch",
			fmt.Sprintf("expected %#v", expected),
			fmt.Sprintf("found %#v", lastStrength[0]),
		)
	}

	units, err := storage.GetUniqueUnits()

	if err != nil {
		t.Fatal(err)
		return
	}

	if (units[0] != "miles") {
		t.Fatal("mimsatch",
			fmt.Sprintf("expected %#v", "miles"),
			fmt.Sprintf("found %#v", units[0]),
		)
	}

	if (units[1] != "lbs") {
		t.Fatal("mimsatch",
			fmt.Sprintf("expected %#v", "lbs"),
			fmt.Sprintf("found %#v", units[1]),
		)
	}

	if len(units) != 2 {
		t.Fatal("mimsatch",
			fmt.Sprintf("expected length %#v", 2),
			fmt.Sprintf("found %#v", len(units)),
		)
	}

}
