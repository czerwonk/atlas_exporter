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
	Rtt      float64
	Asn      int
}

func FromResult(r *measurement.Result) *TracerouteMetric {
	m := &TracerouteMetric{ProbeId: r.PrbId(), HopCount: len(r.TracerouteResults())}
	processLastHop(r, m)

	return m
}

func processLastHop(r *measurement.Result, m *TracerouteMetric) {
	if len(r.TracerouteResults()) == 0 {
		return
	}

	last := r.TracerouteResults()[len(r.TracerouteResults())-1]
	for _, rep := range last.Replies() {
		if rep.From() == r.DstAddr() {
			m.Success = 1
			m.Rtt = rep.Rtt()
		}
	}
}

func (t *TracerouteMetric) Write(w io.Writer, pk string) {
	t.writeMetric(pk, "hops", t.HopCount, w)
	t.writeMetric(pk, "success", t.Success, w)

	if t.Rtt > 0 {
		t.writeMetric(pk, "rtt", t.Rtt, w)
	}
}

func (t *TracerouteMetric) writeMetric(pk string, name string, value interface{}, w io.Writer) {
	const prefix = "atlas_traceroute_"
	fmt.Fprintf(w, prefix+"%s{measurement=\"%s\",probe=\"%d\",asn=\"%d\"} %v\n", name, pk, t.ProbeId, t.Asn, value)
}

func (t *TracerouteMetric) SetAsn(asn int) {
	t.Asn = asn
}
