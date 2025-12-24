package code

import (
	"code/internal/parser"
	"fmt"
	"reflect"
	"sort"
	"strings"
)

var reg *parser.Registry

func init() {
	reg = parser.NewRegistry()
	reg.Register(&parser.JSONParser{}, ".json")
	reg.Register(&parser.YAMLParser{}, ".yaml", ".yml")
}

func GenDiff(path1, path2 string) (string, error) {
	data1, err := reg.ParseFile(path1)
	if err != nil {
		return "", err
	}

	data2, err := reg.ParseFile(path2)
	if err != nil {
		return "", err
	}

	diff := getDiff(data1, data2)
	return diff, nil
}

func getDiff(data1, data2 map[string]interface{}) string {
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

	var sb strings.Builder
	sb.WriteString("{\n")

	for _, k := range keys {
		v1, ok1 := data1[k]
		v2, ok2 := data2[k]

		switch {
		case ok1 && !ok2:
			sb.WriteString(fmt.Sprintf("  - %s: %v\n", k, v1))
		case !ok1 && ok2:
			sb.WriteString(fmt.Sprintf("  + %s: %v\n", k, v2))
		case ok1 && ok2 && !reflect.DeepEqual(v1, v2):
			sb.WriteString(fmt.Sprintf("  - %s: %v\n", k, v1))
			sb.WriteString(fmt.Sprintf("  + %s: %v\n", k, v2))
		default: // одинаковые значения
			sb.WriteString(fmt.Sprintf("    %s: %v\n", k, v1))
		}
	}

	sb.WriteString("}")
	return sb.String()
}
