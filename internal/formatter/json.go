package formatter

import (
	"encoding/json"
	"fmt"

	"code/internal/diff"
)

const (
	jsonTypeRoot      = "root"
	jsonTypeAdded     = "added"
	jsonTypeDeleted   = "deleted"
	jsonTypeChanged   = "changed"
	jsonTypeUnchanged = "unchanged"
	jsonTypeNested    = "nested"
)

type JSONFormatter struct{}

type jsonNode struct {
	Key      string      `json:"key"`
	Type     string      `json:"type"`
	Value1   interface{} `json:"value1,omitempty"`
	Value2   interface{} `json:"value2,omitempty"`
	Children []*jsonNode `json:"children,omitempty"`
}

func (f *JSONFormatter) Format(nodes []*diff.Node) (string, error) {
	root := &jsonNode{
		Key:      "",
		Type:     jsonTypeRoot,
		Children: f.buildJSONNodes(nodes),
	}

	data, err := json.MarshalIndent(root, "", "  ")
	if err != nil {
		return "", fmt.Errorf("marshal diff to json: %w", err)
	}

	return string(data), nil
}

func (f *JSONFormatter) buildJSONNodes(nodes []*diff.Node) []*jsonNode {
	if len(nodes) == 0 {
		return []*jsonNode{}
	}

	res := make([]*jsonNode, 0, len(nodes))
	for _, node := range nodes {
		if node == nil {
			continue
		}
		res = append(res, f.convertNode(node))
	}

	return res
}

func (f *JSONFormatter) convertNode(node *diff.Node) *jsonNode {
	jNode := &jsonNode{
		Key: node.Key,
	}

	switch node.Type {
	case diff.NodeTypeAdded:
		jNode.Type = jsonTypeAdded
		jNode.Value2 = node.Value
	case diff.NodeTypeRemoved:
		jNode.Type = jsonTypeDeleted
		jNode.Value1 = node.Value
	case diff.NodeTypeChanged:
		jNode.Type = jsonTypeChanged
		jNode.Value1 = node.OldValue
		jNode.Value2 = node.NewValue
	case diff.NodeTypeNested:
		jNode.Type = jsonTypeNested
		jNode.Children = f.buildJSONNodes(node.Children)
	case diff.NodeTypeUnchanged:
		jNode.Type = jsonTypeUnchanged
		jNode.Value1 = node.Value
	default:
		jNode.Type = string(node.Type)
		jNode.Value1 = node.Value
	}

	return jNode
}
