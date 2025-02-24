package main

import (
	"fmt"
	"model-consumer/api"
)

func main() {
	fmt.Println("Starting model-consumer API...")
	api.SetupRoutes()
}
