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

const connectionRetryInterval = 30 * time.Second

type streamingStrategy struct {
	stream         *ripeatlas.Stream
	measurements   map[string]*exporter.Measurement
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
		measurements:   make(map[string]*exporter.Measurement),
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
		case <-time.After(connectionRetryInterval):
			s.clearResults(m.ID)
			continue
		}
	}
}

func (s *streamingStrategy) clearResults(id string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.measurements, id)
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

	mes, found := s.measurements[msm]
	if !found {
		var err error
		mes, err = measurementForType(m.Type(), msm, strconv.Itoa(m.Af()), s.cfg)
		if err != nil {
			log.Error(err)
			return
		}

		s.measurements[msm] = mes
	}

	mes.Add(m, probe)
}

func (s *streamingStrategy) MeasurementResults(ctx context.Context, ids []string) ([]*exporter.Measurement, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	result := make([]*exporter.Measurement, 0)
	for _, id := range ids {
		m, found := s.measurements[id]
		if !found {
			continue
		}

		result = append(result, m)
	}

	return result, nil
}
