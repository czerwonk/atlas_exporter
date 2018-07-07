package atlas

import (
	"context"
	"sync"

	"github.com/DNS-OARC/ripeatlas"
	"github.com/DNS-OARC/ripeatlas/measurement"
	"github.com/prometheus/common/log"
)

type requestStrategy struct {
	atlasser ripeatlas.Atlaser
	workers  uint
}

// NewRequestStrategy returns an strategy to retrieve data from Atlas API using requests
func NewRequestStrategy(workers uint) Strategy {
	return requestStrategy{
		atlasser: ripeatlas.Atlaser(ripeatlas.NewHttp()),
		workers:  workers,
	}
}

func (s requestStrategy) MeasurementResults(ctx context.Context, ids []string) ([]*AtlasMeasurement, error) {
	ch := make(chan *AtlasMeasurement)

	wg := sync.WaitGroup{}
	wg.Add(len(ids))

	go func() {
		wg.Wait()
		close(ch)
	}()

	for _, id := range ids {
		go s.getMeasurementForID(ctx, id, ch, &wg)
	}

	res := []*AtlasMeasurement{}
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

func (s *requestStrategy) getMeasurementForID(ctx context.Context, id string, ch chan<- *AtlasMeasurement, wg *sync.WaitGroup) {
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

	r, err := atlasMeasurementForResults(res, id, s.workers)
	if err != nil {
		log.Errorf("failed getting measurement result for %s: %v", id, err)
	}

	ch <- r
}
