package ast

import "hash/fnv"

// VariableNode represents a logical variable
type VariableNode struct {
	Name string
}

func NewVariableNode(name string) *VariableNode {
	return &VariableNode{name}
}

func (n *VariableNode) Hash() uint64 {
	h := fnv.New64a()
	h.Write([]byte("variable"))
	h.Write([]byte(n.Name))
	return h.Sum64()
}

func (n *VariableNode) Equals(other ASTNode) bool {
	variable, ok := other.(*VariableNode)
	return ok && n.Hash() == variable.Hash()
}

func (n *VariableNode) String() string {
	return n.Name
}
