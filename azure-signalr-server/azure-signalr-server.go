package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

// NegotiateResult represents the result of SignalR's negotiate endpoint
type NegotiateResult struct {
	URL         string `json:"url"`
	AccessToken string `json:"accessToken"`
}

// SignalRMessage represents a message sent from the server to Azure SignalR
type SignalRMessage struct {
	Target    string      `json:"target"`
	Arguments interface{} `json:"arguments"`
}

// Variables receiving settings from the command line
var port *uint
var serviceName *string
var hubName *string
var key *string
var interval *int

// Azure SignalR URLs for broadcasting and negotiating
var broadcastURL string
var negotiateURL string

func main() {
	// Define command line arguments
	port = flag.Uint("port", 8080, "Port on which the server should listen")
	serviceName = flag.String("service-name", "mddd19", "Azure SignalR service name")
	hubName = flag.String("hub-name", "mddd19", "Azure SignalR service name")
	key = flag.String("key", "", "Azure SignalR access key")
	interval = flag.Int("interval", 2, "Seconds between broadcasting random values")

	// Parse the command line
	flag.Parse()

	// Allow CORS
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	// Build broadcast and negotiate URL for Azure SignalR
	broadcastURL = fmt.Sprintf("https://%s.service.signalr.net/api/v1/hubs/%s", *serviceName, *hubName)
	negotiateURL = fmt.Sprintf("https://%s.service.signalr.net/client/?hub=%s", *serviceName, *hubName)

	// Start timer that regularly sends out random values
	sendValues()

	// Setup router and start web server
	router := mux.NewRouter()
	router.HandleFunc("/", getIndex).Methods("GET")
	router.HandleFunc(fmt.Sprintf("/%s/negotiate", *hubName), negotiate).Methods("POST")
	fmt.Printf("Server starting, will listen to port %d...\n", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), handlers.CORS(originsOk, headersOk, methodsOk)(router)))
}

// getIndex reads the index.html file and sends it to the client
func getIndex(w http.ResponseWriter, r *http.Request) {
	indexData, _ := ioutil.ReadFile("index.html")
	w.WriteHeader(http.StatusOK)
	w.Write(indexData)
}

// negotiate implements the negotiate endpoint necessary for Azure SignalR
func negotiate(w http.ResponseWriter, r *http.Request) {
	// Create JWT token for Azure SignalR
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"aud": negotiateURL,
		"exp": time.Now().UTC().Add(time.Hour * 1).Unix(),
	})
	tokenString, err := token.SignedString([]byte(*key))
	if err != nil {
		panic(err)
	}

	// Build negotiate result and send it back to the client
	payload := NegotiateResult{
		URL:         negotiateURL,
		AccessToken: tokenString,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(payload)
}

// sendValues triggers the broadcasting of a value via SignalR to all connected clients
func sendValues() {
	ticker := time.NewTicker(time.Duration(*interval) * time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				sendValue()
			}
		}
	}()
}

// sendValue broadcasta a message
func sendValue() {
	// Create JWT token for Azure SignalR
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"aud": broadcastURL,
		"exp": time.Now().UTC().Add(time.Hour * 1).Unix(),
	})
	tokenString, err := token.SignedString([]byte(*key))

	// Generate SignalR message with random value
	arguments := []int{rand.Intn(100)}
	payload := SignalRMessage{
		Target:    "measurement",
		Arguments: arguments,
	}

	// POST message to Azure SignalR, it cares for broadcasting to all clients
	jsonPayload, _ := json.Marshal(payload)
	req, err := http.NewRequest("POST", broadcastURL, bytes.NewBuffer(jsonPayload))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tokenString))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	if resp.StatusCode >= 300 {
		defer resp.Body.Close()

		fmt.Println("Response Status:", resp.Status)
		fmt.Println("Response Headers:", resp.Header)
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println("response Body:", string(body))
	} else {
		fmt.Print("Successfully broadcasted value\n")
	}
}
