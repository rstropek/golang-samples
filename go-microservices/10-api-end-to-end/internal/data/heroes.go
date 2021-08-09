package data

import (
	"time"
)

type Hero struct {
	ID        int64     `json:"id"`             // Unique integer ID for a hero
	CreatedAt time.Time `json:"-"`              // Timestamp for when the hero is added to our list of heroes
	Name      string    `json:"name"`           // Name of hero
	RealName  string    `json:"realName"`       // Hero's real name
	Coolness  int32     `json:"coolness"`       // Coolness factor of hero
	Tags      []string  `json:"tags,omitempty"` // Slice of tags for the hero
	CanFly    bool      `json:"canFly"`         // Indicates whether the hero can fly
}
