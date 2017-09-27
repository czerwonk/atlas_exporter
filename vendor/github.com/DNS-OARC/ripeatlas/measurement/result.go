// Author Jerry Lundström <jerry@dns-oarc.net>
// Copyright (c) 2017, OARC, Inc.
// All rights reserved.
//
// This file is part of ripeatlas.
//
// ripeatlas is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// ripeatlas is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with ripeatlas.  If not, see <http://www.gnu.org/licenses/>.

package measurement

import (
    "encoding/base64"
    "encoding/json"
    "fmt"

    "github.com/DNS-OARC/ripeatlas/measurement/dns"
    "github.com/DNS-OARC/ripeatlas/measurement/http"
    "github.com/DNS-OARC/ripeatlas/measurement/ntp"
    "github.com/DNS-OARC/ripeatlas/measurement/ping"
    "github.com/DNS-OARC/ripeatlas/measurement/sslcert"
    "github.com/DNS-OARC/ripeatlas/measurement/traceroute"
    mdns "github.com/miekg/dns"
)

// A measurement result object, data availability depends on the type
// of measurement and some attributes are shared between measurements.
type Result struct {
    ParseError error

    data struct {
        // DNS and shared data
        Fw         int             `json:"fw"`
        Af         int             `json:"af"`
        DstAddr    string          `json:"dst_addr"`
        DstName    string          `json:"dst_name"`
        Error      json.RawMessage `json:"error"`
        From       string          `json:"from"`
        Lts        int             `json:"lts"`
        MsmId      int             `json:"msm_id"`
        PrbId      int             `json:"prb_id"`
        Proto      string          `json:"proto"`
        Qbuf       string          `json:"qbuf"`
        Result     json.RawMessage `json:"result"`
        Resultsets json.RawMessage `json:"resultset"`
        Retry      int             `json:"retry"`
        Timestamp  int             `json:"timestamp"`
        Type       string          `json:"type"`

        // Ping data (uses shared data)
        Avg     float64 `json:"avg"`
        Dup     int     `json:"dup"`
        Max     float64 `json:"max"`
        Min     float64 `json:"min"`
        Name    string  `json:"name"`
        Rcvd    int     `json:"rcvd"`
        Sent    int     `json:"sent"`
        Size    int     `json:"size"`
        SrcAddr string  `json:"src_addr"`
        Ttl     int     `json:"ttl"`

        // Traceroute data (uses shared and ping data)
        Endtime int `json:"endtime"`
        ParisId int `json:"paris_id"`

        // Http data (uses shared data)
        Uri string `json:"uri"`

        // Ntp data (uses shared and ping data)
        DstPort        string  `json:"dst_port"`
        Li             string  `json:"li"`
        Mode           string  `json:"mode"`
        Poll           float64 `json:"poll"`
        Precision      float64 `json:"precision"`
        RefId          string  `json:"ref-id"`
        RefTs          float64 `json:"ref-ts"`
        RootDelay      float64 `json:"root-delay"`
        RootDispersion float64 `json:"root-dispersion"`
        Stratum        int     `json:"stratum"`
        Version        int     `json:"version"`

        // Sslcert data (uses shared, ping and ntp data)
        Alert  json.RawMessage `json:"alert"`
        Cert   []string        `json:"cert"`
        Method string          `json:"method"`
        Rt     float64         `json:"rt"`
        Ttc    float64         `json:"ttc"`
        Ver    string          `json:"ver"`

        // Wifi data (uses shared)
        Bundle        int               `json:"bundle"`
        MsmName       string            `json:"msm_name"`
        GroupId       int               `json:"group_id"`
        WpaSupplicant map[string]string `json:"wpa_supplicant"`
    }

    dnsError      *dns.Error
    dnsResult     *dns.Result
    dnsResultsets []*dns.Resultset

    pingResults []*ping.Result

    tracerouteResults []*traceroute.Result

    httpResults []*http.Result

    ntpResults []*ntp.Result

    sslcertAlert *sslcert.Alert
}

func (r *Result) UnmarshalJSON(b []byte) error {
    if err := json.Unmarshal(b, &r.data); err != nil {
        return fmt.Errorf("%s for %s", err.Error(), string(b))
    }

    switch r.data.Type {
    case "dns":
        if r.data.Error != nil {
            r.dnsError = &dns.Error{}
            if err := json.Unmarshal(r.data.Error, r.dnsError); err != nil {
                return fmt.Errorf("Unable to process DNS error (fw %d): %s", r.data.Fw, err.Error())
            }
        }
        if r.data.Result != nil {
            r.dnsResult = &dns.Result{}
            if err := json.Unmarshal(r.data.Result, r.dnsResult); err != nil {
                return fmt.Errorf("Unable to process DNS result (fw %d): %s", r.data.Fw, err.Error())
            }
        }
        if r.data.Resultsets != nil {
            if err := json.Unmarshal(r.data.Resultsets, &r.dnsResultsets); err != nil {
                return fmt.Errorf("Unable to process DNS resultset (fw %d): %s", r.data.Fw, err.Error())
            }
        }
    case "ping":
        if r.data.Result != nil {
            if err := json.Unmarshal(r.data.Result, &r.pingResults); err != nil {
                return fmt.Errorf("Unable to process Ping result (fw %d): %s", r.data.Fw, err.Error())
            }
        }
    case "traceroute":
        if r.data.Result != nil {
            if err := json.Unmarshal(r.data.Result, &r.tracerouteResults); err != nil {
                return fmt.Errorf("Unable to process Traceroute result (fw %d): %s", r.data.Fw, err.Error())
            }
        }
    case "http":
        if r.data.Result != nil {
            if err := json.Unmarshal(r.data.Result, &r.httpResults); err != nil {
                return fmt.Errorf("Unable to process HTTP result (fw %d): %s", r.data.Fw, err.Error())
            }
        }
    case "ntp":
        if r.data.Result != nil {
            if err := json.Unmarshal(r.data.Result, &r.ntpResults); err != nil {
                return fmt.Errorf("Unable to process NTP result (fw %d): %s", r.data.Fw, err.Error())
            }
        }
    case "sslcert":
        if r.data.Alert != nil {
            r.sslcertAlert = &sslcert.Alert{}
            if err := json.Unmarshal(r.data.Alert, r.sslcertAlert); err != nil {
                return fmt.Errorf("Unable to process SSL Cert alert (fw %d): %s", r.data.Fw, err.Error())
            }
        }
    }

    return nil
}

// The firmware version used by the probe that generated this result.
func (r *Result) Fw() int {
    return r.data.Fw
}

// IP version: "4" or "6" (optional).
func (r *Result) Af() int {
    return r.data.Af
}

// IP address of the destination (optional).
func (r *Result) DstAddr() string {
    return r.data.DstAddr
}

// Hostname of the destination (optional).
func (r *Result) DstName() string {
    return r.data.DstName
}

// IP address of the source (optional).
func (r *Result) From() string {
    return r.data.From
}

// Last time synchronised. How long ago (in seconds) the clock of the probe
// was found to be in sync with that of a controller. The value -1 is used
// to indicate that the probe does not know whether it is in sync.
func (r *Result) Lts() int {
    return r.data.Lts
}

// Measurement identifier.
func (r *Result) MsmId() int {
    return r.data.MsmId
}

// Source probe ID.
func (r *Result) PrbId() int {
    return r.data.PrbId
}

// Protocol.
func (r *Result) Proto() string {
    return r.data.Proto
}

// Query payload buffer which was sent to the server, UU encoded (optional).
func (r *Result) Qbuf() string {
    return r.data.Qbuf
}

// Decode the Qbuf(), returns a *Msg from the github.com/miekg/dns package
// or nil on error or if Qbuf() is empty.
func (r *Result) UnpackQbuf() (*mdns.Msg, error) {
    if r.data.Qbuf == "" {
        return nil, nil
    }

    b, err := base64.StdEncoding.DecodeString(r.data.Qbuf)
    if err != nil {
        return nil, err
    }

    m := &mdns.Msg{}
    if err := m.Unpack(b); err != nil {
        return nil, err
    }

    return m, nil
}

// Retry count (optional).
func (r *Result) Retry() int {
    return r.data.Retry
}

// Start time, in Unix timestamp.
func (r *Result) Timestamp() int {
    return r.data.Timestamp
}

// The type of measurement within this result.
func (r *Result) Type() string {
    return r.data.Type
}

// Average round-trip time.
func (r *Result) Avg() float64 {
    return r.data.Avg
}

// Number of duplicate packets.
func (r *Result) Dup() int {
    return r.data.Dup
}

// Maximum round-trip time.
func (r *Result) Max() float64 {
    return r.data.Max
}

// Minimum round-trip time.
func (r *Result) Min() float64 {
    return r.data.Min
}

// Name of the destination (deprecated).
func (r *Result) Name() string {
    return r.data.Name
}

// Number of packets received.
func (r *Result) Rcvd() int {
    return r.data.Rcvd
}

// Number of packets sent.
func (r *Result) Sent() int {
    return r.data.Sent
}

// Packet size.
func (r *Result) Size() int {
    return r.data.Size
}

// Source address used by probe (missing due to a bug).
func (r *Result) SrcAddr() string {
    return r.data.SrcAddr
}

// Time-to-live field in the first reply (missing due to a bug).
func (r *Result) Ttl() int {
    return r.data.Ttl
}

// Unix timestamp for end of measurement.
func (r *Result) Endtime() int {
    return r.data.Endtime
}

// Variation for the Paris mode of traceroute.
func (r *Result) ParisId() int {
    return r.data.ParisId
}

// Request uri.
func (r *Result) Uri() string {
    return r.data.Uri
}

// Port name.
func (r *Result) DstPort() string {
    return r.data.DstPort
}

// Leap indicator, values "no", "61", "59", or "unknown".
func (r *Result) Li() string {
    return r.data.Li
}

// "server".
func (r *Result) Mode() string {
    return r.data.Mode
}

// Poll interval in seconds.
func (r *Result) Poll() float64 {
    return r.data.Poll
}

// Precision of the server's clock in seconds.
func (r *Result) Precision() float64 {
    return r.data.Precision
}

// Server's reference clock.
func (r *Result) RefId() string {
    return r.data.RefId
}

// Server's reference timestamp in NTP seconds.
func (r *Result) RefTs() float64 {
    return r.data.RefTs
}

// Round-trip delay from server to stratum 0 time source in seconds.
func (r *Result) RootDelay() float64 {
    return r.data.RootDelay
}

// Total dispersion to stratum 0 time source in seconds.
func (r *Result) RootDispersion() float64 {
    return r.data.RootDispersion
}

// Distance in hops from server to primary time source.
func (r *Result) Stratum() int {
    return r.data.Stratum
}

// NTP protocol version.
func (r *Result) Version() int {
    return r.data.Version
}

// Results of query, not present if "alert" is present (optional).
func (r *Result) Cert() []string {
    return r.data.Cert
}

// "SSL".
func (r *Result) Method() string {
    return r.data.Method
}

// Response time in milli seconds from starting to connect to receving
// the certificates (optional).
func (r *Result) Rt() float64 {
    return r.data.Rt
}

// Time in milli seconds that it took to connect (over TCP) to the
// target (optional).
func (r *Result) Ttc() float64 {
    return r.data.Ttc
}

// (SSL) protocol version.
func (r *Result) Ver() string {
    return r.data.Ver
}

// Wifi bundle (undocumented).
func (r *Result) Bundle() int {
    return r.data.Bundle
}

// Wifi msm_name (undocumented).
func (r *Result) MsmName() string {
    return r.data.MsmName
}

// Wifi group_id (undocumented).
func (r *Result) GroupId() int {
    return r.data.GroupId
}

// Wifi wpa_supplicant (undocumented).
func (r *Result) WpaSupplicant() map[string]string {
    return r.data.WpaSupplicant
}

// DNS error message, nil if the type of measurement is not "dns" (optional).
func (r *Result) DnsError() *dns.Error {
    return r.dnsError
}

// DNS response from the DNS server, nil if the type of measurement is
// not "dns" (optional).
func (r *Result) DnsResult() *dns.Result {
    return r.dnsResult
}

// An array of objects containing the DNS results when querying multiple
// local resolvers, empty if the type of measurement is not "dns" (optional).
func (r *Result) DnsResultsets() []*dns.Resultset {
    return r.dnsResultsets
}

// Ping results, nil if the type of measurement is not "ping" (optional).
func (r *Result) PingResults() []*ping.Result {
    return r.pingResults
}

// Traceroute results, nil if the type of measurement is not
// "traceroute" (optional).
func (r *Result) TracerouteResults() []*traceroute.Result {
    return r.tracerouteResults
}

// HTTP results, nil if the type of measurement is not
// "http" (optional).
func (r *Result) HttpResults() []*http.Result {
    return r.httpResults
}

// NTP results, nil if the type of measurement is not
// "ntp" (optional).
func (r *Result) NtpResults() []*ntp.Result {
    return r.ntpResults
}

// Error sent by server (see RFC 5246, Section 7.2) (from firmware 4720),
// nil if the type of measurement is not "sslcert" (optional).
func (r *Result) SslcertAlert() *sslcert.Alert {
    return r.sslcertAlert
}
