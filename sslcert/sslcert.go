package sslcert

import (
	"github.com/czerwonk/atlas_exporter/config"
	"github.com/czerwonk/atlas_exporter/exporter"
)

const (
	ns  = "atlas"
	sub = "sslcert"
)

// NewMeasurement returns a new instance of `exorter.Measurement` for a SSL measurement
func NewMeasurement(id string, cfg *config.Config) *exporter.Measurement {
	opts := []exporter.MeasurementOpt{}

	if cfg.FilterInvalidResults {
		opts = append(opts, exporter.WithValidator(&exporter.DefaultResultValidator{}))
	}

	return exporter.NewMeasurement(&sslCertExporter{id}, opts...)
}
