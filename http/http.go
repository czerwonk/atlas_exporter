package http

import (
	"strconv"

	"github.com/DNS-OARC/ripeatlas/measurement"
	"github.com/czerwonk/atlas_exporter/exporter"
	"github.com/czerwonk/atlas_exporter/probe"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
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
	labels = []string{"measurement", "probe", "dst_addr", "asn", "ip_version", "uri", "method", "country_code", "lat", "long"}

	successDesc = prometheus.NewDesc(prometheus.BuildFQName(ns, sub, "success"), "Destination was reachable", labels, nil)
	resultDesc = prometheus.NewDesc(prometheus.BuildFQName(ns, sub, "result"), "Code returned from http server", labels, nil)
	httpVerDesc = prometheus.NewDesc(prometheus.BuildFQName(ns, sub, "version"), "HTTP version used for the request", labels, nil)
	bodySizeDesc = prometheus.NewDesc(prometheus.BuildFQName(ns, sub, "body_size"), "Body size in bytes", labels, nil)
	headerSizeDesc = prometheus.NewDesc(prometheus.BuildFQName(ns, sub, "header_size"), "Header size in bytes", labels, nil)
	rttDesc = prometheus.NewDesc(prometheus.BuildFQName(ns, sub, "rtt"), "Round trip time in ms", labels, nil)
	dnsErrDesc = prometheus.NewDesc(prometheus.BuildFQName(ns, sub, "dns_error"), "A DNS error occured (0 if not)", labels, nil)
}

type httpMetricExporter struct {
	id      string
	rttHist prometheus.Histogram
}

// NewExporter creates a exporter for HTTP measurement results
func NewExporter(id string, buckets []float64) exporter.MetricExporter {
	if buckets == nil {
		buckets = []float64{100, 200, 500, 1000}
	}

	hist := prometheus.NewHistogram(prometheus.HistogramOpts{
		Namespace: ns,
		Subsystem: sub,
		Name:      "rtt_hist",
		Buckets:   buckets,
		Help:      "Histogram of round trip times over all HTTP requests",
		ConstLabels: prometheus.Labels{
			"measurement": id,
		},
	})

	return &httpMetricExporter{
		id:      id,
		rttHist: hist,
	}
}

// Export exports metrics for Prometheus
func (m *httpMetricExporter) Export(res *measurement.Result, probe *probe.Probe, ch chan<- prometheus.Metric) {
	for _, h := range res.HttpResults() {
		labelValues := []string{
			m.id,
			strconv.Itoa(probe.ID),
			h.DstAddr(),
			strconv.Itoa(probe.ASNForIPVersion(h.Af())),
			strconv.Itoa(h.Af()),
			res.Uri(),
			h.Method(),
			probe.CountryCode,
			probe.Latitude(),
			probe.Longitude(),
		}

		dnsError := 0
		if len(h.Dnserr()) > 0 {
			dnsError = 1
		}

		httpVer, err := strconv.ParseFloat(h.Ver(), 64)
		if err != nil {
			log.Errorf("error parsing http version %s: %v", h.Ver(), err)
		}

		ch <- prometheus.MustNewConstMetric(resultDesc, prometheus.GaugeValue, float64(h.Res()), labelValues...)
		ch <- prometheus.MustNewConstMetric(httpVerDesc, prometheus.GaugeValue, httpVer, labelValues...)
		ch <- prometheus.MustNewConstMetric(bodySizeDesc, prometheus.GaugeValue, float64(h.Bsize()), labelValues...)
		ch <- prometheus.MustNewConstMetric(headerSizeDesc, prometheus.GaugeValue, float64(h.Hsize()), labelValues...)
		ch <- prometheus.MustNewConstMetric(dnsErrDesc, prometheus.GaugeValue, float64(dnsError), labelValues...)

		if h.Rt() > 0 {
			ch <- prometheus.MustNewConstMetric(successDesc, prometheus.GaugeValue, 1, labelValues...)
			ch <- prometheus.MustNewConstMetric(rttDesc, prometheus.GaugeValue, h.Rt(), labelValues...)
		} else {
			ch <- prometheus.MustNewConstMetric(successDesc, prometheus.GaugeValue, 0, labelValues...)
		}
	}
}

// ExportHistograms exports aggregated metrics for the measurement
func (m *httpMetricExporter) ExportHistograms(res []*measurement.Result, ch chan<- prometheus.Metric) {
	for _, r := range res {
		for _, p := range r.HttpResults() {
			if p.Rt() > 0 {
				m.rttHist.Observe(p.Rt())
			}
		}
	}

	m.rttHist.Collect(ch)
}

// Describe exports metric descriptions for Prometheus
func (m *httpMetricExporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- successDesc
	ch <- resultDesc
	ch <- httpVerDesc
	ch <- bodySizeDesc
	ch <- headerSizeDesc
	ch <- rttDesc
	ch <- dnsErrDesc
	ch <- m.rttHist.Desc()
}

// IsValid returns whether an result is valid or not (e.g. IPv6 measurement and Probe does not support IPv6)
func (m *httpMetricExporter) IsValid(res *measurement.Result, probe *probe.Probe) bool {
	return probe.ASNForIPVersion(res.Af()) > 0
}
