package traceroute

import (
	"strconv"

	"github.com/DNS-OARC/ripeatlas/measurement"
	"github.com/czerwonk/atlas_exporter/exporter"
	"github.com/czerwonk/atlas_exporter/probe"
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
	labels = []string{"measurement", "probe", "dst_addr", "dst_name", "asn", "ip_version", "protocol", "country_code", "lat", "long"}

	successDesc = prometheus.NewDesc(prometheus.BuildFQName(ns, sub, "success"), "Destination was reachable", labels, nil)
	hopDesc = prometheus.NewDesc(prometheus.BuildFQName(ns, sub, "hops"), "Number of hops", labels, nil)
	rttDesc = prometheus.NewDesc(prometheus.BuildFQName(ns, sub, "rtt"), "Round trip time in ms", labels, nil)
}

type tracerouteMetricExporter struct {
	id      string
	rttHist prometheus.Histogram
}

// NewExporter creates a exporter for traceroute measurement results
func NewExporter(id string, buckets []float64) exporter.MetricExporter {
	if buckets == nil {
		buckets = []float64{10, 20, 50, 100}
	}

	hist := prometheus.NewHistogram(prometheus.HistogramOpts{
		Namespace: ns,
		Subsystem: sub,
		Name:      "rtt_hist",
		Buckets:   buckets,
		Help:      "Histogram of round trip times over all traceroute requests",
		ConstLabels: prometheus.Labels{
			"measurement": id,
		},
	})

	return &tracerouteMetricExporter{
		id:      id,
		rttHist: hist,
	}
}

// Export exports a prometheus metric
func (m *tracerouteMetricExporter) Export(res *measurement.Result, probe *probe.Probe, ch chan<- prometheus.Metric) {
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

// ExportHistograms exports aggregated metrics for the measurement
func (m *tracerouteMetricExporter) ExportHistograms(res []*measurement.Result, ch chan<- prometheus.Metric) {
	for _, r := range res {
		success, rtt := processLastHop(r)
		if success == 1 && rtt > 0 {
			m.rttHist.Observe(rtt)
		}
	}

	m.rttHist.Collect(ch)
}

func processLastHop(r *measurement.Result) (success float64, rtt float64) {
	if len(r.TracerouteResults()) == 0 {
		return success, rtt
	}

	last := r.TracerouteResults()[len(r.TracerouteResults())-1]
	for _, rep := range last.Replies() {
		if rep.From() == r.DstAddr() {
			success = 1
		}

		repl := last.Replies()
		if len(repl) > 0 {
			rtt = repl[len(repl)-1].Rtt()
		}
	}

	return success, rtt
}

// Describe exports metric descriptions for Prometheus
func (m *tracerouteMetricExporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- successDesc
	ch <- hopDesc
	ch <- rttDesc
	ch <- m.rttHist.Desc()
}

// IsValid returns whether an result is valid or not (e.g. IPv6 measurement and Probe does not support IPv6)
func (m *tracerouteMetricExporter) IsValid(res *measurement.Result, probe *probe.Probe) bool {
	return probe.ASNForIPVersion(res.Af()) > 0 && len(res.TracerouteResults()) > 1
}
