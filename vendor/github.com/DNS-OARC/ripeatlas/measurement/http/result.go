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

package http

import (
    "encoding/json"
    "fmt"
)

// HTTP result.
type Result struct {
    data struct {
        Af         int             `json:"af"`
        Bsize      int             `json:"bsize"`
        Dnserr     string          `json:"dnserr"`
        DstAddr    string          `json:"dst_addr"`
        Err        string          `json:"err"`
        Header     []string        `json:"header"`
        Hsize      int             `json:"hsize"`
        Method     string          `json:"method"`
        Readtiming json.RawMessage `json:"readtiming"`
        Res        int             `json:"res"`
        Rt         float64         `json:"rt"`
        SrcAddr    string          `json:"src_addr"`
        Subid      int             `json:"subid"`
        Submax     int             `json:"submax"`
        Time       int             `json:"time"`
        Ttc        float64         `json:"ttc"`
        Ttfb       float64         `json:"ttfb"`
        Ttr        float64         `json:"ttr"`
        Ver        string          `json:"ver"`
    }

    readtimings []*Readtiming
}

func (r *Result) UnmarshalJSON(b []byte) error {
    if err := json.Unmarshal(b, &r.data); err != nil {
        return fmt.Errorf("%s for %s", err.Error(), string(b))
    }

    if r.data.Readtiming != nil {
        if err := json.Unmarshal(r.data.Readtiming, &r.readtimings); err != nil {
            return fmt.Errorf("Unable to process Readtiming: %s", err.Error())
        }
    }

    return nil
}

// Address family, 4 or 6.
func (r *Result) Af() int {
    return r.data.Af
}

// Size of body in octets.
func (r *Result) Bsize() int {
    return r.data.Bsize
}

// DNS resolution failed (optional).
func (r *Result) Dnserr() string {
    return r.data.Dnserr
}

// Target address.
func (r *Result) DstAddr() string {
    return r.data.DstAddr
}

// Other failure (optional).
func (r *Result) Err() string {
    return r.data.Err
}

// The last string can be empty to indicate the end of enders or end
// with "[...]" to indicate truncation (optional).
func (r *Result) Header() []string {
    return r.data.Header
}

// Header size in octets.
func (r *Result) Hsize() int {
    return r.data.Hsize
}

// "GET", "HEAD", or "POST".
func (r *Result) Method() string {
    return r.data.Method
}

// Timing results for reply data (optional).
func (r *Result) Readtimings() []*Readtiming {
    return r.readtimings
}

// HTTP result code.
func (r *Result) Res() int {
    return r.data.Res
}

// Time to execute request excluding DNS.
func (r *Result) Rt() float64 {
    return r.data.Rt
}

// Source address used by probe.
func (r *Result) SrcAddr() string {
    return r.data.SrcAddr
}

// Sequence number of this result within a group of results, when
// the 'all' option is used without the 'combine' option (optional).
func (r *Result) Subid() int {
    return r.data.Subid
}

// Total number of results within a group (optional).
func (r *Result) Submax() int {
    return r.data.Submax
}

// Unix timestamp, when the 'all' option is used with the 'combine'
// option (optional).
func (r *Result) Time() int {
    return r.data.Time
}

// Time to connect to the target (in milli seconds) (optional).
func (r *Result) Ttc() float64 {
    return r.data.Ttc
}

// Time to first response byte received by measurent code after starting
// to connect (in milli seconds) (optional).
func (r *Result) Ttfb() float64 {
    return r.data.Ttfb
}

// Time to resolve the DNS name (in milli seconds) (optional).
func (r *Result) Ttr() float64 {
    return r.data.Ttr
}

// Major, minor version of http server.
func (r *Result) Ver() string {
    return r.data.Ver
}
