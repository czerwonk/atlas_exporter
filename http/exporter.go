package http

import (
	"strconv"

	"github.com/DNS-OARC/ripeatlas/measurement"
	"github.com/czerwonk/atlas_exporter/probe"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
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
	dnsErrDesc = prometheus.NewDesc(prometheus.BuildFQName(ns, sub, "dns_error"), "A DNS error occurred (0 if not)", labels, nil)
}

type httpExporter struct {
	id string
}

// Export exports metrics for Prometheus
func (m *httpExporter) Export(res *measurement.Result, probe *probe.Probe, ch chan<- prometheus.Metric) {
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

// Describe exports metric descriptions for Prometheus
func (m *httpExporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- successDesc
	ch <- resultDesc
	ch <- httpVerDesc
	ch <- bodySizeDesc
	ch <- headerSizeDesc
	ch <- rttDesc
	ch <- dnsErrDesc
}
