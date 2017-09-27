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

package ping

import (
    "encoding/json"
    "fmt"
)

// Ping result.
type Result struct {
    data struct {
        X       string  `json:"x"`
        Error   string  `json:"error"`
        Rtt     float64 `json:"rtt"`
        SrcAddr string  `json:"src_Addr"`
        Ttl     int     `json:"ttl"`
        Dup     int     `json:"dup"`
    }
}

func (r *Result) UnmarshalJSON(b []byte) error {
    if err := json.Unmarshal(b, &r.data); err != nil {
        return fmt.Errorf("%s for %s", err.Error(), string(b))
    }
    return nil
}

// On timeout: "*".
func (r *Result) X() string {
    return r.data.X
}

// On error: description of error.
func (r *Result) Error() string {
    return r.data.Error
}

// On reply: round-trip-time in milliseconds.
func (r *Result) Rtt() float64 {
    return r.data.Rtt
}

// On reply: source address if different from the source address in first
// reply (optional).
func (r *Result) SrcAddr() string {
    return r.data.SrcAddr
}

// On reply: time-to-live reply if different from ttl in first reply
// (optional).
func (r *Result) Ttl() int {
    return r.data.Ttl
}

// On reply: signals that the reply is a duplicate (optional).
func (r *Result) Dup() int {
    return r.data.Dup
}
