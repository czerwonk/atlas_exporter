package atlas

import (
	"context"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/czerwonk/atlas_exporter/exporter"
	"github.com/czerwonk/atlas_exporter/probe"

	"github.com/DNS-OARC/ripeatlas/measurement"
	"github.com/czerwonk/atlas_exporter/config"
	"github.com/prometheus/common/log"

	"github.com/DNS-OARC/ripeatlas"
)

const ConnectionRetryInterval = 30 * time.Second

type streamingStrategy struct {
	stream         *ripeatlas.Stream
	results        map[string]*exporter.ResultHandler
	workers        uint
	cfg            *config.Config
	defaultTimeout time.Duration
	mu             sync.Mutex
}

// NewStreamingStrategy returns an strategy using the RIPE Atlas Streaming API
func NewStreamingStrategy(ctx context.Context, cfg *config.Config, workers uint, defaultTimeout time.Duration) Strategy {
	s := &streamingStrategy{
		stream:         ripeatlas.NewStream(),
		workers:        workers,
		defaultTimeout: defaultTimeout,
		cfg:            cfg,
		results:        make(map[string]*exporter.ResultHandler),
	}

	s.start(ctx, cfg.Measurements)
	return s
}

func (s *streamingStrategy) start(ctx context.Context, measurements []config.Measurement) {
	for _, m := range measurements {
		go s.startListening(ctx, m)
	}
}

func (s *streamingStrategy) startListening(ctx context.Context, m config.Measurement) {
	for {
		ch, err := s.subscribe(m.ID)
		if err != nil {
			log.Error(err)
		} else {
			log.Infof("Subscribed to results of measurement #%s", m.ID)
			s.listenForResults(ctx, s.timeoutForMeasurement(m), ch)
		}

		select {
		case <-ctx.Done():
			return
		case <-time.After(ConnectionRetryInterval):
			delete(s.results, m.ID)
			continue
		}
	}
}

func (s *streamingStrategy) timeoutForMeasurement(m config.Measurement) time.Duration {
	if m.Timeout == 0 {
		return s.defaultTimeout
	}

	return m.Timeout
}

func (s *streamingStrategy) subscribe(id string) (<-chan *measurement.Result, error) {
	msm, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}

	ch, err := s.stream.MeasurementResults(ripeatlas.Params{
		"msm": msm,
	})
	if err != nil {
		return nil, err
	}

	return ch, nil
}

func (s *streamingStrategy) listenForResults(ctx context.Context, timeout time.Duration, ch <-chan *measurement.Result) {
	for {
		select {
		case m := <-ch:
			if m.ParseError != nil {
				log.Error(m.ParseError)
			}

			if m.ParseError != nil && strings.HasPrefix(m.ParseError.Error(), "c.On(disconnect)") {
				log.Error(m.ParseError)
				return
			}

			s.processMeasurement(m)
		case <-time.After(timeout):
			log.Errorf("Timeout reached. Trying to reconnect.")
			return
		case <-ctx.Done():
			return
		}
	}
}

func (s *streamingStrategy) processMeasurement(m *measurement.Result) {
	log.Infof("Got result for %d from probe %d", m.MsmId(), m.PrbId())

	probe, err := probeForID(m.PrbId())
	if err != nil {
		log.Error(err)
		return
	}

	s.add(m, probe)
}

func (s *streamingStrategy) add(m *measurement.Result, probe *probe.Probe) {
	s.mu.Lock()
	defer s.mu.Unlock()

	msm := strconv.Itoa(m.MsmId())

	h, found := s.results[msm]
	if !found {
		var err error
		h, err = handlerForType(m.Type(), msm, s.cfg)
		if err != nil {
			log.Error(err)
			return
		}

		s.results[msm] = h
	}

	h.Add(m, probe)
}

func (s *streamingStrategy) MeasurementResults(ctx context.Context, ids []string) ([]*exporter.ResultHandler, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	measurements := make([]*exporter.ResultHandler, 0)
	for _, id := range ids {
		res, found := s.results[id]
		if !found {
			continue
		}

		measurements = append(measurements, res)
		res.Scraped()
	}

	return measurements, nil
}
