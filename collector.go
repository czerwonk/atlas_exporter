package main

import (
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
		for _, res := range m.Results {
			probe := m.Probes[res.PrbId()]

			if c.filterInvalid && !m.Exporter.IsValid(res, probe) {
				continue
			}

			m.Exporter.Export(m.ID, res, probe, ch)
		}
	}
}

// Describe implements Prometheus Collector interface
func (c *collector) Describe(ch chan<- *prometheus.Desc) {
	for _, m := range c.measurements {
		m.Exporter.Describe(ch)
	}
}
