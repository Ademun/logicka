package ast

import (
	"fmt"
	"logicka/lib/lexer"
	"slices"
	"strings"
)

type ASTNode interface {
	Equals(node ASTNode) bool
	Children() []ASTNode
	Contains(node ASTNode) bool
	String() string
}

type GroupingNode struct {
	Expr ASTNode
}

func (g *GroupingNode) Equals(node ASTNode) bool {
	same, ok := node.(*GroupingNode)
	return ok && same.Expr.Equals(g.Expr)
}

func (g *GroupingNode) Children() []ASTNode {
	return []ASTNode{g.Expr}
}

func (g *GroupingNode) Contains(node ASTNode) bool {
	return g.Expr.Equals(node)
}

func (g *GroupingNode) String() string {
	return fmt.Sprintf("(%s)", g.Expr.String())
}

type LiteralNode struct {
	Value bool
}

func (l *LiteralNode) Equals(node ASTNode) bool {
	same, ok := node.(*LiteralNode)
	return ok && same.Value == l.Value
}

func (l *LiteralNode) Children() []ASTNode {
	return []ASTNode{}
}

func (l *LiteralNode) Contains(node ASTNode) bool {
	return false
}

func (l *LiteralNode) String() string {
	return fmt.Sprint(l.Value)
}

type VariableNode struct {
	Name string
}

func (v *VariableNode) Equals(node ASTNode) bool {
	same, ok := node.(*VariableNode)
	return ok && same.Name == v.Name
}

func (v *VariableNode) Children() []ASTNode {
	return []ASTNode{}
}

func (v *VariableNode) Contains(node ASTNode) bool {
	return false
}

func (v *VariableNode) String() string {
	return v.Name
}

type BinaryNode struct {
	Operator    lexer.BooleanTokenType
	Left, Right ASTNode
}

func (b *BinaryNode) Equals(node ASTNode) bool {
	same, ok := node.(*BinaryNode)
	return ok && b.Operator == same.Operator && b.Left.Equals(same.Left) && b.Right.Equals(same.Right)
}

func (b *BinaryNode) Children() []ASTNode {
	return []ASTNode{b.Left, b.Right}
}

func (b *BinaryNode) Contains(node ASTNode) bool {
	return b.Left.Equals(node) || b.Right.Equals(node)
}

func (b *BinaryNode) String() string {
	return fmt.Sprint(b.Left.String(), " ", b.Operator.String(), " ", b.Right.String())
}

// ChainNode represents a flattened chain of binary operations of the same type
type ChainNode struct {
	Operator lexer.BooleanTokenType // CONJ or DISJ
	Operands []ASTNode              // Must have at least 2 operands
}

func (c *ChainNode) Equals(node ASTNode) bool {
	same, ok := node.(*ChainNode)
	if !ok || c.Operator != same.Operator || len(c.Operands) != len(same.Operands) {
		return false
	}

	for i, operand := range c.Operands {
		if !operand.Equals(same.Operands[i]) {
			return false
		}
	}

	return true
}

func (c *ChainNode) Children() []ASTNode {
	return c.Operands
}

func (c *ChainNode) Contains(node ASTNode) bool {
	return slices.ContainsFunc(c.Operands, func(a ASTNode) bool {
		return a.Equals(node)
	})
}

func (c *ChainNode) String() string {
	parts := make([]string, 0)
	for _, operand := range c.Operands {
		if operand == nil {
			continue
		}
		parts = append(parts, operand.String())
	}

	return strings.Join(parts, " "+c.Operator.String()+" ")
}

func (c *ChainNode) Remove(node ASTNode) {
	slices.DeleteFunc(c.Operands, func(a ASTNode) bool {
		return a.Equals(node)
	})
}

// ChainFromBinary creates a ChainNode by flattening nested binary operations
func ChainFromBinary(node *BinaryNode) *ChainNode {
	if node.Operator != lexer.CONJ && node.Operator != lexer.DISJ {
		return nil
	}

	chain := &ChainNode{
		Operator: node.Operator,
		Operands: make([]ASTNode, 0),
	}

	chain.collectOperands(node)
	if len(chain.Operands) < 2 {
		return nil
	}
	return chain
}

// collectOperands recursively collects operands from nested binary nodes
func (c *ChainNode) collectOperands(node ASTNode) {
	switch n := node.(type) {
	case *BinaryNode:
		if n.Operator == c.Operator {
			c.collectOperands(n.Left)
			c.collectOperands(n.Right)
		} else {
			c.Operands = append(c.Operands, &GroupingNode{Expr: node})
		}
	case *GroupingNode:
		// Unwrap groupings and collect their contents
		c.collectOperands(n.Expr)
	default:
		c.Operands = append(c.Operands, node)
	}
}

// ToBinary converts the ChainNode back to nested BinaryNode structure
func (c *ChainNode) ToBinary() ASTNode {
	if len(c.Operands) == 2 {
		return &BinaryNode{
			Operator: c.Operator,
			Left:     c.Operands[0],
			Right:    c.Operands[1],
		}
	}

	// For more than 2 operands, build left-associative tree
	result := &BinaryNode{
		Operator: c.Operator,
		Left:     c.Operands[0],
		Right:    c.Operands[1],
	}

	for i := 2; i < len(c.Operands); i++ {
		result = &BinaryNode{
			Operator: c.Operator,
			Left:     result,
			Right:    c.Operands[i],
		}
	}

	return result
}

// CanFlatten checks if a BinaryNode can be flattened into a ChainNode
func CanFlatten(node *BinaryNode) bool {

	if node.Operator != lexer.CONJ && node.Operator != lexer.DISJ {
		return false
	}
	return hasNestedSameOperator(node, node.Operator)
}

// hasNestedSameOperator checks if the binary tree contains nested operations of the same type
func hasNestedSameOperator(node ASTNode, operator lexer.BooleanTokenType) bool {
	switch n := node.(type) {
	case *BinaryNode:
		if n.Operator == operator {
			// Check if either child is also the same operator (indicating nesting)
			if leftBin, ok := n.Left.(*BinaryNode); ok && leftBin.Operator == operator {
				return true
			}
			if rightBin, ok := n.Right.(*BinaryNode); ok && rightBin.Operator == operator {
				return true
			}
			if leftGr, ok := n.Left.(*GroupingNode); ok {
				if bin, ok1 := leftGr.Expr.(*BinaryNode); ok1 && bin.Operator == operator {
					return true
				}
			}
			if rightGr, ok := n.Right.(*GroupingNode); ok {
				if bin, ok1 := rightGr.Expr.(*BinaryNode); ok1 && bin.Operator == operator {
					return true
				}
			}
			// Also check through groupings
			return hasNestedSameOperator(n.Left, operator) || hasNestedSameOperator(n.Right, operator)
		}
	case *GroupingNode:
		return hasNestedSameOperator(n.Expr, operator)
	}
	return false
}

type UnaryNode struct {
	Operator lexer.BooleanTokenType
	Operand  ASTNode
}

func (u *UnaryNode) Equals(node ASTNode) bool {
	same, ok := node.(*UnaryNode)
	return ok && u.Operator == same.Operator && u.Operand.Equals(same.Operand)
}

func (u *UnaryNode) Children() []ASTNode {
	return []ASTNode{u.Operand}
}

func (u *UnaryNode) Contains(node ASTNode) bool {
	return u.Operand.Equals(node)
}

func (u *UnaryNode) String() string {
	return fmt.Sprint(u.Operator.String(), u.Operand.String())
}

type PredicateNode struct {
	Name string
	Body any
}

func (p *PredicateNode) Equals(node ASTNode) bool {
	return true
}

func (p *PredicateNode) Children() []ASTNode {
	panic("implement me")
}

func (p *PredicateNode) Contains(node ASTNode) bool {
	panic("implement me")
}

func (p *PredicateNode) String() string {
	return p.Name
}

type QuantifierNode struct {
	Type     lexer.BooleanTokenType
	Variable string
	Domain   any
}

func (q *QuantifierNode) Equals(node ASTNode) bool {
	return true
}

func (q *QuantifierNode) Children() []ASTNode {
	panic("implement me")
}

func (q *QuantifierNode) Contains(node ASTNode) bool {
	panic("implement me")
}

func (q *QuantifierNode) String() string {
	return fmt.Sprint(q.Variable)
}
