package api

import (
	"fmt"
	"net/http"
)

// SetupRoutes initializes the API routes
func SetupRoutes() {
	http.HandleFunc("/fetch-prediction", FetchPrediction)
	fmt.Println("model-consumer running on http://localhost:8081")
	http.ListenAndServe(":8081", nil)
}
