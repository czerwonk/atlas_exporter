package exporter

import (
	"github.com/DNS-OARC/ripeatlas/measurement"
	"github.com/prometheus/client_golang/prometheus"
)

// Histogram is the state of a single histogram of a measurement
type Histogram interface {
	ProcessResult(*measurement.Result)
	Hist() prometheus.Histogram
}
