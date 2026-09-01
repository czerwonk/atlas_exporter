// SPDX-License-Identifier: LGPL-3.0-or-later

package atlas

import (
	"context"
	"sync"

	"github.com/czerwonk/atlas_exporter/exporter"

	"github.com/DNS-OARC/ripeatlas"
	"github.com/DNS-OARC/ripeatlas/measurement"
	"github.com/czerwonk/atlas_exporter/config"
	log "github.com/sirupsen/logrus"
)

type requestStrategy struct {
	atlasser ripeatlas.Atlaser
	workers  uint
	cfg      *config.Config
}

// NewRequestStrategy returns an strategy to retrieve data from Atlas API using requests
func NewRequestStrategy(cfg *config.Config, workers uint) Strategy {
	return requestStrategy{
		atlasser: ripeatlas.Atlaser(ripeatlas.NewHttp()),
		cfg:      cfg,
		workers:  workers,
	}
}

func (s requestStrategy) MeasurementResults(ctx context.Context, ids []string) ([]*exporter.Measurement, error) {
	ch := make(chan *exporter.Measurement)

	wg := sync.WaitGroup{}
	wg.Add(len(ids))

	go func() {
		wg.Wait()
		close(ch)
	}()

	for _, id := range ids {
		go s.getMeasurementForID(ctx, id, ch, &wg)
	}

	res := make([]*exporter.Measurement, 0)
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

func (s *requestStrategy) getMeasurementForID(ctx context.Context, id string, ch chan<- *exporter.Measurement, wg *sync.WaitGroup) {
	defer wg.Done()

	resultCh, err := s.atlasser.MeasurementLatest(ripeatlas.Params{"pk": id})
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

	first := res[0]
	ipVersion := exporter.IpVersionForMeasurement(first)
	mes, err := measurementForType(first.Type(), id, ipVersion, s.cfg)
	if err != nil {
		log.Errorln(err)
		return
	}

	probes, err := probesForResults(res, s.workers)
	if err != nil {
		log.Errorln(err)
		return
	}

	for _, r := range res {
		mes.Add(r, probes[r.PrbId()])
	}

	ch <- mes
}
