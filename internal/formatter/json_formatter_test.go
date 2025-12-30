package formatter

import (
	"encoding/json"
	"testing"

	"code/internal/diff"

	"github.com/stretchr/testify/require"
)

func TestJSONFormatter_Format(t *testing.T) {
	tests := []struct {
		name        string
		nodes       []*diff.Node
		expectErr   bool
		errContains string
		assertFunc  func(t *testing.T, root *jsonNode)
	}{
		{
			name: "non empty",
			nodes: []*diff.Node{
				{Type: diff.NodeTypeAdded, Key: "key1", Value: "value1"},
				{
					Type: diff.NodeTypeNested,
					Key:  "parent",
					Children: []*diff.Node{
						{Type: diff.NodeTypeChanged, Key: "child", OldValue: "old", NewValue: "new"},
					},
				},
			},
			assertFunc: func(t *testing.T, root *jsonNode) {
				require.Equal(t, jsonTypeRoot, root.Type)
				require.Len(t, root.Children, 2)

				added := root.Children[0]
				require.Equal(t, "key1", added.Key)
				require.Equal(t, jsonTypeAdded, added.Type)
				require.Nil(t, added.Children)
				require.Nil(t, added.Value1)
				require.Equal(t, "value1", added.Value2)

				parent := root.Children[1]
				require.Equal(t, jsonTypeNested, parent.Type)
				require.Len(t, parent.Children, 1)
				child := parent.Children[0]
				require.Equal(t, jsonTypeChanged, child.Type)
				require.Equal(t, "old", child.Value1)
				require.Equal(t, "new", child.Value2)
			},
		},
		{
			name: "empty nodes",
			assertFunc: func(t *testing.T, root *jsonNode) {
				require.Equal(t, jsonTypeRoot, root.Type)
				require.Empty(t, root.Children)
			},
		},
		{
			name: "marshal error",
			nodes: []*diff.Node{
				{Type: diff.NodeTypeAdded, Key: "invalid", Value: make(chan struct{})},
			},
			expectErr:   true,
			errContains: "marshal diff to json",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			f := &JSONFormatter{}
			result, err := f.Format(tt.nodes)

			if tt.expectErr {
				require.Empty(t, result)
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.errContains)
				return
			}

			require.NoError(t, err)

			var root jsonNode
			require.NoError(t, json.Unmarshal([]byte(result), &root))

			require.NotNil(t, tt.assertFunc, "assertFunc must be provided for non-error cases")
			tt.assertFunc(t, &root)
		})
	}
}
