package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"github.com/gorilla/mux"
)

func TestSetHandler(t *testing.T) {
	// An http request with parameters
	req, err := http.NewRequest("POST", "/set?key=testkey&value=testvalue", nil)
	if err != nil {
		t.Fatal(err)
	}

	// ResponseRecorder
	rr := httptest.NewRecorder()

	// Test router and handler
	r := mux.NewRouter()
	r.HandleFunc("/set", setHandler).Methods("POST")

	// Serve the request to the handler
	r.ServeHTTP(rr, req)

	// Status code check
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status %d; got %d", http.StatusOK, rr.Code)
	}

	// Check if the key-value pair is stored correctly
	store.RLock()
	defer store.RUnlock()
	if store.m["testkey"] != "testvalue" {
		t.Errorf("Expected value 'testvalue'; got '%s'", store.m["testkey"])
	}
}

func TestGetHandler(t *testing.T) {
	// add test key-value pair to the store
	store.Lock()
	store.m["testkey"] = "testvalue"
	store.Unlock()

	// request to get the value
	req, err := http.NewRequest("GET", "/get/testkey", nil)
	if err != nil {
		t.Fatal(err)
	}

	// ResponseRecorder
	rr := httptest.NewRecorder()

	// test router and handler
	r := mux.NewRouter()
	r.HandleFunc("/get/{key}", getHandler).Methods("GET")

	// Serve the request to the handler
	r.ServeHTTP(rr, req)

	// Check the status code
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status %d; got %d", http.StatusOK, rr.Code)
	}

	// Check if the response body contains the correct value
	expected := "testvalue"
	if strings.TrimSpace(rr.Body.String()) != expected {
		t.Errorf("Expected body '%s'; got '%s'", expected, rr.Body.String())
	}
}

func TestSearchHandler(t *testing.T) {
	// Add test key-value pairs to the store
	store.Lock()
	store.m["abc-1"] = "value1"
	store.m["abc-2"] = "value2"
	store.m["xyz-1"] = "value3"
	store.m["xyz-2"] = "value4"
	store.Unlock()

	// request to search for keys with a prefix
	req, err := http.NewRequest("GET", "/search?prefix=abc&suffix=-1", nil)
	if err != nil {
		t.Fatal(err)
	}

	// ResponseRecorder
	rr := httptest.NewRecorder()

	// test router and handler
	r := mux.NewRouter()
	r.HandleFunc("/search", searchHandler).Methods("GET")

	// Serve the request to the handler
	r.ServeHTTP(rr, req)

	// Check the status code
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status %d; got %d", http.StatusOK, rr.Code)
	}

	// Check if the response body contains the expected keys
	expected := "abc-1\n"
	actual := rr.Body.String()
	if actual != expected {
		t.Errorf("Expected body '%s'; got '%s'", expected, actual)
	}
}


