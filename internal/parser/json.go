package parser

import (
	"encoding/json"
	"fmt"
)

type JSONParser struct{}

func (p *JSONParser) Parse(data []byte) (map[string]interface{}, error) {
	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("parse json: %w", err)
	}
	return result, nil
}
