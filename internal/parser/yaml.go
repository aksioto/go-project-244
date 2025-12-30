package parser

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

type YAMLParser struct{}

func (p *YAMLParser) Parse(data []byte) (map[string]interface{}, error) {
	var result map[string]interface{}
	if err := yaml.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("parse yaml: %w", err)
	}
	normalizeValues(result)
	return result, nil
}

// normalizeValues recursively converts YAML-native types to JSON-compatible ones.
// YAML uses ints for whole numbers, while JSON expects float64.
func normalizeValues(m map[string]interface{}) {
	for k, v := range m {
		switch val := v.(type) {
		case int:
			m[k] = float64(val)
		case int64:
			m[k] = float64(val)
		case map[string]interface{}:
			normalizeValues(val)
		case []interface{}:
			normalizeSlice(val)
		}
	}
}

func normalizeSlice(s []interface{}) {
	for i, v := range s {
		switch val := v.(type) {
		case int:
			s[i] = float64(val)
		case int64:
			s[i] = float64(val)
		case map[string]interface{}:
			normalizeValues(val)
		case []interface{}:
			normalizeSlice(val)
		}
	}
}
