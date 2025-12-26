package code

import (
	"code/internal/diff"
	"code/internal/formatter"
	"code/internal/parser"
	"reflect"
	"sort"
)

var reg *parser.Registry

func init() {
	reg = parser.NewRegistry()
	reg.Register(&parser.JSONParser{}, ".json")
	reg.Register(&parser.YAMLParser{}, ".yaml", ".yml")
}

// GenDiff сравнивает два файла и возвращает различия в указанном формате
func GenDiff(path1, path2, format string) (string, error) {
	data1, err := reg.ParseFile(path1)
	if err != nil {
		return "", err
	}

	data2, err := reg.ParseFile(path2)
	if err != nil {
		return "", err
	}

	fmter, err := formatter.GetFormatter(format)
	if err != nil {
		return "", err
	}

	nodes := getNodes(data1, data2)
	return fmter.Format(nodes), nil
}

func getNodes(data1, data2 map[string]interface{}) []*diff.Node {
	keysMap := make(map[string]struct{})
	for k := range data1 {
		keysMap[k] = struct{}{}
	}
	for k := range data2 {
		keysMap[k] = struct{}{}
	}

	keys := make([]string, 0, len(keysMap))
	for k := range keysMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var nodes []*diff.Node
	for _, key := range keys {
		v1, ok1 := data1[key]
		v2, ok2 := data2[key]

		node := getNode(key, v1, ok1, v2, ok2)
		if node != nil {
			nodes = append(nodes, node)
		}
	}

	return nodes
}

func getNode(key string, v1 interface{}, ok1 bool, v2 interface{}, ok2 bool) *diff.Node {
	switch {
	case ok1 && !ok2:
		return &diff.Node{
			Type:  diff.NodeTypeRemoved,
			Key:   key,
			Value: v1,
		}
	case !ok1 && ok2:
		return &diff.Node{
			Type:  diff.NodeTypeAdded,
			Key:   key,
			Value: v2,
		}
	case ok1 && ok2:
		m1, isMap1 := v1.(map[string]interface{})
		m2, isMap2 := v2.(map[string]interface{})
		if isMap1 && isMap2 {
			children := getNodes(m1, m2)
			return &diff.Node{
				Type:     diff.NodeTypeNested,
				Key:      key,
				Children: children,
			}
		}

		if !reflect.DeepEqual(v1, v2) {
			return &diff.Node{
				Type:     diff.NodeTypeChanged,
				Key:      key,
				OldValue: v1,
				NewValue: v2,
			}
		}

		return &diff.Node{Type: diff.NodeTypeUnchanged, Key: key, Value: v1}
	}
	return nil
}
