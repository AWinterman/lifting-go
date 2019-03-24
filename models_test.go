package lifting

import (
	"cloud.google.com/go/civil"
	"database/sql"
	"fmt"
	"testing"
)

func TestConversion(t *testing.T) {
	r := Repetition{
		Exercise:    "squat",
		Effort:      70,
		Volume:      5,
		Weight:      180,
		SessionDate: civil.Date{Year: 2018, Month: 12, Day: 26},
		Units:       "lbs",
		Failure:     false,
		Category:    "strength",
	}

	expected := WorkoutRow{
		Exercise:    "squat",
		Effort:      sql.NullInt64{Int64: 70, Valid: true},
		Volume:      sql.NullFloat64{Float64: 5, Valid: true},
		Weight:      sql.NullInt64{Int64: 180, Valid: true},
		SessionDate: "2018-12-26",
		Units:       "lbs",
		Failure:     false,
		Category:    sql.NullString{String: "strength", Valid: true},
	}

	workout := repetitionToWorkout(r)

	if workout != expected {
		t.Fatal("mimsatch", fmt.Sprintf("expected %#v", expected), fmt.Sprintf("found %#v", workout))
	}

	back, err := workoutToRepetition(workout)

	if err != nil {
		t.Fatal(err)
	}

	if back != r {
		t.Fatal("mimsatch", fmt.Sprintf("expected %#v", r), fmt.Sprintf("found %#v", back))
	}
}
