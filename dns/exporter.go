package dns

import (
	"strconv"

	"github.com/DNS-OARC/ripeatlas/measurement"
	"github.com/DNS-OARC/ripeatlas/measurement/dns"
	"github.com/czerwonk/atlas_exporter/probe"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	labels      []string
	successDesc *prometheus.Desc
	rttDesc     *prometheus.Desc
)

func init() {
	labels = []string{"measurement", "probe", "dst_addr", "asn", "ip_version", "country_code", "lat", "long", "rdata1"}

	successDesc = prometheus.NewDesc(prometheus.BuildFQName(ns, sub, "success"), "Destination was reachable", labels, nil)
	rttDesc = prometheus.NewDesc(prometheus.BuildFQName(ns, sub, "rtt"), "Roundtrip time in ms", labels, nil)
}

type dnsExporter struct {
	id string
}

// Export exports a prometheus metric
func (m *dnsExporter) Export(res *measurement.Result, probe *probe.Probe, ch chan<- prometheus.Metric) {
	var rtt float64
	var answers []*dns.Answer
	if res.DnsResult() != nil {
		rtt = res.DnsResult().Rt()
		answers = res.DnsResult().Answers()
	}

	var rdata1 string
	if len(answers) > 0 {
		if len(answers[0].Rdata()) > 0 {
			rdata1 = answers[0].Rdata()[0]
		}
	}

	labelValues := []string{
		m.id,
		strconv.Itoa(probe.ID),
		res.DstAddr(),
		strconv.Itoa(probe.ASNForIPVersion(res.Af())),
		strconv.Itoa(res.Af()),
		probe.CountryCode,
		probe.Latitude(),
		probe.Longitude(),
		rdata1,
	}

	if rtt > 0 {
		ch <- prometheus.MustNewConstMetric(successDesc, prometheus.GaugeValue, 1, labelValues...)
		ch <- prometheus.MustNewConstMetric(rttDesc, prometheus.GaugeValue, rtt, labelValues...)
	} else {
		ch <- prometheus.MustNewConstMetric(successDesc, prometheus.GaugeValue, 0, labelValues...)
	}
}

// Describe exports metric descriptions for Prometheus
func (m *dnsExporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- successDesc
	ch <- rttDesc
}
