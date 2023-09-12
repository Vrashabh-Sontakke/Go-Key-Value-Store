package main

import (
	"net/http"
	"strings"
	"sync"

	"github.com/gorilla/mux"
)

// store represents a synchronized key-value store.
var store = struct {
	sync.RWMutex
	m map[string]string
}{m: make(map[string]string)}

// setHandler handles the HTTP POST request to set key-value pairs in the store.
func setHandler(w http.ResponseWriter, r *http.Request) {
	key := r.FormValue("key")
	value := r.FormValue("value")

	store.Lock()
	store.m[key] = value
	store.Unlock()
}

// getHandler handles the HTTP GET request to retrieve a value by key from the store.
func getHandler(w http.ResponseWriter, r *http.Request) {
	key := mux.Vars(r)["key"]

	store.RLock()
	value := store.m[key]
	store.RUnlock()

	w.Write([]byte(value))
}

// searchHandler handles the HTTP GET request to search for keys with a given prefix and suffix.
func searchHandler(w http.ResponseWriter, r *http.Request) {
	prefix := r.URL.Query().Get("prefix")
	suffix := r.URL.Query().Get("suffix")

	store.RLock()
	defer store.RUnlock()

	for key := range store.m {
		if strings.HasPrefix(key, prefix) && strings.HasSuffix(key, suffix) {
			w.Write([]byte(key + "\n"))
		}
	}
}

// health check handler
func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/set", setHandler).Methods("POST")
	r.HandleFunc("/get/{key}", getHandler).Methods("GET")
	r.HandleFunc("/search", searchHandler).Methods("GET")

	r.HandleFunc("/healthCheck", healthCheckHandler).Methods("GET")

	http.Handle("/", r)

	http.ListenAndServe(":8080", nil)
}
