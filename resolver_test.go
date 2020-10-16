package dnsbl

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReverseIP(t *testing.T) {

	act := reverseIP("1.2.3.4")
	require.Equal(t, "4.3.2.1", act)

}

func TestIsValidIP(t *testing.T) {

	testcases := []struct {
		name  string
		ip    string
		valid bool
	}{
		{"should be a valid ip", "1.2.3.4", true},
		{"should be invalid ip", "1.2.x.4", false},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			act := isValidIP4(tc.ip)
			require.Equal(t, tc.valid, act)
		})
	}
}
