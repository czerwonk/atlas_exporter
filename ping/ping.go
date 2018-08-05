package ping

import (
	"strconv"

	"github.com/DNS-OARC/ripeatlas/measurement"
	"github.com/czerwonk/atlas_exporter/exporter"
	"github.com/czerwonk/atlas_exporter/probe"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	ns  = "atlas"
	sub = "ping"
)

var (
	labels         []string
	rttHistDesc    *prometheus.Desc
	successDesc    *prometheus.Desc
	minLatencyDesc *prometheus.Desc
	maxLatencyDesc *prometheus.Desc
	avgLatencyDesc *prometheus.Desc
	sentDesc       *prometheus.Desc
	rcvdDesc       *prometheus.Desc
	dupDesc        *prometheus.Desc
	ttlDesc        *prometheus.Desc
	sizeDesc       *prometheus.Desc
)

type pingMetricExporter struct {
	id      string
	rttHist prometheus.Histogram
}

// NewExporter creates a exporter for PING measurement results
func NewExporter(id string, buckets []float64) exporter.MetricExporter {
	if buckets == nil {
		buckets = []float64{10, 20, 50, 100}
	}

	hist := prometheus.NewHistogram(prometheus.HistogramOpts{
		Namespace: ns,
		Subsystem: sub,
		Name:      "rtt_hist",
		Buckets:   buckets,
		Help:      "Histogram of round trip times over all ICMP requests",
		ConstLabels: prometheus.Labels{
			"measurement": id,
		},
	})

	return &pingMetricExporter{
		id:      id,
		rttHist: hist,
	}
}

func init() {
	labels = []string{"measurement", "probe", "dst_addr", "dst_name", "asn", "ip_version", "country_code", "lat", "long"}
	successDesc = prometheus.NewDesc(prometheus.BuildFQName(ns, sub, "success"), "Destination was reachable", labels, nil)
	minLatencyDesc = prometheus.NewDesc(prometheus.BuildFQName(ns, sub, "min_latency"), "Minimum latency", labels, nil)
	maxLatencyDesc = prometheus.NewDesc(prometheus.BuildFQName(ns, sub, "max_latency"), "Maximum latency", labels, nil)
	avgLatencyDesc = prometheus.NewDesc(prometheus.BuildFQName(ns, sub, "avg_latency"), "Average latency", labels, nil)
	sentDesc = prometheus.NewDesc(prometheus.BuildFQName(ns, sub, "sent"), "Number of sent icmp requests", labels, nil)
	rcvdDesc = prometheus.NewDesc(prometheus.BuildFQName(ns, sub, "received"), "Number of received icmp repsponses", labels, nil)
	dupDesc = prometheus.NewDesc(prometheus.BuildFQName(ns, sub, "dup"), "Number of duplicate icmp repsponses", labels, nil)
	ttlDesc = prometheus.NewDesc(prometheus.BuildFQName(ns, sub, "ttl"), "Time-to-live field in the response", labels, nil)
	sizeDesc = prometheus.NewDesc(prometheus.BuildFQName(ns, sub, "size"), "Size of ICMP packet", labels, nil)
}

// Export exports a prometheus metric
func (m *pingMetricExporter) Export(res *measurement.Result, probe *probe.Probe, ch chan<- prometheus.Metric) {
	labelValues := []string{
		m.id,
		strconv.Itoa(probe.ID),
		res.DstAddr(),
		res.DstName(),
		strconv.Itoa(probe.ASNForIPVersion(res.Af())),
		strconv.Itoa(res.Af()),
		probe.CountryCode,
		probe.Latitude(),
		probe.Longitude(),
	}

	if res.Min() > 0 {
		ch <- prometheus.MustNewConstMetric(successDesc, prometheus.GaugeValue, 1, labelValues...)
		ch <- prometheus.MustNewConstMetric(minLatencyDesc, prometheus.GaugeValue, res.Min(), labelValues...)
		ch <- prometheus.MustNewConstMetric(maxLatencyDesc, prometheus.GaugeValue, res.Max(), labelValues...)
		ch <- prometheus.MustNewConstMetric(avgLatencyDesc, prometheus.GaugeValue, res.Avg(), labelValues...)
	} else {
		ch <- prometheus.MustNewConstMetric(successDesc, prometheus.GaugeValue, 0, labelValues...)
	}

	ch <- prometheus.MustNewConstMetric(sentDesc, prometheus.GaugeValue, float64(res.Sent()), labelValues...)
	ch <- prometheus.MustNewConstMetric(rcvdDesc, prometheus.GaugeValue, float64(res.Rcvd()), labelValues...)
	ch <- prometheus.MustNewConstMetric(dupDesc, prometheus.GaugeValue, float64(res.Dup()), labelValues...)
	ch <- prometheus.MustNewConstMetric(ttlDesc, prometheus.GaugeValue, float64(res.Ttl()), labelValues...)
	ch <- prometheus.MustNewConstMetric(sizeDesc, prometheus.GaugeValue, float64(res.Size()), labelValues...)
}

// ExportHistograms exports aggregated metrics for the measurement
func (m *pingMetricExporter) ExportHistograms(res []*measurement.Result, ch chan<- prometheus.Metric) {
	for _, r := range res {
		for _, p := range r.PingResults() {
			if p.Rtt() > 0 {
				m.rttHist.Observe(p.Rtt())
			}
		}
	}

	m.rttHist.Collect(ch)
}

// Describe exports metric descriptions for Prometheus
func (m *pingMetricExporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- successDesc
	ch <- minLatencyDesc
	ch <- maxLatencyDesc
	ch <- avgLatencyDesc
	ch <- sentDesc
	ch <- rcvdDesc
	ch <- dupDesc
	ch <- ttlDesc
	ch <- sizeDesc
	ch <- m.rttHist.Desc()
}

// IsValid returns whether an result is valid or not (e.g. IPv6 measurement and Probe does not support IPv6)
func (m *pingMetricExporter) IsValid(res *measurement.Result, probe *probe.Probe) bool {
	return probe.ASNForIPVersion(res.Af()) > 0
}
