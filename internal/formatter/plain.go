package formatter

import (
	"code/internal/diff"
	"fmt"
	"strings"
)

type PlainFormatter struct{}

func (f *PlainFormatter) Format(nodes []*diff.Node) (string, error) {
	var lines []string
	f.collectPlainLines(nodes, "", &lines)

	if len(lines) == 0 {
		return "", nil
	}

	return strings.Join(lines, "\n"), nil
}

func (f *PlainFormatter) collectPlainLines(nodes []*diff.Node, parentPath string, lines *[]string) {
	for _, node := range nodes {
		currentPath := buildPath(parentPath, node.Key)

		switch node.Type {
		case diff.NodeTypeAdded:
			line := fmt.Sprintf(
				"Property '%s' was added with value: %s",
				currentPath,
				formatPlainValue(node.Value),
			)
			*lines = append(*lines, line)

		case diff.NodeTypeRemoved:
			line := fmt.Sprintf(
				"Property '%s' was removed",
				currentPath,
			)
			*lines = append(*lines, line)

		case diff.NodeTypeChanged:
			line := fmt.Sprintf(
				"Property '%s' was updated. From %s to %s",
				currentPath,
				formatPlainValue(node.OldValue),
				formatPlainValue(node.NewValue),
			)
			*lines = append(*lines, line)

		case diff.NodeTypeUnchanged:
			// Unchanged properties are not rendered in the plain format
			continue

		case diff.NodeTypeNested:
			// Process nested nodes recursively
			f.collectPlainLines(node.Children, currentPath, lines)
		}
	}
}

func buildPath(parentPath, key string) string {
	if parentPath == "" {
		return key
	}
	return parentPath + "." + key
}

func formatPlainValue(v interface{}) string {
	if v == nil {
		return "null"
	}

	// If the value is a map, output [complex value]
	if _, ok := v.(map[string]interface{}); ok {
		return "[complex value]"
	}

	// Arrays and slices are also considered complex values
	if _, ok := v.([]interface{}); ok {
		return "[complex value]"
	}

	switch val := v.(type) {
	case string:
		// Wrap strings in single quotes
		return fmt.Sprintf("'%s'", val)
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
