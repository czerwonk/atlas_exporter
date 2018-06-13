package dns

import (
	"strconv"

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

// DNSMetricExporter exports metrics for DNS measurement results
type DNSMetricExporter struct {
}

// Export exports a prometheus metric
func (m *DNSMetricExporter) Export(id string, res *measurement.Result, probe *probe.Probe, ch chan<- prometheus.Metric) {
	labelValues := []string{
		id,
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

// Describe exports metric descriptions for Prometheus
func (m *DNSMetricExporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- successDesc
	ch <- rttDesc
}

// IsValid returns whether an result is valid or not (e.g. IPv6 measurement and Probe does not support IPv6)
func (m *DNSMetricExporter) IsValid(res *measurement.Result, probe *probe.Probe) bool {
	return probe.ASNForIPVersion(res.Af()) > 0
}
