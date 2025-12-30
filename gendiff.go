package code

import (
	"code/internal/diff"
	"code/internal/formatter"
	"code/internal/parser"
	"code/internal/utils"
	"errors"
	"fmt"
	"reflect"
)

var (
	ErrEmptyPath = errors.New("file path cannot be empty")
)

type Differ struct {
	fileParser *parser.FileParser
}

type Option func(*Differ)

func NewDiffer(opts ...Option) *Differ {
	d := &Differ{
		fileParser: defaultParsers(),
	}

	for _, opt := range opts {
		opt(d)
	}

	return d
}

func WithFileParser(fp *parser.FileParser) Option {
	return func(d *Differ) {
		d.fileParser = fp
	}
}

func defaultParsers() *parser.FileParser {
	p := parser.NewFileParser()
	p.Add(&parser.JSONParser{}, ".json")
	p.Add(&parser.YAMLParser{}, ".yaml", ".yml")
	return p
}

// GenDiff compares two files and returns the diff in the requested format.
func GenDiff(path1, path2, format string) (string, error) {
	return NewDiffer().GetDiff(path1, path2, format)
}

func (d *Differ) GetDiff(path1, path2, format string) (string, error) {
	if path1 == "" {
		return "", fmt.Errorf("first file: %w", ErrEmptyPath)
	}
	if path2 == "" {
		return "", fmt.Errorf("second file: %w", ErrEmptyPath)
	}

	data1, err := d.fileParser.Parse(path1)
	if err != nil {
		return "", fmt.Errorf("parse first file %q: %w", path1, err)
	}

	data2, err := d.fileParser.Parse(path2)
	if err != nil {
		return "", fmt.Errorf("parse second file %q: %w", path2, err)
	}

	fmter, err := formatter.GetFormatter(format)
	if err != nil {
		return "", fmt.Errorf("get formatter: %w", err)
	}

	nodes := d.getNodes(data1, data2)

	result, err := fmter.Format(nodes)
	if err != nil {
		return "", fmt.Errorf("format diff: %w", err)
	}

	return result, nil
}

func (d *Differ) getNodes(data1, data2 map[string]interface{}) []*diff.Node {
	keys := utils.MergedSortedKeys(data1, data2)

	var nodes []*diff.Node
	for _, key := range keys {
		v1, ok1 := data1[key]
		v2, ok2 := data2[key]

		node := d.getNode(key, v1, ok1, v2, ok2)
		if node != nil {
			nodes = append(nodes, node)
		}
	}

	return nodes
}

func (d *Differ) getNode(key string, v1 interface{}, ok1 bool, v2 interface{}, ok2 bool) *diff.Node {
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
			children := d.getNodes(m1, m2)
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
