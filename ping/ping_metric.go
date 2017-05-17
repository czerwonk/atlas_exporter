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
}

func FromResult(r *measurement.Result) *PingMetric {
	return &PingMetric{ProbeId: r.PrbId(), Max: r.Max(), Min: r.Min(), Rcvd: r.Rcvd(), Avg: r.Avg(), Sent: r.Sent()}
}

func (p *PingMetric) Write(w io.Writer, pk string) {
	const prefix = "atlas_ping_"
	if p.Min > 0 {
		fmt.Fprintf(w, prefix+"success{measurement=\"%s\",probe=\"%d\"} 1\n", pk, p.ProbeId)
		fmt.Fprintf(w, prefix+"min_latency{measurement=\"%s\",probe=\"%d\"} %f\n", pk, p.ProbeId, p.Min)
		fmt.Fprintf(w, prefix+"max_latency{measurement=\"%s\",probe=\"%d\"} %f\n", pk, p.ProbeId, p.Max)
		fmt.Fprintf(w, prefix+"avg_latency{measurement=\"%s\",probe=\"%d\"} %f\n", pk, p.ProbeId, p.Avg)
	} else {
		fmt.Fprintf(w, prefix+"success{measurement=\"%s\",probe=\"%d\"} 0\n", pk, p.ProbeId)
	}

	fmt.Fprintf(w, prefix+"sent{measurement=\"%s\",probe=\"%d\"} %d\n", pk, p.ProbeId, p.Sent)
	fmt.Fprintf(w, prefix+"received{measurement=\"%s\",probe=\"%d\"} %d\n", pk, p.ProbeId, p.Rcvd)
}
