package metric

import (
	"github.com/prometheus/client_golang/prometheus"
)

// MetricCollector implements the Prometheus Collector interface
type MetricCollector struct {
	Pk      string
	Metrics []MetricExporter
}

// NewMetricCollector creates a new collector
func NewMetricCollector(pk string, metrics []MetricExporter) prometheus.Collector {
	return &MetricCollector{Pk: pk, Metrics: metrics}
}

// Collect implements Prometheus Collector interface
func (c *MetricCollector) Collect(ch chan<- prometheus.Metric) {
	for _, m := range c.Metrics {
		m.Export(ch, c.Pk)
	}
}

// Describe implements Prometheus Collector interface
func (c *MetricCollector) Describe(ch chan<- *prometheus.Desc) {
	for _, m := range c.Metrics {
		m.Describe(ch)
	}
}
