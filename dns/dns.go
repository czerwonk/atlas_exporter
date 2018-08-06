package dns

import (
	"github.com/czerwonk/atlas_exporter/config"
	"github.com/czerwonk/atlas_exporter/exporter"
)

const (
	ns  = "atlas"
	sub = "dns"
)

func NewResultHandler(id string, cfg *config.Config) *exporter.ResultHandler {
	opts := []exporter.ResultHandlerOpt{
		exporter.WithHistograms(newRttHistogram(id, cfg.HistogramBuckets.DNS.Rtt)),
	}

	if cfg.FilterInvalidResults {
		opts = append(opts, exporter.WithValidator(&exporter.DefaultResultValidator{}))
	}

	return exporter.NewResultHandler(&dnsExporter{id}, opts...)
}
