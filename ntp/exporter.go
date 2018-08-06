package ntp

import (
	"strconv"

	"github.com/DNS-OARC/ripeatlas/measurement"
	"github.com/czerwonk/atlas_exporter/probe"
	"github.com/prometheus/client_golang/prometheus"
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
	labels = []string{"measurement", "probe", "dst_addr", "dst_name", "asn", "ip_version", "country_code", "lat", "long"}

	pollDesc = prometheus.NewDesc(prometheus.BuildFQName(ns, sub, "poll"), "Poll", labels, nil)
	precisionDesc = prometheus.NewDesc(prometheus.BuildFQName(ns, sub, "precision"), "Precision", labels, nil)
	roolDelayDesc = prometheus.NewDesc(prometheus.BuildFQName(ns, sub, "root_delay"), "Root delay", labels, nil)
	rootDispersionDesc = prometheus.NewDesc(prometheus.BuildFQName(ns, sub, "root_dispersion"), "Root dispersion", labels, nil)
	ntpVersionDesc = prometheus.NewDesc(prometheus.BuildFQName(ns, sub, "ntp_version"), "NTP Version", labels, nil)
}

type ntpExporter struct {
	id string
}

// Export exports a prometheus metric
func (m *ntpExporter) Export(res *measurement.Result, probe *probe.Probe, ch chan<- prometheus.Metric) {
	labelValues := []string{
		m.id,
		strconv.Itoa(probe.ID),
		res.DstAddr(),
		res.DstName(),
		strconv.Itoa(probe.ASNForIPVersion(res.Af())),
		strconv.Itoa(res.Af()),
		probe.CountryCode,
		probe.Latitude(),
		probe.Longitude(),
	}

	ch <- prometheus.MustNewConstMetric(pollDesc, prometheus.GaugeValue, res.Poll(), labelValues...)
	ch <- prometheus.MustNewConstMetric(precisionDesc, prometheus.GaugeValue, res.Precision(), labelValues...)
	ch <- prometheus.MustNewConstMetric(roolDelayDesc, prometheus.GaugeValue, res.RootDelay(), labelValues...)
	ch <- prometheus.MustNewConstMetric(rootDispersionDesc, prometheus.GaugeValue, res.RootDispersion(), labelValues...)
	ch <- prometheus.MustNewConstMetric(ntpVersionDesc, prometheus.GaugeValue, float64(res.Version()), labelValues...)
}

// Describe exports metric descriptions for Prometheus
func (m *ntpExporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- pollDesc
	ch <- precisionDesc
	ch <- roolDelayDesc
	ch <- rootDispersionDesc
	ch <- ntpVersionDesc
}
