package main

import (
    "net/http"

    "github.com/julienschmidt/httprouter"
)

func (app *application) routes() *httprouter.Router {
    // Initialize a new httprouter router instance.
    router := httprouter.New()

    // Error handler
    router.NotFound = http.HandlerFunc(app.notFoundResponse)
    router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

    // Register the relevant methods, URL patterns and handler functions.
    router.HandlerFunc(http.MethodGet, "/healthcheck", app.healthcheckHandler)
    router.HandlerFunc(http.MethodPost, "/people", app.createPerson)
    router.HandlerFunc(http.MethodGet, "/people", app.getPeople)
    router.HandlerFunc(http.MethodGet, "/people/:id", app.getPerson)

    // Return the httprouter instance.
    return router
}
