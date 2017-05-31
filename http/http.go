package http

import (
	"strconv"

	"github.com/DNS-OARC/ripeatlas/measurement"
	"github.com/DNS-OARC/ripeatlas/measurement/http"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	ns  = "atlas"
	sub = "http"
)

var (
	labels         []string
	resultDesc     *prometheus.Desc
	httpVerDesc    *prometheus.Desc
	bodySizeDesc   *prometheus.Desc
	headerSizeDesc *prometheus.Desc
	rttDesc        *prometheus.Desc
	dnsErrDesc     *prometheus.Desc
	successDesc    *prometheus.Desc
)

func init() {
	labels = []string{"measurement", "probe", "dst_addr", "asn", "ip_version", "uri", "method"}

	successDesc = prometheus.NewDesc(prometheus.BuildFQName(ns, sub, "success"), "Destination was reachable", labels, nil)
	resultDesc = prometheus.NewDesc(prometheus.BuildFQName(ns, sub, "result"), "Code returned from http server", labels, nil)
	httpVerDesc = prometheus.NewDesc(prometheus.BuildFQName(ns, sub, "version"), "HTTP version used for the request", labels, nil)
	bodySizeDesc = prometheus.NewDesc(prometheus.BuildFQName(ns, sub, "body_size"), "Body size in bytes", labels, nil)
	headerSizeDesc = prometheus.NewDesc(prometheus.BuildFQName(ns, sub, "header_size"), "Header size in bytes", labels, nil)
	rttDesc = prometheus.NewDesc(prometheus.BuildFQName(ns, sub, "rtt"), "Round trip time in ms", labels, nil)
	dnsErrDesc = prometheus.NewDesc(prometheus.BuildFQName(ns, sub, "dns_error"), "Round trip time in ms", labels, nil)
}

// HttpMetricExporter exports metrics for HTTP measurement results
type HttpMetricExporter struct {
	ProbeId     int
	DstAddr     string
	Uri         string
	ReturnCode  int
	HttpVersion float64
	BodySize    int
	HeaderSize  int
	Method      string
	Rtt         float64
	DnsError    int
	Asn         int
	IpVersion   int
}

// FromResult creates  metric exporter for HTTP measurement result
func FromResult(r *measurement.Result) *HttpMetricExporter {
	m := &HttpMetricExporter{ProbeId: r.PrbId(), Uri: r.Uri()}

	if len(r.HttpResults()) > 0 {
		h := r.HttpResults()[0]
		m.fillFromHttpResult(h)
	}

	return m
}

func (m *HttpMetricExporter) fillFromHttpResult(h *http.Result) {
	m.IpVersion = h.Af()
	m.DstAddr = h.DstAddr()
	m.ReturnCode = h.Res()
	m.BodySize = h.Bsize()
	m.HeaderSize = h.Hsize()
	m.Method = h.Method()
	m.HttpVersion, _ = strconv.ParseFloat(h.Ver(), 64)
	m.Rtt = h.Rt()

	if len(h.Dnserr()) > 0 {
		m.DnsError = 1
	}
}

// Export exports metrics for Prometheus
func (m *HttpMetricExporter) Export(ch chan<- prometheus.Metric, pk string) {
	labelValues := []string{pk, strconv.Itoa(m.ProbeId), m.DstAddr, strconv.Itoa(m.Asn), strconv.Itoa(m.IpVersion), m.Uri, m.Method}

	ch <- prometheus.MustNewConstMetric(resultDesc, prometheus.GaugeValue, float64(m.ReturnCode), labelValues...)
	ch <- prometheus.MustNewConstMetric(httpVerDesc, prometheus.GaugeValue, m.HttpVersion, labelValues...)
	ch <- prometheus.MustNewConstMetric(bodySizeDesc, prometheus.GaugeValue, float64(m.BodySize), labelValues...)
	ch <- prometheus.MustNewConstMetric(headerSizeDesc, prometheus.GaugeValue, float64(m.HeaderSize), labelValues...)
	ch <- prometheus.MustNewConstMetric(dnsErrDesc, prometheus.GaugeValue, float64(m.DnsError), labelValues...)

	if m.Rtt > 0 {
		ch <- prometheus.MustNewConstMetric(successDesc, prometheus.GaugeValue, 1, labelValues...)
		ch <- prometheus.MustNewConstMetric(rttDesc, prometheus.GaugeValue, m.Rtt, labelValues...)
	} else {
		ch <- prometheus.MustNewConstMetric(successDesc, prometheus.GaugeValue, 0, labelValues...)
	}
}

// Describe exports metric descriptions for Prometheus
func (m *HttpMetricExporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- successDesc
	ch <- resultDesc
	ch <- httpVerDesc
	ch <- bodySizeDesc
	ch <- headerSizeDesc
	ch <- rttDesc
	ch <- dnsErrDesc
}

// SetAsn sets AS number for measurement result
func (m *HttpMetricExporter) SetAsn(asn int) {
	m.Asn = asn
}

// IsValid returns whether an result is valid or not (e.g. IPv6 measurement and Probe does not support IPv6)
func (m *HttpMetricExporter) IsValid() bool {
	return m.Asn > 0
}
