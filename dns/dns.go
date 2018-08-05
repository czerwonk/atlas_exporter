package dns

import (
	"strconv"

	"github.com/czerwonk/atlas_exporter/exporter"
	"github.com/czerwonk/atlas_exporter/probe"

	"github.com/DNS-OARC/ripeatlas/measurement"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	ns  = "atlas"
	sub = "dns"
)

var (
	labels      []string
	successDesc *prometheus.Desc
	rttDesc     *prometheus.Desc
)

func init() {
	labels = []string{"measurement", "probe", "dst_addr", "asn", "ip_version", "country_code", "lat", "long"}

	successDesc = prometheus.NewDesc(prometheus.BuildFQName(ns, sub, "success"), "Destination was reachable", labels, nil)
	rttDesc = prometheus.NewDesc(prometheus.BuildFQName(ns, sub, "rtt"), "Roundtrip time in ms", labels, nil)
}

type dnsMetricExporter struct {
	id      string
	rttHist prometheus.Histogram
}

// NewExporter creates a exporter for DNS measurement results
func NewExporter(id string, buckets []float64) exporter.MetricExporter {
	if buckets == nil {
		buckets = prometheus.LinearBuckets(10, 10, 100)
	}

	hist := prometheus.NewHistogram(prometheus.HistogramOpts{
		Namespace: ns,
		Subsystem: sub,
		Name:      "rtt_hist",
		Buckets:   buckets,
		Help:      "Histogram of round trip times over all DNS requests",
		ConstLabels: prometheus.Labels{
			"measurement": id,
		},
	})

	return &dnsMetricExporter{
		id:      id,
		rttHist: hist,
	}
}

// Export exports a prometheus metric
func (m *dnsMetricExporter) Export(res *measurement.Result, probe *probe.Probe, ch chan<- prometheus.Metric) {
	labelValues := []string{
		m.id,
		strconv.Itoa(probe.ID),
		res.DstAddr(),
		strconv.Itoa(probe.ASNForIPVersion(res.Af())),
		strconv.Itoa(res.Af()),
		probe.CountryCode,
		probe.Latitude(),
		probe.Longitude(),
	}

	var rtt float64
	if res.DnsResult() != nil {
		rtt = res.DnsResult().Rt()
	}

	if rtt > 0 {
		ch <- prometheus.MustNewConstMetric(successDesc, prometheus.GaugeValue, 1, labelValues...)
		ch <- prometheus.MustNewConstMetric(rttDesc, prometheus.GaugeValue, rtt, labelValues...)
	} else {
		ch <- prometheus.MustNewConstMetric(successDesc, prometheus.GaugeValue, 0, labelValues...)
	}
}

// ExportHistograms exports aggregated metrics for the measurement
func (m *dnsMetricExporter) ExportHistograms(res []*measurement.Result, ch chan<- prometheus.Metric) {
	for _, r := range res {
		if r.DnsResult() == nil {
			continue
		}

		if r.DnsResult().Rt() > 0 {
			m.rttHist.Observe(r.DnsResult().Rt())
		}
	}

	m.rttHist.Collect(ch)
}

// Describe exports metric descriptions for Prometheus
func (m *dnsMetricExporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- successDesc
	ch <- rttDesc
	ch <- m.rttHist.Desc()
}

// IsValid returns whether an result is valid or not (e.g. IPv6 measurement and Probe does not support IPv6)
func (m *dnsMetricExporter) IsValid(res *measurement.Result, probe *probe.Probe) bool {
	return probe.ASNForIPVersion(res.Af()) > 0
}
