package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/user/gab-conference-api/gabdownloader"
)

var sessionReferences gabdownloader.SessionReferences
var sessions gabdownloader.Sessions

func main() {
	sessionReferences = gabdownloader.GetAndCacheSessionReferencesFromGitHub()
	sessions = gabdownloader.GetAndCacheSessionsFromGitHub(sessionReferences)

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
	router.HandleFunc("/sessions", getSessions).Methods("GET")
	router.HandleFunc("/sessions/{id}", getSession).Methods("GET")
	router.HandleFunc("/slots", getSlots).Methods("GET")
	fmt.Printf("Server starting, will listen to port %d...\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), handlers.CORS(originsOk, headersOk, methodsOk)(router)))
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func getPing(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode("Pong")
}

type sessionResponse struct {
	ID         string `json:"id"`
	DetailsURL string `json:"detailsUrl"`
}

func getSessions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// check whether parameter `full` is present
	params := r.URL.Query()
	_, prs := params["full"]
	if !prs {
		response := make([]sessionResponse, len(sessionReferences))
		for i, s := range sessionReferences {
			response[i] = sessionResponse{ID: s.Name, DetailsURL: fmt.Sprintf("/sessions/%s", s.Name)}
		}

		json.NewEncoder(w).Encode(response)
	} else {
		json.NewEncoder(w).Encode(sessions)
	}
}

func writeError(sc int, err error, w http.ResponseWriter, enc *json.Encoder) {
	w.WriteHeader(sc)
	enc.Encode(map[string]string{"error": err.Error()})
}

func getSession(w http.ResponseWriter, r *http.Request) {
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

	// search for session based on ID
	for _, item := range sessions {
		if item.ID == paramIDStr {
			// we have found the session based on ID
			w.WriteHeader(http.StatusOK)
			enc.Encode(item)
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
}

type slot struct {
	ID    uint   `json:"id"`
	Start string `json:"start"`
	End   string `json:"end"`
}

var slots = []slot{
	{ID: 1, Start: "09:15", End: "09:50"},
	{ID: 2, Start: "09:55", End: "10:45"},
	{ID: 3, Start: "11:10", End: "12:00"},
	{ID: 4, Start: "13:00", End: "13:50"},
	{ID: 5, Start: "14:05", End: "14:55"},
	{ID: 6, Start: "15:10", End: "16:00"},
}

func getSlots(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(slots)
}
