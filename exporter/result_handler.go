package exporter

import (
	"github.com/DNS-OARC/ripeatlas/measurement"
	"github.com/czerwonk/atlas_exporter/probe"
	"github.com/prometheus/client_golang/prometheus"
)

type ResultHandlerOpt func(r *ResultHandler)

func WithHistograms(h ...Histogram) ResultHandlerOpt {
	return func(r *ResultHandler) {
		r.histograms = append(r.histograms, h...)
	}
}

func WithValidator(v ResultValidator) ResultHandlerOpt {
	return func(r *ResultHandler) {
		r.validator = v
	}
}

type ResultHandler struct {
	latest          map[int]*measurement.Result
	sinceLastscrape []*measurement.Result
	probes          map[int]*probe.Probe
	histograms      []Histogram
	exporter        Exporter
	validator       ResultValidator
}

func NewResultHandler(exporter Exporter, opts ...ResultHandlerOpt) *ResultHandler {
	r := &ResultHandler{
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

func (r *ResultHandler) Add(m *measurement.Result, probe *probe.Probe) {
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

func (r *ResultHandler) Scraped() {
	r.sinceLastscrape = make([]*measurement.Result, 0)
}

func (r *ResultHandler) Describe(ch chan<- *prometheus.Desc) {
	r.exporter.Describe(ch)

	for _, h := range r.histograms {
		h.Hist().Describe(ch)
	}
}

func (r *ResultHandler) Collect(ch chan<- prometheus.Metric) {
	for _, v := range r.latest {
		r.exporter.Export(v, r.probes[v.PrbId()], ch)
	}

	for _, h := range r.histograms {
		h.Hist().Collect(ch)
	}
}
