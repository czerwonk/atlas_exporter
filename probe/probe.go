package probe

import (
	"encoding/json"
)

// Probe holds information about a single Atlas probe
type Probe struct {
	Id   int `json:"id"`
	Asn4 int `json:"asn_v4"`
	Asn6 int `json:"asn_v6"`
}

// FromJson parses json and returns a probe
func FromJson(body []byte) (*Probe, error) {
	var p Probe
	err := json.Unmarshal(body, &p)

	if err != nil {
		return nil, err
	}

	return &p, nil
}

// ASNForIPVersion return the ASN for the given IP Version
func (p *Probe) ASNForIPVersion(v int) int {
	if v == 6 {
		return p.Asn6
	}

	return p.Asn4
}
