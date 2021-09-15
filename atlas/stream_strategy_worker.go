package atlas

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/DNS-OARC/ripeatlas"
	"github.com/DNS-OARC/ripeatlas/measurement"
	"github.com/czerwonk/atlas_exporter/config"
	log "github.com/sirupsen/logrus"
)

const connectionRetryInterval = 30 * time.Second

type streamStrategyWorker struct {
	resultCh    chan<- *measurement.Result
	resetCh     chan<- *config.Measurement
	measurement config.Measurement
	timeout     time.Duration
}

func (w *streamStrategyWorker) run(ctx context.Context) error {
	for {
		ch, err := w.subscribe()
		if err != nil {
			log.Error(err)
		} else {
			log.Infof("Subscribed to results of measurement #%s", w.measurement.ID)
			w.listenForResults(ctx, w.timeout, ch)
		}

		w.resetCh <- &w.measurement

		select {
		case <-ctx.Done():
			return nil
		case <-time.After(connectionRetryInterval):
			continue
		}
	}
}

func (w *streamStrategyWorker) subscribe() (<-chan *measurement.Result, error) {
	stream := ripeatlas.NewStream()

	msm, err := strconv.Atoi(w.measurement.ID)
	if err != nil {
		return nil, err
	}

	ch, err := stream.MeasurementResults(ripeatlas.Params{
		"msm": msm,
	})
	if err != nil {
		return nil, err
	}

	return ch, nil
}

func (w *streamStrategyWorker) listenForResults(ctx context.Context, timeout time.Duration, ch <-chan *measurement.Result) {
	for {
		select {
		case m := <-ch:
			if m == nil {
				continue
			}

			if m.ParseError != nil {
				log.Error(m.ParseError)
			}

			if m.ParseError != nil && strings.HasPrefix(m.ParseError.Error(), "c.On(disconnect)") {
				log.Error(m.ParseError)
				return
			}

			w.resultCh <- m
		case <-time.After(timeout):
			log.Errorf("Timeout reached for measurement #%s. Trying to reconnect.", w.measurement.ID)
			return
		case <-ctx.Done():
			return
		}
	}
}
