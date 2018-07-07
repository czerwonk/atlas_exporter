package atlas

import (
	"context"
	"strconv"

	"github.com/DNS-OARC/ripeatlas/measurement"

	"github.com/DNS-OARC/ripeatlas"
)

type streamingStrategy struct {
	stream  *ripeatlas.Stream
	results map[string][]*measurement.Result
	workers uint
}

// NewStreamingStrategy returns an strategy using the RIPE Atlas Streaming API
func NewStreamingStrategy(ctx context.Context, ids []string, workers uint) (Strategy, error) {
	s := &streamingStrategy{
		stream:  ripeatlas.NewStream(),
		workers: workers,
		results: make(map[string][]*measurement.Result),
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
			msm := strconv.Itoa(m.MsmId())
			res := s.results[msm]
			s.results[msm] = append(res, m)
		case <-ctx.Done():
			return
		}
	}
}

func (s *streamingStrategy) MeasurementResults(ctx context.Context, ids []string) ([]*AtlasMeasurement, error) {
	res := make([]*AtlasMeasurement, 0)
	for k, v := range s.results {
		r, err := atlasMeasurementForResults(v, k, s.workers)
		if err != nil {
			return nil, err
		}

		res = append(res, r)
	}

	return res, nil
}
