package main

import (
	"log"

	"errors"
	"fmt"
	"time"

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
	if r.Type() == "ping" {
		out <- ping.FromResult(r)
		return
	}

	if r.Type() == "traceroute" {
		out <- traceroute.FromResult(r)
		return
	}

	log.Printf("Type %s is not yet supported\n", r.Type())
	out <- nil
}
