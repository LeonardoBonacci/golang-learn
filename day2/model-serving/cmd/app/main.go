package main

import (
	"fmt"
	"log"
	"model-serving/api"
	"net/http"
)

func main() {
	// Define the server address
	port := 8080
	addr := fmt.Sprintf(":%d", port)

	// Define the HTTP handler
	http.HandleFunc("/prediction", api.GetPrediction)

	// Start the server
	log.Printf("Model-serving API is running on http://localhost%s\n", addr)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
