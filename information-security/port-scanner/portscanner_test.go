package main

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/goleak"
)

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)
}

func TestOpenPorts(t *testing.T) {
	tests := []struct {
		name      string
		target    string
		portStart int
		portEnd   int

		want        []int
		wantVerbose string
		wantErr     error
	}{
		{
			name:        "ip",
			target:      "209.216.230.240",
			portStart:   440,
			portEnd:     445,
			want:        []int{443},
			wantVerbose: "",
			wantErr:     nil,
		},
		{
			name:        "url",
			target:      "www.stackoverflow.com",
			portStart:   79,
			portEnd:     82,
			want:        []int{80},
			wantVerbose: "",
			wantErr:     nil,
		},
		{
			name:        "url multiple ports",
			target:      "scanme.nmap.org",
			portStart:   20,
			portEnd:     80,
			want:        []int{22, 80},
			wantVerbose: "",
			wantErr:     nil,
		},
		{
			name:        "verbose ip no hostname returned single port",
			target:      "104.26.10.78",
			portStart:   440,
			portEnd:     450,
			want:        []int{443},
			wantVerbose: "Open ports for 104.26.10.78\nPORT     SERVICE\n443      https",
			wantErr:     nil,
		},
		{
			name:        "verbose ip hostname returned single port",
			target:      "137.74.187.104",
			portStart:   440,
			portEnd:     450,
			want:        []int{443},
			wantVerbose: "Open ports for hackthissite.org (137.74.187.104)\nPORT     SERVICE\n443      https",
			wantErr:     nil,
		},
		{
			name:        "verbose hostname returned multiple ports",
			target:      "scanme.nmap.org",
			portStart:   20,
			portEnd:     80,
			want:        []int{22, 80},
			wantVerbose: "Open ports for scanme.nmap.org (45.33.32.156)\nPORT     SERVICE\n22       ssh\n80       http",
			wantErr:     nil,
		},
		{
			name:        "invalid hostname",
			target:      "scanme.nmap",
			portStart:   22,
			portEnd:     42,
			want:        nil,
			wantVerbose: "",
			wantErr:     errors.New("invalid hostname"),
		},
		{
			name:        "invalid ip address",
			target:      "266.255.9.10",
			portStart:   22,
			portEnd:     42,
			want:        nil,
			wantVerbose: "",
			wantErr:     errors.New("invalid IP address"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := OpenPorts(tt.target, tt.portStart, tt.portEnd)
			if tt.wantErr != nil {
				assert.EqualError(t, err, tt.wantErr.Error())
			} else {
				assert.NoError(t, err)
			}
			if tt.want != nil {
				assert.Equal(t, tt.want, got.Ports())
			}
			if tt.wantVerbose != "" {
				assert.Equal(t, tt.wantVerbose, got.Verbose())
			}
		})
	}
}
