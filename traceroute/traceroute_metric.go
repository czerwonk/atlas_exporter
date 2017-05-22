package traceroute

import (
	"fmt"
	"io"

	"github.com/DNS-OARC/ripeatlas/measurement"
)

type TracerouteMetric struct {
	ProbeId   int
	DstAddr   string
	DstName   string
	HopCount  int
	Success   int
	Rtt       float64
	Asn       int
	IpVersion int
}

func FromResult(r *measurement.Result) *TracerouteMetric {
	m := &TracerouteMetric{ProbeId: r.PrbId(), DstAddr: r.DstAddr(), DstName: r.DstName(), HopCount: len(r.TracerouteResults()), IpVersion: r.Af()}
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

func (m *TracerouteMetric) Write(w io.Writer, pk string) {
	m.writeMetric(pk, "hops", m.HopCount, w)
	m.writeMetric(pk, "success", m.Success, w)

	if m.Rtt > 0 {
		m.writeMetric(pk, "rtt", m.Rtt, w)
	}
}

func (m *TracerouteMetric) writeMetric(pk string, name string, value interface{}, w io.Writer) {
	const prefix = "atlas_traceroute_"
	fmt.Fprintf(w, prefix+"%s{measurement=\"%s\",probe=\"%d\",dst_addr=\"%s\",dst_name=\"%s\",asn=\"%d\",ip_version=\"%d\"} %v\n", name, pk, m.ProbeId, m.DstAddr, m.DstName, m.Asn, m.IpVersion, value)
}

func (m *TracerouteMetric) SetAsn(asn int) {
	m.Asn = asn
}

func (m *TracerouteMetric) Isvalid() bool {
	return (m.Success == 1 || m.HopCount > 1) && m.Asn > 0
}
