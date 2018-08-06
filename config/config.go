package config

import (
	"fmt"
	"io"
	"io/ioutil"
	"time"

	yaml "gopkg.in/yaml.v2"
)

// Config represents the configuration for the exporter
type Config struct {
	// Measurements is the ids of measurements used as source for metrics generation
	Measurements      []Measurement    `yaml:"measurements"`
	HistogramBrackets HistogramBuckets `yaml:"histogram_buckets"`
}

// Measurement represents config options for one measurement
type Measurement struct {
	ID      string        `id`
	Timeout time.Duration `timeout`
}

// HistogramBrackets represents histogram brackets for different measurement types
type HistogramBuckets struct {
	DNS        []float64 `yaml:"dns,omitempty"`
	HTTP       []float64 `yaml:"http,omitempty"`
	Ping       []float64 `yaml:"ping,omitempty"`
	Traceroute []float64 `yaml:"traceroute,omitempty"`
}

func (c *Config) MeasurementIDs() []string {
	ids := make([]string, len(c.Measurements))
	for i, m := range c.Measurements {
		ids[i] = m.ID
	}

	return ids
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
