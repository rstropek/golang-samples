package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/rstropek/golang-samples/web-api/customerrepository"
	"github.com/shopspring/decimal"
)

func getCustomers(w http.ResponseWriter, r *http.Request) {
	custArray := repo.GetCustomersArray()
	orderBy := r.FormValue("orderBy")
	if len(orderBy) > 0 {
		if orderBy != "companyName" {
			http.Error(w, "Currently, we can only order by companyName", http.StatusBadRequest)
			return
		}

		sort.Sort(customerrepository.ByCompanyName(custArray))
	}

	// Return all customers
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(custArray)
}

func getCustomer(w http.ResponseWriter, r *http.Request) {
	// Get customer ID from path
	cid, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "Invalid customer ID format", http.StatusBadRequest)
		return
	}

	// Check if customer with given ID exists
	if c, ok := repo.GetCustomerByID(cid); ok {
		// Return customer
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(c)
		return
	}

	// Customer hasn't been found
	http.NotFound(w, r)
}

// newUUID returns a new UUID and ignores potential errors
func newUUID() uuid.UUID {
	r, _ := uuid.NewUUID()
	return r
}

func addCustomer(w http.ResponseWriter, r *http.Request) {
	// Decode customer data from request body
	var c = customerrepository.Customer{}
	if json.NewDecoder(r.Body).Decode(&c) != nil {
		http.Error(w, "Could not deserialize customer from HTTP body", http.StatusBadRequest)
		return
	}

	// Make sure that incoming custer data is sane
	if c.CustomerID != uuid.Nil {
		http.Error(w, "CustomerID must be empty", http.StatusBadRequest)
		return
	}

	if len(c.CompanyName) == 0 {
		http.Error(w, "Company name must not be empty", http.StatusBadRequest)
		return
	}

	if len(c.ContactName) == 0 {
		http.Error(w, "Contact name must not be empty", http.StatusBadRequest)
		return
	}

	if len(c.Country) != 3 {
		http.Error(w, "Country name must be three characters long (use ISO 3166-1 Alpha-3 code)", http.StatusBadRequest)
		return
	}

	if decimal.NewFromInt(0).GreaterThan(c.HourlyRate) {
		http.Error(w, "Hourly rate must be >= 0", http.StatusBadRequest)
		return
	}

	// Assign new customer ID
	c.CustomerID = newUUID()

	// Add customer to our list
	repo.AddCustomer(c)

	// Return customer
	w.Header().Set("Location", fmt.Sprintf("/customers/%s", c.CustomerID))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(c)
}

func deleteCustomer(w http.ResponseWriter, r *http.Request) {
	// Get customer ID from path
	cid, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "Invalid customer ID format", http.StatusBadRequest)
		return
	}

	// Delete customer
	if repo.DeleteCustomerByID(cid) {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// Customer hasn't been found
	http.NotFound(w, r)
}

func patchCustomer(w http.ResponseWriter, r *http.Request) {
	// Get customer ID from path
	cid, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "Invalid customer ID format", http.StatusBadRequest)
		return
	}

	// Decode customer data from request body
	var c = customerrepository.Customer{}
	if json.NewDecoder(r.Body).Decode(&c) != nil {
		http.Error(w, "Could not deserialize customer from HTTP body", http.StatusBadRequest)
		return
	}

	// If customer ID was specified, it must match the customer ID from path
	if c.CustomerID != uuid.Nil && cid != c.CustomerID {
		http.Error(w, "Cannot update customer ID", http.StatusBadRequest)
		return
	}

	if len(c.Country) > 0 && len(c.Country) != 3 {
		http.Error(w, "Country name must be three characters long (use ISO 3166-1 Alpha-3 code)", http.StatusBadRequest)
		return
	}

	if c.HourlyRate != decimal.NewFromInt(0) && decimal.NewFromInt(0).GreaterThan(c.HourlyRate) {
		http.Error(w, "Hourly rate must be >= 0", http.StatusBadRequest)
		return
	}

	if cNew, ok := repo.PatchCustomer(cid, c); ok {
		// Return updated customer data
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(cNew)
		return
	}

	// Customer hasn't been found
	http.NotFound(w, r)
}
