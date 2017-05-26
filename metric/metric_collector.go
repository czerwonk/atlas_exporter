package metric

import (
	"github.com/prometheus/client_golang/prometheus"
)

type MetricCollector struct {
	Pk      string
	Metrics []Metric
}

func NewMetricCollector(pk string, metrics []Metric) prometheus.Collector {
	return &MetricCollector{Pk: pk, Metrics: metrics}
}

func (c *MetricCollector) Collect(ch chan<- prometheus.Metric) {
	for _, m := range c.Metrics {
		m.GetMetrics(ch, c.Pk)
	}
}

func (c *MetricCollector) Describe(ch chan<- *prometheus.Desc) {
	for _, m := range c.Metrics {
		m.Describe(ch)
	}
}
