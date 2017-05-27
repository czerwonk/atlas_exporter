package dns

import (
	"strconv"

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
	labels = make([]string, 0)
	labels = append(labels, "measurement", "probe", "dst_addr", "asn", "ip_version")

	successDesc = prometheus.NewDesc(prometheus.BuildFQName(ns, sub, "success"), "Destination was reachable", labels, nil)
	rttDesc = prometheus.NewDesc(prometheus.BuildFQName(ns, sub, "rtt"), "Roundtrip time in ms", labels, nil)
}

// Metrics for DNS measurement results
type DnsMetricExporter struct {
	ProbeId   int
	DstAddr   string
	Rtt       float64
	Asn       int
	Success   int
	IpVersion int
}

// Creates metric exporter for DNS measurement result
func FromResult(r *measurement.Result) *DnsMetricExporter {
	var rtt float64
	if r.DnsResult() != nil {
		rtt = r.DnsResult().Rt()
	}

	var success int
	if r.DnsError() == nil {
		success = 1
	}

	return &DnsMetricExporter{ProbeId: r.PrbId(), DstAddr: r.DstAddr(), Rtt: rtt, Success: success, IpVersion: r.Af()}
}

// Exports metrics for prometheus
func (m *DnsMetricExporter) GetMetrics(ch chan<- prometheus.Metric, pk string) {
	labelValues := make([]string, 0)
	labelValues = append(labelValues, pk, strconv.Itoa(m.ProbeId), m.DstAddr, strconv.Itoa(m.Asn), strconv.Itoa(m.IpVersion))

	if m.Rtt > 0 {
		ch <- prometheus.MustNewConstMetric(successDesc, prometheus.GaugeValue, 1, labelValues...)
		ch <- prometheus.MustNewConstMetric(rttDesc, prometheus.GaugeValue, m.Rtt, labelValues...)
	} else {
		ch <- prometheus.MustNewConstMetric(successDesc, prometheus.GaugeValue, 0, labelValues...)
	}
}

// Exports metric descriptions for prometheus
func (m *DnsMetricExporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- successDesc
	ch <- rttDesc
}

// Sets AS number for measurement result
func (m *DnsMetricExporter) SetAsn(asn int) {
	m.Asn = asn
}

// Gets whether an result is valid (e.g. IPv6 measurement and Probe does not support IPv6)
func (m *DnsMetricExporter) Isvalid() bool {
	return m.Asn > 0
}
