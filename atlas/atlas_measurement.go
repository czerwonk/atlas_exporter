package atlas

import (
	"github.com/DNS-OARC/ripeatlas/measurement"
	"github.com/czerwonk/atlas_exporter/exporter"
	"github.com/czerwonk/atlas_exporter/probe"
)

// AtlasMeasurement contains all measurement results for one measurement including probe information
type AtlasMeasurement struct {
	ID       string
	Results  []*measurement.Result
	Exporter exporter.MetricExporter
	Probes   map[int]*probe.Probe
}

// LatestResults returns the lastest result for each probe
func (a *AtlasMeasurement) LatestResults() []*measurement.Result {
	res := make([]*measurement.Result, 0)
	byProbe := make(map[int]*measurement.Result)

	for i := len(a.Results) - 1; len(byProbe) < len(a.Probes) && i >= 0; i-- {
		m := a.Results[i]
		if _, found := byProbe[m.PrbId()]; found {
			continue
		}

		res = append(res, m)
		byProbe[m.PrbId()] = m
	}

	return res
}
