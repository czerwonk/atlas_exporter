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

package measurement

import (
    "encoding/json"
    "fmt"
)

type MeasurementStatus struct {
    ParseError error

    data struct {
        Id   int    `json:"id"`
        Name string `json:"name"`
    }
}

func (m *MeasurementStatus) UnmarshalJSON(b []byte) error {
    if err := json.Unmarshal(b, &m.data); err != nil {
        return fmt.Errorf("%s for %s", err.Error(), string(b))
    }

    return nil
}

// Numeric ID of this status.
func (m *MeasurementStatus) Id() int {
    return m.data.Id
}

// Human-readable description of this status.
func (m *MeasurementStatus) Name() string {
    return m.data.Name
}
