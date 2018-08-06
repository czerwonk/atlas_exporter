package traceroute

import (
	"github.com/DNS-OARC/ripeatlas/measurement"
	"github.com/czerwonk/atlas_exporter/config"
	"github.com/czerwonk/atlas_exporter/exporter"
)

const (
	ns  = "atlas"
	sub = "traceroute"
)

// NewMeasurement returns a new instance of `exorter.Measurement` for a traceroute measurement
func NewMeasurement(id, ipVersion string, cfg *config.Config) *exporter.Measurement {
	opts := []exporter.MeasurementOpt{
		exporter.WithHistograms(newRttHistogram(id, ipVersion, cfg.HistogramBuckets.Traceroute.Rtt)),
	}

	if cfg.FilterInvalidResults {
		opts = append(opts, exporter.WithValidator(&tracerouteResultValidator{}))
	}

	return exporter.NewMeasurement(&tracerouteExporter{id}, opts...)
}

func processLastHop(r *measurement.Result) (success float64, rtt float64) {
	if len(r.TracerouteResults()) == 0 {
		return success, rtt
	}

	last := r.TracerouteResults()[len(r.TracerouteResults())-1]
	for _, rep := range last.Replies() {
		if rep.From() == r.DstAddr() {
			success = 1
		}

		repl := last.Replies()
		if len(repl) > 0 {
			rtt = repl[len(repl)-1].Rtt()
		}
	}

	return success, rtt
}
