package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"model-consumer/db"
	"model-consumer/models"
	"net/http"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

// FetchPrediction handles POST /fetch-prediction
func FetchPrediction(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Read environment variables (fallback to default)
	modelServingHost := os.Getenv("MODEL_SERVING_HOST")
	if modelServingHost == "" {
		modelServingHost = "localhost"
	}

	modelServingPort := os.Getenv("MODEL_SERVING_PORT")
	if modelServingPort == "" {
		modelServingPort = "8080"
	}

	// Decode request body
	var req models.PredictionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Prepare data for MongoDB (generate timestamp)
	record := models.PredictionRecord{
		FooID:     req.FooID,
		Timestamp: time.Now(),
	}

	// Get MongoDB collection
	collection := db.GetCollection("predictions")
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	// Query MongoDB for records with the same fooId
	filter := bson.M{"foo_id": req.FooID}
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		log.Printf("Failed to retrieve past predictions: %v", err)
		http.Error(w, "Error querying past predictions", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	// Decode results into slice
	var pastPredictions []models.PredictionRecord
	if err := cursor.All(ctx, &pastPredictions); err != nil {
		log.Printf("Error decoding past predictions: %v", err)
		http.Error(w, "Error processing past predictions", http.StatusInternalServerError)
		return
	}

	// Log retrieved past predictions (or process them as needed)
	log.Printf("Found %d past predictions for fooId %d in the last minute", len(pastPredictions), req.FooID)

	// Query MongoDB for records with the same fooId from the last minute
	//	collection := db.GetCollection("predictions")
	//	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Save request asynchronously in MongoDB
	log.Println("Storing request in MongoDB")
	go func(record models.PredictionRecord) {
		collection := db.GetCollection("predictions")
		log.Println("Storing request in MongoDB:", collection)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		_, err := collection.InsertOne(ctx, record)
		if err != nil {
			log.Println("Failed to store request in MongoDB:", err)
		} else {
			log.Println("Stored request in MongoDB:", record)
		}
	}(record)

	// Get model-serving response
	servingResponse, err := CallModelServing(req.FooID, len(pastPredictions))
	if err != nil {
		// Check if it's an HTTP error (contains status code)
		var httpErrMsg string
		var statusCode int

		if errors.As(err, &httpErrMsg) {
			statusCode = http.StatusInternalServerError
			log.Printf("Error calling model-serving: HTTP %d - %s", statusCode, httpErrMsg)
			http.Error(w, err.Error(), statusCode)
		} else {
			log.Printf("Error calling model-serving: %v", err)
			http.Error(w, "Failed to get prediction", http.StatusInternalServerError)
		}
		return
	}

	log.Printf("Received prediction: %+v\n", servingResponse)

	// Transform data
	transformedResponse := models.PredictionResponse{
		FooID:       servingResponse.FooID,
		Value:       servingResponse.Value,
		Certainty:   servingResponse.Confidence,
		NumRequests: len(pastPredictions),
	}

	// Return response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(transformedResponse)
}

// CallModelServing calls the model-serving service
func CallModelServing(fooID, frequency int) (*models.ServingResponse, error) {
	// Get host & port from environment variables
	host := os.Getenv("MODEL_SERVING_HOST")
	port := os.Getenv("MODEL_SERVING_PORT")

	// Construct the request URL
	url := fmt.Sprintf("http://%s:%s/prediction?fooId=%d&frequency=%d", host, port, fooID, frequency)

	// Send the GET request
	client := &http.Client{Timeout: 2 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		log.Printf("Error calling model-serving: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	// Check for HTTP errors
	if resp.StatusCode != http.StatusOK {
		log.Printf("model-serving returned status: %d", resp.StatusCode)
		return nil, fmt.Errorf("received non-OK status: %d", resp.StatusCode)
	}

	// Parse JSON response
	var servingResponse models.ServingResponse
	if err := json.NewDecoder(resp.Body).Decode(&servingResponse); err != nil {
		log.Printf("Failed to parse response: %v", err)
		return nil, err
	}

	return &servingResponse, nil
}
