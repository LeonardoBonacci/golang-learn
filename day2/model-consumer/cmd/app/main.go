package main

import (
	"fmt"
	"model-consumer/api"
	"os"
)

func main() {
	fmt.Println("Starting model-consumer API...")

	// Print environment variables
	modelServingHost := os.Getenv("MODEL_SERVING_HOST")
	modelServingPort := os.Getenv("MODEL_SERVING_PORT")
	fmt.Printf("Pointing to model-serving at %s:%s\n", modelServingHost, modelServingPort)

	api.SetupRoutes()
}
