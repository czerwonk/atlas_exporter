package metric

import (
	"github.com/prometheus/client_golang/prometheus"
)

// Metric exporter for measurement result (single probe)
type MetricExporter interface {
	// Exports metrics for prometheus
	GetMetrics(ch chan<- prometheus.Metric, pk string)

	// Exports metric descriptions for prometheus
	Describe(ch chan<- *prometheus.Desc)

	// Sets AN number for measurement result
	SetAsn(asn int)

	// Gets whether an result is valid (e.g. IPv6 measurement and Probe does not support IPv6)
	Isvalid() bool
}
