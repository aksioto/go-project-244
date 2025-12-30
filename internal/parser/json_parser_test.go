package parser

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type jsonParserCase struct {
	name      string
	data      []byte
	expected  map[string]interface{}
	expectErr bool
}

func TestJSONParser_Parse(t *testing.T) {
	parser := &JSONParser{}

	tests := []jsonParserCase{
		{
			name: "valid nested json",
			data: []byte(`{"foo": 1, "bar": {"baz": true}}`),
			expected: map[string]interface{}{
				"foo": float64(1),
				"bar": map[string]interface{}{"baz": true},
			},
		},
		{
			name:      "invalid json",
			data:      []byte("{invalid"),
			expectErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			result, err := parser.Parse(tt.data)
			if tt.expectErr {
				require.Error(t, err)
				require.Nil(t, result)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.expected, result)
		})
	}
}
