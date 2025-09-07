package visitor

import (
	"logicka/lib/ast"
)

type Visitor[T any] interface {
	Visit(node *ast.ASTNode) T
}

type EvaluationContext struct {
	Variables map[string]bool
}
