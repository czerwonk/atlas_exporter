// SPDX-License-Identifier: LGPL-3.0-or-later

package dns

import (
	"strconv"

	"github.com/DNS-OARC/ripeatlas/measurement"
	"github.com/czerwonk/atlas_exporter/probe"
	mdns "github.com/miekg/dns"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	labels      []string
	successDesc *prometheus.Desc
	rttDesc     *prometheus.Desc
	answerDesc  *prometheus.Desc
)

func init() {
	labels = []string{"measurement", "probe", "dst_addr", "asn", "ip_version", "country_code", "lat", "long"}

	successDesc = prometheus.NewDesc(prometheus.BuildFQName(ns, sub, "success"), "Destination was reachable", labels, nil)
	rttDesc = prometheus.NewDesc(prometheus.BuildFQName(ns, sub, "rtt"), "Roundtrip time in ms", labels, nil)
	answerDesc = prometheus.NewDesc(
		prometheus.BuildFQName(ns, sub, "answer"),
		"DNS answer IP for query",
		[]string{
			"measurement",
			"probe",
			"resolver",
			"asn",
			"ip_version",
			"country_code",
			"lat",
			"long",
			"qname",
			"rr_type",
			"answer_ip",
		},
		nil,
	)
}

type dnsExporter struct {
	id string
}

// Export exports a prometheus metric
func (m *dnsExporter) Export(res *measurement.Result, p *probe.Probe, ch chan<- prometheus.Metric) {
	if rs := res.DnsResultsets(); len(rs) > 0 {
		for _, s := range rs {
			if s == nil {
				continue
			}

			labelValues := []string{
				m.id,
				strconv.Itoa(p.ID),
				s.DstAddr(),
				strconv.Itoa(p.ASNForIPVersion(s.Af())),
				strconv.Itoa(s.Af()),
				p.CountryCode,
				p.Latitude(),
				p.Longitude(),
			}

			if s.DnsError() != nil || s.Result() == nil {
				ch <- prometheus.MustNewConstMetric(successDesc, prometheus.GaugeValue, 0, labelValues...)
				continue
			}

			r := s.Result()
			rtt := r.Rt()

			if msg, err := r.UnpackAbuf(); err == nil && msg != nil {
				for _, ans := range msg.Answer {
					switch rr := ans.(type) {
					case *mdns.A:
						ch <- prometheus.MustNewConstMetric(
							answerDesc,
							prometheus.GaugeValue,
							1,
							[]string{
								m.id,
								strconv.Itoa(p.ID),
								s.DstAddr(),
								strconv.Itoa(p.ASNForIPVersion(s.Af())),
								strconv.Itoa(s.Af()),
								p.CountryCode,
								p.Latitude(),
								p.Longitude(),
								rr.Hdr.Name,
								"A",
								rr.A.String(),
							}...,
						)
					case *mdns.AAAA:
						ch <- prometheus.MustNewConstMetric(
							answerDesc,
							prometheus.GaugeValue,
							1,
							[]string{
								m.id,
								strconv.Itoa(p.ID),
								s.DstAddr(),
								strconv.Itoa(p.ASNForIPVersion(s.Af())),
								strconv.Itoa(s.Af()),
								p.CountryCode,
								p.Latitude(),
								p.Longitude(),
								rr.Hdr.Name,
								"AAAA",
								rr.AAAA.String(),
							}...,
						)
					}
				}
			}

			if rtt > 0 {
				ch <- prometheus.MustNewConstMetric(successDesc, prometheus.GaugeValue, 1, labelValues...)
				ch <- prometheus.MustNewConstMetric(rttDesc, prometheus.GaugeValue, rtt, labelValues...)
			} else {
				ch <- prometheus.MustNewConstMetric(successDesc, prometheus.GaugeValue, 0, labelValues...)
			}
		}
		return
	}

	labelValues := []string{
		m.id,
		strconv.Itoa(p.ID),
		res.DstAddr(),
		strconv.Itoa(p.ASNForIPVersion(res.Af())),
		strconv.Itoa(res.Af()),
		p.CountryCode,
		p.Latitude(),
		p.Longitude(),
	}

	var rtt float64
	dnsRes := res.DnsResult()
	if dnsRes != nil {
		rtt = dnsRes.Rt()

		if msg, err := dnsRes.UnpackAbuf(); err == nil && msg != nil {
			for _, ans := range msg.Answer {
				switch rr := ans.(type) {
				case *mdns.A:
					ch <- prometheus.MustNewConstMetric(
						answerDesc,
						prometheus.GaugeValue,
						1,
						[]string{
							m.id,
							strconv.Itoa(p.ID),
							res.DstAddr(),
							strconv.Itoa(p.ASNForIPVersion(res.Af())),
							strconv.Itoa(res.Af()),
							p.CountryCode,
							p.Latitude(),
							p.Longitude(),
							rr.Hdr.Name,
							"A",
							rr.A.String(),
						}...,
					)
				case *mdns.AAAA:
					ch <- prometheus.MustNewConstMetric(
						answerDesc,
						prometheus.GaugeValue,
						1,
						[]string{
							m.id,
							strconv.Itoa(p.ID),
							res.DstAddr(),
							strconv.Itoa(p.ASNForIPVersion(res.Af())),
							strconv.Itoa(res.Af()),
							p.CountryCode,
							p.Latitude(),
							p.Longitude(),
							rr.Hdr.Name,
							"AAAA",
							rr.AAAA.String(),
						}...,
					)
				}
			}
		}
	}

	if rtt > 0 {
		ch <- prometheus.MustNewConstMetric(successDesc, prometheus.GaugeValue, 1, labelValues...)
		ch <- prometheus.MustNewConstMetric(rttDesc, prometheus.GaugeValue, rtt, labelValues...)
	} else {
		ch <- prometheus.MustNewConstMetric(successDesc, prometheus.GaugeValue, 0, labelValues...)
	}
}

// Describe exports metric descriptions for Prometheus
func (m *dnsExporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- successDesc
	ch <- rttDesc
	ch <- answerDesc
}
