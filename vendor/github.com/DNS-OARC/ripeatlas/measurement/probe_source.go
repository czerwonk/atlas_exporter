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

type ProbeSource struct {
    ParseError error

    data struct {
        Requested   int    `json:"requested"`
        Type        string `json:"type"`
        Value       string `json:"value"`
        TagsInclude string `json:"tags_include"`
        TagsExclude string `json:"tags_exclude"`
    }
}

func (p *ProbeSource) UnmarshalJSON(b []byte) error {
    if err := json.Unmarshal(b, &p.data); err != nil {
        return fmt.Errorf("%s for %s", err.Error(), string(b))
    }

    return nil
}

// Number of probes that have to be added or removed.
func (p *ProbeSource) Requested() int {
    return p.data.Requested
}

// Probe selector. Options are: `area` allows a compass quarter of the world, `asn` selects an Autonomous System, `country` selects a country, `msm` selects the probes used in another measurement, `prefix` selects probes based on prefix, `probes` selects probes directly.
func (p *ProbeSource) Type() string {
    return p.data.Type
}

// Value for given selector type. `area`: ['WW','West','North-Central','South-Central','North-East','South-East']. `asn`: ASN (integer). `country`: two-letter country code according to ISO 3166-1 alpha-2, e.g. GR. `msm`: measurement id (integer). `prefix`: prefix in CIDR notation, e.g. 193.0.0/16. `probes`: comma-separated list of probe IDs.
func (p *ProbeSource) Value() string {
    return p.data.Value
}

// Comma-separated list of probe tags. Only probes with all these tags attached will be selected from this participation request.
func (p *ProbeSource) TagsInclude() string {
    return p.data.TagsInclude
}

// Comma-separated list of probe tags. Probes with any of these tags attached will be excluded from this participation request.
func (p *ProbeSource) TagsExclude() string {
    return p.data.TagsExclude
}
