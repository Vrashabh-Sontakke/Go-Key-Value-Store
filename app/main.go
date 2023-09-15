package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// store represents a synchronized key-value store.
var store = struct {
	sync.RWMutex
	m map[string]string
}{m: make(map[string]string)}

// Define Prometheus metrics
var (
	// Latency metrics
	requestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "http_request_duration_seconds",
			Help: "HTTP request duration in seconds",
		},
		[]string{"endpoint"},
	)

	// HTTP status code metrics
	httpStatusCodes = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_status_codes_total",
			Help: "Total HTTP status codes per endpoint",
		},
		[]string{"endpoint", "status"},
	)

	// Total number of keys metric
	totalKeys = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "total_keys",
			Help: "Total number of keys in the DB",
		},
	)
)

func init() {
	// Register Prometheus metrics
	prometheus.MustRegister(requestDuration, httpStatusCodes, totalKeys)
}

// setHandler handles the HTTP POST request to set key-value pairs in the store.
func setHandler(w http.ResponseWriter, r *http.Request) {
	defer observeRequestDuration("setHandler", time.Now())
	key := r.FormValue("key")
	value := r.FormValue("value")

	store.Lock()
	store.m[key] = value
	store.Unlock()

	// Record HTTP status code
	recordHTTPStatusCode("setHandler", http.StatusOK)
}

// getHandler handles the HTTP GET request to retrieve a value by key from the store.
func getHandler(w http.ResponseWriter, r *http.Request) {
	defer observeRequestDuration("getHandler", time.Now())
	key := mux.Vars(r)["key"]

	store.RLock()
	value := store.m[key]
	store.RUnlock()

	w.Write([]byte(value))

	// Record HTTP status code
	recordHTTPStatusCode("getHandler", http.StatusOK)
}

// searchHandler handles the HTTP GET request to search for keys with a given prefix and suffix.
func searchHandler(w http.ResponseWriter, r *http.Request) {
	defer observeRequestDuration("searchHandler", time.Now())
	prefix := r.URL.Query().Get("prefix")
	suffix := r.URL.Query().Get("suffix")

	store.RLock()
	defer store.RUnlock()

	for key := range store.m {
		if strings.HasPrefix(key, prefix) && strings.HasSuffix(key, suffix) {
			w.Write([]byte(key + "\n"))
		}
	}

	// Record HTTP status code
	recordHTTPStatusCode("searchHandler", http.StatusOK)
}

// health check handler
func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

// Helper function to record request duration
func observeRequestDuration(endpoint string, start time.Time) {
	duration := time.Since(start).Seconds()
	requestDuration.WithLabelValues(endpoint).Observe(duration)
}

// Helper function to record HTTP status codes
func recordHTTPStatusCode(endpoint string, status int) {
	httpStatusCodes.WithLabelValues(endpoint, strconv.Itoa(status)).Inc()
}

// Helper function to update the total number of keys
func updateTotalKeys() {
	store.RLock()
	defer store.RUnlock()
	totalKeys.Set(float64(len(store.m)))
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/set", setHandler).Methods("POST")
	r.HandleFunc("/get/{key}", getHandler).Methods("GET")
	r.HandleFunc("/search", searchHandler).Methods("GET")

	r.HandleFunc("/healthCheck", healthCheckHandler).Methods("GET")

	http.Handle("/", r)

	// Start a goroutine to update the total number of keys regularly
	go func() {
		ticker := time.NewTicker(1 * time.Minute)
		for {
			select {
			case <-ticker.C:
				updateTotalKeys()
			}
		}
	}()

	fmt.Println("Server is listening on :8080...")

	// Expose Prometheus metrics on /metrics endpoint
	http.Handle("/metrics", promhttp.Handler())

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}