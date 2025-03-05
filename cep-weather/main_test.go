package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestValidCEP(t *testing.T) {
	// Create a request to pass to our handler
	req, err := http.NewRequest("GET", "/weather/01001001", nil)
	assert.NoError(t, err)

	// Create a ResponseRecorder to record the response
	w := httptest.NewRecorder()

	// Create a router and serve the request
	r := mux.NewRouter()
	r.HandleFunc("/weather/{cep}", weatherHandler)
	r.ServeHTTP(w, req)

	// Check the status code
	assert.Equal(t, http.StatusOK, w.Code)

	// You might want to add more specific checks on the response body
}

func TestInvalidCEPFormat(t *testing.T) {
	req, err := http.NewRequest("GET", "/weather/123", nil)
	assert.NoError(t, err)

	w := httptest.NewRecorder()

	r := mux.NewRouter()
	r.HandleFunc("/weather/{cep}", weatherHandler)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	assert.Contains(t, w.Body.String(), "invalid zipcode")
}

func TestNonExistentCEP(t *testing.T) {
	req, err := http.NewRequest("GET", "/weather/99999999", nil)
	assert.NoError(t, err)

	w := httptest.NewRecorder()

	r := mux.NewRouter()
	r.HandleFunc("/weather/{cep}", weatherHandler)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "can not find zipcode")
}