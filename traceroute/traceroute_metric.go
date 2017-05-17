package traceroute

import (
	"fmt"
	"io"

	"github.com/DNS-OARC/ripeatlas/measurement"
)

type TracerouteMetric struct {
	ProbeId  int
	HopCount int
	Success  int
}

func FromResult(r *measurement.Result) *TracerouteMetric {
	success := getTracerouteSuccess(r)
	return &TracerouteMetric{ProbeId: r.PrbId(), HopCount: len(r.TracerouteResults()), Success: success}
}

func getTracerouteSuccess(r *measurement.Result) int {
	success := 0

	if len(r.TracerouteResults()) > 0 {
		last := r.TracerouteResults()[len(r.TracerouteResults())-1]
		for _, rep := range last.Replies() {
			if rep.From() == r.DstAddr() {
				success = 1
			}
		}
	}

	return success
}

func (t *TracerouteMetric) Write(w io.Writer, pk string) {
	const prefix = "atlas_traceroute_"
	fmt.Fprintf(w, prefix+"hops{measurement=\"%s\",probe=\"%d\"} %d\n", pk, t.ProbeId, t.HopCount)
	fmt.Fprintf(w, prefix+"success{measurement=\"%s\",probe=\"%d\"} %d\n", pk, t.ProbeId, t.Success)
}
