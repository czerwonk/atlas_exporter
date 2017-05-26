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
	labels = make([]string, 0)
	labels = append(labels, "measurement", "probe", "dst_addr", "dst_name", "asn", "ip_version")

	successDesc = prometheus.NewDesc(prometheus.BuildFQName(ns, sub, "success"), "Destination was reachable", labels, nil)
	hopDesc = prometheus.NewDesc(prometheus.BuildFQName(ns, sub, "hops"), "Number of hops", labels, nil)
	rttDesc = prometheus.NewDesc(prometheus.BuildFQName(ns, sub, "rtt"), "Round trip time in ms", labels, nil)
}

type TracerouteMetric struct {
	ProbeId   int
	DstAddr   string
	DstName   string
	HopCount  int
	Success   int
	Rtt       float64
	Asn       int
	IpVersion int
}

func FromResult(r *measurement.Result) *TracerouteMetric {
	m := &TracerouteMetric{ProbeId: r.PrbId(), DstAddr: r.DstAddr(), DstName: r.DstName(), HopCount: len(r.TracerouteResults()), IpVersion: r.Af()}
	processLastHop(r, m)

	return m
}

func processLastHop(r *measurement.Result, m *TracerouteMetric) {
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

func (m *TracerouteMetric) GetMetrics(ch chan<- prometheus.Metric, pk string) {
	labelValues := make([]string, 0)
	labelValues = append(labelValues, pk, strconv.Itoa(m.ProbeId), m.DstAddr, m.DstName, strconv.Itoa(m.Asn), strconv.Itoa(m.IpVersion))

	ch <- prometheus.MustNewConstMetric(successDesc, prometheus.GaugeValue, float64(m.Success), labelValues...)
	ch <- prometheus.MustNewConstMetric(hopDesc, prometheus.GaugeValue, float64(m.HopCount), labelValues...)

	if m.Rtt > 0 {
		ch <- prometheus.MustNewConstMetric(rttDesc, prometheus.GaugeValue, m.Rtt, labelValues...)
	}
}

func (m *TracerouteMetric) Describe(ch chan<- *prometheus.Desc) {
	ch <- successDesc
	ch <- hopDesc
	ch <- rttDesc
}

func (m *TracerouteMetric) SetAsn(asn int) {
	m.Asn = asn
}

func (m *TracerouteMetric) Isvalid() bool {
	return (m.Success == 1 || m.HopCount > 1) && m.Asn > 0
}
