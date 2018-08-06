package ntp

import (
	"github.com/czerwonk/atlas_exporter/config"
	"github.com/czerwonk/atlas_exporter/exporter"
)

const (
	ns  = "atlas"
	sub = "ntp"
)

// NewMeasurement returns a new instance of `exorter.Measurement` for a NTP measurement
func NewMeasurement(id string, cfg *config.Config) *exporter.Measurement {
	opts := []exporter.MeasurementOpt{}

	if cfg.FilterInvalidResults {
		opts = append(opts, exporter.WithValidator(&exporter.DefaultResultValidator{}))
	}

	return exporter.NewMeasurement(&ntpExporter{id}, opts...)
}
