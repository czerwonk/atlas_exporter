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

/*
Package ripeatlas implements bindings for RIPE Atlas.

The Atlaser is the interface to access RIPE Atlas and there are a few
different ways to do so, for example read measurement results from a
JSON file:
    a := ripeatlas.Atlaser(ripeatlas.NewFile())
    c, err := a.MeasurementResults(ripeatlas.Params{"file": name})
    if err != nil {
        ...
    }
    for r := range c {
        if r.ParseError != nil {
            ...
        }
        fmt.Printf("%d %s\n", r.MsmId(), r.Type())
    }

See File for file access, Http for REST API access and Stream for Streaming
API access.
*/
package ripeatlas

import (
    "github.com/DNS-OARC/ripeatlas/measurement"
)

// Params is used to give parameters to the different access methods.
type Params map[string]interface{}

// Atlaser is the interface for accessing RIPE Atlas, designed after
// the REST API (https://atlas.ripe.net/docs/api/v2/reference/).
type Atlaser interface {
    MeasurementLatest(p Params) (<-chan *measurement.Result, error)
    MeasurementResults(p Params) (<-chan *measurement.Result, error)
}
