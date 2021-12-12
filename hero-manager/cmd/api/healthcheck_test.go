package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHealthcheck(t *testing.T) {
	// Setup mock configuration for "Test" environment
	app := &application{
		config: config{
			env: "Test",
		},
	}

	// Test healthcheck handler function in isolation
	rr := httptest.NewRecorder()
	r, err := http.NewRequest(http.MethodGet, "/v1/healthcheck", nil)
	if err != nil {
		t.Fatal(err)
	}

	app.healthcheckHandler(rr, r)

	// Get and verify result
	rs := rr.Result()
	if rs.StatusCode != http.StatusOK {
		t.Errorf("want %d; got %d", http.StatusOK, rs.StatusCode)
	}

	healthResult := make(map[string]string)
	err = json.NewDecoder(rs.Body).Decode(&healthResult)
	if err != nil {
		t.Fatal(err)
		return
	}

	assert.Equal(t, "Test", healthResult["environment"])
}

func TestHealthcheckEndToEnd(t *testing.T) {
	// Setup mock configuration for "Test" environment
	app := &application{
		config: config{
			env: "Test",
		},
	}

	// Run HTTPS server on random port
	ts := httptest.NewTLSServer(app.routes())
	defer ts.Close()

	rs, err := ts.Client().Get(ts.URL + "/v1/healthcheck")
	if err != nil {
		t.Fatal(err)
	}

	if rs.StatusCode != http.StatusOK {
		t.Errorf("want %d; got %d", http.StatusOK, rs.StatusCode)
	}

	healthResult := make(map[string]string)
	err = json.NewDecoder(rs.Body).Decode(&healthResult)
	if err != nil {
		t.Fatal(err)
		return
	}

	assert.Equal(t, "Test", healthResult["environment"])
}
