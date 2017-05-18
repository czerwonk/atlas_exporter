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
	Asn      string
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
	t.writeMetric(pk, "hops", t.HopCount, w)
	t.writeMetric(pk, "success", t.Success, w)
}

func (t *TracerouteMetric) writeMetric(pk string, name string, value interface{}, w io.Writer) {
	const prefix = "atlas_traceroute_"
	fmt.Fprintf(w, prefix+"%s{measurement=\"%s\",probe=\"%d\",asn=\"%s\"} %v\n", name, pk, t.ProbeId, t.Asn, value)
}

func (t *TracerouteMetric) SetAsn(asn string) {
	t.Asn = asn
}
