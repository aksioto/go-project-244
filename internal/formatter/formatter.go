package formatter

import (
	"code/internal/diff"
)

// Formatter интерфейс для форматирования различий
type Formatter interface {
	Format(nodes []*diff.Node) string
}

// GetFormatter возвращает форматер по имени
func GetFormatter(name string) Formatter {
	switch name {
	case "stylish":
		return &StylishFormatter{}
	default:
		return &StylishFormatter{} // По умолчанию stylish
	}
}

// StylishFormatter реализует форматирование в стиле stylish
type StylishFormatter struct{}

func (f *StylishFormatter) Format(nodes []*diff.Node) string {
	return Stylish(nodes)
}
