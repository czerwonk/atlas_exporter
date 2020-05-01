package probe

import (
	"encoding/json"
	"strconv"
)

const IPv6 int = 6

// Probe holds information about a single Atlas probe
type Probe struct {
	ID          int    `json:"id"`
	Asn4        int    `json:"asn_v4"`
	Asn6        int    `json:"asn_v6"`
	CountryCode string `json:"country_code"`
	Geometry    struct {
		Coordinates []float64 `json:"coordinates"`
	} `json:"geometry"`
}

// FromJSON parses json and returns a probe
func FromJSON(body []byte) (*Probe, error) {
	var p Probe
	err := json.Unmarshal(body, &p)

	if err != nil {
		return nil, err
	}

	return &p, nil
}

// ASNForIPVersion return the ASN for the given IP Version
func (p *Probe) ASNForIPVersion(v int) int {
	if v == IPv6 {
		return p.Asn6
	}

	return p.Asn4
}

// Longitude of the geo location of the probe
func (p *Probe) Longitude() string {
	if len(p.Geometry.Coordinates) == 0 {
		return ""
	}

	return strconv.FormatFloat(p.Geometry.Coordinates[0], 'f', 4, 64)
}

// Latitude of the geo location of the probe
func (p *Probe) Latitude() string {
	if len(p.Geometry.Coordinates) == 0 {
		return ""
	}

	return strconv.FormatFloat(p.Geometry.Coordinates[1], 'f', 4, 64)
}
