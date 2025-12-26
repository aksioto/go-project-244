package formatter

import (
	"code/internal/diff"
	"fmt"
	"sort"
	"strings"
)

// Plain форматирует дерево различий в плоском формате
func Plain(nodes []*diff.Node) string {
	var lines []string
	collectPlainLines(nodes, "", &lines)

	if len(lines) == 0 {
		return ""
	}

	return strings.Join(lines, "\n")
}

func collectPlainLines(nodes []*diff.Node, parentPath string, lines *[]string) {
	// Сортируем узлы для стабильного вывода
	sortedNodes := make([]*diff.Node, len(nodes))
	copy(sortedNodes, nodes)
	sort.Slice(sortedNodes, func(i, j int) bool {
		return sortedNodes[i].Key < sortedNodes[j].Key
	})

	for _, node := range sortedNodes {
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
			// Неизменённые свойства не выводятся в plain формате
			continue

		case diff.NodeTypeNested:
			// Рекурсивно обрабатываем вложенные узлы
			collectPlainLines(node.Children, currentPath, lines)
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

	// Если значение - это map (объект), выводим [complex value]
	if _, ok := v.(map[string]interface{}); ok {
		return "[complex value]"
	}

	switch val := v.(type) {
	case string:
		// Строки оборачиваем в одинарные кавычки
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
