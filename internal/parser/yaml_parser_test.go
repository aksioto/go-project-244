package parser

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type yamlParserCase struct {
	name      string
	data      []byte
	expected  map[string]interface{}
	expectErr bool
}

func TestYAMLParser_Parse(t *testing.T) {
	parser := &YAMLParser{}

	tests := []yamlParserCase{
		{
			name: "valid yaml with normalization",
			data: []byte(`
foo: 1
bar:
  baz: true
nested:
  value: 3
list:
  - 2
  - name: item
deep_list:
  - - 4
    - sublist:
        value: 5
`),
			expected: map[string]interface{}{
				"foo": float64(1),
				"bar": map[string]interface{}{"baz": true},
				"nested": map[string]interface{}{
					"value": float64(3),
				},
				"list": []interface{}{
					float64(2),
					map[string]interface{}{"name": "item"},
				},
				"deep_list": []interface{}{
					[]interface{}{
						float64(4),
						map[string]interface{}{
							"sublist": map[string]interface{}{
								"value": float64(5),
							},
						},
					},
				},
			},
		},
		{
			name:      "invalid yaml",
			data:      []byte(":invalid"),
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

func TestNormalizeValuesAndSlice(t *testing.T) {
	m := map[string]interface{}{
		"int":   1,
		"int64": int64(2),
		"map": map[string]interface{}{
			"nestedInt": 3,
		},
		"slice": []interface{}{
			4,
			int64(5),
			map[string]interface{}{
				"deep": 6,
			},
			[]interface{}{7, int64(8)},
		},
	}

	normalizeValues(m)

	require.Equal(t, float64(1), m["int"])
	require.Equal(t, float64(2), m["int64"])

	nestedMap, ok := m["map"].(map[string]interface{})
	require.True(t, ok)
	require.Equal(t, float64(3), nestedMap["nestedInt"])

	slice, ok := m["slice"].([]interface{})
	require.True(t, ok)
	require.Equal(t, float64(4), slice[0])
	require.Equal(t, float64(5), slice[1])

	deepMap, ok := slice[2].(map[string]interface{})
	require.True(t, ok)
	require.Equal(t, float64(6), deepMap["deep"])

	nestedSlice, ok := slice[3].([]interface{})
	require.True(t, ok)
	require.Equal(t, float64(7), nestedSlice[0])
	require.Equal(t, float64(8), nestedSlice[1])
}
