package http

import (
	"github.com/DNS-OARC/ripeatlas/measurement"
	"github.com/czerwonk/atlas_exporter/exporter"
	"github.com/prometheus/client_golang/prometheus"
)

type rttHistogram struct {
	rtt prometheus.Histogram
}

func newRttHistogram(id, ipVersion string, buckets []float64) exporter.Histogram {
	if buckets == nil {
		buckets = []float64{100, 200, 500, 1000}
	}

	return &rttHistogram{
		rtt: prometheus.NewHistogram(prometheus.HistogramOpts{
			Namespace: ns,
			Subsystem: sub,
			Name:      "rtt_hist",
			Buckets:   buckets,
			Help:      "Histogram of round trip times over all HTTP requests",
			ConstLabels: prometheus.Labels{
				"measurement": id,
				"ip_version":  ipVersion,
			},
		}),
	}
}

func (h *rttHistogram) ProcessResult(r *measurement.Result) {
	for _, p := range r.HttpResults() {
		if p.Rt() > 0 {
			h.rtt.Observe(p.Rt())
		}
	}
}

func (h *rttHistogram) Hist() prometheus.Histogram {
	return h.rtt
}
