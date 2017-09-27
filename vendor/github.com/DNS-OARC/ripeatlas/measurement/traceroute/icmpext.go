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

package traceroute

import (
    "encoding/json"
    "fmt"
)

// Traceroute ICMP extension reply
type Icmpext struct {
    data struct {
        Version int `json:"version"`
        Rfc4884 int `json:"rfc4884"`
        /* TODO: Even RIPE NCC have not implemented this:
           "obj" -- elements of the object (array)
               "class" -- RFC4884 class (int)
               "type" -- RFC4884 type (int)
               "mpls" -- [optional] MPLS data, RFC4950, shown when class is "1" and type is "1" (array)
                   "exp" -- for experimental use (int)
                   "label" -- mpls label (int)
                   "s" -- bottom of stack (int)
                   "ttl" -- time to live value (int)
        */
        Objects []interface{} `json:"obj"`
    }
}

func (i *Icmpext) UnmarshalJSON(b []byte) error {
    if err := json.Unmarshal(b, &i.data); err != nil {
        return fmt.Errorf("%s for %s", err.Error(), string(b))
    }
    return nil
}

// RFC4884 version.
func (i *Icmpext) Version() int {
    return i.data.Version
}

// "1" if length indication is present, "0" otherwise.
func (i *Icmpext) Rfc4884() int {
    return i.data.Rfc4884
}

// ICMP extension objects.
func (i *Icmpext) Objects() []interface{} {
    return i.data.Objects
}
