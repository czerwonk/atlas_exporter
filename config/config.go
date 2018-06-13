package config

import (
	"fmt"
	"io"
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

// Config represents the configuration for the exporter
type Config struct {
	// Measurements is the ids of measurements used as source for metrics generation
	Measurements []string `yaml:"measurements"`
}

// Load loads a config from a reader
func Load(r io.Reader) (*Config, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("could not load config: %v", err)
	}

	c := &Config{}
	if len(b) == 0 {
		return c, nil
	}

	err = yaml.Unmarshal(b, c)
	if err != nil {
		return nil, fmt.Errorf("could not parse config: %v", err)
	}

	return c, err
}
