package diff

type NodeType string

const (
	NodeTypeAdded     NodeType = "added"
	NodeTypeRemoved   NodeType = "removed"
	NodeTypeChanged   NodeType = "changed"
	NodeTypeUnchanged NodeType = "unchanged"
	NodeTypeNested    NodeType = "nested"
)

// Node represents a node in the diff tree.
type Node struct {
	Type     NodeType    `json:"type"`
	Key      string      `json:"key"`
	Value    interface{} `json:"value,omitempty"`
	OldValue interface{} `json:"oldValue,omitempty"`
	NewValue interface{} `json:"newValue,omitempty"`
	Children []*Node     `json:"children,omitempty"`
}
