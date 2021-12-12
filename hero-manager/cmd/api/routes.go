package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)
	router.HandlerFunc(http.MethodGet, "/v1/heroes", app.listHeroesHandler)
	router.HandlerFunc(http.MethodPost, "/v1/heroes", app.createHeroHandler)
	router.HandlerFunc(http.MethodPut, "/v1/heroes/:id", app.updateHeroHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/heroes/:id", app.deleteHeroHandler)
	router.HandlerFunc(http.MethodPost, "/v1/generate", app.generateDemoDataHandler)
	router.HandlerFunc(http.MethodGet, "/v1/heroes/:id", app.showHeroHandler)

	c := alice.New()
	c.Append(app.recoverPanic)
	c.Append(app.enableCORS)
	chain := c.Then(router)

	return chain
}
