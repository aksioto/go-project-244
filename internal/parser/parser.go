package parser

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Parser interface {
	Parse(data []byte) (map[string]interface{}, error)
}

type Registry struct {
	parsers map[string]Parser
}

func NewRegistry() *Registry {
	return &Registry{
		parsers: make(map[string]Parser),
	}
}

func (r *Registry) Register(p Parser, exts ...string) {
	for _, ext := range exts {
		r.parsers[strings.ToLower(ext)] = p
	}
}

func (r *Registry) ParseFile(path string) (map[string]interface{}, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("cannot get absolute path for %s: %w", path, err)
	}
	absPath = filepath.Clean(absPath)

	data, err := os.ReadFile(absPath)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", absPath, err)
	}

	ext := strings.ToLower(filepath.Ext(absPath))
	parser, ok := r.parsers[ext]
	if !ok {
		return nil, fmt.Errorf("unsupported file format: %s", ext)
	}

	return parser.Parse(data)
}
