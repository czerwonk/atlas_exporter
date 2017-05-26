package metric

import (
	"github.com/prometheus/client_golang/prometheus"
)

type Metric interface {
	GetMetrics(ch chan<- prometheus.Metric, pk string)
	Describe(ch chan<- *prometheus.Desc)

	SetAsn(asn int)
	Isvalid() bool
}
