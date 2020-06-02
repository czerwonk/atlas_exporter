package config

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	tests := []struct {
		name      string
		value     string
		expected  Config
		wantsFail bool
	}{
		{
			name:  "empty config",
			value: ``,
			expected: Config{
				FilterInvalidResults: true,
			},
		},
		{
			name: "valid config",
			value: `
measurements:
  - id: 123
  - id: 456`,
			expected: Config{
				Measurements: []Measurement{
					{ID: "123"},
					{ID: "456"},
				},
				FilterInvalidResults: true,
			},
		},
		{
			name: "valid config with brackets",
			value: `
measurements:
  - id: 123
histogram_buckets:
  dns:
    rtt: [ 1.0, 2.0 ]
  http: 
    rtt: [ 3.0, 4.0 ]
  ping: 
    rtt: [ 5.0, 6.0 ]
  traceroute: 
    rtt: [ 7.0, 8.0 ]`,
			expected: Config{
				Measurements: []Measurement{
					{ID: "123"},
				},
				HistogramBuckets: HistogramBuckets{
					DNS: RttHistogramBucket{
						Rtt: []float64{1, 2},
					},
					HTTP: RttHistogramBucket{
						Rtt: []float64{3, 4},
					},
					Ping: RttHistogramBucket{
						Rtt: []float64{5, 6},
					},
					Traceroute: RttHistogramBucket{
						Rtt: []float64{7, 8},
					},
				},
				FilterInvalidResults: true,
			},
		},
		{
			name:      "invalid config",
			value:     `measurements: { 123, 456 }`,
			wantsFail: true,
		},
		{
			name: "valid config with timeout",
			value: `
measurements:
  - id: 123
    timeout: 30s`,
			expected: Config{
				Measurements: []Measurement{
					{ID: "123", Timeout: 30 * time.Second},
				},
				FilterInvalidResults: true,
			},
		},
		{
			name: "valid config with filter override",
			value: `
filter_invalid_results: false`,
			expected: Config{
				FilterInvalidResults: false,
			},
		},
		{
			name:      "invalid config",
			value:     `measurements: { 123, 456 }`,
			wantsFail: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(te *testing.T) {
			r := strings.NewReader(test.value)
			c, err := Load(r)
			if err != nil {
				if !test.wantsFail {
					te.Fatalf("unecpected error: %v", err)
				}

				return
			}

			assert.Equal(te, test.expected, *c)
		})
	}
}
