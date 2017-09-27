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
    "encoding/json"
    "fmt"
)

type Resultset struct {
    data struct {
        Af        int             `json:"af"`
        DstAddr   string          `json:"dst_addr"`
        DstName   string          `json:"dst_name"`
        Error     json.RawMessage `json:"error"`
        Lts       int             `json:"lts"`
        Proto     string          `json:"proto"`
        Qbuf      string          `json:"qbuf"`
        Result    json.RawMessage `json:"result"`
        Retry     int             `json:"retry"`
        Timestamp int             `json:"timestamp"`
    }

    error  *Error
    result *Result
}

func (r *Resultset) UnmarshalJSON(b []byte) error {
    if err := json.Unmarshal(b, &r.data); err != nil {
        return fmt.Errorf("%s for %s", err.Error(), string(b))
    }

    if r.data.Error != nil {
        r.error = &Error{}
        if err := json.Unmarshal(r.data.Error, r.error); err != nil {
            return fmt.Errorf("Unable to process DNS error: %s", err.Error())
        }
    }
    if r.data.Result != nil {
        r.result = &Result{}
        if err := json.Unmarshal(r.data.Result, r.result); err != nil {
            return fmt.Errorf("Unable to process DNS result: %s", err.Error())
        }
    }

    return nil
}

// IP version: "4" or "6" (optional).
func (r *Resultset) Af() int {
    return r.data.Af
}

// IP address of the destination (optional).
func (r *Resultset) DstAddr() string {
    return r.data.DstAddr
}

// Hostname of the destination (optional).
func (r *Resultset) DstName() string {
    return r.data.DstName
}

// DNS error message, nil if not present.
func (r *Resultset) DnsError() *Error {
    return r.error
}

// Last time synchronised. How long ago (in seconds) the clock of the probe
// was found to be in sync with that of a controller. The value -1 is used
// to indicate that the probe does not know whether it is in sync.
func (r *Resultset) Lts() int {
    return r.data.Lts
}

// Protocol, "TCP" or "UDP".
func (r *Resultset) Proto() string {
    return r.data.Proto
}

// Query payload buffer which was sent to the server, UU encoded (optional).
func (r *Resultset) Qbuf() string {
    return r.data.Qbuf
}

// DNS response from the DNS server, nil if not present.
func (r *Resultset) Result() *Result {
    return r.result
}

// Retry count (optional).
func (r *Resultset) Retry() int {
    return r.data.Retry
}

// Start time, in Unix timestamp.
func (r *Resultset) Timestamp() int {
    return r.data.Timestamp
}
