package ast

import (
	"fmt"
	"hash/fnv"
	"logicka/lib/lexer"
)

// QuantifierNode represents quantified expressions
type QuantifierNode struct {
	Type     lexer.BooleanTokenType
	Variable *VariableNode
	Domain   string
}

func NewQuantifierNode(quantifierType lexer.BooleanTokenType, variable *VariableNode, domain string) *QuantifierNode {
	return &QuantifierNode{Type: quantifierType, Variable: variable, Domain: domain}
}

func (n *QuantifierNode) Hash() uint64 {
	h := fnv.New64a()
	h.Write([]byte("quantifier"))
	h.Write([]byte(n.Type.String()))
	h.Write([]byte(n.Variable.String()))
	h.Write([]byte(n.Domain))
	return h.Sum64()
}

func (n *QuantifierNode) Equals(other ASTNode) bool {
	quantifier, ok := other.(*QuantifierNode)
	return ok && n.Hash() == quantifier.Hash()
}

func (n *QuantifierNode) String() string {
	if n.Domain != "" {
		return fmt.Sprintf("(%s%s ∈ %s)", n.Type.String(), n.Variable.String(), n.Domain)
	}
	return fmt.Sprintf("%s%s", n.Type.String(), n.Variable.String())
}
