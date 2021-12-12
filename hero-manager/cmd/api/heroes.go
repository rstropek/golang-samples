package main

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"heroes.rainerstropek.com/internal/data"
	"heroes.rainerstropek.com/internal/validator"
)

func (app *application) createHeroHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name      string    `json:"name"`
		FirstSeen time.Time `json:"firstSeen"`
		CanFly    bool      `json:"canFly"`
		RealName  string    `json:"realName,omitempty"`
		Abilities []string  `json:"abilities"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	hero := &data.Hero{
		Name:      input.Name,
		FirstSeen: input.FirstSeen,
		CanFly:    input.CanFly,
		RealName:  input.RealName,
		Abilities: input.Abilities,
	}

	v := validator.New()

	if data.ValidateHero(v, hero); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Heroes.Insert(hero)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/movies/%d", hero.ID))

	err = app.writeJSON(w, http.StatusCreated, hero, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) generateDemoDataHandler(w http.ResponseWriter, r *http.Request) {
	for i := 0; i < 20; i++ {
		hero := &data.Hero{
			Name:      gofakeit.Name(),
			FirstSeen: gofakeit.DateRange(time.Date(1900, time.Month(1), 1, 0, 0, 0, 0, time.UTC), time.Now()),
			CanFly:    gofakeit.Bool(),
			RealName:  gofakeit.Name(),
			Abilities: []string{"foo", "bar"},
		}

		err := app.models.Heroes.Insert(hero)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}
	}

	headers := make(http.Header)
	err := app.writeJSON(w, http.StatusCreated, "created", headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) showHeroHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	hero, err := app.models.Heroes.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, hero, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateHeroHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	hero, err := app.models.Heroes.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	var input struct {
		Name      string    `json:"name"`
		FirstSeen time.Time `json:"firstSeen"`
		CanFly    bool      `json:"canFly"`
		RealName  string    `json:"realName,omitempty"`
		Abilities []string  `json:"abilities"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	hero.FirstSeen = input.FirstSeen
	hero.Name = input.Name
	hero.CanFly = input.CanFly
	hero.RealName = input.RealName
	hero.Abilities = input.Abilities

	v := validator.New()

	if data.ValidateHero(v, hero); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Heroes.Update(hero)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, hero, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteHeroHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.models.Heroes.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, "successfully deleted", nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) listHeroesHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name      string
		Abilities []string
		data.Filters
	}

	v := validator.New()
	qs := r.URL.Query()
	input.Name = app.readString(qs, "name", "")
	input.Name = fmt.Sprintf("%%%s%%", input.Name)
	input.Abilities = app.readCSV(qs, "abilities", []string{})
	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)
	input.Filters.Sort = app.readString(qs, "sort", "id")
	input.Filters.SortSafelist = []string{"id", "name", "realname"}

	if data.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	heroes, err := app.models.Heroes.GetAll(input.Name, input.Abilities, input.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, heroes, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
