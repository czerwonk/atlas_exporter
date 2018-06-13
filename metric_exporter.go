package main

import (
	"github.com/DNS-OARC/ripeatlas/measurement"
	"github.com/czerwonk/atlas_exporter/probe"
	"github.com/prometheus/client_golang/prometheus"
)

// MetricExporter defines a set of metrics for an ATLAS measurement type
type MetricExporter interface {

	// Export exports a prometheus metric
	Export(id string, res *measurement.Result, probe *probe.Probe, ch chan<- prometheus.Metric)

	// Describes metrics exported for this measurement type
	Describe(ch chan<- *prometheus.Desc)

	// IsValid returns if a meaurement result is valid (can be filtered when needed)
	IsValid(res *measurement.Result, probe *probe.Probe) bool
}
