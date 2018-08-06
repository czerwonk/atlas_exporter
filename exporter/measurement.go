package exporter

import (
	"github.com/DNS-OARC/ripeatlas/measurement"
	"github.com/czerwonk/atlas_exporter/probe"
	"github.com/prometheus/client_golang/prometheus"
)

// MeasurementOpt are options to apply to the `Measurement`
type MeasurementOpt func(r *Measurement)

// WithHistograms adds histograms to the measurement
func WithHistograms(h ...Histogram) MeasurementOpt {
	return func(r *Measurement) {
		r.histograms = append(r.histograms, h...)
	}
}

// WithValidator sets an validator to validate results for a measurement
func WithValidator(v ResultValidator) MeasurementOpt {
	return func(r *Measurement) {
		r.validator = v
	}
}

// Measurement handles measurement results and converts to metrics
type Measurement struct {
	latest          map[int]*measurement.Result
	sinceLastscrape []*measurement.Result
	probes          map[int]*probe.Probe
	histograms      []Histogram
	exporter        Exporter
	validator       ResultValidator
}

// NewMeasurement returns a new instance of `Measurement`
func NewMeasurement(exporter Exporter, opts ...MeasurementOpt) *Measurement {
	r := &Measurement{
		latest:          make(map[int]*measurement.Result),
		sinceLastscrape: make([]*measurement.Result, 0),
		probes:          make(map[int]*probe.Probe),
		histograms:      make([]Histogram, 0),
		exporter:        exporter,
	}

	for _, opt := range opts {
		opt(r)
	}

	return r
}

// Add adds an result to a measurement
func (r *Measurement) Add(m *measurement.Result, probe *probe.Probe) {
	if r.validator != nil && !r.validator.IsValid(m, probe) {
		return
	}

	r.sinceLastscrape = append(r.sinceLastscrape, m)
	r.latest[m.PrbId()] = m
	r.probes[m.PrbId()] = probe

	for _, h := range r.histograms {
		h.ProcessResult(m)
	}
}

// Scraped is called when the scrape happened to clear some internal states
func (r *Measurement) Scraped() {
	r.sinceLastscrape = make([]*measurement.Result, 0)
}

// Describe describes all metrics for the `Measurement`
func (r *Measurement) Describe(ch chan<- *prometheus.Desc) {
	r.exporter.Describe(ch)

	for _, h := range r.histograms {
		h.Hist().Describe(ch)
	}
}

// Collect collects metrics for the `Measurement`
func (r *Measurement) Collect(ch chan<- prometheus.Metric) {
	for _, v := range r.latest {
		r.exporter.Export(v, r.probes[v.PrbId()], ch)
	}

	for _, h := range r.histograms {
		h.Hist().Collect(ch)
	}
}
