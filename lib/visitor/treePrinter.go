package visitor

import (
	"fmt"
	"logicka/lib/ast"
	"strings"
)

type TreePrinter struct {
	indentLevel int
}

func NewTreePrinter() *TreePrinter {
	return &TreePrinter{}
}

func (t *TreePrinter) Visit(node ast.ASTNode) {
	switch n := node.(type) {
	case *ast.GroupingNode:
		t.visitGrouping(n)
	case *ast.LiteralNode:
		t.visitLiteral(n)
	case *ast.VariableNode:
		t.visitVariable(n)
	case *ast.BinaryNode:
		t.visitBinary(n)
	case *ast.UnaryNode:
		t.visitUnary(n)
	case *ast.PredicateNode:
		t.visitPredicate(n)
	case *ast.QuantifierNode:
		t.visitQuantifier(n)
	default:
		panic(fmt.Sprintf("unknown node type: %T", node))
	}
	fmt.Println(node.String())
}

func (t *TreePrinter) visitGrouping(node *ast.GroupingNode) {
	t.printIndent("Grouping")
	t.indentLevel++
	t.Visit(node.Expr)
	t.indentLevel--
}

func (t *TreePrinter) visitLiteral(node *ast.LiteralNode) {
	t.printIndent(fmt.Sprintf("Literal: %t", node.Value))
}

func (t *TreePrinter) visitVariable(node *ast.VariableNode) {
	t.printIndent(fmt.Sprintf("Variable: %s", node.Name))
}

func (t *TreePrinter) visitBinary(node *ast.BinaryNode) {
	t.printIndent(fmt.Sprintf("Binary: %s", node.Operator.String()))
	t.indentLevel++
	t.Visit(node.Left)
	t.Visit(node.Right)
	t.indentLevel--
}

func (t *TreePrinter) visitUnary(node *ast.UnaryNode) {
	t.printIndent(fmt.Sprintf("Unary: %s", node.Operator.String()))
	t.indentLevel++
	t.Visit(node.Operand)
	t.indentLevel--
}

func (t *TreePrinter) visitPredicate(node *ast.PredicateNode) {
	t.printIndent(fmt.Sprintf("Predicate: %s", node.Name))
}

func (t *TreePrinter) visitQuantifier(node *ast.QuantifierNode) {
	t.printIndent(fmt.Sprintf("Quantifier: %s %s",
		node.Type.String(),
		node.Variable))
}

func (t *TreePrinter) printIndent(text string) {
	fmt.Printf("%s%s\n", strings.Repeat("  ", t.indentLevel), text)
}
