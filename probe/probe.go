package probe

import (
	"encoding/json"
)

type Probe struct {
	Id   int `json:"id"`
	Asn4 int `json:"asn_v4"`
	Asn6 int `json:"asn_v6"`
}

func FromJson(body []byte) (*Probe, error) {
	var p Probe
	err := json.Unmarshal(body, &p)

	if err != nil {
		return nil, err
	}

	return &p, nil
}
