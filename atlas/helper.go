package atlas

import (
	"fmt"
	"sync"

	"github.com/DNS-OARC/ripeatlas/measurement"
	"github.com/czerwonk/atlas_exporter/config"
	"github.com/czerwonk/atlas_exporter/dns"
	"github.com/czerwonk/atlas_exporter/exporter"
	"github.com/czerwonk/atlas_exporter/http"
	"github.com/czerwonk/atlas_exporter/ntp"
	"github.com/czerwonk/atlas_exporter/ping"
	"github.com/czerwonk/atlas_exporter/probe"
	"github.com/czerwonk/atlas_exporter/sslcert"
	"github.com/czerwonk/atlas_exporter/traceroute"
)

func probesForResults(res []*measurement.Result, workers uint) (map[int]*probe.Probe, error) {
	probes := make(map[int]*probe.Probe)

	in := startProducer(res)
	out := make(chan *probe.Probe)
	errCh := make(chan error)

	go func() {
		startConsumers(in, out, errCh, int(workers))
	}()

	for {
		select {
		case err := <-errCh:
			return nil, err
		case p, more := <-out:
			if !more {
				return probes, nil
			}

			probes[p.ID] = p
		}
	}
}

func startProducer(res []*measurement.Result) chan int {
	ch := make(chan int)

	go func() {
		added := make(map[int]bool)

		for _, m := range res {
			if _, found := added[m.PrbId()]; found {
				continue
			}

			ch <- m.PrbId()
			added[m.PrbId()] = true
		}

		close(ch)
	}()

	return ch
}

func startConsumers(idChan chan int, out chan<- *probe.Probe, errCh chan<- error, workers int) {
	wg := sync.WaitGroup{}
	wg.Add(workers)

	for i := 0; i < workers; i++ {
		go func() {
			defer wg.Done()
			for id := range idChan {
				p, err := probeForID(id)
				if err != nil {
					errCh <- err
					continue
				}
				out <- p
			}
		}()
	}

	wg.Wait()
	close(out)
}

func probeForID(id int) (*probe.Probe, error) {
	p, found := cache.Get(id)
	if found {
		return p, nil
	}

	p, err := probe.Get(id)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve probe information for probe %d: %v", id, err)
	}

	cache.Add(id, p)
	return p, nil
}

func measurementForType(t, id, ipVersion string, cfg *config.Config) (*exporter.Measurement, error) {
	switch t {
	case "ping":
		return ping.NewMeasurement(id, ipVersion, cfg), nil
	case "traceroute":
		return traceroute.NewMeasurement(id, ipVersion, cfg), nil
	case "ntp":
		return ntp.NewMeasurement(id, cfg), nil
	case "dns":
		return dns.NewMeasurement(id, ipVersion, cfg), nil
	case "http":
		return http.NewMeasurement(id, ipVersion, cfg), nil
	case "sslcert":
		return sslcert.NewMeasurement(id, cfg), nil
	}

	return nil, fmt.Errorf("type %s is not supported yet", t)
}
