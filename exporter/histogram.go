package exporter

import (
	"github.com/DNS-OARC/ripeatlas/measurement"
	"github.com/prometheus/client_golang/prometheus"
)

type Histogram interface {
	ProcessResult(*measurement.Result)
	Hist() prometheus.Histogram
}
