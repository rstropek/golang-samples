package data

import (
	"time"

	"github.com/rstropek/golang-samples/api-end-to-end/internal/validator"
)

type Hero struct {
	ID        int64     `json:"id"`             // Unique integer ID for a hero
	CreatedAt time.Time `json:"-"`              // Timestamp for when the hero is added to our list of heroes
	Name      string    `json:"name"`           // Name of hero
	RealName  string    `json:"realName"`       // Hero's real name
	Coolness  int32     `json:"coolness"`       // Coolness factor of hero
	Tags      []string  `json:"tags,omitempty"` // Slice of tags for the hero
	CanFly    CanFly    `json:"canFly"`         // Indicates whether the hero can fly
}

func ValidateHero(v *validator.Validator, movie *Hero) {
    v.Check(movie.Name != "", "name", "must not be empty")
    v.Check(len(movie.Name) <= 100, "name", "must not be more than 100 bytes long")

    v.Check(movie.Coolness >= 0 && movie.Coolness <= 9, "coolness", "must be between 0 and 9")

    v.Check(movie.Tags != nil, "tags", "must be provided")
    v.Check(len(movie.Tags) >= 1, "tags", "must contain at least 1 tag")
    v.Check(len(movie.Tags) <= 5, "tags", "must not contain more than 5 tags")
    v.Check(validator.Unique(movie.Tags), "tags", "must not contain duplicate values")
}
