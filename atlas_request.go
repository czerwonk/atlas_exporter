package main

import (
	"log"

	"github.com/DNS-OARC/ripeatlas"
	"github.com/DNS-OARC/ripeatlas/measurement"
	"github.com/czerwonk/atlas_exporter/ping"
	"github.com/czerwonk/atlas_exporter/traceroute"
)

func getMeasurement(id string) ([]Metric, error) {
	a := ripeatlas.Atlaser(ripeatlas.NewHttp())
	c, err := a.MeasurementLatest(ripeatlas.Params{"pk": id})

	if err != nil {
		return nil, err
	}

	res := make([]Metric, 0)

	for r := range c {
		if r.ParseError != nil {
			return nil, err
		}

		m := getMetric(r)

		if m != nil {
			res = append(res, m)
		}
	}

	return res, nil
}

func getMetric(r *measurement.Result) Metric {
	if r.Type() == "ping" {
		return ping.FromResult(r)
	}

	if r.Type() == "traceroute" {
		return traceroute.FromResult(r)
	}

	log.Printf("Type %s is not yet supported\n", r.Type())
	return nil
}
