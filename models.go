package lifting

import (
	"cloud.google.com/go/civil"
	"database/sql"
	"time"
)

type (
	// Repetition represents a single repetition of an exercise, e.g. squats, a run, etc.
	// we store like this so we can easily represent failure on a given set, and so
	// we have the flexibility to represent whatever you might do with our body.
	Repetition struct {
		ID *int
		// exercise represents whatever exercise you completed
		Exercise string
		// What day did you do the activity?
		SessionDate civil.Date
		// how do you measure what you did? Miles, pounds, kg?
		Units string
		// did you fail in attempting the exercise?
		Failure bool
		// Effor on a scale from 0 - 100, 0 being asleep, 30 being moderate, 50 is hard,
		// 70 is very hard, 100 is almost failed
		Effort int
		// Volume how much did you do? 1 mile? 8 Squats? etc.
		Volume float64
		// Did you add weight to the exercise? If so how much?
		Weight int
		// How long did the exercise take you, if it is relevant?
		Elapsed civil.Time
		// What kind of workout it is-- aerobic/recovery, strength, power endurance?
		Category string
		// If non zero, this object indicates sets of the above specifications
		// were performed.
		Sets int
		// anything about the specific workout not captured in other parameters
		Comment string
	}

	// WorkoutRow represents The SQL database format for the Repetition
	WorkoutRow struct {
		ID          *int
		Exercise    string
		Effort      sql.NullInt64
		Volume      sql.NullFloat64
		Weight      sql.NullInt64
		Category    sql.NullString
		Elapsed     sql.NullString `db:"duration"`
		SessionDate string         `db:"session_date"`
		Units       string
		Failure     bool
		Comment     sql.NullString
		Sets        int
	}

	// CategoryQuery represents how we pull by category out of the database
	CategoryQuery struct {
		Category      string
		Count, Offset int
	}
)

// ParseSessionDateString defines how to go from a string to a civil date. If it
// is not successful a non null error is returned
func ParseSessionDateString(date string) (civil.Date, error) {
	asTime, e := time.Parse(time.RFC3339, date)
	if e == nil {
		return civil.DateOf(asTime), e
	}
	return civil.ParseDate(date)
}

// RepetitionToWorkout transforms from a repetition to a workout row.
func RepetitionToWorkout(r Repetition) WorkoutRow {
	effort := sql.NullInt64{Valid: false}
	volume := sql.NullFloat64{Valid: false}
	weight := sql.NullInt64{Valid: false}
	Category := sql.NullString{Valid: false}
	elapsed := sql.NullString{Valid: false}
	comment := sql.NullString{Valid: false}
	if r.Effort != 0 {
		effort = sql.NullInt64{Int64: int64(r.Effort), Valid: true}
	}
	if r.Volume != 0 {
		volume = sql.NullFloat64{Float64: float64(r.Volume), Valid: true}
	}
	if r.Weight != 0 {
		weight = sql.NullInt64{Int64: int64(r.Weight), Valid: true}
	}
	if r.Category != "" {
		Category = sql.NullString{String: r.Category, Valid: true}
	}

	if (r.Elapsed != civil.Time{}) {
		elapsed = sql.NullString{String: r.Elapsed.String(), Valid: true}
	}

	if r.Comment != "" {
		comment = sql.NullString{String: r.Comment, Valid: true}
	}

	return WorkoutRow{
		ID:          r.ID,
		Exercise:    r.Exercise,
		Effort:      effort,
		Volume:      volume,
		Weight:      weight,
		Elapsed:     elapsed,
		SessionDate: r.SessionDate.String(),
		Units:       r.Units,
		Failure:     r.Failure,
		Category:    Category,
		Sets:        r.Sets,
		Comment:     comment,
	}
}

// WorkoutToRepetition converts from workout to repetition.
func WorkoutToRepetition(w WorkoutRow) (Repetition, error) {
	var (
		rep         Repetition
		effort      int
		volume      float64
		weight      int
		sessionDate civil.Date
		elapsed     civil.Time
		Category    string
		err         error
		comment     string
	)

	sessionDate, err = ParseSessionDateString(w.SessionDate)

	if err != nil {
		return rep, err
	}

	if w.Effort.Valid {
		effort = int(w.Effort.Int64)
	}

	if w.Volume.Valid {
		volume = w.Volume.Float64
	}

	if w.Weight.Valid {
		weight = int(w.Weight.Int64)
	}
	if w.Category.Valid {
		Category = w.Category.String
	}

	if w.Elapsed.Valid {
		elapsed, err = civil.ParseTime(w.Elapsed.String)
		if err != nil {
			return rep, err
		}
	}

	if w.Comment.Valid {
		comment = w.Comment.String
	}

	rep = Repetition{
		ID:          w.ID,
		Exercise:    w.Exercise,
		Effort:      effort,
		Volume:      volume,
		Weight:      weight,
		Elapsed:     elapsed,
		SessionDate: sessionDate,
		Units:       w.Units,
		Failure:     w.Failure,
		Category:    Category,
		Sets:        w.Sets,
		Comment:     comment,
	}
	return rep, nil
}
