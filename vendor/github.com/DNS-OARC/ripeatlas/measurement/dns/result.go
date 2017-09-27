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

package dns

import (
    "encoding/base64"
    "encoding/json"
    "fmt"

    mdns "github.com/miekg/dns"
)

// Response from the DNS server.
type Result struct {
    data struct {
        Ancount int             `json:"ANCOUNT"`
        Arcount int             `json:"ARCOUNT"`
        Id      int             `json:"ID"`
        Nscount int             `json:"NSCOUNT"`
        Qdcount int             `json:"QDCOUNT"`
        Abuf    string          `json:"abuf"`
        Answers json.RawMessage `json:"answers"`
        Rt      float64         `json:"rt"`
        Size    int             `json:"size"`
        SrcAddr string          `json:"src_addr"`
        Subid   int             `json:"subid"`
        Submax  int             `json:"submax"`
    }

    answers []*Answer
}

func (r *Result) UnmarshalJSON(b []byte) error {
    if err := json.Unmarshal(b, &r.data); err != nil {
        return fmt.Errorf("%s for %s", err.Error(), string(b))
    }

    if r.data.Answers != nil {
        if err := json.Unmarshal(r.data.Answers, &r.answers); err != nil {
            return fmt.Errorf("Unable to process DNS answers: %s", err.Error())
        }
    }

    return nil
}

// Answer count.
func (r *Result) Ancount() int {
    return r.data.Ancount
}

// Additional record count.
func (r *Result) Arcount() int {
    return r.data.Arcount
}

// Query ID.
func (r *Result) Id() int {
    return r.data.Id
}

// Name server count.
func (r *Result) Nscount() int {
    return r.data.Nscount
}

// Number of queries.
func (r *Result) Qdcount() int {
    return r.data.Qdcount
}

// Answer payload buffer from the server, UU encoded.
func (r *Result) Abuf() string {
    return r.data.Abuf
}

// First two records from the response decoded by the probe, if they are
// TXT or SOA; other RR can be decoded from Abuf() using UnpackAbuf().
func (r *Result) Answers() []*Answer {
    return r.answers
}

// Response time in milli seconds (optional).
func (r *Result) Rt() float64 {
    return r.data.Rt
}

// Response size (optional).
func (r *Result) Size() int {
    return r.data.Size
}

// The source IP address added by the probe (optional).
func (r *Result) SrcAddr() string {
    return r.data.SrcAddr
}

// Sequence number of this result within a group of results, available
// if the resolution is done by the probe's local resolver (optional).
func (r *Result) Subid() int {
    return r.data.Subid
}

// Total number of results within a group (optional).
func (r *Result) Submax() int {
    return r.data.Submax
}

// Decode the Abuf(), returns a *Msg from the github.com/miekg/dns package
// or nil on error or if Abuf() is empty.
func (r *Result) UnpackAbuf() (*mdns.Msg, error) {
    if r.data.Abuf == "" {
        return nil, nil
    }

    b, err := base64.StdEncoding.DecodeString(r.data.Abuf)
    if err != nil {
        return nil, err
    }

    m := &mdns.Msg{}
    if err := m.Unpack(b); err != nil {
        return nil, err
    }

    return m, nil
}
