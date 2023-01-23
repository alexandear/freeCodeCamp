package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/goleak"
)

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)
}

func TestOpenPorts(t *testing.T) {
	t.Parallel()
	for name, tc := range map[string]struct {
		target    string
		portStart int
		portEnd   int

		wantPorts   []int
		wantVerbose string
	}{
		"ip": {
			target:      "209.216.230.240",
			portStart:   440,
			portEnd:     445,
			wantPorts:   []int{443},
			wantVerbose: "",
		},
		"url": {
			target:      "www.stackoverflow.com",
			portStart:   79,
			portEnd:     82,
			wantPorts:   []int{80},
			wantVerbose: "",
		},
		"url multiple ports": {
			target:      "scanme.nmap.org",
			portStart:   20,
			portEnd:     80,
			wantPorts:   []int{22, 80},
			wantVerbose: "",
		},
		"verbose ip no hostname returned single port": {
			target:      "104.26.10.78",
			portStart:   440,
			portEnd:     450,
			wantPorts:   []int{443},
			wantVerbose: "Open ports for 104.26.10.78\nPORT     SERVICE\n443      https",
		},
		"verbose ip hostname returned single port": {
			target:      "137.74.187.104",
			portStart:   440,
			portEnd:     450,
			wantPorts:   []int{443},
			wantVerbose: "Open ports for hackthissite.org (137.74.187.104)\nPORT     SERVICE\n443      https",
		},
		"verbose hostname returned multiple ports": {
			target:      "scanme.nmap.org",
			portStart:   20,
			portEnd:     80,
			wantPorts:   []int{22, 80},
			wantVerbose: "Open ports for scanme.nmap.org (45.33.32.156)\nPORT     SERVICE\n22       ssh\n80       http",
		},
	} {
		t.Run(name, func(t *testing.T) {
			got, err := OpenPorts(tc.target, tc.portStart, tc.portEnd)

			assert.NoError(t, err)
			assert.Equal(t, tc.wantPorts, got.Ports())
			if tc.wantVerbose != "" {
				assert.Equal(t, tc.wantVerbose, got.Verbose())
			}
		})
	}

	t.Run("invalid hostname", func(t *testing.T) {
		_, err := OpenPorts("scanme.nmap", 22, 42)

		assert.EqualError(t, err, "invalid hostname")
	})

	t.Run("invalid ip", func(t *testing.T) {
		_, err := OpenPorts("266.255.9.10", 22, 42)

		assert.EqualError(t, err, "invalid IP address")
	})
}
