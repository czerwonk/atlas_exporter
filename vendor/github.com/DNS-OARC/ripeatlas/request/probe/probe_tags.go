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

type ProbeTags struct {
    ParseError error

    data struct {
        Name string `json:"name"`
        Slug string `json:"slug"`
    }
}

func (p *ProbeTags) UnmarshalJSON(b []byte) error {
    if err := json.Unmarshal(b, &p.data); err != nil {
        return fmt.Errorf("%s for %s", err.Error(), string(b))
    }

    return nil
}

// tagname.
func (p *ProbeTags) Name() string {
    return p.data.Name
}

// tag as a slug, the tagname in lowercase and hyphenated (for use in URLs).
func (p *ProbeTags) Slug() string {
    return p.data.Slug
}
