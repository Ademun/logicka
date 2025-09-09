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

func (t *TreePrinter) Print(node ast.ASTNode) error {
	_, err := Accept[interface{}](node, t)
	return err
}

func (t *TreePrinter) VisitGrouping(node *ast.GroupingNode) (interface{}, error) {
	t.printIndent("Grouping")
	t.indentLevel++
	_, err := Accept[interface{}](node.Expr, t)
	t.indentLevel--
	return nil, err
}

func (t *TreePrinter) VisitLiteral(node *ast.LiteralNode) (interface{}, error) {
	t.printIndent(fmt.Sprintf("Literal: %t", node.Value))
	return nil, nil
}

func (t *TreePrinter) VisitVariable(node *ast.VariableNode) (interface{}, error) {
	t.printIndent(fmt.Sprintf("Variable: %s", node.Name))
	return nil, nil
}

func (t *TreePrinter) VisitBinary(node *ast.BinaryNode) (interface{}, error) {
	t.printIndent(fmt.Sprintf("Binary: %s", node.Operator.String()))
	t.indentLevel++
	_, err1 := Accept[interface{}](node.Left, t)
	_, err2 := Accept[interface{}](node.Right, t)
	t.indentLevel--
	if err1 != nil {
		return nil, err1
	}
	return nil, err2
}

func (t *TreePrinter) VisitUnary(node *ast.UnaryNode) (interface{}, error) {
	t.printIndent(fmt.Sprintf("Unary: %s", node.Operator.String()))
	t.indentLevel++
	_, err := Accept[interface{}](node.Operand, t)
	t.indentLevel--
	return nil, err
}

func (t *TreePrinter) VisitPredicate(node *ast.PredicateNode) (interface{}, error) {
	t.printIndent(fmt.Sprintf("Predicate: %s", node.Name))
	return nil, nil
}

func (t *TreePrinter) VisitQuantifier(node *ast.QuantifierNode) (interface{}, error) {
	t.printIndent(fmt.Sprintf("Quantifier: %s %s",
		node.Type.String(),
		node.Variable))
	return nil, nil
}

func (t *TreePrinter) printIndent(text string) {
	fmt.Printf("%s%s\n", strings.Repeat("  ", t.indentLevel), text)
}
