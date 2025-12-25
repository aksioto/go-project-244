package diff

type NodeType string

const (
	NodeTypeAdded     NodeType = "added"
	NodeTypeRemoved   NodeType = "removed"
	NodeTypeChanged   NodeType = "changed"
	NodeTypeUnchanged NodeType = "unchanged"
	NodeTypeNested    NodeType = "nested"
)

// Node представляет узел в дереве различий
type Node struct {
	Type     NodeType
	Key      string
	Value    interface{}
	OldValue interface{}
	NewValue interface{}
	Children []*Node
}
