package main

import (
	"net/http"
	"strings"
	"sync"

	"github.com/gorilla/mux"
)

var store = struct {
	sync.RWMutex
	m map[string]string
}{m: make(map[string]string)}

func setHandler(w http.ResponseWriter, r *http.Request) {
	key := r.FormValue("key")
	value := r.FormValue("value")

	store.Lock()
	store.m[key] = value
	store.Unlock()
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	key := mux.Vars(r)["key"]

	store.RLock()
	value := store.m[key]
	store.RUnlock()

	w.Write([]byte(value))
}

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

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/set", setHandler).Methods("POST")
	r.HandleFunc("/get/{key}", getHandler).Methods("GET")
	r.HandleFunc("/search", searchHandler).Methods("GET")
	http.Handle("/", r)

	http.ListenAndServe(":8080", nil)
}
