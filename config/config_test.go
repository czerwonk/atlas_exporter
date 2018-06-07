package config

import (
	"strings"
	"testing"

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
			name:     "empty config",
			value:    ``,
			expected: Config{},
		},
		{
			name:  "valid config",
			value: `measurements: [ "123", "456" ]`,
			expected: Config{
				Measurements: []string{
					"123", "456",
				},
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
