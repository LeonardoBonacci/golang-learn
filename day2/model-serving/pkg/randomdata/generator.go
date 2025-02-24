package randomdata

import (
	"math/rand"
	"time"
)

// Prediction represents the generated model output
type Prediction struct {
	FooID      int     `json:"foo_id"`
	Timestamp  string  `json:"timestamp"`
	Value      float64 `json:"value"`
	Confidence float64 `json:"confidence"`
	Frequency  int     `json:"frequency"`
}

// RandomPrediction generates a pseudo-random prediction based on fooId and frequency
func RandomPrediction(fooId, frequency int) Prediction {
	rand.Seed(time.Now().UnixNano() + int64(fooId)) // Seed randomness with fooId

	// The frequency affects the range of randomness in value
	scaledValue := (rand.Float64() * 100) / float64(frequency+1) // Avoid divide by zero
	confidence := 0.5 + rand.Float64()*0.5                       // Confidence between 0.5 - 1.0

	return Prediction{
		FooID:      fooId,
		Timestamp:  time.Now().Format(time.RFC3339),
		Value:      scaledValue,
		Confidence: confidence,
		Frequency:  frequency,
	}
}
