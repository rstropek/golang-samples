package main

import (
    "net/http"

    "github.com/julienschmidt/httprouter"
)

func (app *application) routes() *httprouter.Router {
    // Create a new httprouter router instance.
    router := httprouter.New()

    // Register routes
    router.NotFound = http.HandlerFunc(app.notFoundResponse)
    router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

    router.HandlerFunc(http.MethodGet, "/healthcheck", app.healthcheckHandler)
    router.HandlerFunc(http.MethodGet, "/crash", app.crash)
    router.HandlerFunc(http.MethodPost, "/heroes", app.createHeroHandler)
    router.HandlerFunc(http.MethodGet, "/heroes/:id", app.showHeroHandler)

    // Return the httprouter instance.
    return router
}
