// Define a package named main.
package main

import (
	// Importing necessary packages for encoding, error handling, formatting, logging, HTTP handling, string conversion, synchronization, and Gorilla Mux for routing.
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
	"sync/atomic"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

// Define a struct to represent a lunch order.
type lunchOrder struct {
	ID     uint64 `json:"id,omitempty"` // Order ID, omitted if zero.
	Day    byte   `json:"day"`          // Day of the week as a byte (0-6).
	Person string `json:"person"`       // Person who placed the order.
	MealID uint16 `json:"mealId"`       // ID of the meal ordered.
	Notes  string `json:"notes"`        // Additional notes for the order.
}

// Method to validate a lunchOrder instance.
func (lo *lunchOrder) isValid() error {
	switch {
	case len(lo.Person) == 0:
		return errors.New("person must be set") // Ensure 'Person' is not empty.
	case lo.Day > 6:
		return errors.New("day must be between 0 (Monday) and 6 (Sunday)") // Validate day value.
	case lo.MealID < 1 || lo.MealID > 6:
		return errors.New("invalid meal id") // Check if MealID is within a valid range.
	default:
		return nil // If all checks pass, return nil (no error).
	}
}

// Define a struct to represent a meal.
type meal struct {
	ID    uint64  `json:"id,omitempty"` // Meal ID.
	Desc  string  `json:"desc"`         // Description of the meal.
	Price float32 `json:"price"`        // Price of the meal.
}

// Define and initialize an array of meal instances.
var meals = [...]meal{
	// Array of predefined meals.
	{ID: 1, Desc: "Spaghetti Pomodoro", Price: 5.90},
	{ID: 2, Desc: "Wiener Schnitzel", Price: 8.50},
	{ID: 3, Desc: "Pizza Speciale", Price: 7.20},
	{ID: 4, Desc: "Greek Salat", Price: 5.30},
	{ID: 5, Desc: "Vegetable Soup", Price: 2.10},
}

// Initialize a slice to store lunch orders.
var orders = []lunchOrder{
	// Pre-populated lunch orders.
	{ID: 1, Day: 0, Person: "Tom", MealID: 1, Notes: "With lots of Parmesan"},
	{ID: 2, Day: 1, Person: "Tom", MealID: 2},
	{ID: 3, Day: 2, Person: "Tom", MealID: 4, Notes: "Without onions"},
	{ID: 4, Day: 2, Person: "Tom", MealID: 5},
	{ID: 5, Day: 3, Person: "Tom", MealID: 3},
	{ID: 6, Day: 4, Person: "Tom", MealID: 6},
	{ID: 7, Day: 0, Person: "Jane", MealID: 2, Notes: "Potatoes instead of Fries"},
	{ID: 8, Day: 1, Person: "Jane", MealID: 3},
	{ID: 9, Day: 2, Person: "Jane", MealID: 5},
	{ID: 10, Day: 2, Person: "Jane", MealID: 6},
	{ID: 11, Day: 3, Person: "Jane", MealID: 1},
	{ID: 12, Day: 4, Person: "Jane", MealID: 4, Notes: "No Olives"},
}

// Variable to store the current highest order ID.
var id uint64 = 12

// Mutex for synchronizing access to the orders slice.
var collectionMutex sync.Mutex

// Main function of the program.
func main() {
	var port uint16 = 8080 // Port number for the server.

	// Set up CORS (Cross-Origin Resource Sharing) settings.
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	// Initialize a new Gorilla Mux router.
	router := mux.NewRouter()

	// Define route handlers for different endpoints.
	router.HandleFunc("/ping", getPing).Methods("GET")
	router.HandleFunc("/meals", getMeals).Methods("GET")
	router.HandleFunc("/lunchOrders", getLunchOrders).Methods("GET")
	router.HandleFunc("/lunchOrders/{id}", getLunchOrder).Methods("GET")
	router.HandleFunc("/lunchOrders", addLunchOrder).Methods("POST")

	// Start the server and listen on the specified port.
	fmt.Printf("Server starting, will listen to port %d...\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), handlers.CORS(originsOk, headersOk, methodsOk)(router)))
}

// Helper function to write error responses.
func writeError(sc int, err error, w http.ResponseWriter, enc *json.Encoder) {
	w.WriteHeader(sc) // Set HTTP status code.
	enc.Encode(map[string]string{"error": err.Error()}) // Encode the error message.
}

// Handler for the /ping endpoint.
func getPing(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode("Pong") // Respond with "Pong".
}

// Handler for the /meals endpoint.
func getMeals(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json") // Set content type to JSON.
	json.NewEncoder(w).Encode(meals) // Encode and send the meals array.
}

// Handler for the /lunchOrders endpoint.
func getLunchOrders(w http.ResponseWriter, r *http.Request) {
	queryValues := r.URL.Query() // Extract query parameters.
	person := queryValues.Get("person") // Get the 'person' query parameter.
	filtered := len(person) != 0 // Check if filtering is needed.
	result := []lunchOrder{} // Initialize a slice to store the filtered results.

	// Iterate over orders and filter based on the 'person' parameter.
	for _, order := range orders {
		if !filtered || order.Person == person {
			result = append(result, order) // Append matching orders to the result.
		}
	}

	w.Header().Set("Content-Type", "application/json") // Set content type to JSON.
	json.NewEncoder(w).Encode(result) // Encode and send the result.
}

// Handler for the /lunchOrders/{id} endpoint.
func getLunchOrder(w http.ResponseWriter, r *http.Request) {
	// Set response content type to JSON.
	w.Header().Set("Content-Type", "application/json")

	// Initialize a JSON encoder for the response.
	enc := json.NewEncoder(w)

	// Extract the 'id' parameter from the route.
	params := mux.Vars(r)
	paramIDStr, prs := params["id"]
	if !prs {
		writeError(http.StatusBadRequest, errors.New("missing id"), w, enc) // Handle missing ID.
		return
	}

	// Convert the ID parameter to a uint64.
	paramID, convErr := strconv.ParseUint(paramIDStr, 10, 64)
	if convErr != nil {
		writeError(http.StatusBadRequest, convErr, w, enc) // Handle invalid ID format.
		return
	}

	// Search for the order with the given ID.
	for _, item := range orders {
		if item.ID == paramID {
			w.WriteHeader(http.StatusOK) // Set HTTP status to OK.
			enc.Encode(item) // Send the found order.
			return
		}
	}

	// If no order is found, set status to Not Found.
	w.WriteHeader(http.StatusNotFound)
}

// Handler for adding a new lunch order.
func addLunchOrder(w http.ResponseWriter, r *http.Request) {
	// Set response content type to JSON.
	w.Header().Set("Content-Type", "application/json")

	// Variable to store the incoming order.
	var order lunchOrder

	// Initialize a JSON encoder for the response.
	enc := json.NewEncoder(w)

	// Decode the incoming JSON from the request body into the order variable.
	err := json.NewDecoder(r.Body).Decode(&order)
	if err != nil {
		writeError(http.StatusBadRequest, err, w, enc) // Handle JSON decoding errors.
		return
	}

	// Validate the decoded order.
	if err = order.isValid(); err != nil {
		writeError(http.StatusBadRequest, err, w, enc) // Handle validation errors.
		return
	}

	// Atomically increment the unique order ID counter and assign it to the new order.
	atomic.AddUint64(&id, 1)
	order.ID = atomic.LoadUint64(&id)

	// Add the new order to the orders slice, ensuring thread safety with a mutex.
	collectionMutex.Lock()
	orders = append(orders, order)
	collectionMutex.Unlock()

	// Set response headers and status, then send the order in response.
	w.Header().Set("Location", fmt.Sprintf("/lunchOrders/%d", order.ID))
	w.WriteHeader(http.StatusCreated)
	enc.Encode(order)
}