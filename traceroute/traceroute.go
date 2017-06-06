package traceroute

import (
	"strconv"

	"github.com/DNS-OARC/ripeatlas/measurement"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	ns  = "atlas"
	sub = "traceroute"
)

var (
	labels      []string
	successDesc *prometheus.Desc
	hopDesc     *prometheus.Desc
	rttDesc     *prometheus.Desc
)

func init() {
	labels = []string{"measurement", "probe", "dst_addr", "dst_name", "asn", "ip_version", "protocol"}

	successDesc = prometheus.NewDesc(prometheus.BuildFQName(ns, sub, "success"), "Destination was reachable", labels, nil)
	hopDesc = prometheus.NewDesc(prometheus.BuildFQName(ns, sub, "hops"), "Number of hops", labels, nil)
	rttDesc = prometheus.NewDesc(prometheus.BuildFQName(ns, sub, "rtt"), "Round trip time in ms", labels, nil)
}

// TracerouteMetricExporter exports metrics for traceroute measurement results
type TracerouteMetricExporter struct {
	ProbeId   int
	DstAddr   string
	DstName   string
	HopCount  int
	Success   int
	Rtt       float64
	Asn       int
	IpVersion int
	Protocol  string
}

// FromResult creates metric exporter for traceroute measurement result
func FromResult(r *measurement.Result) *TracerouteMetricExporter {
	m := &TracerouteMetricExporter{ProbeId: r.PrbId(), DstAddr: r.DstAddr(), DstName: r.DstName(), HopCount: len(r.TracerouteResults()), IpVersion: r.Af(), Protocol: r.Proto()}
	processLastHop(r, m)

	return m
}

func processLastHop(r *measurement.Result, m *TracerouteMetricExporter) {
	if len(r.TracerouteResults()) == 0 {
		return
	}

	last := r.TracerouteResults()[len(r.TracerouteResults())-1]
	for _, rep := range last.Replies() {
		if rep.From() == r.DstAddr() {
			m.Success = 1
			m.Rtt = rep.Rtt()
		}
	}
}

// Export exports metrics for Prometheus
func (m *TracerouteMetricExporter) Export(ch chan<- prometheus.Metric, pk string) {
	labelValues := []string{pk, strconv.Itoa(m.ProbeId), m.DstAddr, m.DstName, strconv.Itoa(m.Asn), strconv.Itoa(m.IpVersion), m.Protocol}

	ch <- prometheus.MustNewConstMetric(successDesc, prometheus.GaugeValue, float64(m.Success), labelValues...)
	ch <- prometheus.MustNewConstMetric(hopDesc, prometheus.GaugeValue, float64(m.HopCount), labelValues...)

	if m.Rtt > 0 {
		ch <- prometheus.MustNewConstMetric(rttDesc, prometheus.GaugeValue, m.Rtt, labelValues...)
	}
}

// Describe exports metric descriptions for Prometheus
func (m *TracerouteMetricExporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- successDesc
	ch <- hopDesc
	ch <- rttDesc
}

// SetAsn sets AS number for measurement result
func (m *TracerouteMetricExporter) SetAsn(asn int) {
	m.Asn = asn
}

// IsValid returns whether an result is valid or not (e.g. IPv6 measurement and Probe does not support IPv6)
func (m *TracerouteMetricExporter) IsValid() bool {
	return (m.Success == 1 || m.HopCount > 1) && m.Asn > 0
}
