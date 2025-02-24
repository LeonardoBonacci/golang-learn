package api

import (
	"context"
	"encoding/json"
	"log"
	"model-consumer/db"
	"model-consumer/models"
	"net/http"
	"os"
	"time"
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

	// Generate current timestamp
	now := time.Now()
	oneMinuteAgo := now.Add(-1 * time.Minute)

	// Query MongoDB for records with the same fooId from the last minute
	collection := db.GetCollection("predictions")
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	filter := map[string]interface{}{
		"foo_id": req.FooID,
		"timestamp": map[string]interface{}{
			"$gte": oneMinuteAgo,
			"$lte": now,
		},
	}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		log.Printf("Failed to retrieve past predictions: %v", err)
		http.Error(w, "Error querying past predictions", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

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

	// Parse model-serving response
	fullResponse := struct {
		FooID      int
		Value      float64
		Confidence float64
	}{
		FooID:      123,
		Value:      45.67,
		Confidence: 0.95,
	}

	// Transform data
	transformedResponse := models.PredictionResponse{
		FooID:     fullResponse.FooID,
		Value:     fullResponse.Value,
		Certainty: fullResponse.Confidence,
	}

	// Return response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(transformedResponse)
}
