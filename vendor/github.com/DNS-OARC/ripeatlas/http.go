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

package ripeatlas

import (
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "net/url"
    "strings"

    "github.com/DNS-OARC/ripeatlas/measurement"
)

// A Http reads RIPE Atlas data from the RIPE Atlas REST API.
type Http struct {
}

const (
    MeasurementsUrl = "https://atlas.ripe.net/api/v2/measurements"
)

// NewHttp returns a new Atlaser for reading from the RIPE Atlas REST API.
func NewHttp() *Http {
    return &Http{}
}

func (h *Http) get(url string, fragmented bool) (<-chan *measurement.Result, error) {
    r, err := http.Get(url)
    if err != nil {
        return nil, fmt.Errorf("http.Get(%s): %s", url, err.Error())
    }

    ch := make(chan *measurement.Result)
    go func() {
        d := json.NewDecoder(r.Body)
        defer r.Body.Close()

        if fragmented {
            for {
                var r measurement.Result
                if err := d.Decode(&r); err == io.EOF {
                    break
                } else if err != nil {
                    r.ParseError = fmt.Errorf("json.Decode(%s): %s", url, err.Error())
                    ch <- &r
                    break
                }
                ch <- &r
            }
        } else {
            var r []*measurement.Result
            if err := d.Decode(&r); err != nil {
                if err != io.EOF {
                    r := &measurement.Result{ParseError: fmt.Errorf("json.Decode(%s): %s", url, err.Error())}
                    ch <- r
                }
            } else {
                for _, i := range r {
                    ch <- i
                }
            }
        }
        close(ch)
    }()

    return ch, nil
}

// MeasurementLatest gets the latest measurement results, as described
// by the Params, and sends them to the returned channel.
//
// Params available are:
//
// "pk": string - The measurement id to read results from (required).
//
// "fragmented": bool - If true, use the fragmented/stream format when reading.
func (h *Http) MeasurementLatest(p Params) (<-chan *measurement.Result, error) {
    var pk string
    var fragmented bool

    for k, v := range p {
        switch k {
        case "pk":
            v, ok := v.(string)
            if !ok {
                return nil, fmt.Errorf("Invalid %s parameter, must be string", k)
            }
            pk = v
        case "fragmented":
            v, ok := v.(bool)
            if !ok {
                return nil, fmt.Errorf("Invalid %s parameter, must be bool", k)
            }
            fragmented = v
        default:
            return nil, fmt.Errorf("Invalid parameter %s", k)
        }
    }

    if pk == "" {
        return nil, fmt.Errorf("Required parameter pk missing")
    }

    url := fmt.Sprintf("%s/%s/latest", MeasurementsUrl, url.PathEscape(pk))
    if fragmented {
        url += "?format=txt"
    } else {
        url += "?format=json"
    }

    return h.get(url, fragmented)
}

// MeasurementResults gets the measurement results, as described by the Params,
// and sends them to the returned channel.
//
// Params available are:
//
// "pk": string - The measurement id to read results from (required).
//
// "start": int64 - Get the results starting at the given UNIX timestamp.
//
// "stop": int64 - Get the results up to the given UNIX timestamp.
//
// "probe_ids": none - Unimplemented
//
// "anchors-only": none - Unimplemented
//
// "public-only": none - Unimplemented
//
// "fragmented": bool - If true, use the fragmented/stream format when reading.
func (h *Http) MeasurementResults(p Params) (<-chan *measurement.Result, error) {
    var qstr []string
    var pk string
    var fragmented bool

    for k, v := range p {
        switch k {
        case "pk":
            v, ok := v.(string)
            if !ok {
                return nil, fmt.Errorf("Invalid %s parameter, must be string", k)
            }
            pk = v
        case "start":
            fallthrough
        case "stop":
            v, ok := v.(int64)
            if !ok {
                return nil, fmt.Errorf("Invalid %s parameter, must be int64", k)
            }
            qstr = append(qstr, fmt.Sprintf("%s=%d", k, v))
        case "probe_ids":
            fallthrough
        case "anchors-only":
            fallthrough
        case "public-only":
            return nil, fmt.Errorf("Unimplemented parameter %s", k)
        case "fragmented":
            v, ok := v.(bool)
            if !ok {
                return nil, fmt.Errorf("Invalid %s parameter, must be bool", k)
            }
            fragmented = v
        default:
            return nil, fmt.Errorf("Invalid parameter %s", k)
        }
    }

    if pk == "" {
        return nil, fmt.Errorf("Required parameter pk missing")
    }

    url := fmt.Sprintf("%s/%s/results", MeasurementsUrl, url.PathEscape(pk))
    if fragmented {
        url += "?format=txt"
    } else {
        url += "?format=json"
    }
    if len(qstr) > 0 {
        url += "&" + strings.Join(qstr, "&")
    }

    return h.get(url, fragmented)
}
