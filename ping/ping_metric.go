package ping

import (
	"fmt"
	"io"

	"github.com/DNS-OARC/ripeatlas/measurement"
)

type PingMetric struct {
	ProbeId   int
	DstAddr   string
	DstName   string
	Min       float64
	Max       float64
	Avg       float64
	Sent      int
	Rcvd      int
	Dup       int
	Ttl       int
	Size      int
	Asn       int
	IpVersion int
}

func FromResult(r *measurement.Result) *PingMetric {
	return &PingMetric{ProbeId: r.PrbId(), DstAddr: r.DstAddr(), DstName: r.DstName(), Max: r.Max(), Min: r.Min(), Rcvd: r.Rcvd(), Avg: r.Avg(), Sent: r.Sent(), Dup: r.Dup(), Ttl: r.Ttl(), Size: r.Size(), IpVersion: r.Af()}
}

func (m *PingMetric) Write(w io.Writer, pk string) {
	if m.Min > 0 {
		m.writeMetric(pk, "success", 1, w)
		m.writeMetric(pk, "min_latency", m.Min, w)
		m.writeMetric(pk, "max_latency", m.Max, w)
		m.writeMetric(pk, "avg_latency", m.Avg, w)
	} else {
		m.writeMetric(pk, "success", 0, w)
	}

	m.writeMetric(pk, "sent", m.Sent, w)
	m.writeMetric(pk, "received", m.Rcvd, w)
	m.writeMetric(pk, "dup", m.Dup, w)
	m.writeMetric(pk, "ttl", m.Ttl, w)
	m.writeMetric(pk, "size", m.Size, w)
}

func (m *PingMetric) writeMetric(pk string, name string, value interface{}, w io.Writer) {
	const prefix = "atlas_ping_"
	fmt.Fprintf(w, prefix+"%s{measurement=\"%s\",probe=\"%d\",dst_addr=\"%s\",dst_name=\"%s\",asn=\"%d\",ip_version=\"%d\"} %v\n", name, pk, m.ProbeId, m.DstAddr, m.DstName, m.Asn, m.IpVersion, value)
}

func (m *PingMetric) SetAsn(asn int) {
	m.Asn = asn
}

func (m *PingMetric) Isvalid() bool {
	return m.Asn > 0
}
