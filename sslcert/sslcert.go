package sslcert

import (
	"strconv"

	"github.com/DNS-OARC/ripeatlas/measurement"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	ns  = "atlas"
	sub = "sslcert"
)

var (
	labels               []string
	rttDesc              *prometheus.Desc
	sslVerDesc           *prometheus.Desc
	successDesc          *prometheus.Desc
	alertLevelDesc       *prometheus.Desc
	alertDescriptionDesc *prometheus.Desc
)

func init() {
	labels = []string{"measurement", "probe", "dst_addr", "asn", "ip_version"}

	successDesc = prometheus.NewDesc(prometheus.BuildFQName(ns, sub, "success"), "Destination was reachable", labels, nil)
	sslVerDesc = prometheus.NewDesc(prometheus.BuildFQName(ns, sub, "version"), "SSL/TLS version used for the request", labels, nil)
	rttDesc = prometheus.NewDesc(prometheus.BuildFQName(ns, sub, "rtt"), "Round trip time in ms", labels, nil)
	alertLevelDesc = prometheus.NewDesc(prometheus.BuildFQName(ns, sub, "alert_level"), "Status of the SSL/TLS certificate (0 = valid)", labels, nil)
	alertDescriptionDesc = prometheus.NewDesc(prometheus.BuildFQName(ns, sub, "alert_description"), "Description for the alert level (see RIPIE Atlas documentation)", labels, nil)
}

// SslCertMetricExporter exports metrics for SSL certificate measurement results
type SslCertMetricExporter struct {
	ProbeId          int
	DstAddr          string
	SslVersion       float64
	Rtt              float64
	AlertLevel       int
	AlertDescription int
	Asn              int
	IpVersion        int
}

// FromResult creates  metric exporter for SSL certificate measurement result
func FromResult(r *measurement.Result) *SslCertMetricExporter {
	m := &SslCertMetricExporter{ProbeId: r.PrbId(), DstAddr: r.DstAddr(), IpVersion: r.Af(), Rtt: r.Rt()}
	m.SslVersion, _ = strconv.ParseFloat(r.Ver(), 64)

	if r.SslcertAlert() != nil {
		m.AlertLevel = r.SslcertAlert().Level()
		m.AlertDescription = r.SslcertAlert().Description()
	}

	return m
}

// Export exports metrics for Prometheus
func (m *SslCertMetricExporter) Export(ch chan<- prometheus.Metric, pk string) {
	labelValues := []string{pk, strconv.Itoa(m.ProbeId), m.DstAddr, strconv.Itoa(m.Asn), strconv.Itoa(m.IpVersion)}

	ch <- prometheus.MustNewConstMetric(sslVerDesc, prometheus.GaugeValue, m.SslVersion, labelValues...)
	ch <- prometheus.MustNewConstMetric(alertLevelDesc, prometheus.GaugeValue, float64(m.AlertLevel), labelValues...)
	ch <- prometheus.MustNewConstMetric(alertDescriptionDesc, prometheus.GaugeValue, float64(m.AlertDescription), labelValues...)

	if m.Rtt > 0 {
		ch <- prometheus.MustNewConstMetric(successDesc, prometheus.GaugeValue, 1, labelValues...)
		ch <- prometheus.MustNewConstMetric(rttDesc, prometheus.GaugeValue, m.Rtt, labelValues...)
	} else {
		ch <- prometheus.MustNewConstMetric(successDesc, prometheus.GaugeValue, 0, labelValues...)
	}
}

// Describe exports metric descriptions for Prometheus
func (m *SslCertMetricExporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- successDesc
	ch <- rttDesc
	ch <- sslVerDesc
	ch <- alertLevelDesc
	ch <- alertDescriptionDesc
}

// SetAsn sets AS number for measurement result
func (m *SslCertMetricExporter) SetAsn(asn int) {
	m.Asn = asn
}

// IsValid returns whether an result is valid or not (e.g. IPv6 measurement and Probe does not support IPv6)
func (m *SslCertMetricExporter) IsValid() bool {
	return m.Asn > 0
}
