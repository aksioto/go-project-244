package formatter

import (
	"code/internal/diff"
	"errors"
	"fmt"
	"strings"
)

const (
	FormatStylish = "stylish"
	FormatPlain   = "plain"
	FormatJSON    = "json"
)

var SupportedFormats = []string{
	FormatStylish,
	FormatPlain,
	FormatJSON,
}

var ErrUnknownFormat = errors.New("unknown format")

// Formatter formats a diff tree into a string representation.
type Formatter interface {
	Format(nodes []*diff.Node) (string, error)
}

// GetFormatter returns a formatter implementation by its name.
func GetFormatter(name string) (Formatter, error) {
	if name == "" {
		return &StylishFormatter{}, nil
	}

	switch name {
	case FormatStylish:
		return &StylishFormatter{}, nil
	case FormatPlain:
		return &PlainFormatter{}, nil
	case FormatJSON:
		return &JSONFormatter{}, nil
	default:
		return nil, fmt.Errorf("%w: %s (allowed: %s)", ErrUnknownFormat, name, strings.Join(SupportedFormats, ", "))
	}
}
