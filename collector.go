package main

import (
	"github.com/czerwonk/atlas_exporter/exporter"
	"github.com/prometheus/client_golang/prometheus"
)

type collector struct {
	measurements []*exporter.Measurement
}

func newCollector(measurements []*exporter.Measurement) *collector {
	return &collector{
		measurements: measurements,
	}
}

// Collect implements Prometheus Collector interface
func (c *collector) Collect(ch chan<- prometheus.Metric) {
	for _, m := range c.measurements {
		m.Collect(ch)
	}
}

// Describe implements Prometheus Collector interface
func (c *collector) Describe(ch chan<- *prometheus.Desc) {
	for _, m := range c.measurements {
		m.Describe(ch)
	}
}
