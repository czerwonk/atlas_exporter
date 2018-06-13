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

package ripeatlas

import (
    "encoding/json"
    "fmt"

    "github.com/DNS-OARC/ripeatlas/measurement"
)

type Measurement struct {
    ParseError error

    data struct {
        Id                    int             `json:"id"`
        Result                string          `json:"result"`
        GroupId               int             `json:"group_id"`
        Af                    int             `json:"af"`
        IsOneoff              bool            `json:"is_oneoff"`
        IsPublic              bool            `json:"is_public"`
        Description           string          `json:"description"`
        Spread                int             `json:"spread"`
        ResolveOnProbe        bool            `json:"resolve_on_probe"`
        StartTime             int             `json:"start_time"`
        StopTime              int             `json:"stop_time"`
        Type                  string          `json:"type"`
        Status                json.RawMessage `json:"status"`
        IsAllScheduled        bool            `json:"is_all_scheduled"`
        ParticipantCount      int             `json:"participant_count"`
        TargetAsn             int             `json:"target_asn"`
        TargetIp              string          `json:"target_ip"`
        CreationTime          int             `json:"creation_time"`
        InWifiGroup           bool            `json:"in_wifi_group"`
        ResolvedIps           []string        `json:"resolved_ips"`
        ProbesRequested       int             `json:"probes_requested"`
        ProbesScheduled       int             `json:"probes_scheduled"`
        Group                 string          `json:"group"`
        Probes                json.RawMessage `json:"probes"`
        ProbeSources          json.RawMessage `json:"probe_sources"`
        ParticipationRequests json.RawMessage `json:"participation_requests"`
        Interval              int             `json:"interval"`
        Packets               int             `json:"packets"`
        Size                  int             `json:"size"`
        PacketInterval        int             `json:"packet_interval"`
        IncludeProbeId        bool            `json:"include_probe_id"`
        Port                  int             `json:"port"`
        FirstHop              int             `json:"first_hop"`
        MaxHops               int             `json:"max_hops"`
        Paris                 int             `json:"paris"`
        Protocol              string          `json:"protocol"`
        ResponseTimeout       int             `json:"response_timeout"`
        DuplicateTimeout      int             `json:"duplicate_timeout"`
        HopByHopOptionSize    int             `json:"hop_by_hop_option_size"`
        DestinationOptionSize int             `json:"destination_option_size"`
        DontFragment          bool            `json:"dont_fragment"`
        UdpPayloadSize        int             `json:"udp_payload_size"`
        UseProbeResolver      bool            `json:"use_probe_resolver"`
        SetRdBit              bool            `json:"set_rd_bit"`
        PrependProbeId        bool            `json:"prepend_probe_id"`
        Retry                 int             `json:"retry"`
        IncludeQbuf           bool            `json:"include_qbuf"`
        SetNsidBit            bool            `json:"set_nsid_bit"`
        IncludeAbuf           bool            `json:"include_abuf"`
        QueryClass            string          `json:"query_class"`
        QueryArgument         string          `json:"query_argument"`
        QueryType             string          `json:"query_type"`
        SetCdBit              bool            `json:"set_cd_bit"`
        SetDoBit              bool            `json:"set_do_bit"`
        UseMacros             bool            `json:"use_macros"`
        Timeout               int             `json:"timeout"`
        ExtendedTiming        bool            `json:"extended_timing"`
        MoreExtendedTiming    bool            `json:"more_extended_timing"`
        HeaderBytes           int             `json:"header_bytes"`
        Method                string          `json:"method"`
        Path                  string          `json:"path"`
        QueryString           string          `json:"query_string"`
        UserAgent             string          `json:"user_agent"`
        MaxBytesRead          int             `json:"max_bytes_read"`
        Version               string          `json:"version"`
        Hostname              string          `json:"hostname"`
        Ipv4                  bool            `json:"ipv4"`
        Ipv6                  bool            `json:"ipv6"`
        Cert                  string          `json:"cert"`
        ExtraWait             int             `json:"extra_wait"`
        Ssid                  string          `json:"ssid"`
        KeyMgmt               string          `json:"key_mgmt"`
        Eap                   string          `json:"eap"`
        Identity              string          `json:"identity"`
        AnonymousIdentity     string          `json:"anonymous_identity"`
        Phase2                string          `json:"phase2"`
        Rssi                  bool            `json:"rssi"`
    }

    status *measurement.MeasurementStatus

    probes []*measurement.Probe

    probeSources []*measurement.ProbeSource

    participationRequests []*measurement.ProbeSource
}

func (m *Measurement) UnmarshalJSON(b []byte) error {
    if err := json.Unmarshal(b, &m.data); err != nil {
        return fmt.Errorf("%s for %s", err.Error(), string(b))
    }

    if m.data.Status != nil {
        if err := json.Unmarshal(m.data.Status, &m.status); err != nil {
            return fmt.Errorf("Unable to process Measurement Status: %s", err.Error())
        }
    }

    if m.data.Probes != nil {
        if err := json.Unmarshal(m.data.Probes, &m.probes); err != nil {
            return fmt.Errorf("Unable to process Measurement Probes: %s", err.Error())
        }
    }

    if m.data.ProbeSources != nil {
        if err := json.Unmarshal(m.data.ProbeSources, &m.probeSources); err != nil {
            return fmt.Errorf("Unable to process Measurement ProbeSources: %s", err.Error())
        }
    }

    if m.data.ParticipationRequests != nil {
        if err := json.Unmarshal(m.data.ParticipationRequests, &m.participationRequests); err != nil {
            return fmt.Errorf("Unable to process Measurement ParticipationRequests: %s", err.Error())
        }
    }

    return nil
}

// [core] .
func (m *Measurement) Id() int {
    return m.data.Id
}

// [core] The URL that contains the results of this measurement.
func (m *Measurement) Result() string {
    return m.data.Result
}

// [core] The ID of the measurement group. This ID references a measurement acting as group master.
func (m *Measurement) GroupId() int {
    return m.data.GroupId
}

// [core] [Not for wifi] IPv4 of IPv6 Address family of the measurement.
// [wifi] IPv4 of IPv6 Address family of the measurement.
func (m *Measurement) Af() int {
    return m.data.Af
}

// [core] Flag indicating this is a one-off measurement.
func (m *Measurement) IsOneoff() bool {
    return m.data.IsOneoff
}

// [core] Flag indicating this measurement is a publicly available.
func (m *Measurement) IsPublic() bool {
    return m.data.IsPublic
}

// [core] User-defined description of the measurement.
func (m *Measurement) Description() string {
    return m.data.Description
}

// [core] Distribution of probes' measurements throughout the interval (default is half the interval, maximum 400 seconds).
func (m *Measurement) Spread() int {
    return m.data.Spread
}

// [core] Flag that, when set to true, indicates that a name should be resolved (using DNS) on the probe. Otherwise it will be resolved on the RIPE Atlas servers.
func (m *Measurement) ResolveOnProbe() bool {
    return m.data.ResolveOnProbe
}

// [core] Configured start time (as a unix timestamp).
func (m *Measurement) StartTime() int {
    return m.data.StartTime
}

// [core] Actual end time of measurement (as a unix timestamp).
func (m *Measurement) StopTime() int {
    return m.data.StopTime
}

// [core] Returns the type of the measurement.
func (m *Measurement) Type() string {
    return m.data.Type
}

// [core] Returns a JSON object containing `id` and `name` (0: Specified, 1: Scheduled, 2: Ongoing, 4: Stopped, 5: Forced to stop, 6: No suitable probes, 7: Failed, 8: Archived).
func (m *Measurement) Status() *measurement.MeasurementStatus {
    return m.status
}

// [core] Returns true if all probe requests have made it through the scheduling process..
func (m *Measurement) IsAllScheduled() bool {
    return m.data.IsAllScheduled
}

// [core] Number of participating probes.
func (m *Measurement) ParticipantCount() int {
    return m.data.ParticipantCount
}

// [core] The ASN the IP the target is in.
func (m *Measurement) TargetAsn() int {
    return m.data.TargetAsn
}

// [core] The IP Address of the target of the measurement.
func (m *Measurement) TargetIp() string {
    return m.data.TargetIp
}

// [core] The creation date and time of the measurement (Defaults to unix timestamp format).
func (m *Measurement) CreationTime() int {
    return m.data.CreationTime
}

// [core] Flag indicating this measurement belongs to a wifi measurement group.
func (m *Measurement) InWifiGroup() bool {
    return m.data.InWifiGroup
}

// [core] The list of IP addresses returned for the fqdn in the `target` field by the backend infra-structure resolvers.
func (m *Measurement) ResolvedIps() []string {
    return m.data.ResolvedIps
}

// [core] Number of probes requested, but not necessarily granted to this measurement.
func (m *Measurement) ProbesRequested() int {
    return m.data.ProbesRequested
}

// [core] Number of probes actually scheduled for this measurement.
func (m *Measurement) ProbesScheduled() int {
    return m.data.ProbesScheduled
}

// [core] The API URL of the measurement group..
func (m *Measurement) Group() string {
    return m.data.Group
}

// [core] probes involved in this measurement.
func (m *Measurement) Probes() []*measurement.Probe {
    return m.probes
}

// [core] .
func (m *Measurement) ProbeSources() []*measurement.ProbeSource {
    return m.probeSources
}

// [core] .
func (m *Measurement) ParticipationRequests() []*measurement.ProbeSource {
    return m.participationRequests
}

// [dns] Interval between samples from a single probe. Defaults to 240 seconds..
// [http] Interval between samples from a single probe. Defaults to 1800 seconds..
// [ntp] Interval between samples from a single probe. Defaults to 1800 seconds..
// [ping] Interval between samples from a single probe. Defaults to 240 seconds..
// [sslcert] Interval between samples from a single probe. Defaults to 900 seconds..
// [traceroute] Interval between samples from a single probe. Defaults to 900 seconds..
// [wifi] Interval between samples from a single probe. Defaults to 900 seconds.
func (m *Measurement) Interval() int {
    return m.data.Interval
}

// [ntp] The number of packets send in a measurement execution. Value must be between 1 and 16. Default is 3.
// [ping] The number of packets send in a measurement execution. Value must be between 1 and 16. Default is 3.
// [traceroute] The number of packets send in a measurement execution. Value must be between 1 and 16. Default is 3.
func (m *Measurement) Packets() int {
    return m.data.Packets
}

// [ping] size of the data part of the packet, i.e. excluding any IP and ICMP headers. Value must be between 1 and 2048.
// [traceroute] size of the data part of the packet, i.e. excluding any IP, ICMP, UDP or TCP headers. Value must be between 0 and 2048.
func (m *Measurement) Size() int {
    return m.data.Size
}

// [ping] Time between packets in milliseconds. Value must be between 2 and 300000.
func (m *Measurement) PacketInterval() int {
    return m.data.PacketInterval
}

// [ping] Include the probe ID (encoded as ASCII digits) as part of the payload.
func (m *Measurement) IncludeProbeId() bool {
    return m.data.IncludeProbeId
}

// [http] The target port number Defaults to 80.
// [sslcert] The target port number. Defaults to 443.
// [traceroute] The target port number (TCP only). Defaults to 80.
func (m *Measurement) Port() int {
    return m.data.Port
}

// [traceroute] TTL (time to live) of the first hop.
func (m *Measurement) FirstHop() int {
    return m.data.FirstHop
}

// [traceroute] Traceroute measurement stops after the hop at which the TTL reaches this value.
func (m *Measurement) MaxHops() int {
    return m.data.MaxHops
}

// [traceroute] The number of paris traceroute variations to try. Zero disables paris traceroute. Value must be between 0 and 64.
func (m *Measurement) Paris() int {
    return m.data.Paris
}

// [dns] Protocol used in measurement. Defaults to UDP.
// [traceroute] Protocol used in measurement.
func (m *Measurement) Protocol() string {
    return m.data.Protocol
}

// [traceroute] Response timeout for one packet.
func (m *Measurement) ResponseTimeout() int {
    return m.data.ResponseTimeout
}

// [traceroute] Time to wait (in milliseconds) for a duplicate response after receiving the first response.
func (m *Measurement) DuplicateTimeout() int {
    return m.data.DuplicateTimeout
}

// [traceroute] Size of an IPv6 hop-by-hop option header filled with NOPs.
func (m *Measurement) HopByHopOptionSize() int {
    return m.data.HopByHopOptionSize
}

// [traceroute] Size of an IPv6 destination option header filled with NOPs.
func (m *Measurement) DestinationOptionSize() int {
    return m.data.DestinationOptionSize
}

// [traceroute] Do not fragment outgoing packets.
func (m *Measurement) DontFragment() bool {
    return m.data.DontFragment
}

// [dns] Set the DNS0 option for UDP payload size to this value, between 512 and 4096.Defaults to 512).
func (m *Measurement) UdpPayloadSize() int {
    return m.data.UdpPayloadSize
}

// [dns] Send the DNS query to the probe's local resolvers (instead of an explicitly specified target).
func (m *Measurement) UseProbeResolver() bool {
    return m.data.UseProbeResolver
}

// [dns] Flag indicating Recursion Desired bit was set.
func (m *Measurement) SetRdBit() bool {
    return m.data.SetRdBit
}

// [dns] Each probe prepends its probe number and a timestamp to the DNS query argument to make it unique.
func (m *Measurement) PrependProbeId() bool {
    return m.data.PrependProbeId
}

// [dns] Number of times to retry.
func (m *Measurement) Retry() int {
    return m.data.Retry
}

// [dns] include the raw DNS query data in the result. Defaults to false.
func (m *Measurement) IncludeQbuf() bool {
    return m.data.IncludeQbuf
}

// [dns] Flag indicating Name Server Identifier (RFC5001) was set.
func (m *Measurement) SetNsidBit() bool {
    return m.data.SetNsidBit
}

// [dns] include the raw DNS answer data in the result. Defaults to true.
func (m *Measurement) IncludeAbuf() bool {
    return m.data.IncludeAbuf
}

// [dns] The `class` part of the query used in the measurement.
func (m *Measurement) QueryClass() string {
    return m.data.QueryClass
}

// [dns] The `argument` part of the query used in the measurement.
func (m *Measurement) QueryArgument() string {
    return m.data.QueryArgument
}

// [dns] The `type` part of the query used in the measurement.
func (m *Measurement) QueryType() string {
    return m.data.QueryType
}

// [dns] Flag indicating DNSSEC Checking Disabled (RFC4035) was set.
func (m *Measurement) SetCdBit() bool {
    return m.data.SetCdBit
}

// [dns] Flag indicating DNSSEC OK (RFC3225) was set.
func (m *Measurement) SetDoBit() bool {
    return m.data.SetDoBit
}

// [dns] Allow the use of $p (probe ID), $r (random 16-digit hex string) and $t (timestamp) in the query_argument.
func (m *Measurement) UseMacros() bool {
    return m.data.UseMacros
}

// [dns] Timeout in milliseconds (default: 5000).
// [ntp] Per packet timeout in milliseconds.
func (m *Measurement) Timeout() int {
    return m.data.Timeout
}

// [http] Enable time-to-resolve, time-to-connect and time-to-first-byte measurements.
func (m *Measurement) ExtendedTiming() bool {
    return m.data.ExtendedTiming
}

// [http] Include fields added by extended_timing and adds readtiming which reports for each read system call when it happened and how much data was delivered.
func (m *Measurement) MoreExtendedTiming() bool {
    return m.data.MoreExtendedTiming
}

// [http] Maximum number of bytes in the reponse header, defaults to 0.
func (m *Measurement) HeaderBytes() int {
    return m.data.HeaderBytes
}

// [http] http verb of the measurement request.
func (m *Measurement) Method() string {
    return m.data.Method
}

// [http] Path of the requested URL.
func (m *Measurement) Path() string {
    return m.data.Path
}

// [http] Optional query parameters of the requested URL.
func (m *Measurement) QueryString() string {
    return m.data.QueryString
}

// [http] user agent header field sent in the http request. Always set to 'RIPE Atlas: https//atlas.ripe.net'.
func (m *Measurement) UserAgent() string {
    return m.data.UserAgent
}

// [http] .
func (m *Measurement) MaxBytesRead() int {
    return m.data.MaxBytesRead
}

// [http] http version of measurement request.
func (m *Measurement) Version() string {
    return m.data.Version
}

// [sslcert] Server Name Indication (SNI) hostname.
func (m *Measurement) Hostname() string {
    return m.data.Hostname
}

// [wifi] Flag indicating IPv4 measurements are attempted in this group.
func (m *Measurement) Ipv4() bool {
    return m.data.Ipv4
}

// [wifi] Flag indicating IPv6 measurements are attempted in this group.
func (m *Measurement) Ipv6() bool {
    return m.data.Ipv6
}

// [wifi] Certificate in PEM format.
func (m *Measurement) Cert() string {
    return m.data.Cert
}

// [wifi] Wait this amount of time before executing measurement commands..
func (m *Measurement) ExtraWait() int {
    return m.data.ExtraWait
}

// [wifi] Wifi SSID to connect to. Max. 32 characters.
func (m *Measurement) Ssid() string {
    return m.data.Ssid
}

// [wifi] Authentication mechanism used for the wifi connection. For WPA-PSK `psk` field is also required,for WPA-EAP `eap` and `password` fields are required.
func (m *Measurement) KeyMgmt() string {
    return m.data.KeyMgmt
}

// [wifi] Extensible Authentication Protocol type. Currently only `TTLS` is available.
func (m *Measurement) Eap() string {
    return m.data.Eap
}

// [wifi] Username used for wifi connection. Used for both outer and inner connection if anonymous_identity is omitted.
func (m *Measurement) Identity() string {
    return m.data.Identity
}

// [wifi] Username used for outer connection. If omitted the `identity` field is used for the outer connection.
func (m *Measurement) AnonymousIdentity() string {
    return m.data.AnonymousIdentity
}

// [wifi] Connection and Authentication directives for the inner connection. Only used for WPA-EAP. Currently only EAP-MSCHAPv2 is available.
func (m *Measurement) Phase2() string {
    return m.data.Phase2
}

// [wifi] Flag indicating that BSSID radio signal strength will be measured and stored.
func (m *Measurement) Rssi() bool {
    return m.data.Rssi
}
