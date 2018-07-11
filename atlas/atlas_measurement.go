package atlas

import (
	"github.com/DNS-OARC/ripeatlas/measurement"
	"github.com/czerwonk/atlas_exporter/probe"
)

// AtlasMeasurement contains all measurement results for one measurement including probe information
type AtlasMeasurement struct {
	ID       string
	Results  []*measurement.Result
	Exporter MetricExporter
	Probes   map[int]*probe.Probe
}
