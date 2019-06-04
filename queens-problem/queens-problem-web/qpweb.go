package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	qpbas "github.com/rstropek/golang-samples/queens-problem/queens-problem-bitarray-solver"
)

func main() {
	port := flag.Uint("p", 8080, "Port number for webserver")
	flag.Parse()

	router := mux.NewRouter()
	router.HandleFunc("/ping", getPing).Methods("GET")
	router.HandleFunc("/solve", solve).Methods("POST")
	fmt.Printf("Server starting, will listen to port %d...\n", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), router))
}

func getPing(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode("Pong")
}

type solverParameters struct {
	SideLength byte `json:"sideLength"`
}

type solverResult struct {
	NumberOfResults int    `json:"numberOfResults"`
	CalculationTime string `json:"calculationTime"`
}

func (sp *solverParameters) isValid() error {
	switch {
	case sp.SideLength > 12:
		return errors.New("Side length must be less or equal 12")
	default:
		return nil
	}
}

func writeError(sc int, err error, w http.ResponseWriter, enc *json.Encoder) {
	w.WriteHeader(sc)
	enc.Encode(map[string]string{"error": err.Error()})
}

func solve(w http.ResponseWriter, r *http.Request) {
	enc := json.NewEncoder(w)
	var parameters solverParameters

	err := json.NewDecoder(r.Body).Decode(&parameters)
	if err != nil {
		writeError(http.StatusBadRequest, err, w, enc)
		return
	}

	if err = parameters.isValid(); err != nil {
		writeError(http.StatusBadRequest, err, w, enc)
		return
	}

	result := qpbas.FindSolutions(parameters.SideLength)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	enc.Encode(solverResult{
		NumberOfResults: len(result.Solutions),
		CalculationTime: fmt.Sprintf("%v", result.CalculationTime),
	})
}
