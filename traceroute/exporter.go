package traceroute

import (
	"strconv"

	"github.com/DNS-OARC/ripeatlas/measurement"
	"github.com/czerwonk/atlas_exporter/probe"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	labels      []string
	successDesc *prometheus.Desc
	hopDesc     *prometheus.Desc
	rttDesc     *prometheus.Desc
)

func init() {
	labels = []string{"measurement", "probe", "dst_addr", "dst_name", "asn", "ip_version", "protocol", "country_code", "lat", "long"}

	successDesc = prometheus.NewDesc(prometheus.BuildFQName(ns, sub, "success"), "Destination was reachable", labels, nil)
	hopDesc = prometheus.NewDesc(prometheus.BuildFQName(ns, sub, "hops"), "Number of hops", labels, nil)
	rttDesc = prometheus.NewDesc(prometheus.BuildFQName(ns, sub, "rtt"), "Round trip time in ms", labels, nil)
}

type tracerouteExporter struct {
	id string
}

// Export exports a prometheus metric
func (m *tracerouteExporter) Export(res *measurement.Result, probe *probe.Probe, ch chan<- prometheus.Metric) {
	labelValues := []string{
		m.id,
		strconv.Itoa(probe.ID),
		res.DstAddr(),
		res.DstName(),
		strconv.Itoa(probe.ASNForIPVersion(res.Af())),
		strconv.Itoa(res.Af()),
		res.Proto(),
		probe.CountryCode,
		probe.Latitude(),
		probe.Longitude(),
	}

	success, rtt := processLastHop(res)
	hops := float64(len(res.TracerouteResults()))
	ch <- prometheus.MustNewConstMetric(successDesc, prometheus.GaugeValue, success, labelValues...)
	ch <- prometheus.MustNewConstMetric(hopDesc, prometheus.GaugeValue, hops, labelValues...)

	if rtt > 0 {
		ch <- prometheus.MustNewConstMetric(rttDesc, prometheus.GaugeValue, rtt, labelValues...)
	}
}

// Describe exports metric descriptions for Prometheus
func (m *tracerouteExporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- successDesc
	ch <- hopDesc
	ch <- rttDesc
}
