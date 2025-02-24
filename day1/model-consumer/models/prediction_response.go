package models

// PredictionResponse represents the subset of data we return
type PredictionResponse struct {
	FooID     int     `json:"foo_id"`
	Value     float64 `json:"value"`
	Certainty float64 `json:"certainty"`
}
