package models

import "time"

// PredictionRecord represents the stored MongoDB record
type PredictionRecord struct {
	FooID     int       `bson:"foo_id" json:"foo_id"`
	Timestamp time.Time `bson:"timestamp" json:"timestamp"`
}
