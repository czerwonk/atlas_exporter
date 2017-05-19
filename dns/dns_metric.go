package dns

import (
	"fmt"
	"io"

	"github.com/DNS-OARC/ripeatlas/measurement"
)

type DnsMetric struct {
	ProbeId   int
	Rtt       float64
	Asn       int
	Success   int
	IpVersion int
}

func FromResult(r *measurement.Result) *DnsMetric {
	var rtt float64
	if r.DnsResult() != nil {
		rtt = r.DnsResult().Rt()
	}

	var success int
	if r.DnsError() == nil {
		success = 1
	}

	return &DnsMetric{ProbeId: r.PrbId(), Rtt: rtt, Success: success, IpVersion: r.Af()}
}

func (m *DnsMetric) Write(w io.Writer, pk string) {
	m.writeMetric(pk, "rtt", m.Rtt, w)
	m.writeMetric(pk, "success", m.Success, w)
}

func (m *DnsMetric) writeMetric(pk string, name string, value interface{}, w io.Writer) {
	const prefix = "atlas_dns_"
	fmt.Fprintf(w, prefix+"%s{measurement=\"%s\",probe=\"%d\",asn=\"%d\",ip_version=\"%d\"} %v\n", name, pk, m.ProbeId, m.Asn, m.IpVersion, value)
}

func (m *DnsMetric) SetAsn(asn int) {
	m.Asn = asn
}

func (m *DnsMetric) Isvalid() bool {
	return m.Asn > 0
}
