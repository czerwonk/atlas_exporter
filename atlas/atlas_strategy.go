package atlas

import "context"

// Strategy defines an strategy to retrieve data for generating metrics
type Strategy interface {
	// MeasurementResults gets results for a list of measurements
	MeasurementResults(ctx context.Context, ids []string) ([]*AtlasMeasurement, error)
}
