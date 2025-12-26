package parser

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var (
	ErrUnsupportedFormat = errors.New("unsupported file format")
	ErrReadFile          = errors.New("read file error")
	ErrAbsPath           = errors.New("cannot resolve file path")
)

type Parser interface {
	Parse(data []byte) (map[string]interface{}, error)
}

type Registry struct {
	parsers        map[string]Parser
	allowedFormats []string
}

func NewRegistry() *Registry {
	return &Registry{
		parsers: make(map[string]Parser),
	}
}

func (r *Registry) Register(p Parser, exts ...string) {
	for _, ext := range exts {
		ext = strings.ToLower(ext)
		r.parsers[ext] = p
		r.allowedFormats = append(r.allowedFormats, ext)
	}
}

func (r *Registry) ParseFile(path string) (map[string]interface{}, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrAbsPath, path)
	}
	absPath = filepath.Clean(absPath)

	data, err := os.ReadFile(absPath)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrReadFile, absPath)
	}

	ext := strings.ToLower(filepath.Ext(absPath))
	parser, ok := r.parsers[ext]
	if !ok {
		return nil, fmt.Errorf("%w: %s (allowed: %s)", ErrUnsupportedFormat, ext, r.getAllowedFormats())
	}

	return parser.Parse(data)
}

func (r *Registry) getAllowedFormats() string {
	return strings.Join(r.allowedFormats, ", ")
}
