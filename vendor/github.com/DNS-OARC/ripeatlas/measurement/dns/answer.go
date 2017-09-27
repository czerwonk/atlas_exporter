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

// First two records from the response decoded by the probe, if they are
// TXT or SOA; other RR can be decoded from "abuf"
type Answer struct {
    data struct {
        Mname  string      `json:"MNAME"`
        Name   string      `json:"NAME"`
        Rdata  interface{} `json:"RDATA"`
        Rname  string      `json:"RNAME"`
        Serial int         `json:"SERIAL"`
        Ttl    int         `json:"TTL"`
        Type   string      `json:"TYPE"`
    }

    rdata []string
}

func (a *Answer) UnmarshalJSON(b []byte) error {
    if err := json.Unmarshal(b, &a.data); err != nil {
        return fmt.Errorf("%s for %s", err.Error(), string(b))
    }

    switch a.data.Rdata.(type) {
    case string:
        a.rdata = []string{a.data.Rdata.(string)}
    case []interface{}:
        for _, i := range a.data.Rdata.([]interface{}) {
            switch i.(type) {
            case string:
                a.rdata = append(a.rdata, i.(string))
            default:
                return fmt.Errorf("Element within RDATA field unsupported type %T", i)
            }
        }
    case nil:
        return nil
    default:
        return fmt.Errorf("RDATA field unsupported type %T", a.data.Rdata)
    }

    return nil
}

// If the type is "SOA" this will have the original or primary domain name.
func (a *Answer) Mname() string {
    return a.data.Mname
}

// Domain name.
func (a *Answer) Name() string {
    return a.data.Name
}

// If the type is "TXT" this will have the value of that record.
func (a *Answer) Rdata() []string {
    return a.rdata
}

// If the type is "SOA" this will have the mailbox.
func (a *Answer) Rname() string {
    return a.data.Rname
}

// If the type is "SOA" this will have the zone serial number.
func (a *Answer) Serial() int {
    return a.data.Serial
}

// If the type is "SOA" this will have the time to live.
func (a *Answer) Ttl() int {
    return a.data.Ttl
}

// Resource record type ("SOA" or "TXT").
func (a *Answer) Type() string {
    return a.data.Type
}
