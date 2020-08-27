package customerhandlers

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

// ObjectResultWriter writes a given object to the HTTP response
type ObjectResultWriter interface {
	WriteObjectResult(w http.ResponseWriter, object interface{})
}

// CustomerHandlers represents functions handling HTTP requests for customers management web api
type CustomerHandlers struct {
	repo customerrepository.CustomerRepository
	orw  ObjectResultWriter
}

// NewCustomerHandlers creates a customer handler object
func NewCustomerHandlers(repo customerrepository.CustomerRepository, orw ObjectResultWriter) CustomerHandlers {
	return CustomerHandlers{
		repo: repo,
		orw:  orw,
	}
}

// GetCustomers returns all customers
func (ch CustomerHandlers) GetCustomers(w http.ResponseWriter, r *http.Request) {
	custArray := ch.repo.GetCustomersArray()
	orderBy := r.FormValue("orderBy")
	if len(orderBy) > 0 {
		if orderBy != "companyName" {
			http.Error(w, "Currently, we can only order by companyName", http.StatusBadRequest)
			return
		}

		sort.Sort(customerrepository.ByCompanyName(custArray))
	}

	// Return all customers
	ch.orw.WriteObjectResult(w, custArray)
}

// GetCustomer returns a single customer based on a given customer ID
func (ch CustomerHandlers) GetCustomer(w http.ResponseWriter, r *http.Request) {
	// Get customer ID from path
	cid, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "Invalid customer ID format", http.StatusBadRequest)
		return
	}

	// Check if customer with given ID exists
	if c, ok := ch.repo.GetCustomerByID(cid); ok {
		// Return customer
		ch.orw.WriteObjectResult(w, c)
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

// AddCustomer adds a customer
func (ch CustomerHandlers) AddCustomer(w http.ResponseWriter, r *http.Request) {
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
	ch.repo.AddCustomer(c)

	// Return customer
	w.Header().Set("Location", fmt.Sprintf("/customers/%s", c.CustomerID))
	w.WriteHeader(http.StatusCreated)
	ch.orw.WriteObjectResult(w, c)
}

// DeleteCustomer deletes a customer based on a given ID
func (ch CustomerHandlers) DeleteCustomer(w http.ResponseWriter, r *http.Request) {
	// Get customer ID from path
	cid, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "Invalid customer ID format", http.StatusBadRequest)
		return
	}

	// Delete customer
	if ch.repo.DeleteCustomerByID(cid) {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// Customer hasn't been found
	http.NotFound(w, r)
}

// PatchCustomer patches a customer based on a given ID and new field values
func (ch CustomerHandlers) PatchCustomer(w http.ResponseWriter, r *http.Request) {
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

	if cNew, ok := ch.repo.PatchCustomer(cid, c); ok {
		// Return updated customer data
		ch.orw.WriteObjectResult(w, cNew)
		return
	}

	// Customer hasn't been found
	http.NotFound(w, r)
}
