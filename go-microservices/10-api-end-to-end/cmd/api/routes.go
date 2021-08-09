package main

import (
    "net/http"

    "github.com/julienschmidt/httprouter"
)

func (app *application) routes() *httprouter.Router {
    // Create a new httprouter router instance.
    router := httprouter.New()

    // Register routes
    router.HandlerFunc(http.MethodGet, "/healthcheck", app.healthcheckHandler)
    router.HandlerFunc(http.MethodPost, "/heroes", app.createHeroHandler)
    router.HandlerFunc(http.MethodGet, "/heroes/:id", app.showHeroHandler)

    // Return the httprouter instance.
    return router
}
