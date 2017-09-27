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

package sslcert

import (
    "encoding/json"
    "fmt"
)

var level = map[int]string{
    1: "warning",
    2: "fatal",
}

var description = map[int]string{
    0:   "close_notify",
    10:  "unexpected_message",
    20:  "bad_record_mac",
    21:  "decryption_failed_RESERVED",
    22:  "record_overflow",
    30:  "decompression_failure",
    40:  "handshake_failure",
    41:  "no_certificate_RESERVED",
    42:  "bad_certificate",
    43:  "unsupported_certificate",
    44:  "certificate_revoked",
    45:  "certificate_expired",
    46:  "certificate_unknown",
    47:  "illegal_parameter",
    48:  "unknown_ca",
    49:  "access_denied",
    50:  "decode_error",
    51:  "decrypt_error",
    60:  "export_restriction_RESERVED",
    70:  "protocol_version",
    71:  "insufficient_security",
    80:  "internal_error",
    90:  "user_canceled",
    100: "no_renegotiation",
    110: "unsupported_extension",
}

// SSL Certificate error alert (see RFC 5246, Section 7.2).
type Alert struct {
    data struct {
        Level       int `json:"level"`
        Description int `json:"description"`
    }
}

func (r *Alert) UnmarshalJSON(b []byte) error {
    if err := json.Unmarshal(b, &r.data); err != nil {
        return fmt.Errorf("%s for %s", err.Error(), string(b))
    }
    return nil
}

// AlertLevel.
func (r *Alert) Level() int {
    return r.data.Level
}

func (r *Alert) LevelString() (string, error) {
    s, ok := level[r.data.Level]
    if !ok {
        return "", fmt.Errorf("Unknown level %d", r.data.Level)
    }
    return s, nil
}

// AlertDescription.
func (r *Alert) Description() int {
    return r.data.Description
}

func (r *Alert) DescriptionString() (string, error) {
    s, ok := description[r.data.Description]
    if !ok {
        return "", fmt.Errorf("Unknown description %d", r.data.Description)
    }
    return s, nil
}
