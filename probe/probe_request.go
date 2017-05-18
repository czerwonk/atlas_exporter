package probe

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

const url = "https://atlas.ripe.net/api/v2/probes/"

func Get(id int) (*Probe, error) {
	c := &http.Client{}
	u := fmt.Sprintf("%s%d", url, id)

	resp, err := c.Get(u)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	return FromJson(body)
}
