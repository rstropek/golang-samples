package data

import (
	"database/sql"
	"errors"
)

// Error returned when looking up a hero that doesn't exist in our database.
var (
	ErrRecordNotFound = errors.New("record not found")
)

type HeroesRepository interface {
	Insert(hero *Hero) error
	Get(id int64) (*Hero, error)
	Update(hero *Hero) error
	Delete(id int64) error
	GetAll(name string, abilities []string, filters Filters) ([]*Hero, error)
}

type Models struct {
	Heroes HeroesRepository
}

func NewModels(db *sql.DB) Models {
	return Models{
		Heroes: HeroModel{DB: db},
	}
}
