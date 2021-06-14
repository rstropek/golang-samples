package main

import (
	"fmt"
	"net/http"

	cr "github.com/rstropek/golang-samples/go-microservices/05b-web-api/internal/customerrepository"
	"github.com/rstropek/golang-samples/go-microservices/05b-web-api/internal/validator"
	"github.com/shopspring/decimal"
)

func (app *application) createPerson(w http.ResponseWriter, r *http.Request) {
	// Decode customer data from request body
	var c = cr.Customer{}
	if err := app.readJSON(w, r, &c); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// Make sure that incoming custer data is sane
	v := validator.New()
	v.Check(validator.IsEmptyUuid(c.CustomerID), "CustomerID", "must be empty")
	v.Check(validator.IsNotEmptyString(c.CompanyName), "Company name", "must not be empty")
	v.Check(validator.IsNotEmptyString(c.ContactName), "Contact name", "must not be empty")
	v.Check(validator.HasLen(c.Country, 3), "Country", "must be three characters long (use ISO 3166-1 Alpha-3 code)")
	v.Check(validator.IsGreaterThan(c.HourlyRate, decimal.NewFromInt(0)), "Hourly rate", "must be greater than 0")
	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	// Assign new customer ID
	c.CustomerID = app.newUUID()

	// Add customer to our list
	app.repository.AddCustomer(c)

	// Return customer
	app.writeJSON(w, http.StatusCreated, c, http.Header{
		"Location": []string{fmt.Sprintf("/customers/%s", c.CustomerID)},
	})
}

func (app *application) getPerson(w http.ResponseWriter, r *http.Request) {
	id, err := app.readUuidParam(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()
	v.Check(!validator.IsEmptyUuid(id), "id", "must not be empty")
	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	c, found := app.repository.GetCustomerByID(id)

	if !found {
		app.notFoundResponse(w, r)
	}

	app.writeJSON(w, http.StatusOK, c, nil)
}


func (app *application) getPeople(w http.ResponseWriter, r *http.Request) {
	c := app.repository.GetCustomersArray()
	app.writeJSON(w, http.StatusOK, c, nil)
}
