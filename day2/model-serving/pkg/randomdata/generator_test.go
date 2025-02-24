package randomdata

import (
	"testing"
	"time"
)

func TestRandomPrediction(t *testing.T) {
	fooId := 42
	frequency := 5

	prediction := RandomPrediction(fooId, frequency)

	// Check if FooID matches input
	if prediction.FooID != fooId {
		t.Errorf("Expected FooID %d, got %d", fooId, prediction.FooID)
	}

	// Check if Frequency matches input
	if prediction.Frequency != frequency {
		t.Errorf("Expected Frequency %d, got %d", frequency, prediction.Frequency)
	}

	// Check if Value is within expected range
	if prediction.Value < 0 || prediction.Value > (100/float64(frequency+1)) {
		t.Errorf("Value out of expected range: got %f", prediction.Value)
	}

	// Check if Confidence is within expected range (0.5 - 1.0)
	if prediction.Confidence < 0.5 || prediction.Confidence > 1.0 {
		t.Errorf("Confidence out of range: got %f", prediction.Confidence)
	}

	// Check if Timestamp is valid
	_, err := time.Parse(time.RFC3339, prediction.Timestamp)
	if err != nil {
		t.Errorf("Invalid timestamp format: %s", prediction.Timestamp)
	}
}
