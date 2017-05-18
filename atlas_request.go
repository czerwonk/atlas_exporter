package main

import (
	"log"

	"errors"
	"fmt"
	"time"

	"github.com/DNS-OARC/ripeatlas"
	"github.com/DNS-OARC/ripeatlas/measurement"
	"github.com/czerwonk/atlas_exporter/ping"
	"github.com/czerwonk/atlas_exporter/probe"
	"github.com/czerwonk/atlas_exporter/traceroute"
)

func getMeasurement(id string) ([]Metric, error) {
	a := ripeatlas.Atlaser(ripeatlas.NewHttp())
	c, err := a.MeasurementLatest(ripeatlas.Params{"pk": id})

	if err != nil {
		return nil, err
	}

	res := make([]Metric, 0)
	ch := make(chan Metric)

	count := 0
	for r := range c {
		if r.ParseError != nil {
			return nil, err
		}

		go convertToMetric(r, ch)
		count++
	}

	for i := 0; i < count; i++ {
		select {
		case m := <-ch:
			if m != nil {
				res = append(res, m)
			}
		case <-time.After(60 * time.Second):
			return nil, errors.New(fmt.Sprintln("Timeout exceeded!"))
		}
	}

	return res, nil
}

func convertToMetric(r *measurement.Result, out chan Metric) {
	var m Metric

	if r.Type() == "ping" {
		m = ping.FromResult(r)
	}

	if r.Type() == "traceroute" {
		m = traceroute.FromResult(r)
	}

	if m != nil {
		setAsnForMetric(r, m)
	} else {
		log.Printf("Type %s is not yet supported\n", r.Type())
	}

	out <- m
}
func setAsnForMetric(r *measurement.Result, m Metric) {
	p, err := probe.Get(r.PrbId())

	if err != nil {
		log.Printf("Could not get information for probe %d: %v\n", r.PrbId(), err)
		return
	}

	if r.Af() == 4 {
		m.SetAsn(p.Asn4)
	} else {
		m.SetAsn(p.Asn6)
	}
}
