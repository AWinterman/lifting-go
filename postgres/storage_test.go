package postgres

import (
	"fmt"
	"testing"

	"cloud.google.com/go/civil"
	"github.com/awinterman/lifting"
)

func TestLiftingStorage(t *testing.T) {
	storage, err := CreateStorage("user=testing dbname=test_lifting password=testing", nil)
	//defer storage.Drop()
	if err != nil {
		t.Error(err)
	}
	duration := civil.Time{Minute: 30}
	if err != nil {
		t.Error(err)
	}

	reps := []lifting.Repetition{
		lifting.Repetition{
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
		lifting.Repetition{
			Exercise:    "squat",
			Effort:      70,
			Volume:      5,
			Weight:      180,
			SessionDate: civil.Date{Year: 2018, Month: 12, Day: 26},
			Units:       "lbs",
			Failure:     false,
			Category:    "strength",
		},
		lifting.Repetition{
			Exercise:    "squat",
			Effort:      70,
			Volume:      5,
			Weight:      180,
			SessionDate: civil.Date{Year: 2018, Month: 12, Day: 26},
			Units:       "lbs",
			Failure:     false,
			Category:    "strength",
		},
		lifting.Repetition{
			Exercise:    "squat",
			Effort:      70,
			Volume:      5,
			Weight:      180,
			SessionDate: civil.Date{Year: 2018, Month: 12, Day: 26},
			Units:       "lbs",
			Failure:     false,
			Category:    "strength",
		},
		lifting.Repetition{
			SessionDate: civil.Date{Year: 2018, Month: 12, Day: 26},
			Exercise: "yoga",
			Category: "mobility",
		},
		lifting.Repetition{
			Exercise:    "overhead press",
			Effort:      90,
			Volume:      5,
			Weight:      95,
			SessionDate: civil.Date{Year: 2018, Month: 12, Day: 26},
			Units:       "lbs",
			Failure:     false,
			Category:    "strength",
		},
		lifting.Repetition{
			Exercise:    "overhead press",
			Effort:      90,
			Volume:      5,
			Weight:      95,
			SessionDate: civil.Date{Year: 2018, Month: 12, Day: 26},
			Units:       "lbs",
			Failure:     false,
			Category:    "strength",
		},
		lifting.Repetition{
			Exercise:    "overhead press",
			Effort:      90,
			Volume:      5,
			Weight:      95,
			SessionDate: civil.Date{Year: 2018, Month: 12, Day: 26},
			Units:       "lbs",
			Failure:     false,
			Category:    "strength",
		},
	}
	err = storage.Load(reps)
	if err != nil {
		t.Errorf("loading %v failed %v", reps, err)
		return
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

	rep, err := storage.GetByID(*last[0].ID)

	if err != nil {
		t.Fatal(err)
	}

	if rep == nil {
		t.Fatal("mimsatch",
			fmt.Sprintf("expected %#v", last[0]),
			fmt.Sprintf("found %#v", rep),
		)

	}

	last[0].ID = nil
	rep.ID = nil

	if *rep != last[0] {
		t.Fatal("mimsatch",
			fmt.Sprintf("expected %#v", last[0]),
			fmt.Sprintf("found %#v", rep),
		)

	}

	last[0].ID = nil
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

	// lastStrength, err := storage.GetByCategory("aerobic", 1, 0)
    //
	// if err != nil {
	// 	t.Fatal(err)
	// 	return
	// }
    //
	// if len(lastStrength) == 0 {
	// 	t.Error("nothing returned by get getByCategory")
	// 	return
	// }
	// lastStrength[0].ID = nil
	// expected = reps[0]
	// expected.Sets = 1
	// if expected != lastStrength[0] {
	// 	t.Fatal("mimsatch",
	// 		fmt.Sprintf("expected %#v", expected),
	// 		fmt.Sprintf("found %#v", lastStrength[0]),
	// 	)
	// }

	units, err := storage.GetUniqueUnits()

	if err != nil {
		t.Fatal(err)
		return
	}

	if len(units) != 2 {
		t.Fatal("mimsatch",
			fmt.Sprintf("expected %#v", 2),
			fmt.Sprintf("found %#v", len(units)),
		)
	}


	if units[1] != "miles" {
		t.Fatal("mimsatch",
			fmt.Sprintf("expected %#v", "miles"),
			fmt.Sprintf("found %#v", units[0]),
		)
	}

	if units[0] != "lbs" {
		t.Fatal("mimsatch",
			fmt.Sprintf("expected %#v", "lbs"),
			fmt.Sprintf("found %#v", units[1]),
		)
	}

}
