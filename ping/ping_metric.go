package ping

import (
	"fmt"
	"io"

	"github.com/DNS-OARC/ripeatlas/measurement"
)

type PingMetric struct {
	ProbeId int
	Min     float64
	Max     float64
	Avg     float64
	Sent    int
	Rcvd    int
	Dup     int
	Ttl     int
	Size    int
	Asn     int
}

func FromResult(r *measurement.Result) *PingMetric {
	return &PingMetric{ProbeId: r.PrbId(), Max: r.Max(), Min: r.Min(), Rcvd: r.Rcvd(), Avg: r.Avg(), Sent: r.Sent(), Dup: r.Dup(), Ttl: r.Ttl(), Size: r.Size()}
}

func (p *PingMetric) Write(w io.Writer, pk string) {
	if p.Min > 0 {
		p.writeMetric(pk, "success", 1, w)
		p.writeMetric(pk, "min_latency", p.Min, w)
		p.writeMetric(pk, "max_latency", p.Max, w)
		p.writeMetric(pk, "avg_latency", p.Avg, w)
	} else {
		p.writeMetric(pk, "success", 0, w)
	}

	p.writeMetric(pk, "sent", p.Sent, w)
	p.writeMetric(pk, "received", p.Rcvd, w)
	p.writeMetric(pk, "dup", p.Dup, w)
	p.writeMetric(pk, "ttl", p.Ttl, w)
	p.writeMetric(pk, "size", p.Size, w)
}

func (p *PingMetric) writeMetric(pk string, name string, value interface{}, w io.Writer) {
	const prefix = "atlas_ping_"
	fmt.Fprintf(w, prefix+"%s{measurement=\"%s\",probe=\"%d\",asn=\"%d\"} %v\n", name, pk, p.ProbeId, p.Asn, value)
}

func (p *PingMetric) SetAsn(asn int) {
	p.Asn = asn
}

func (p *PingMetric) Isvalid() bool {
	return p.Asn > 0
}
