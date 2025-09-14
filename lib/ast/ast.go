// Package ast provides Abstract Syntax Tree nodes and operations for logical expressions.
package ast

import (
	"errors"
	"fmt"
	"slices"
	"strings"

	"logicka/lib/lexer"
)

var (
	// ErrUnknownNodeType indicates an AST node type that cannot be processed
	ErrUnknownNodeType = errors.New("unknown AST node type")

	// ErrUnknownOperator indicates an unsupported operator
	ErrUnknownOperator = errors.New("unknown operator")

	// ErrInvalidChain indicates a chain node with insufficient operands
	ErrInvalidChain = errors.New("chain must have at least 2 operands")
)

// ASTNode represents a node in the abstract syntax tree for logical expressions.
type ASTNode interface {
	// Equals returns true if this node is equivalent to the other node
	Equals(other ASTNode) bool

	// String returns a string representation of the node
	String() string
}

// Traversable represents nodes that can be traversed (have children).
type Traversable interface {
	ASTNode
	Children() []ASTNode
}

// Container represents nodes that can contain other nodes.
type Container interface {
	ASTNode
	Contains(node ASTNode) bool
}

// Mutable represents nodes that can be modified.
type Mutable interface {
	ASTNode
	// Remove removes the specified node if it exists
	Remove(node ASTNode) bool
}

// GroupingNode represents parenthesized expressions.
type GroupingNode struct {
	Expr ASTNode
}

func NewGroupingNode(expr ASTNode) *GroupingNode {
	return &GroupingNode{Expr: expr}
}

func (g *GroupingNode) Equals(other ASTNode) bool {
	node, ok := other.(*GroupingNode)
	return ok && g.Expr.Equals(node.Expr)
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

// LiteralNode represents boolean literals (true/false).
type LiteralNode struct {
	Value bool
}

func NewLiteralNode(value bool) *LiteralNode {
	return &LiteralNode{Value: value}
}

func (l *LiteralNode) Equals(other ASTNode) bool {
	node, ok := other.(*LiteralNode)
	return ok && l.Value == node.Value
}

func (l *LiteralNode) Children() []ASTNode {
	return nil
}

func (l *LiteralNode) Contains(node ASTNode) bool {
	return l.Equals(node)
}

func (l *LiteralNode) String() string {
	if l.Value {
		return "true"
	}
	return "false"
}

// VariableNode represents logical variables.
type VariableNode struct {
	Name string
}

func NewVariableNode(name string) *VariableNode {
	return &VariableNode{Name: name}
}

func (v *VariableNode) Equals(other ASTNode) bool {
	node, ok := other.(*VariableNode)
	return ok && v.Name == node.Name
}

func (v *VariableNode) Children() []ASTNode {
	return nil
}

func (v *VariableNode) Contains(node ASTNode) bool {
	return v.Equals(node)
}

func (v *VariableNode) String() string {
	return v.Name
}

// BinaryNode represents binary logical operations.
type BinaryNode struct {
	Operator    lexer.BooleanTokenType
	Left, Right ASTNode
}

func NewBinaryNode(operator lexer.BooleanTokenType, left, right ASTNode) *BinaryNode {
	return &BinaryNode{
		Operator: operator,
		Left:     left,
		Right:    right,
	}
}

func (b *BinaryNode) Equals(other ASTNode) bool {
	node, ok := other.(*BinaryNode)
	return ok &&
		b.Operator == node.Operator &&
		b.Left.Equals(node.Left) &&
		b.Right.Equals(node.Right)
}

func (b *BinaryNode) Children() []ASTNode {
	return []ASTNode{b.Left, b.Right}
}

func (b *BinaryNode) Contains(node ASTNode) bool {
	return b.Left.Equals(node) || b.Right.Equals(node)
}

func (b *BinaryNode) String() string {
	return fmt.Sprintf("%s %s %s",
		b.Left.String(),
		b.Operator.String(),
		b.Right.String())
}

// ChainNode represents a flattened chain of binary operations of the same type.
// This optimization reduces tree depth for associative operations.
type ChainNode struct {
	Operator lexer.BooleanTokenType // CONJ or DISJ
	Operands []ASTNode              // Must have at least 2 operands
}

func NewChainNode(operator lexer.BooleanTokenType, operands ...ASTNode) (*ChainNode, error) {
	if len(operands) < 2 {
		return nil, fmt.Errorf("%w: got %d operands", ErrInvalidChain, len(operands))
	}

	return &ChainNode{
		Operator: operator,
		Operands: slices.Clone(operands), // Defensive copy
	}, nil
}

func (c *ChainNode) Equals(other ASTNode) bool {
	node, ok := other.(*ChainNode)
	if !ok || c.Operator != node.Operator || len(c.Operands) != len(node.Operands) {
		return false
	}

	for i, operand := range c.Operands {
		if !operand.Equals(node.Operands[i]) {
			return false
		}
	}
	return true
}

func (c *ChainNode) Children() []ASTNode {
	return slices.Clone(c.Operands) // Defensive copy
}

func (c *ChainNode) Contains(node ASTNode) bool {
	return slices.ContainsFunc(c.Operands, func(operand ASTNode) bool {
		return operand.Equals(node)
	})
}

func (c *ChainNode) Remove(node ASTNode) bool {
	originalLen := len(c.Operands)
	c.Operands = slices.DeleteFunc(c.Operands, func(operand ASTNode) bool {
		return operand.Equals(node)
	})
	return len(c.Operands) != originalLen
}

func (c *ChainNode) String() string {
	if len(c.Operands) == 0 {
		return ""
	}

	parts := make([]string, 0, len(c.Operands))
	for _, operand := range c.Operands {
		if operand != nil {
			parts = append(parts, operand.String())
		}
	}

	return "[" + strings.Join(parts, " "+c.Operator.String()+" ") + "]"
}

// ToBinary converts the ChainNode to a left-associative binary tree.
func (c *ChainNode) ToBinary() ASTNode {
	if len(c.Operands) < 2 {
		return nil
	}

	if len(c.Operands) == 2 {
		return NewBinaryNode(c.Operator, c.Operands[0], c.Operands[1])
	}

	// Build left-associative tree: ((a op b) op c) op d...
	result := NewBinaryNode(c.Operator, c.Operands[0], c.Operands[1])
	for i := 2; i < len(c.Operands); i++ {
		result = NewBinaryNode(c.Operator, result, c.Operands[i])
	}

	return result
}

// UnaryNode represents unary logical operations (primarily negation).
type UnaryNode struct {
	Operator lexer.BooleanTokenType
	Operand  ASTNode
}

func NewUnaryNode(operator lexer.BooleanTokenType, operand ASTNode) *UnaryNode {
	return &UnaryNode{
		Operator: operator,
		Operand:  operand,
	}
}

func (u *UnaryNode) Equals(other ASTNode) bool {
	node, ok := other.(*UnaryNode)
	return ok &&
		u.Operator == node.Operator &&
		u.Operand.Equals(node.Operand)
}

func (u *UnaryNode) Children() []ASTNode {
	return []ASTNode{u.Operand}
}

func (u *UnaryNode) Contains(node ASTNode) bool {
	return u.Operand.Equals(node)
}

func (u *UnaryNode) String() string {
	return fmt.Sprintf("%s%s", u.Operator.String(), u.Operand.String())
}

// PredicateNode represents predicate logic expressions (future extension).
type PredicateNode struct {
	Name string
	Args []ASTNode // Changed from 'any' to []ASTNode for type safety
}

func NewPredicateNode(name string, args ...ASTNode) *PredicateNode {
	return &PredicateNode{
		Name: name,
		Args: slices.Clone(args),
	}
}

func (p *PredicateNode) Equals(other ASTNode) bool {
	node, ok := other.(*PredicateNode)
	if !ok || p.Name != node.Name || len(p.Args) != len(node.Args) {
		return false
	}

	for i, arg := range p.Args {
		if !arg.Equals(node.Args[i]) {
			return false
		}
	}
	return true
}

func (p *PredicateNode) Children() []ASTNode {
	return slices.Clone(p.Args)
}

func (p *PredicateNode) Contains(node ASTNode) bool {
	return slices.ContainsFunc(p.Args, func(arg ASTNode) bool {
		return arg.Equals(node)
	})
}

func (p *PredicateNode) String() string {
	if len(p.Args) == 0 {
		return p.Name
	}

	argStrs := make([]string, len(p.Args))
	for i, arg := range p.Args {
		argStrs[i] = arg.String()
	}

	return fmt.Sprintf("%s(%s)", p.Name, strings.Join(argStrs, ", "))
}

// QuantifierNode represents quantified expressions (∀, ∃).
type QuantifierNode struct {
	Type     lexer.BooleanTokenType
	Variable string
	Domain   ASTNode // Changed from 'any' to ASTNode
	Body     ASTNode
}

func NewQuantifierNode(qType lexer.BooleanTokenType, variable string, domain, body ASTNode) *QuantifierNode {
	return &QuantifierNode{
		Type:     qType,
		Variable: variable,
		Domain:   domain,
		Body:     body,
	}
}

func (q *QuantifierNode) Equals(other ASTNode) bool {
	node, ok := other.(*QuantifierNode)
	return ok &&
		q.Type == node.Type &&
		q.Variable == node.Variable &&
		q.Domain.Equals(node.Domain) &&
		q.Body.Equals(node.Body)
}

func (q *QuantifierNode) Children() []ASTNode {
	return []ASTNode{q.Domain, q.Body}
}

func (q *QuantifierNode) Contains(node ASTNode) bool {
	return q.Domain.Equals(node) || q.Body.Equals(node)
}

func (q *QuantifierNode) String() string {
	return fmt.Sprintf("%s %s: %s . %s",
		q.Type.String(),
		q.Variable,
		q.Domain.String(),
		q.Body.String())
}

// Utility functions for common AST operations

// IsTrue returns true if the node represents a true literal.
func IsTrue(node ASTNode) bool {
	literal, ok := node.(*LiteralNode)
	return ok && literal.Value
}

// IsFalse returns true if the node represents a false literal.
func IsFalse(node ASTNode) bool {
	literal, ok := node.(*LiteralNode)
	return ok && !literal.Value
}

// IsNegation returns true if the node is a negation operation.
func IsNegation(node ASTNode) bool {
	unary, ok := node.(*UnaryNode)
	return ok && unary.Operator == lexer.NEG
}

// IsNegationOf returns true if left is the negation of right.
func IsNegationOf(left, right ASTNode) bool {
	unary, ok := left.(*UnaryNode)
	return ok && unary.Operator == lexer.NEG && unary.Operand.Equals(right)
}

// CanFlatten determines if a binary node can be converted to a chain node.
func CanFlatten(node *BinaryNode) bool {
	if node.Operator != lexer.CONJ && node.Operator != lexer.DISJ {
		return false
	}
	return hasNestedSameOperator(node, node.Operator)
}

// hasNestedSameOperator checks for nested operations of the same type.
func hasNestedSameOperator(node ASTNode, operator lexer.BooleanTokenType) bool {
	switch n := node.(type) {
	case *BinaryNode:
		if n.Operator == operator {
			return isSameOperatorBinary(n.Left, operator) ||
				isSameOperatorBinary(n.Right, operator) ||
				hasNestedSameOperator(n.Left, operator) ||
				hasNestedSameOperator(n.Right, operator)
		}
	case *GroupingNode:
		return hasNestedSameOperator(n.Expr, operator)
	}
	return false
}

// isSameOperatorBinary checks if a node is a binary node with the specified operator.
func isSameOperatorBinary(node ASTNode, operator lexer.BooleanTokenType) bool {
	switch n := node.(type) {
	case *BinaryNode:
		return n.Operator == operator
	case *ChainNode:
		return n.Operator == operator
	case *GroupingNode:
		if bin, ok := n.Expr.(*BinaryNode); ok {
			return bin.Operator == operator
		}
		if chain, ok := n.Expr.(*ChainNode); ok {
			return chain.Operator == operator
		}
	}
	return false
}

// ChainFromBinary creates a ChainNode by flattening nested binary operations.
func ChainFromBinary(node *BinaryNode) *ChainNode {
	if !CanFlatten(node) {
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

// collectOperands recursively collects operands from nested binary nodes.
func (c *ChainNode) collectOperands(node ASTNode) {
	switch n := node.(type) {
	case *BinaryNode:
		if n.Operator == c.Operator {
			c.collectOperands(n.Left)
			c.collectOperands(n.Right)
		} else {
			c.Operands = append(c.Operands, NewGroupingNode(node))
		}
	case *ChainNode:
		if n.Operator == c.Operator {
			c.Operands = append(c.Operands, n.Operands...)
		} else {
			c.Operands = append(c.Operands, NewGroupingNode(node))
		}
	case *GroupingNode:
		c.collectOperands(n.Expr)
	default:
		c.Operands = append(c.Operands, node)
	}
}
