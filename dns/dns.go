package dns

import (
	"github.com/czerwonk/atlas_exporter/config"
	"github.com/czerwonk/atlas_exporter/exporter"
)

const (
	ns  = "atlas"
	sub = "dns"
)

// NewMeasurement returns a new instance of `exorter.Measurement` for a DNS measurement
func NewMeasurement(id, ipVersion string, cfg *config.Config) *exporter.Measurement {
	opts := []exporter.MeasurementOpt{
		exporter.WithHistograms(newRttHistogram(id, ipVersion, cfg.HistogramBuckets.DNS.Rtt)),
	}

	if cfg.FilterInvalidResults {
		opts = append(opts, exporter.WithValidator(&exporter.DefaultResultValidator{}))
	}

	return exporter.NewMeasurement(&dnsExporter{id}, opts...)
}
