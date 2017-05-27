package ntp

import (
	"strconv"

	"github.com/DNS-OARC/ripeatlas/measurement"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	ns  = "atlas"
	sub = "ntp"
)

var (
	labels             []string
	pollDesc           *prometheus.Desc
	precisionDesc      *prometheus.Desc
	roolDelayDesc      *prometheus.Desc
	rootDispersionDesc *prometheus.Desc
	ntpVersionDesc     *prometheus.Desc
)

func init() {
	labels = make([]string, 0)
	labels = append(labels, "measurement", "probe", "dst_addr", "dst_name", "asn", "ip_version")

	pollDesc = prometheus.NewDesc(prometheus.BuildFQName(ns, sub, "poll"), "Poll", labels, nil)
	precisionDesc = prometheus.NewDesc(prometheus.BuildFQName(ns, sub, "precision"), "Precision", labels, nil)
	roolDelayDesc = prometheus.NewDesc(prometheus.BuildFQName(ns, sub, "root_delay"), "Root delay", labels, nil)
	rootDispersionDesc = prometheus.NewDesc(prometheus.BuildFQName(ns, sub, "root_dispersion"), "Root dispersion", labels, nil)
	ntpVersionDesc = prometheus.NewDesc(prometheus.BuildFQName(ns, sub, "ntp_version"), "NTP Version", labels, nil)
}

// NtpMetricExporter exports metrics for NTP measurement results
type NtpMetricExporter struct {
	ProbeId        int
	DstAddr        string
	DstName        string
	Poll           float64
	Precision      float64
	RootDelay      float64
	RootDispersion float64
	Version        int
	Asn            int
	IpVersion      int
}

// FromResult creates  metric exporter for NTP measurement result
func FromResult(r *measurement.Result) *NtpMetricExporter {
	return &NtpMetricExporter{ProbeId: r.PrbId(), DstAddr: r.DstAddr(), DstName: r.DstName(), Poll: r.Poll(), Precision: r.Precision(), RootDelay: r.RootDelay(), RootDispersion: r.RootDispersion(), Version: r.Version(), IpVersion: r.Af()}
}

// Export exports metrics for prometheus
func (m *NtpMetricExporter) Export(ch chan<- prometheus.Metric, pk string) {
	labelValues := make([]string, 0)
	labelValues = append(labelValues, pk, strconv.Itoa(m.ProbeId), m.DstAddr, m.DstName, strconv.Itoa(m.Asn), strconv.Itoa(m.IpVersion))

	ch <- prometheus.MustNewConstMetric(pollDesc, prometheus.GaugeValue, m.Poll, labelValues...)
	ch <- prometheus.MustNewConstMetric(precisionDesc, prometheus.GaugeValue, m.Precision, labelValues...)
	ch <- prometheus.MustNewConstMetric(roolDelayDesc, prometheus.GaugeValue, m.RootDelay, labelValues...)
	ch <- prometheus.MustNewConstMetric(rootDispersionDesc, prometheus.GaugeValue, m.RootDispersion, labelValues...)
	ch <- prometheus.MustNewConstMetric(ntpVersionDesc, prometheus.GaugeValue, float64(m.Version), labelValues...)
}

// Describe exports metric descriptions for prometheus
func (m *NtpMetricExporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- pollDesc
	ch <- precisionDesc
	ch <- roolDelayDesc
	ch <- rootDispersionDesc
	ch <- ntpVersionDesc
}

// SetAsn sets AS number for measurement result
func (m *NtpMetricExporter) SetAsn(asn int) {
	m.Asn = asn
}

// IsValid returns whether an result is valid or not (e.g. IPv6 measurement and Probe does not support IPv6)
func (m *NtpMetricExporter) Isvalid() bool {
	return m.Asn > 0
}
