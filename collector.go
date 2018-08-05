package main

import (
	"github.com/DNS-OARC/ripeatlas/measurement"
	"github.com/czerwonk/atlas_exporter/atlas"
	"github.com/prometheus/client_golang/prometheus"
)

type collector struct {
	measurements  []*atlas.AtlasMeasurement
	filterInvalid bool
}

func newCollector(measurements []*atlas.AtlasMeasurement, filterInvalid bool) *collector {
	return &collector{
		measurements:  measurements,
		filterInvalid: filterInvalid,
	}
}

// Collect implements Prometheus Collector interface
func (c *collector) Collect(ch chan<- prometheus.Metric) {
	for _, m := range c.measurements {
		for _, res := range c.filterInvalids(m.LatestResults(), m) {
			probe := m.Probes[res.PrbId()]

			if c.filterInvalid && !m.Exporter.IsValid(res, probe) {
				continue
			}

			m.Exporter.Export(res, probe, ch)
		}

		m.Exporter.ExportHistograms(c.filterInvalids(m.Results, m), ch)
	}
}

func (c *collector) filterInvalids(results []*measurement.Result, m *atlas.AtlasMeasurement) []*measurement.Result {
	if !c.filterInvalid {
		return results
	}

	valids := make([]*measurement.Result, 0)

	for _, res := range results {
		probe := m.Probes[res.PrbId()]

		if !m.Exporter.IsValid(res, probe) {
			continue
		}

		valids = append(valids, res)
	}

	return valids
}

// Describe implements Prometheus Collector interface
func (c *collector) Describe(ch chan<- *prometheus.Desc) {
	for _, m := range c.measurements {
		m.Exporter.Describe(ch)
	}
}
