package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCrackSHA1Password(t *testing.T) {
	for name, tc := range map[string]struct {
		passwordHash string
		ifUseSalts   bool
		expected     string
	}{
		"hash1": {
			passwordHash: "18c28604dd31094a8d69dae60f1bcd347f1afc5a",
			ifUseSalts:   false,
			expected:     "superman",
		},
		"hash2": {
			passwordHash: "5d70c3d101efd9cc0a69f4df2ddf33b21e641f6a",
			ifUseSalts:   false,
			expected:     "q1w2e3r4t5",
		},
		"hash3": {
			passwordHash: "b80abc2feeb1e37c66477b0824ac046f9e2e84a0",
			ifUseSalts:   false,
			expected:     "bubbles1",
		},
		"hash4": {
			passwordHash: "80540a46a2c1a0eae58d9868f01c32bdcec9a010",
			ifUseSalts:   false,
			expected:     "01071988",
		},
		"hash_salted_1": {
			passwordHash: "53d8b3dc9d39f0184144674e310185e41a87ffd5",
			ifUseSalts:   true,
			expected:     "superman",
		},
		"hash_salted_2": {
			passwordHash: "da5a4e8cf89539e66097acd2f8af128acae2f8ae",
			ifUseSalts:   true,
			expected:     "q1w2e3r4t5",
		},
		"hash_salted_3": {
			passwordHash: "ea3f62d498e3b98557f9f9cd0d905028b3b019e1",
			ifUseSalts:   true,
			expected:     "bubbles1",
		},
		"hash_salted_4": {
			passwordHash: "05bbf26a28148f531cf57872df546961d1ed0861",
			ifUseSalts:   true,
			expected:     "01071988",
		},
		"not_in_database": {
			passwordHash: "03810a46a2c1a0eae58d9332f01c32bdcec9a01a",
			ifUseSalts:   false,
			expected:     "PASSWORD NOT IN DATABASE",
		},
	} {
		t.Run(name, func(t *testing.T) {
			actual := CrackSHA1Hash(tc.passwordHash, tc.ifUseSalts)
			assert.Equal(t, tc.expected, actual)
		})
	}
}
