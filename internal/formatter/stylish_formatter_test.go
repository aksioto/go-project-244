package formatter

import (
	"strings"
	"testing"

	"code/internal/diff"

	"github.com/stretchr/testify/require"
)

func TestFormatSimpleValue(t *testing.T) {
	tests := []struct {
		name     string
		value    interface{}
		expected string
	}{
		{name: "nil", value: nil, expected: "null"},
		{name: "string", value: "text", expected: "text"},
		{name: "bool true", value: true, expected: "true"},
		{name: "bool false", value: false, expected: "false"},
		{name: "int float64", value: float64(10), expected: "10"},
		{name: "fraction float64", value: 3.14, expected: "3.14"},
		{name: "default", value: []int{1, 2}, expected: "[1 2]"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			result := formatSimpleValue(tt.value)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestMakeIndent(t *testing.T) {
	tests := []struct {
		name     string
		depth    int
		marker   string
		expected string
	}{
		{name: "depth one plus", depth: 1, marker: "+", expected: "  + "},
		{name: "zero depth space", depth: 0, marker: " ", expected: "  "},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.expected, makeIndent(tt.depth, tt.marker))
		})
	}
}

func TestMakeIndentForMap(t *testing.T) {
	tests := []struct {
		name     string
		depth    int
		expected string
	}{
		{name: "depth zero", depth: 0, expected: ""},
		{name: "depth three", depth: 3, expected: "            "},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.expected, makeIndentForMap(tt.depth))
		})
	}
}

func TestStylishFormatter_Format(t *testing.T) {
	tests := []struct {
		name           string
		nodes          []*diff.Node
		expectContains []string
		expectExact    string
	}{
		{
			name: "non empty",
			nodes: []*diff.Node{
				{Type: diff.NodeTypeAdded, Key: "key1", Value: "value1"},
				{Type: diff.NodeTypeRemoved, Key: "key2", Value: 42.0},
				{Type: diff.NodeTypeUnchanged, Key: "key3", Value: true},
			},
			expectContains: []string{"+ key1: value1", "- key2: 42", "  key3: true"},
		},
		{
			name:        "empty nodes",
			expectExact: "{\n}",
		},
		{
			name: "complex nested",
			nodes: []*diff.Node{
				{
					Type: diff.NodeTypeAdded,
					Key:  "object",
					Value: map[string]interface{}{
						"inner":  map[string]interface{}{"foo": "bar"},
						"number": float64(3),
					},
				},
				{
					Type:     diff.NodeTypeChanged,
					Key:      "count",
					OldValue: float64(1),
					NewValue: 2.0,
				},
				{
					Type: diff.NodeTypeNested,
					Key:  "group",
					Children: []*diff.Node{
						{Type: diff.NodeTypeRemoved, Key: "old", Value: false},
						{Type: diff.NodeTypeAdded, Key: "new", Value: "value"},
					},
				},
			},
			expectExact: strings.Join([]string{
				"{",
				"  + object: {",
				"        inner: {",
				"            foo: bar",
				"        }",
				"        number: 3",
				"    }",
				"  - count: 1",
				"  + count: 2",
				"    group: {",
				"      - old: false",
				"      + new: value",
				"    }",
				"}",
			}, "\n"),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			f := &StylishFormatter{}
			result, err := f.Format(tt.nodes)

			require.NoError(t, err)
			if tt.expectExact != "" {
				require.Equal(t, tt.expectExact, result)
				return
			}

			for _, substr := range tt.expectContains {
				require.Contains(t, result, substr)
			}
		})
	}
}
