package http

import (
	"github.com/czerwonk/atlas_exporter/config"
	"github.com/czerwonk/atlas_exporter/exporter"
)

const (
	ns  = "atlas"
	sub = "http"
)

func NewResultHandler(id string, cfg *config.Config) *exporter.ResultHandler {
	opts := []exporter.ResultHandlerOpt{
		exporter.WithHistograms(newRttHistogram(id, cfg.HistogramBuckets.HTTP.Rtt)),
	}

	if cfg.FilterInvalidResults {
		opts = append(opts, exporter.WithValidator(&exporter.DefaultResultValidator{}))
	}

	return exporter.NewResultHandler(&httpExporter{id}, opts...)
}
