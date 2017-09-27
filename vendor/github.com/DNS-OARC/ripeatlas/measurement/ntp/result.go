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

package ntp

import (
    "encoding/json"
    "fmt"
)

// NTP result.
type Result struct {
    data struct {
        FinalTs    float64 `json:"final-ts"`
        Offset     float64 `json:"offset"`
        OriginTs   float64 `json:"origin-ts"`
        ReceiveTs  float64 `json:"receive-ts"`
        Rtt        float64 `json:"rtt"`
        TransmitTs float64 `json:"transmit-ts"`
    }
}

func (r *Result) UnmarshalJSON(b []byte) error {
    if err := json.Unmarshal(b, &r.data); err != nil {
        return fmt.Errorf("%s for %s", err.Error(), string(b))
    }
    return nil
}

// NTP time the response from the server is received.
func (r *Result) FinalTs() float64 {
    return r.data.FinalTs
}

// Clock offset between client and server in seconds.
func (r *Result) Offset() float64 {
    return r.data.Offset
}

// NTP time the request was sent.
func (r *Result) OriginTs() float64 {
    return r.data.OriginTs
}

// NTP time the server received the request.
func (r *Result) ReceiveTs() float64 {
    return r.data.ReceiveTs
}

// Round trip time between client and server in seconds.
func (r *Result) Rtt() float64 {
    return r.data.Rtt
}

// NTP time the server sent the response.
func (r *Result) TransmitTs() float64 {
    return r.data.TransmitTs
}
