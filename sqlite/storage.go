package sqlite

import (
	"fmt"

	"cloud.google.com/go/civil"
	"github.com/jmoiron/sqlx"
	"github.com/awinterman/lifting"
	_ "github.com/mattn/go-sqlite3" // how they told me to do it i guess
)

const (
	workoutSchema = `
            CREATE TABLE IF NOT EXISTS workout (
               id integer primary key,
               exercise varchar NOT NULL,
               effort decimal,
               volume int,
               weight decimal,
               duration interval,
               session_date date NOT NULL,
               failure boolean default false,
               units str NOT NULL,
			   category text,
			   comment text
            );
        `

	drop = `
            DROP TABLE IF EXISTS workout;
        `
	namedInsert = `INSERT INTO workout(
            exercise, effort, volume, weight, duration, session_date, units, failure, category, comment
            ) values (
            :exercise, :effort, :volume, :weight, :duration, :session_date, :units, :failure, :category, :comment
			)`

	//
	namedDelete = `DELETE FROM workout WHERE id = :id`
	namedUpdate = `UPDATE workout
			SET exercise = :exercise,
				effort = :effort,
				volume = :volume,
				 weight = :weight, 
				 duration = :duration, 
				 session_date = :session_date, 
				 units = :units, 
				 failure = :failure, 
				 category = :category,
				 comment = :comment
			WHERE
				id = :id

            `

	uniquecategory = `SELECT DISTINCT category FROM workout`
	uniqueExercise = `SELECT DISTINCT exercise FROM workout`
	uniqueUnits    = `SELECT DISTINCT units FROM workout WHERE units != ""`

	getlast = `
            SELECT 
            id, exercise, effort, volume, weight, duration, session_date, units, failure, category, comment
            FROM workout 
            ORDER BY session_date desc, id desc LIMIT ? OFFSET ?`
	getBetween = `
            SELECT 
            id, exercise, effort, volume, weight, duration, session_date, units, failure, category, comment
            FROM workout WHERE session_date BETWEEN ? and ? 
            ORDER BY session_date DESC, id DESC`
	getByID = `
            SELECT 
            id, exercise, effort, volume, weight, duration, session_date, units, failure, category, comment
            FROM workout WHERE id = ?`
	getByCategory = `
			WITH vars AS (SELECT :category as category)
            SELECT 
				MAX(id) as id,
				exercise, 
				MIN(effort) as effort,
				MAX(volume) as volume,
				MAX(weight) as weight,
				MAX(duration) as duration,
				MAX(session_date) as session_date,
				count(*) as sets,
				units,
				false as failure,
				workout.category as category
			FROM workout INNER JOIN vars ON(workout.category LIKE '%'||vars.category||'%')
			GROUP BY exercise, workout.category, units
			ORDER BY 
				workout.category = vars.category, workout.category like vars.category||'%', workout.category like '%'||vars.category||'%',
				session_date DESC, 
				id DESC 
			LIMIT :count OFFSET :offset`
)

// SqliteStorage is a sqlite implementation of the Storage interface
type SqliteStorage struct {
	Path string
	db   *sqlx.DB
}

// CreateStorage sets up the database resources
func CreateStorage(Path string, db *sqlx.DB) (*SqliteStorage, error) {
	var s = SqliteStorage{Path: Path, db: db}

	if s.db == nil {
		db, err := sqlx.Connect("sqlite3", s.Path)
		s.db = db
		s.db.Ping()
		if err != nil {
			return nil, err
		}
		stmt, err := s.db.Prepare(workoutSchema)
		if err != nil {
			return nil, err
		}
		_, err = stmt.Exec()
		if err != nil {
			return nil, err
		}
		return &s, stmt.Close()
	}
	return &s, nil
}

// Drop drops the database
func (s *SqliteStorage) Drop() error {
	_, err := s.db.Exec(drop)
	return err
}

//Load the repetitions into the database
func (s *SqliteStorage) Load(repetitions []lifting.Repetition) error {
	tx, err := s.db.Beginx()
	if err != nil {
		return err
	}

	for _, rep := range repetitions {
		workout := lifting.RepetitionToWorkout(rep)
		var err error

		if workout.ID == nil {
			_, err = tx.NamedExec(
				namedInsert,
				&workout,
			)
		} else {
			_, err = tx.NamedExec(
				namedUpdate,
				&workout,
			)

		}
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}

// Delete removes the corresponding database row
func (s *SqliteStorage) Delete(id int) error {
	tx, err := s.db.Beginx()
	if err != nil {
		return err
	}

	workout := lifting.WorkoutRow{ID: &id}

	_, err = tx.NamedExec(
		namedDelete,
		workout,
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (s *SqliteStorage) getCollectionWithStruct(query string, arg interface{}) ([]lifting.Repetition, error) {
	var (
		rs []lifting.Repetition
		w  lifting.WorkoutRow
		r  lifting.Repetition
	)

	rows, err := s.db.NamedQuery(query, arg)
	if err != nil {
		panic(err)
	}

	for rows.Next() {
		err = rows.StructScan(&w)
		if err != nil {
			return rs, err
		}
		r, err = lifting.WorkoutToRepetition(w)
		rs = append(rs, r)
		if err != nil {
			return rs, err
		}

	}
	return rs, nil
}

func (s *SqliteStorage) getCollection(query string, args ...interface{}) ([]lifting.Repetition, error) {
	var rs []lifting.Repetition
	ws := []lifting.WorkoutRow{}
	err := s.db.Select(&ws, query, args...)
	if err != nil {
		panic(err)
	}
	rs = make([]lifting.Repetition, len(ws))
	for i, w := range ws {
		rs[i], err = lifting.WorkoutToRepetition(w)
		if err != nil {
			return rs, err
		}
	}
	return rs, nil
}

//GetUniqueCategories retrieves what categories have been input
func (s *SqliteStorage) GetUniqueCategories() ([]string, error) {
	categorys := make([]string, 0)
	err := s.db.Select(&categorys, uniquecategory)
	if err != nil {
		return categorys, err
	}
	return categorys, nil
}

//GetUniqueUnits retrieves what categories have been input
func (s *SqliteStorage) GetUniqueUnits() ([]string, error) {
	categorys := make([]string, 0)
	err := s.db.Select(&categorys, uniqueUnits)
	if err != nil {
		return categorys, err
	}
	return categorys, nil
}

//GetUniqueExercises retrieves what Exercises have been input
func (s *SqliteStorage) GetUniqueExercises() ([]string, error) {
	r := make([]string, 0)
	err := s.db.Select(&r, uniqueExercise)
	if err != nil {
		return r, err
	}
	return r, nil
}

//GetByCategory retrieves reps in a given category
func (s *SqliteStorage) GetByCategory(category string, count, offset int) ([]lifting.Repetition, error) {
	return s.getCollectionWithStruct(getByCategory, lifting.CategoryQuery{
		Category: category, Count: count, Offset: offset,
	})
}

// GetLast retrieves data in order
func (s *SqliteStorage) GetLast(count, offset int) ([]lifting.Repetition, error) {
	return s.getCollection(getlast, count, offset)
}

// GetByID finds a particular repetition
func (s *SqliteStorage) GetByID(id int) (*lifting.Repetition, error) {
	reps, err := s.getCollection(getByID, id)

	if err != nil {
		return nil, err
	}

	if len(reps) == 0 {
		return nil, nil
	}

	if len(reps) > 1 {
		return nil, fmt.Errorf("Multiple return values for SqliteStorage#GetByID, %v", reps)
	}

	return &reps[0], nil
}

// GetBetween returns the reps between the start and end date.
func (s *SqliteStorage) GetBetween(start, end civil.Date) ([]lifting.Repetition, error) {
	return s.getCollection(getBetween, start, end)
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
