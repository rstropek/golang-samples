package customerhandlers

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/rstropek/golang-samples/web-api/customerrepository"
    "net/http"
    "net/http/httptest"
    "testing"
)

type testResponseWriter struct {}

func (r testResponseWriter) WriteObjectResult(w http.ResponseWriter, object interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(object)
}

func TestGetCustomers(t *testing.T) {
	// Here we use the existing customer repository. In practice, you would probably
	// use a mocking framework like https://github.com/stretchr/testify. However, proper
	// mocking for unit tests is out of scope here.
	repo := customerrepository.NewCustomerRepository()
	repo.AddCustomer(customerrepository.Customer{CompanyName: "Foo Bar"})
	ch := NewCustomerHandlers(repo, testResponseWriter{})

    // Create a request to pass to our handler
    req, _ := http.NewRequest("GET", "/", nil)

    // Create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response
    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(ch.GetCustomers)

    // Our handlers satisfy http.Handler, so we can call their ServeHTTP method 
    // directly and pass in our Request and ResponseRecorder
    handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	assert.Equal(t, http.StatusOK, rr.Code)

	// Check content type
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))

	// Check the JSON result
	result := make([]customerrepository.Customer, 0)
	json.NewDecoder(rr.Body).Decode(&result)
	assert.Equal(t, 1, len(result))
	assert.Equal(t, "Foo Bar", result[0].CompanyName)
}