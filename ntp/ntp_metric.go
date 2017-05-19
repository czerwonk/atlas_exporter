package ntp

import (
	"fmt"
	"io"

	"github.com/DNS-OARC/ripeatlas/measurement"
)

type NtpMetric struct {
	ProbeId        int
	Poll           float64
	Precision      float64
	RootDelay      float64
	RootDispersion float64
	Version        int
	Asn            int
	IpVersion      int
}

func FromResult(r *measurement.Result) *NtpMetric {
	return &NtpMetric{ProbeId: r.PrbId(), Poll: r.Poll(), Precision: r.Precision(), RootDelay: r.RootDelay(), RootDispersion: r.RootDispersion(), Version: r.Version(), IpVersion: r.Af()}
}

func (m *NtpMetric) Write(w io.Writer, pk string) {
	m.writeMetric(pk, "poll", m.Poll, w)
	m.writeMetric(pk, "precision", m.Precision, w)
	m.writeMetric(pk, "root_delay", m.RootDelay, w)
	m.writeMetric(pk, "root_dispersion", m.RootDispersion, w)
	m.writeMetric(pk, "ntp_version", m.Version, w)
}

func (m *NtpMetric) writeMetric(pk string, name string, value interface{}, w io.Writer) {
	const prefix = "atlas_ntp_"
	fmt.Fprintf(w, prefix+"%s{measurement=\"%s\",probe=\"%d\",asn=\"%d\",ip_version=\"%d\"} %v\n", name, pk, m.ProbeId, m.Asn, m.IpVersion, value)
}

func (m *NtpMetric) SetAsn(asn int) {
	m.Asn = asn
}

func (m *NtpMetric) Isvalid() bool {
	return m.Asn > 0
}
