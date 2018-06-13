package main

import (
	"github.com/prometheus/client_golang/prometheus"
)

type collector struct {
	measurements  []*atlasMeasurement
	filterInvalid bool
}

func newCollector(measurements []*atlasMeasurement, filterInvalid bool) *collector {
	return &collector{
		measurements:  measurements,
		filterInvalid: filterInvalid,
	}
}

// Collect implements Prometheus Collector interface
func (c *collector) Collect(ch chan<- prometheus.Metric) {
	for _, m := range c.measurements {
		for _, res := range m.results {
			probe := m.probes[res.PrbId()]

			if c.filterInvalid && !m.exporter.IsValid(res, probe) {
				continue
			}

			m.exporter.Export(m.id, res, probe, ch)
		}
	}
}

// Describe implements Prometheus Collector interface
func (c *collector) Describe(ch chan<- *prometheus.Desc) {
	for _, m := range c.measurements {
		m.exporter.Describe(ch)
	}
}
