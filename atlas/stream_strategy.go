package atlas

import (
	"context"
	"strconv"
	"sync"

	"github.com/DNS-OARC/ripeatlas/measurement"

	"github.com/DNS-OARC/ripeatlas"
)

type streamingStrategy struct {
	stream  *ripeatlas.Stream
	results map[string]map[int]*measurement.Result
	workers uint
	mu      sync.RWMutex
}

// NewStreamingStrategy returns an strategy using the RIPE Atlas Streaming API
func NewStreamingStrategy(ctx context.Context, ids []string, workers uint) (Strategy, error) {
	s := &streamingStrategy{
		stream:  ripeatlas.NewStream(),
		workers: workers,
		results: make(map[string]map[int]*measurement.Result),
	}

	err := s.start(ctx, ids)
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (s *streamingStrategy) start(ctx context.Context, ids []string) error {
	for _, id := range ids {
		msm, err := strconv.Atoi(id)
		if err != nil {
			return err
		}

		ch, err := s.stream.MeasurementResults(ripeatlas.Params{
			"msm": msm,
		})
		if err != nil {
			return err
		}

		go s.listenForResults(ctx, ch)
	}

	return nil
}

func (s *streamingStrategy) listenForResults(ctx context.Context, ch <-chan *measurement.Result) {
	for {
		select {
		case m := <-ch:
			s.addOrReplace(m)
		case <-ctx.Done():
			return
		}
	}
}

func (s *streamingStrategy) addOrReplace(m *measurement.Result) {
	s.mu.Lock()
	defer s.mu.Unlock()

	msm := strconv.Itoa(m.MsmId())

	_, found := s.results[msm]
	if !found {
		s.results[msm] = make(map[int]*measurement.Result)
	}

	s.results[msm][m.PrbId()] = m
}

func (s *streamingStrategy) MeasurementResults(ctx context.Context, ids []string) ([]*AtlasMeasurement, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	measurements := make([]*AtlasMeasurement, 0)
	for _, id := range ids {
		m, found := s.results[id]
		if !found {
			continue
		}

		res := make([]*measurement.Result, 0)
		for _, v := range m {
			res = append(res, v)
		}

		r, err := atlasMeasurementForResults(res, id, s.workers)
		if err != nil {
			return nil, err
		}

		measurements = append(measurements, r)
	}

	return measurements, nil
}
