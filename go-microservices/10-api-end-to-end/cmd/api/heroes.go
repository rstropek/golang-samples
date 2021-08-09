package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/rstropek/golang-samples/api-end-to-end/internal/data"
	"github.com/rstropek/golang-samples/api-end-to-end/internal/validator"
)

func (app *application) createHeroHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name     string      `json:"name"`
		RealName string      `json:"realName"`
		Coolness int32       `json:"coolness"`
		Tags     []string    `json:"tags"`
		CanFly   data.CanFly `json:"canFly"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// You could consider creating deep copies using JSON marshal,
	// ProtoBuf marshal, deepcopier, etc.
	hero := &data.Hero{
		Name:     input.Name,
		RealName: input.RealName,
		Coolness: input.Coolness,
		Tags:     input.Tags,
		CanFly:   input.CanFly,
	}

	v := validator.New()
	if data.ValidateHero(v, hero); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	fmt.Fprintf(w, "%+v\n", input)
}

func (app *application) showHeroHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	hero := data.Hero{
		ID:        id,
		CreatedAt: time.Now(),
		Name:      "Homelander",
		RealName:  "John",
		Coolness:  9,
		Tags:      []string{"The Boys", "Evil"},
		CanFly:    true,
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"hero": hero}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
