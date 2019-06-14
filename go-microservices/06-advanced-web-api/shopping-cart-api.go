package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type updateQuantitiesParameters struct {
	ID       int    `json:"id"`
	Quantity uint16 `json:"quantity"`
}

type applyCouponParameters struct {
	Coupon string `json:"coupon"`
}

var cartChannels CartChannels

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

	cartChannels = CreateCartChannels()
	go Cart(cartChannels)

	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	router := mux.NewRouter()
	router.HandleFunc("/ping", getPing).Methods("GET")
	router.HandleFunc("/cart/reset", resetCart).Methods("POST")
	router.HandleFunc("/cart", getCart).Methods("GET")
	router.HandleFunc("/cart/quantities", updateQuantities).Methods("POST")
	router.HandleFunc("/cart/applyCoupon", applyCoupon).Methods("POST")
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

func resetCart(w http.ResponseWriter, r *http.Request) {
	reset := ReadCartOp{Resp: make(chan []CartItem)}
	cartChannels.resets <- reset
	cart := <-reset.Resp

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cart)
}

func getCart(w http.ResponseWriter, r *http.Request) {
	read := ReadCartOp{Resp: make(chan []CartItem)}
	cartChannels.reads <- read
	cart := <-read.Resp

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cart)
}

func updateQuantities(w http.ResponseWriter, r *http.Request) {
	// content type will always be json
	w.Header().Set("Content-Type", "application/json")

	// resulting update struct
	var uqs []updateQuantitiesParameters

	// create an encoder for writing response
	enc := json.NewEncoder(w)

	// decode incoming json from body
	err := json.NewDecoder(r.Body).Decode(&uqs)
	if err != nil {
		writeError(http.StatusBadRequest, err, w, enc)
		return
	}

	for _, uq := range uqs {
		update := UpdateQuantityOp{ID: uq.ID, Quantity: uq.Quantity, Resp: make(chan bool)}
		cartChannels.updates <- update
		success := <-update.Resp

		if !success {
			writeError(http.StatusBadRequest, fmt.Errorf("Shopping cart item with ID %d not found", uq.ID), w, enc)
			return
		}
	}

	read := ReadCartOp{Resp: make(chan []CartItem)}
	cartChannels.reads <- read
	cart := <-read.Resp

	// return success
	w.WriteHeader(http.StatusOK)
	enc.Encode(cart)
}

func applyCoupon(w http.ResponseWriter, r *http.Request) {
	// content type will always be json
	w.Header().Set("Content-Type", "application/json")

	var acp applyCouponParameters

	// create an encoder for writing response
	enc := json.NewEncoder(w)

	// decode incoming json from body
	err := json.NewDecoder(r.Body).Decode(&acp)
	if err != nil {
		writeError(http.StatusBadRequest, err, w, enc)
		return
	}

	coupon := ApplyCouponOp{Coupon: acp.Coupon, Resp: make(chan bool)}
	cartChannels.coupons <- coupon
	success := <-coupon.Resp

	if !success {
		writeError(http.StatusBadRequest, fmt.Errorf("Coupon %s is not valid", acp.Coupon), w, enc)
		return
	}

	read := ReadCartOp{Resp: make(chan []CartItem)}
	cartChannels.reads <- read
	cart := <-read.Resp

	// return success
	w.WriteHeader(http.StatusOK)
	enc.Encode(cart)
}
