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

package request

import (
    "encoding/json"
    "fmt"

    "github.com/DNS-OARC/ripeatlas/request/probe"
)

type Probe struct {
    ParseError error

    data struct {
        AddressV4      string          `json:"address_v4"`
        AddressV6      string          `json:"address_v6"`
        AsnV4          int             `json:"asn_v4"`
        AsnV6          int             `json:"asn_v6"`
        CountryCode    string          `json:"country_code"`
        Description    string          `json:"description"`
        FirstConnected int             `json:"first_connected"`
        Geometry       json.RawMessage `json:"geometry"`
        Id             int             `json:"id"`
        IsAnchor       bool            `json:"is_anchor"`
        IsPublic       bool            `json:"is_public"`
        LastConnected  int             `json:"last_connected"`
        PrefixV4       string          `json:"prefix_v4"`
        PrefixV6       string          `json:"prefix_v6"`
        Status         json.RawMessage `json:"status"`
        StatusSince    int             `json:"status_since"`
        Tags           json.RawMessage `json:"tags"`
        TotalUptime    int             `json:"total_uptime"`
        Type           string          `json:"type"`
    }

    geometry *probe.Number

    status *probe.ProbeStatus

    tags []*probe.ProbeTags
}

func (p *Probe) UnmarshalJSON(b []byte) error {
    if err := json.Unmarshal(b, &p.data); err != nil {
        return fmt.Errorf("%s for %s", err.Error(), string(b))
    }

    if p.data.Geometry != nil {
        if err := json.Unmarshal(p.data.Geometry, &p.geometry); err != nil {
            return fmt.Errorf("Unable to process Probe Geometry: %s", err.Error())
        }
    }

    if p.data.Status != nil {
        if err := json.Unmarshal(p.data.Status, &p.status); err != nil {
            return fmt.Errorf("Unable to process Probe Status: %s", err.Error())
        }
    }

    if p.data.Tags != nil {
        if err := json.Unmarshal(p.data.Tags, &p.tags); err != nil {
            return fmt.Errorf("Unable to process Probe Tags: %s", err.Error())
        }
    }

    return nil
}

// The last IPv4 address that was known to be held by this probe, or null if there is no known address. Note: a probe that connects over IPv6 may fail to report its IPv4 address, meaning that this field can sometimes be null even though the probe may have working IPv4.
func (p *Probe) AddressV4() string {
    return p.data.AddressV4
}

// The last IPv6 address that was known to be held by this probe, or null if there is no known address..
func (p *Probe) AddressV6() string {
    return p.data.AddressV6
}

// The IPv4 ASN if any.
func (p *Probe) AsnV4() int {
    return p.data.AsnV4
}

// The IPv6 ASN if any.
func (p *Probe) AsnV6() int {
    return p.data.AsnV6
}

// An ISO-3166-1 alpha-2 code indicating the country that this probe is located in, as derived from the user supplied longitude and latitude.
func (p *Probe) CountryCode() string {
    return p.data.CountryCode
}

// User defined description of the probe.
func (p *Probe) Description() string {
    return p.data.Description
}

// When the probe connected for the first time (UTC Time and date in ISO-8601/ECMA 262 format).
func (p *Probe) FirstConnected() int {
    return p.data.FirstConnected
}

// A GeoJSON point object containing the user-supplied location of this probe. The longitude and latitude are contained within the `coordinates` array.
func (p *Probe) Geometry() *probe.Number {
    return p.geometry
}

// The id of the probe.
func (p *Probe) Id() int {
    return p.data.Id
}

// Whether or not this probe is a RIPE Atlas Anchor.
func (p *Probe) IsAnchor() bool {
    return p.data.IsAnchor
}

// If a probe is not public then certain details, including exact IP addresses, are not returned..
func (p *Probe) IsPublic() bool {
    return p.data.IsPublic
}

// When the probe connected for the last time (UTC Time and date in ISO-8601/ECMA 262 format).
func (p *Probe) LastConnected() int {
    return p.data.LastConnected
}

// The IPv4 prefix if any.
func (p *Probe) PrefixV4() string {
    return p.data.PrefixV4
}

// The IPv6 prefix if any.
func (p *Probe) PrefixV6() string {
    return p.data.PrefixV6
}

// A JSON object containing id: The connection status ID for this probe (integer [0-3]), name: The connection status (string [Never Connected, Connected, Disconnected, Abandoned]), since: The datetime of the last change in connection status.
func (p *Probe) Status() *probe.ProbeStatus {
    return p.status
}

// A datetime field that can hold a datetime both as a timestamp and as a JSON datetime.
func (p *Probe) StatusSince() int {
    return p.data.StatusSince
}

// .
func (p *Probe) Tags() []*probe.ProbeTags {
    return p.tags
}

// Accumulated uptime for this probe in seconds.
func (p *Probe) TotalUptime() int {
    return p.data.TotalUptime
}

// The type of the object.
func (p *Probe) Type() string {
    return p.data.Type
}
