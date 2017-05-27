package ping

import (
	"strconv"

	"github.com/DNS-OARC/ripeatlas/measurement"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	ns  = "atlas"
	sub = "ping"
)

var (
	labels         []string
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

// PingMetricExporter exports metrics for PING measurement results
type PingMetricExporter struct {
	ProbeId   int
	DstAddr   string
	DstName   string
	Min       float64
	Max       float64
	Avg       float64
	Sent      int
	Rcvd      int
	Dup       int
	Ttl       int
	Size      int
	Asn       int
	IpVersion int
}

func init() {
	labels = make([]string, 0)
	labels = append(labels, "measurement", "probe", "dst_addr", "dst_name", "asn", "ip_version")

	successDesc = prometheus.NewDesc(prometheus.BuildFQName(ns, sub, "success"), "Destination was reachable", labels, nil)
	minLatencyDesc = prometheus.NewDesc(prometheus.BuildFQName(ns, sub, "min_latency"), "Minimum latency", labels, nil)
	maxLatencyDesc = prometheus.NewDesc(prometheus.BuildFQName(ns, sub, "max_latency"), "Maximum latency", labels, nil)
	avgLatencyDesc = prometheus.NewDesc(prometheus.BuildFQName(ns, sub, "avg_latency"), "Average latency", labels, nil)
	sentDesc = prometheus.NewDesc(prometheus.BuildFQName(ns, sub, "sent"), "Number of sent icmp requests", labels, nil)
	rcvdDesc = prometheus.NewDesc(prometheus.BuildFQName(ns, sub, "received"), "Number of received icmp repsponses", labels, nil)
	dupDesc = prometheus.NewDesc(prometheus.BuildFQName(ns, sub, "dup"), "Number of duplicate icmp repsponses", labels, nil)
	ttlDesc = prometheus.NewDesc(prometheus.BuildFQName(ns, sub, "ttl"), "Time to live remaining", labels, nil)
	sizeDesc = prometheus.NewDesc(prometheus.BuildFQName(ns, sub, "size"), "Size of ICMP packet", labels, nil)
}

// FromResult creates  metric exporter for PING measurement result
func FromResult(r *measurement.Result) *PingMetricExporter {
	return &PingMetricExporter{ProbeId: r.PrbId(), DstAddr: r.DstAddr(), DstName: r.DstName(), Max: r.Max(), Min: r.Min(), Rcvd: r.Rcvd(), Avg: r.Avg(), Sent: r.Sent(), Dup: r.Dup(), Ttl: r.Ttl(), Size: r.Size(), IpVersion: r.Af()}
}

// Export exports metrics for prometheus
func (m *PingMetricExporter) Export(ch chan<- prometheus.Metric, pk string) {
	labelValues := make([]string, 0)
	labelValues = append(labelValues, pk, strconv.Itoa(m.ProbeId), m.DstAddr, m.DstName, strconv.Itoa(m.Asn), strconv.Itoa(m.IpVersion))

	if m.Min > 0 {
		ch <- prometheus.MustNewConstMetric(successDesc, prometheus.GaugeValue, 1, labelValues...)
		ch <- prometheus.MustNewConstMetric(minLatencyDesc, prometheus.GaugeValue, m.Min, labelValues...)
		ch <- prometheus.MustNewConstMetric(maxLatencyDesc, prometheus.GaugeValue, m.Max, labelValues...)
		ch <- prometheus.MustNewConstMetric(avgLatencyDesc, prometheus.GaugeValue, m.Avg, labelValues...)
	} else {
		ch <- prometheus.MustNewConstMetric(successDesc, prometheus.GaugeValue, 0, labelValues...)
	}

	ch <- prometheus.MustNewConstMetric(sentDesc, prometheus.GaugeValue, float64(m.Sent), labelValues...)
	ch <- prometheus.MustNewConstMetric(rcvdDesc, prometheus.GaugeValue, float64(m.Rcvd), labelValues...)
	ch <- prometheus.MustNewConstMetric(dupDesc, prometheus.GaugeValue, float64(m.Dup), labelValues...)
	ch <- prometheus.MustNewConstMetric(ttlDesc, prometheus.GaugeValue, float64(m.Ttl), labelValues...)
	ch <- prometheus.MustNewConstMetric(sizeDesc, prometheus.GaugeValue, float64(m.Size), labelValues...)
}

// Describe exports metric descriptions for prometheus
func (m *PingMetricExporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- successDesc
	ch <- minLatencyDesc
	ch <- maxLatencyDesc
	ch <- avgLatencyDesc
	ch <- sentDesc
	ch <- rcvdDesc
	ch <- dupDesc
	ch <- ttlDesc
	ch <- sizeDesc
}

// SetAsn sets AS number for measurement result
func (m *PingMetricExporter) SetAsn(asn int) {
	m.Asn = asn
}

// IsValid returns whether an result is valid or not (e.g. IPv6 measurement and Probe does not support IPv6)
func (m *PingMetricExporter) IsValid() bool {
	return m.Asn > 0
}
