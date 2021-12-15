package main

import (
	"net/http"

	"heroes.rainerstropek.com/internal/middleware"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)
	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)

	jwtMiddleware := middleware.NewJwtMiddleware(
		app.config.azure.tenantId,
		[]string{"api://4fac0887-b94f-4ea9-a8d3-06c7bca2a7bd"}) // Make scope configurable if you need to

	protectedrouter := httprouter.New()
	protectedrouter.NotFound = http.HandlerFunc(app.notFoundResponse)
	protectedrouter.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)
	protectedrouter.HandlerFunc(http.MethodGet, "/v1/heroes", app.listHeroesHandler)
	protectedrouter.HandlerFunc(http.MethodPost, "/v1/heroes", app.createHeroHandler)
	protectedrouter.HandlerFunc(http.MethodPut, "/v1/heroes/:id", app.updateHeroHandler)
	protectedrouter.HandlerFunc(http.MethodDelete, "/v1/heroes/:id", app.deleteHeroHandler)
	protectedrouter.HandlerFunc(http.MethodPost, "/v1/generate", app.generateDemoDataHandler)
	protectedrouter.HandlerFunc(http.MethodGet, "/v1/heroes/:id", app.showHeroHandler)
	protectedrouter.HandlerFunc(http.MethodGet, "/v1/claims", middleware.ClaimsHandler)
	router.NotFound = jwtMiddleware.CheckJWT(protectedrouter)

	c := alice.New()
	c.Append(app.recoverPanic)
	c.Append(app.enableCORS)
	chain := c.Then(router)

	return chain
}
