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
)

var SupportedFormats = []string{
	FormatStylish,
	FormatPlain,
}

var ErrUnknownFormat = errors.New("unknown format")

// Formatter интерфейс для форматирования различий
type Formatter interface {
	Format(nodes []*diff.Node) string
}

// GetFormatter возвращает форматер по имени
func GetFormatter(name string) (Formatter, error) {
	if name == "" {
		return &StylishFormatter{}, nil
	}

	switch name {
	case FormatStylish:
		return &StylishFormatter{}, nil
	case FormatPlain:
		return &PlainFormatter{}, nil
	default:
		return nil, fmt.Errorf("%w: %s (allowed: %s)", ErrUnknownFormat, name, strings.Join(SupportedFormats, ", "))
	}
}

// StylishFormatter реализует форматирование в стиле stylish
type StylishFormatter struct{}

func (f *StylishFormatter) Format(nodes []*diff.Node) string {
	return Stylish(nodes)
}

type PlainFormatter struct{}

func (f *PlainFormatter) Format(nodes []*diff.Node) string {
	return Plain(nodes)
}
