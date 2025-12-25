package formatter

import (
	"code/internal/diff"
	"fmt"
	"sort"
	"strings"
)

const indentSize = 4

// Stylish форматирует Diff в стиле stylish
func Stylish(nodes []*diff.Node) string {
	if len(nodes) == 0 {
		return "{\n}"
	}

	var sb strings.Builder
	sb.WriteString("{\n")
	formatNodes(&sb, nodes, 1)
	sb.WriteString("}")
	return sb.String()
}

func formatNodes(sb *strings.Builder, nodes []*diff.Node, depth int) {
	for _, node := range nodes {
		switch node.Type {
		case diff.NodeTypeAdded:
			formatValue(sb, depth, "+", node.Key, node.Value)

		case diff.NodeTypeRemoved:
			formatValue(sb, depth, "-", node.Key, node.Value)

		case diff.NodeTypeChanged:
			formatValue(sb, depth, "-", node.Key, node.OldValue)
			formatValue(sb, depth, "+", node.Key, node.NewValue)

		case diff.NodeTypeUnchanged:
			formatValue(sb, depth, " ", node.Key, node.Value)

		case diff.NodeTypeNested:
			lineIndent := makeIndent(depth, " ")

			sb.WriteString(fmt.Sprintf(
				"%s%s: {\n",
				lineIndent,
				node.Key,
			))

			formatNodes(sb, node.Children, depth+1)

			sb.WriteString(fmt.Sprintf("%s}\n", makeIndent(depth, " ")))
		}
	}
}

func formatValue(sb *strings.Builder, depth int, marker string, key string, value interface{}) {
	lineIndent := makeIndent(depth, marker)

	if m, ok := value.(map[string]interface{}); ok {
		sb.WriteString(fmt.Sprintf("%s%s: {\n", lineIndent, key))

		formatMap(sb, m, depth+1)

		closeIndent := makeIndentForMap(depth)
		sb.WriteString(fmt.Sprintf("%s}\n", closeIndent))
		return
	}

	sb.WriteString(fmt.Sprintf("%s%s: %s\n", lineIndent, key, formatSimpleValue(value)))
}

func formatMap(sb *strings.Builder, m map[string]interface{}, depth int) {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	valIndent := makeIndentForMap(depth)

	for _, k := range keys {
		v := m[k]

		if nested, ok := v.(map[string]interface{}); ok {
			sb.WriteString(fmt.Sprintf("%s%s: {\n", valIndent, k))

			formatMap(sb, nested, depth+1)

			sb.WriteString(fmt.Sprintf("%s}\n", valIndent))
		} else {
			sb.WriteString(fmt.Sprintf("%s%s: %s\n", valIndent, k, formatSimpleValue(v)))
		}
	}
}

func formatSimpleValue(v interface{}) string {
	if v == nil {
		return "null"
	}

	switch val := v.(type) {
	case string:
		return val
	case bool:
		if val {
			return "true"
		}
		return "false"
	case float64:
		if val == float64(int64(val)) {
			return fmt.Sprintf("%.0f", val)
		}
		return fmt.Sprintf("%g", val)
	default:
		return fmt.Sprintf("%v", v)
	}
}

func makeIndent(depth int, marker string) string {
	spacesCount := depth*indentSize - 2
	if spacesCount < 0 {
		spacesCount = 0
	}
	return strings.Repeat(" ", spacesCount) + marker + " "
}

func makeIndentForMap(depth int) string {
	spacesCount := depth * indentSize
	return strings.Repeat(" ", spacesCount)
}
