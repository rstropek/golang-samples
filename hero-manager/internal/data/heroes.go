package data

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/lib/pq"
	"heroes.rainerstropek.com/internal/validator"
)

type Hero struct {
	ID        int64     `json:"id"`
	FirstSeen time.Time `json:"firstSeen"`
	Name      string    `json:"name"`
	CanFly    bool      `json:"canFly"`
	RealName  string    `json:"realName,omitempty"`
	Abilities []string  `json:"-"`
	Version   int32     `json:"version"`
}

func (h Hero) MarshalJSON() ([]byte, error) {
	var abilities string

	if h.Abilities != nil {
		abilities = strings.Join(h.Abilities, ", ")
	}

	type HeroAlias Hero

	aux := struct {
		HeroAlias
		Abilities string `json:"abilities,omitempty"`
	}{
		HeroAlias: HeroAlias(h),
		Abilities: abilities,
	}

	return json.Marshal(aux)
}

func ValidateHero(v *validator.Validator, hero *Hero) {
	v.Check(hero.Name != "", "name", "must be provided")
	v.Check(len(hero.Name) <= 100, "name", "must not be more than 100 bytes long")

	v.Check(hero.Abilities != nil, "abilities", "must be provided")
	v.Check(len(hero.Abilities) >= 1, "abilities", "must contain at least 1 ability")
	v.Check(len(hero.Abilities) <= 5, "abilities", "must not contain more than 5 abilities")
	v.Check(validator.Unique(hero.Abilities), "abilities", "must not contain duplicate values")
}

type HeroModel struct {
	DB *sql.DB
}

func (m HeroModel) Insert(hero *Hero) error {
	query := `
        INSERT INTO heroes (first_seen, name, can_fly, realname, abilities) 
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id, version`
	args := []interface{}{hero.FirstSeen, hero.Name, hero.CanFly, hero.RealName, pq.Array(hero.Abilities)}
	return m.DB.QueryRow(query, args...).Scan(&hero.ID, &hero.Version)
}

func (m HeroModel) Get(id int64) (*Hero, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
        SELECT id, first_seen, name, can_fly, realname, abilities, version
        FROM heroes
        WHERE id = $1`

	var hero Hero

	err := m.DB.QueryRow(query, id).Scan(
		&hero.ID,
		&hero.FirstSeen,
		&hero.Name,
		&hero.CanFly,
		&hero.RealName,
		pq.Array(&hero.Abilities),
		&hero.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &hero, nil
}

func (m HeroModel) Update(hero *Hero) error {
	query := `
        UPDATE heroes
        SET first_seen = $1, name = $2, can_fly = $3, realname = $4, abilities = $5, version = version + 1
        WHERE id = $6
        RETURNING version`

	args := []interface{}{
		hero.FirstSeen,
		hero.Name,
		hero.CanFly,
		hero.RealName,
		pq.Array(hero.Abilities),
		hero.ID,
	}

	return m.DB.QueryRow(query, args...).Scan(&hero.Version)
}

func (m HeroModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `
        DELETE FROM heroes
        WHERE id = $1`

	result, err := m.DB.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}

func (m HeroModel) GetAll(name string, abilities []string, filters Filters) ([]*Hero, error) {
	query := fmt.Sprintf(`
        SELECT id, first_seen, name, can_fly, realname, abilities, version
        FROM heroes
        WHERE (LOWER(name) LIKE LOWER($1) OR $1 = '') 
        AND (abilities @> $2 OR $2 = '{}')     
        ORDER BY %s ASC, id ASC
        LIMIT $3 OFFSET $4`, filters.Sort)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []interface{}{name, pq.Array(abilities), filters.limit(), filters.offset()}
	rows, err := m.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	heroes := []*Hero{}

	for rows.Next() {
		var hero Hero

		err := rows.Scan(
			&hero.ID,
			&hero.FirstSeen,
			&hero.Name,
			&hero.CanFly,
			&hero.RealName,
			pq.Array(&hero.Abilities),
			&hero.Version,
		)
		if err != nil {
			return nil, err
		}

		heroes = append(heroes, &hero)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return heroes, nil
}
