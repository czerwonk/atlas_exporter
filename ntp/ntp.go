package ntp

import (
	"github.com/czerwonk/atlas_exporter/config"
	"github.com/czerwonk/atlas_exporter/exporter"
)

const (
	ns  = "atlas"
	sub = "ntp"
)

func NewResultHandler(id string, cfg *config.Config) *exporter.ResultHandler {
	opts := []exporter.ResultHandlerOpt{}

	if cfg.FilterInvalidResults {
		opts = append(opts, exporter.WithValidator(&exporter.DefaultResultValidator{}))
	}

	return exporter.NewResultHandler(&ntpExporter{id}, opts...)
}
