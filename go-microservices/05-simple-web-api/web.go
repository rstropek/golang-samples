package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"path"
	"strconv"

	"github.com/gorilla/mux"
)

// Go has good JSON support (see also https://blog.golang.org/json-and-go).
// Note that this example uses Tags (see also https://medium.com/golangspec/tags-in-golang-3e5db0b8ef3e)
// for turning structs into JSON.

// Person represents a person
type Person struct {
	ID        int      `json:"id,omitempty"`
	Firstname string   `json:"firstname,omitempty"`
	Lastname  string   `json:"lastname,omitempty"`
	Address   *Address `json:"address,omitempty"`
}

// Address represents an address of a person
type Address struct {
	City  string `json:"city,omitempty"`
	State string `json:"state,omitempty"`
}

func (p *Person) isValid() error {
	switch {
	case len(p.Firstname) == 0:
		return errors.New("Firstname has to be set")
	case len(p.Lastname) == 0:
		return errors.New("Lastname has to be set")
	default:
		return nil
	}
}

// Define a global variable holding people. Note that this is a slice
// (see also https://gobyexample.com/slices).
var people []Person

func main() {
	// Note the use of the built-in `append` function
	people = append(people, Person{ID: 1, Firstname: "John", Lastname: "Doe", Address: &Address{City: "City X", State: "State X"}})
	people = append(people, Person{ID: 2, Firstname: "Koko", Lastname: "Doe", Address: &Address{City: "City Z", State: "State Y"}})
	people = append(people, Person{ID: 3, Firstname: "Francis", Lastname: "Sunday"})

	// Note the use of the Gorilla MUX here (see also https://github.com/gorilla/mux)
	router := mux.NewRouter()
	router.HandleFunc("/people", GetPeople).Methods("GET")
	router.HandleFunc("/people", CreatePerson).Methods("POST")
	router.HandleFunc("/people/{id:[0-9]+}", GetPerson).Methods("GET")
	router.HandleFunc("/people/{id:[0-9]+}", DeletePerson).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8080", router))
}

// GetPeople returns a list of people
func GetPeople(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(people)
}

// GetPerson returns a single person identified by the ID in the URL path
func GetPerson(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	// Find person based on ID
	id, _ := strconv.Atoi(params["id"])
	for _, item := range people {
		if item.ID == id {
			json.NewEncoder(w).Encode(item)
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
}

// CreatePerson creates a new person and adds it to the list of people
func CreatePerson(w http.ResponseWriter, r *http.Request) {
	// JSON decode body
	var person Person
	_ = json.NewDecoder(r.Body).Decode(&person)
	if err := person.isValid(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	// Get highest ID of all people
	person.ID = 0
	for _, p := range people {
		if p.ID > person.ID {
			person.ID = p.ID
		}
	}

	// Append person to people
	people = append(people, person)

	// Return the (possibly modified) person
	w.Header().Add("Location", path.Join(r.URL.Path, strconv.Itoa(person.ID)))
	json.NewEncoder(w).Encode(person)
}

// DeletePerson deletes a person based on its ID in the URL path
func DeletePerson(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	// Find person based on its ID
	// Note the range iterator here (see also https://gobyexample.com/range)
	id, _ := strconv.Atoi(params["id"])
	for index, item := range people {
		if item.ID == id {
			// Delete item using its index
			people = append(people[:index], people[index+1:]...)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
}
