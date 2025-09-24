package ast

import (
	"logicka/lib/lexer"
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
	case *IdentifierExpr, *BracedExpression:
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

func (n *QuantifierExpr) String() string {
	var postfix string
	if n.Domain != nil {
		postfix = "∈ " + n.Domain.String()
	}
	return n.Type.String() + n.Variable.String() + postfix + ": " + n.Body.String()
}

func (n *QuantifierExpr) Children() []Expr {
	return []Expr{n.Body}
}

func (n *QuantifierExpr) expr() {}
