package ast

import (
	"hash/fnv"
	"logicka/lib/utils"
	"strings"
)

type FunctionDeclExpr struct {
	Name string
	Args []Expr
	Body Expr
}

func NewFunctionDeclExpr(name string, args []Expr, body Expr) (*FunctionDeclExpr, error) {
	for _, arg := range args {
		if _, ok := arg.(*IdentifierExpr); !ok {
			return nil, NewValidationError("Arguments of function declaration must be identifiers")
		}
	}

	return &FunctionDeclExpr{
		Name: name,
		Args: args,
		Body: body,
	}, nil
}

func (n *FunctionDeclExpr) Hash() uint64 {
	h := fnv.New64a()
	h.Write([]byte("function_decl"))
	for _, arg := range n.Args {
		h.Write(utils.Uint64ToBytes(arg.Hash()))
	}
	h.Write(utils.Uint64ToBytes(n.Body.Hash()))
	return h.Sum64()
}

func (n *FunctionDeclExpr) Equals(other Node) bool {
	return n.Hash() == other.Hash()
}

func (n *FunctionDeclExpr) String() string {
	var results []string
	for _, arg := range n.Args {
		results = append(results, arg.String())
	}
	return n.Name + "(" + strings.Join(results, ", ") + ")" + " := " + n.Body.String()
}

func (n *FunctionDeclExpr) expr() {}
