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

// HTTP timing result.
type Readtiming struct {
    data struct {
        T float64 `json:"t"`
        O int     `json:"o,string"` // TODO: Not documented as string
    }
}

func (r *Readtiming) UnmarshalJSON(b []byte) error {
    if err := json.Unmarshal(b, &r.data); err != nil {
        return fmt.Errorf("%s for %s", err.Error(), string(b))
    }
    return nil
}

// Time since starting to connect when data is received (in milli seconds).
func (r *Readtiming) T() float64 {
    return r.data.T
}

// Offset in stream of reply data.
func (r *Readtiming) O() int {
    return r.data.O
}
