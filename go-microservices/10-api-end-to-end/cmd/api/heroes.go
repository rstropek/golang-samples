package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/rstropek/golang-samples/api-end-to-end/internal/data"
)

func (app *application) createHeroHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "create a new hero")
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

	err = app.writeJSON(w, http.StatusOK, hero, nil)
	if err != nil {
		app.logger.Println(err)
		http.Error(w, "The server encountered a problem and could not process your request", http.StatusInternalServerError)
	}
}
