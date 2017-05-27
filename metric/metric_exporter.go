package metric

import (
	"github.com/prometheus/client_golang/prometheus"
)

// MetricExporter exports metrics for Atlas measurement results
type MetricExporter interface {

	// Export exports metrics for prometheus
	Export(ch chan<- prometheus.Metric, pk string)

	// Describe exports metric descriptions for prometheus
	Describe(ch chan<- *prometheus.Desc)

	// SetAsn sets AS number for measurement result
	SetAsn(asn int)

	// IsValid returns whether an result is valid or not (e.g. IPv6 measurement and Probe does not support IPv6)
	Isvalid() bool
}
