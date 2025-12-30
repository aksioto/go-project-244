package formatter

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

type getFormatterCase struct {
	name      string
	format    string
	expectErr bool
}

func TestGetFormatter(t *testing.T) {
	tests := []getFormatterCase{
		{name: "default", format: "", expectErr: false},
		{name: "stylish", format: "stylish", expectErr: false},
		{name: "plain", format: "plain", expectErr: false},
		{name: "json", format: "json", expectErr: false},
		{name: "invalid", format: "xml", expectErr: true},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			f, err := GetFormatter(tt.format)
			if tt.expectErr {
				require.Error(t, err)
				require.Nil(t, f)
				require.True(t, errors.Is(err, ErrUnknownFormat))
				return
			}

			require.NoError(t, err)
			require.NotNil(t, f)
		})
	}
}
