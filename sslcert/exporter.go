// SPDX-License-Identifier: LGPL-3.0-or-later

package sslcert

import (
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"strconv"

	"github.com/DNS-OARC/ripeatlas/measurement"
	"github.com/czerwonk/atlas_exporter/probe"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	labels               []string
	rttDesc              *prometheus.Desc
	sslVerDesc           *prometheus.Desc
	successDesc          *prometheus.Desc
	alertLevelDesc       *prometheus.Desc
	alertDescriptionDesc *prometheus.Desc
)

func init() {
	labels = []string{"measurement", "probe", "dst_addr", "asn", "ip_version", "country_code", "lat", "long", "cert_fingerprint", "cert_issuer"}

	successDesc = prometheus.NewDesc(prometheus.BuildFQName(ns, sub, "success"), "Destination was reachable", labels, nil)
	sslVerDesc = prometheus.NewDesc(prometheus.BuildFQName(ns, sub, "version"), "SSL/TLS version used for the request", labels, nil)
	rttDesc = prometheus.NewDesc(prometheus.BuildFQName(ns, sub, "rtt"), "Round trip time in ms", labels, nil)
	alertLevelDesc = prometheus.NewDesc(prometheus.BuildFQName(ns, sub, "alert_level"), "Status of the SSL/TLS certificate (0 = valid)", labels, nil)
	alertDescriptionDesc = prometheus.NewDesc(prometheus.BuildFQName(ns, sub, "alert_description"), "Description for the alert level (see RIPE Atlas documentation)", labels, nil)
}

type sslCertExporter struct {
	id string
}

func fingerprintFromResult(res *measurement.Result) string {
	certs := res.Cert()
	if len(certs) == 0 {
		return ""
	}

	if block, _ := pem.Decode([]byte(certs[0])); block != nil {
		sum := sha256.Sum256(block.Bytes)
		return fmt.Sprintf("%x", sum)
	}

	if der, err := base64.StdEncoding.DecodeString(certs[0]); err == nil {
		sum := sha256.Sum256(der)
		return fmt.Sprintf("%x", sum)
	}

	return ""
}

func issuerOrgFromResult(res *measurement.Result) string {
	certs := res.Cert()
	if len(certs) == 0 {
		return "unknown"
	}

	for _, raw := range certs {
		var der []byte

		if block, _ := pem.Decode([]byte(raw)); block != nil {
			der = block.Bytes
		} else {
			// base64 DER
			b, err := base64.StdEncoding.DecodeString(raw)
			if err != nil {
				continue
			}
			der = b
		}

		cert, err := x509.ParseCertificate(der)
		if err != nil {
			continue
		}

		if len(cert.Issuer.Organization) > 0 && cert.Issuer.Organization[0] != "" {
			return cert.Issuer.Organization[0]
		}

		// если O пустой — пробуем CN
		if cn := cert.Issuer.CommonName; cn != "" {
			return cn
		}

		// совсем ничего — unknown
		return "unknown"
	}

	return "unknown"
}

// Export exports a prometheus metric
func (m *sslCertExporter) Export(res *measurement.Result, probe *probe.Probe, ch chan<- prometheus.Metric) {
	fp := fingerprintFromResult(res)
	issuer := issuerOrgFromResult(res)

	labelValues := []string{
		m.id,
		strconv.Itoa(probe.ID),
		res.DstAddr(),
		strconv.Itoa(probe.ASNForIPVersion(res.Af())),
		strconv.Itoa(res.Af()),
		probe.CountryCode,
		probe.Latitude(),
		probe.Longitude(),
		fp,
		issuer,
	}

	ver, _ := strconv.ParseFloat(res.Ver(), 64)
	ch <- prometheus.MustNewConstMetric(sslVerDesc, prometheus.GaugeValue, ver, labelValues...)

	var alertLevel, alertDescription float64
	if res.SslcertAlert() != nil {
		alertLevel = float64(res.SslcertAlert().Level())
		alertDescription = float64(res.SslcertAlert().Description())
	}
	ch <- prometheus.MustNewConstMetric(alertLevelDesc, prometheus.GaugeValue, alertLevel, labelValues...)
	ch <- prometheus.MustNewConstMetric(alertDescriptionDesc, prometheus.GaugeValue, alertDescription, labelValues...)

	if res.Rt() > 0 {
		ch <- prometheus.MustNewConstMetric(successDesc, prometheus.GaugeValue, 1, labelValues...)
		ch <- prometheus.MustNewConstMetric(rttDesc, prometheus.GaugeValue, res.Rt(), labelValues...)
	} else {
		ch <- prometheus.MustNewConstMetric(successDesc, prometheus.GaugeValue, 0, labelValues...)
	}
}

// Describe exports metric descriptions for Prometheus
func (m *sslCertExporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- successDesc
	ch <- rttDesc
	ch <- sslVerDesc
	ch <- alertLevelDesc
	ch <- alertDescriptionDesc
}
