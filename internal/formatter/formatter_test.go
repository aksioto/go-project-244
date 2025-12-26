package formatter

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGetFormatter(t *testing.T) {
	tests := []struct {
		name      string
		format    string
		expectErr bool
	}{
		{"default", "", false},
		{"stylish", "stylish", false},
		{"plain", "plain", false},
		{"invalid", "xml", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, err := GetFormatter(tt.format)
			if tt.expectErr {
				require.Error(t, err)
				require.Nil(t, f)
			} else {
				require.NoError(t, err)
				require.NotNil(t, f)
			}
		})
	}
}
