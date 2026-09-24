package main

import "testing"

func TestResolveConfigFilePath(t *testing.T) {
	tests := []struct {
		name     string
		flag     string
		env      string
		expected string
	}{
		{name: "neither set", expected: ""},
		{name: "flag only", flag: "/etc/flag.yml", expected: "/etc/flag.yml"},
		{name: "env only", env: "/etc/env.yml", expected: "/etc/env.yml"},
		{name: "flag takes precedence over env", flag: "/etc/flag.yml", env: "/etc/env.yml", expected: "/etc/flag.yml"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Setenv("CONFIG", test.env)

			actual := resolveConfigFilePath(test.flag)
			if actual != test.expected {
				t.Errorf("resolveConfigFilePath(%q) with CONFIG=%q: expected %q, got %q", test.flag, test.env, test.expected, actual)
			}
		})
	}
}
