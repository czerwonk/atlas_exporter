// SPDX-License-Identifier: LGPL-3.0-or-later

package dns

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
		buckets = []float64{10, 20, 50, 100}
	}

	return &rttHistogram{
		rtt: prometheus.NewHistogram(prometheus.HistogramOpts{
			Namespace: ns,
			Subsystem: sub,
			Name:      "rtt_hist",
			Buckets:   buckets,
			Help:      "Histogram of round trip times over all DNS requests",
			ConstLabels: prometheus.Labels{
				"measurement": id,
				"ip_version":  ipVersion,
			},
		}),
	}
}

func (h *rttHistogram) ProcessResult(r *measurement.Result) {
	if rs := r.DnsResultsets(); len(rs) > 0 {
		for _, s := range rs {
			if s == nil || s.Result() == nil {
				continue
			}

			rt := s.Result().Rt()
			if rt > 0 {
				h.rtt.Observe(rt)
			}
		}
		return
	}

	if dr := r.DnsResult(); dr != nil {
		if rt := dr.Rt(); rt > 0 {
			h.rtt.Observe(rt)
		}
	}
}

func (h *rttHistogram) Hist() prometheus.Histogram {
	return h.rtt
}
