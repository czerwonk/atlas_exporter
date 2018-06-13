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

package probe

import (
    "encoding/json"
    "fmt"
)

type ProbeStatus struct {
    ParseError error

    data struct {
        Status string `json:"status"`
        Id     int    `json:"id"`
        Name   string `json:"name"`
    }
}

func (p *ProbeStatus) UnmarshalJSON(b []byte) error {
    if err := json.Unmarshal(b, &p.data); err != nil {
        return fmt.Errorf("%s for %s", err.Error(), string(b))
    }

    return nil
}

// .
func (p *ProbeStatus) Status() string {
    return p.data.Status
}

// .
func (p *ProbeStatus) Id() int {
    return p.data.Id
}

// .
func (p *ProbeStatus) Name() string {
    return p.data.Name
}
