package formatter

import (
	"testing"

	"code/internal/diff"

	"github.com/stretchr/testify/require"
)

func TestFormatPlainValue(t *testing.T) {
	tests := []struct {
		name     string
		value    interface{}
		expected string
	}{
		{name: "nil", value: nil, expected: "null"},
		{name: "map", value: map[string]interface{}{"a": 1.0}, expected: "[complex value]"},
		{name: "string", value: "text", expected: "'text'"},
		{name: "bool true", value: true, expected: "true"},
		{name: "bool false", value: false, expected: "false"},
		{name: "int float64", value: float64(5), expected: "5"},
		{name: "fraction float64", value: 3.14, expected: "3.14"},
		{name: "slice treated as complex", value: []interface{}{"a", "b"}, expected: "[complex value]"},
		{name: "default type", value: 123, expected: "123"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			result := formatPlainValue(tt.value)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestBuildPath(t *testing.T) {
	tests := []struct {
		name       string
		parentPath string
		key        string
		expected   string
	}{
		{name: "root level", parentPath: "", key: "key", expected: "key"},
		{name: "nested level", parentPath: "parent.child", key: "key", expected: "parent.child.key"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			result := buildPath(tt.parentPath, tt.key)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestPlainFormatter_Format(t *testing.T) {
	tests := []struct {
		name              string
		nodes             []*diff.Node
		expectContains    []string
		expectNotContains []string
		expectExact       string
	}{
		{
			name: "non empty",
			nodes: []*diff.Node{
				{Type: diff.NodeTypeAdded, Key: "key1", Value: "value1"},
				{Type: diff.NodeTypeRemoved, Key: "key2", Value: 42.0},
				{Type: diff.NodeTypeChanged, Key: "key3", OldValue: "old", NewValue: "new"},
			},
			expectContains: []string{
				"Property 'key1' was added with value: 'value1'",
				"Property 'key2' was removed",
				"Property 'key3' was updated. From 'old' to 'new'",
			},
		},
		{
			name: "nested and complex values",
			nodes: []*diff.Node{
				{
					Type: diff.NodeTypeNested,
					Key:  "group",
					Children: []*diff.Node{
						{
							Type:  diff.NodeTypeAdded,
							Key:   "child",
							Value: map[string]interface{}{"a": 1},
						},
						{
							Type:     diff.NodeTypeChanged,
							Key:      "value",
							OldValue: nil,
							NewValue: false,
						},
					},
				},
				{
					Type: diff.NodeTypeUnchanged,
					Key:  "ignored",
					Value: map[string]interface{}{
						"nested": true,
					},
				},
			},
			expectContains: []string{
				"Property 'group.child' was added with value: [complex value]",
				"Property 'group.value' was updated. From null to false",
			},
			expectNotContains: []string{"ignored"},
		},
		{
			name:        "empty nodes",
			expectExact: "",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			f := &PlainFormatter{}
			result, err := f.Format(tt.nodes)

			require.NoError(t, err)
			if tt.expectExact != "" || tt.nodes == nil {
				require.Equal(t, tt.expectExact, result)
				return
			}

			for _, substr := range tt.expectContains {
				require.Contains(t, result, substr)
			}

			for _, substr := range tt.expectNotContains {
				require.NotContains(t, result, substr)
			}
		})
	}
}
