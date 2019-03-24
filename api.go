package lifting

import (
	"cloud.google.com/go/civil"
)

// Storage is an interface for the storage class
type Storage interface {
	Load(repetitions []Repetition) error
	GetLast(count, offset int) ([]Repetition, error)
	GetByID(id int) (*Repetition, error)
	GetBetween(start, end civil.Date) ([]Repetition, error)
	GetUniqueCategories() ([]string, error)
	GetByCategory(label string, count, offset int) ([]Repetition, error)
	GetUniqueExercises() ([]string, error)
	GetUniqueUnits() ([]string, error)
}
