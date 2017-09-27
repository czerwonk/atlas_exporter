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

// Traceroute result.
type Result struct {
    data struct {
        Hop    int             `json:"hop"`
        Error  string          `json:"error"`
        Result json.RawMessage `json:"result"`
    }

    replies []*Reply
}

func (r *Result) UnmarshalJSON(b []byte) error {
    if err := json.Unmarshal(b, &r.data); err != nil {
        return fmt.Errorf("%s for %s", err.Error(), string(b))
    }

    if r.data.Result != nil {
        if err := json.Unmarshal(r.data.Result, &r.replies); err != nil {
            return fmt.Errorf("Unable to process Replies: %s", err.Error())
        }
    }

    return nil
}

// Hop number.
func (r *Result) Hop() int {
    return r.data.Hop
}

// When an error occurs trying to send a packet. In that case there will
// not be a result structure (optional).
func (r *Result) Error() string {
    return r.data.Error
}

// Traceroute replies (called "result" in RIPE Atlas API documentation).
func (r *Result) Replies() []*Reply {
    return r.replies
}
