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

// Error message.
type Error struct {
    data struct {
        Timeout     int    `json:"timeout"`
        Getaddrinfo string `json:"getaddrinfo"`
    }
}

func (e *Error) UnmarshalJSON(b []byte) error {
    if err := json.Unmarshal(b, &e.data); err != nil {
        return fmt.Errorf("%s for %s", err.Error(), string(b))
    }
    return nil
}

// Query timeout.
func (e *Error) Timeout() int {
    return e.data.Timeout
}

// Error message.
func (e *Error) Getaddrinfo() string {
    return e.data.Getaddrinfo
}
