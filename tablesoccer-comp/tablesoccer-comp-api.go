package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync/atomic"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

// Match represents a table soccer match
// Note that scores are pointers so that we can identify missing values.
// Otherwise missing values would become zero.
type Match struct {
	ID      uint64 `json:"id,omitempty"`
	Player1 string `json:"p1"`
	Player2 string `json:"p2"`
	Score1  byte   `json:"s1"`
	Score2  byte   `json:"s2"`
}

func (m *Match) isValid() error {
	switch {
	case len(m.Player1) == 0 || len(m.Player2) == 0:
		return errors.New("Player 1 and 2 have to be set")
	case m.Score1 > 10 || m.Score2 > 10:
		return errors.New("Both scores must be >= 0 and <= 10")
	default:
		return nil
	}
}

var matches = []Match{
	{ID: 1, Player1: "Tom", Player2: "Jane", Score1: 5, Score2: 10},
	{ID: 2, Player1: "Jane", Player2: "Phil", Score1: 10, Score2: 8},
	{ID: 3, Player1: "Jane", Player2: "Tom", Score1: 10, Score2: 6},
	{ID: 4, Player1: "Tom", Player2: "Phil", Score1: 8, Score2: 10},
	{ID: 5, Player1: "Phil", Player2: "Jane", Score1: 9, Score2: 10},
	{ID: 6, Player1: "Jane", Player2: "Phil", Score1: 10, Score2: 7},
	{ID: 7, Player1: "Tom", Player2: "Jane", Score1: 3, Score2: 10},
}
var id uint64 = 7

// our main function
func main() {
	var port uint16 = 8080
	args := os.Args[1:]
	if len(args) > 0 {
		argsPort, convErr := strconv.ParseUint(args[0], 10, 16)
		if convErr != nil {
			fmt.Fprintf(os.Stderr, "%s is not a valid port number", args[0])
			return
		}

		port = uint16(argsPort)
	} else {
		fmt.Printf("No port given in command line, using default (%d).\n", port)
	}

	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	router := mux.NewRouter()
	router.HandleFunc("/ping", getPing).Methods("GET")
	router.HandleFunc("/matches", getMatches).Methods("GET")
	router.HandleFunc("/matches/{id}", getMatch).Methods("GET")
	router.HandleFunc("/matches", addMatch).Methods("POST")
	fmt.Printf("Server starting, will listen to port %d...\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), handlers.CORS(originsOk, headersOk, methodsOk)(router)))
}

func writeError(sc int, err error, w http.ResponseWriter, enc *json.Encoder) {
	w.WriteHeader(sc)
	enc.Encode(map[string]string{"error": err.Error()})
}

func getPing(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode("Pong")
}

func getMatches(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(matches)
}

func getMatch(w http.ResponseWriter, r *http.Request) {
	// content type will always be json
	w.Header().Set("Content-Type", "application/json")

	// create an encoder for writing response
	enc := json.NewEncoder(w)

	// check whether ID is present
	params := mux.Vars(r)
	paramIDStr, prs := params["id"]
	if !prs {
		writeError(http.StatusBadRequest, errors.New("Missing id"), w, enc)
		return
	}

	// try to parse ID as an int
	paramID, convErr := strconv.ParseUint(paramIDStr, 10, 64)
	if convErr != nil {
		writeError(http.StatusBadRequest, convErr, w, enc)
		return
	}

	// search for match based on ID
	for _, item := range matches {
		if item.ID == paramID {
			// we have found the match based on ID
			w.WriteHeader(http.StatusOK)
			enc.Encode(item)
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
}

func addMatch(w http.ResponseWriter, r *http.Request) {
	// content type will always be json
	w.Header().Set("Content-Type", "application/json")

	// resulting match struct
	var match Match

	// create an encoder for writing response
	enc := json.NewEncoder(w)

	// decode incoming json from body
	err := json.NewDecoder(r.Body).Decode(&match)
	if err != nil {
		writeError(http.StatusBadRequest, err, w, enc)
		return
	}

	// Validate match
	if err = match.isValid(); err != nil {
		writeError(http.StatusBadRequest, err, w, enc)
		return
	}

	// increment unique match id counter and add it to match
	atomic.AddUint64(&id, 1)
	match.ID = atomic.LoadUint64(&id)

	// add match to matches
	matches = append(matches, match)

	// return success
	w.Header().Set("Location", fmt.Sprintf("/matches/%d", match.ID))
	w.WriteHeader(http.StatusCreated)
	enc.Encode(match)
}
