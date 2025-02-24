package models

// ServingResponse represents the model serving response object
type ServingResponse struct {
	FooID      int     `json:"foo_id"`
	Timestamp  string  `json:"timestamp"`
	Value      float64 `json:"value"`
	Confidence float64 `json:"confidence"`
	Frequency  int     `json:"frequency"`
}
