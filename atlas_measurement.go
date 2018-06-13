package main

import (
	"context"
	"fmt"
	"sync"

	"github.com/czerwonk/atlas_exporter/dns"
	"github.com/czerwonk/atlas_exporter/http"
	"github.com/czerwonk/atlas_exporter/ntp"
	"github.com/czerwonk/atlas_exporter/ping"
	"github.com/czerwonk/atlas_exporter/probe"
	"github.com/czerwonk/atlas_exporter/sslcert"
	"github.com/czerwonk/atlas_exporter/traceroute"

	"github.com/DNS-OARC/ripeatlas"
	"github.com/DNS-OARC/ripeatlas/measurement"
	"github.com/prometheus/common/log"
)

var atlasser ripeatlas.Atlaser

func init() {
	atlasser = ripeatlas.Atlaser(ripeatlas.NewHttp())
}

type atlasMeasurement struct {
	id       string
	results  []*measurement.Result
	exporter MetricExporter
	probes   map[int]*probe.Probe
}

func getMeasurements(ctx context.Context, ids []string) ([]*atlasMeasurement, error) {
	ch := make(chan *atlasMeasurement)

	wg := sync.WaitGroup{}
	wg.Add(len(ids))

	go func() {
		wg.Wait()
		close(ch)
	}()

	for _, id := range ids {
		go getMeasurementForID(ctx, id, ch, &wg)
	}

	res := []*atlasMeasurement{}
	for {
		select {
		case m, more := <-ch:
			if !more {
				return res, nil
			}

			res = append(res, m)
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
}

func getMeasurementForID(ctx context.Context, id string, ch chan<- *atlasMeasurement, wg *sync.WaitGroup) {
	defer wg.Done()

	resultCh, err := atlasser.MeasurementLatest(ripeatlas.Params{"pk": id})
	if err != nil {
		log.Errorf("could not retrieve measurement results for %s: %v", id, err)
		return
	}

	res := []*measurement.Result{}
	for m := range resultCh {
		if m.ParseError != nil {
			log.Errorf("failed parsing measurement result for %s: %v", id, m.ParseError)
			return
		}

		res = append(res, m)
	}

	if len(res) == 0 {
		return
	}

	probes, err := probesForResults(res)
	if err != nil {
		log.Errorf("could not retrieve probe information for measurement %s: %v", id, err)
		return
	}

	exporter, err := exporterForType(res[0].Type())
	if err != nil {
		log.Errorf("could determine exporter for measurement %s: %v", id, err)
		return
	}

	ch <- &atlasMeasurement{
		id:       id,
		results:  res,
		probes:   probes,
		exporter: exporter,
	}
}

func probesForResults(res []*measurement.Result) (map[int]*probe.Probe, error) {
	probes := make(map[int]*probe.Probe)

	for _, m := range res {
		id := m.PrbId()
		p, found := cache.Get(id)
		if found {
			probes[id] = p
			continue
		}

		p, err := probe.Get(id)
		if err != nil {
			return nil, fmt.Errorf("could not retrieve probe information for probe %d: %v", id, err)
		}

		cache.Add(id, p)
		probes[id] = p
	}

	return probes, nil
}

func exporterForType(t string) (MetricExporter, error) {
	switch t {
	case "ping":
		return &ping.PingMetricExporter{}, nil
	case "traceroute":
		return &traceroute.TracerouteMetricExporter{}, nil
	case "ntp":
		return &ntp.NTPMetricExporter{}, nil
	case "dns":
		return &dns.DNSMetricExporter{}, nil
	case "http":
		return &http.HTTPMetricExporter{}, nil
	case "sslcert":
		return &sslcert.SslCertMetricExporter{}, nil
	}

	return nil, fmt.Errorf("type %s is not supported yet", t)
}
