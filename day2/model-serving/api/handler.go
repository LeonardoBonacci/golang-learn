package api

import (
	"encoding/json"
	"log"
	"model-serving/pkg/randomdata"
	"net/http"
	"strconv"
)

// GetPrediction handles GET /prediction?fooId=123&frequency=10
func GetPrediction(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Parse fooId from query parameters
	fooIdStr := r.URL.Query().Get("fooId")
	if fooIdStr == "" {
		http.Error(w, "Missing fooId parameter", http.StatusBadRequest)
		return
	}

	fooId, err := strconv.Atoi(fooIdStr)
	if err != nil {
		http.Error(w, "Invalid fooId parameter", http.StatusBadRequest)
		return
	}

	// Parse frequency from query parameters (default to 1 if not provided)
	frequencyStr := r.URL.Query().Get("frequency")
	frequency := 1
	if frequencyStr != "" {
		frequency, err = strconv.Atoi(frequencyStr)
		if err != nil {
			http.Error(w, "Invalid frequency parameter", http.StatusBadRequest)
			return
		}
	}

	// Log requests
	log.Printf("Incoming fooId %d and frequency %d", fooId, frequency)

	// Generate a prediction using fooId and frequency
	prediction := randomdata.RandomPrediction(fooId, frequency)

	// Encode and return JSON response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(prediction)
}
