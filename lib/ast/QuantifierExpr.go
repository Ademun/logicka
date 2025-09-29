package ast

import (
	"hash/fnv"
	"logicka/lib/lexer"
	"logicka/lib/utils"
)

type QuantifierExpr struct {
	Type     lexer.TokenType
	Variable Expr
	Domain   Expr
	Body     Expr
}

func NewQuantifierExpr(qtype lexer.TokenType, variable Expr, domain, body Expr) (*QuantifierExpr, error) {
	if _, ok := variable.(*IdentifierExpr); !ok {
		return nil, NewValidationError("Quantifier variable must be an identifier")
	}

	switch domain.(type) {
	case *IdentifierExpr, *BracedExpr:
		break
	default:
		return nil, NewValidationError("Quantifier domain must be an identifier or a set literal")
	}
	return &QuantifierExpr{
		Type:     qtype,
		Variable: variable,
		Domain:   domain,
		Body:     body,
	}, nil
}

func (n *QuantifierExpr) Hash() uint64 {
	h := fnv.New64a()
	h.Write([]byte("quantifier"))
	h.Write(utils.Uint64ToBytes(n.Variable.Hash()))
	h.Write(utils.Uint64ToBytes(n.Domain.Hash()))
	h.Write(utils.Uint64ToBytes(n.Body.Hash()))
	return h.Sum64()
}

func (n *QuantifierExpr) Equals(other Node) bool {
	return n.Hash() == other.Hash()
}

func (n *QuantifierExpr) String() string {
	var postfix string
	if n.Domain != nil {
		postfix = "∈ " + n.Domain.String()
	}
	return n.Type.String() + n.Variable.String() + postfix + ": " + n.Body.String()
}

func (n *QuantifierExpr) expr() {}
