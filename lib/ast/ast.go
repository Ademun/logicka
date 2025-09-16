// Package ast provides Abstract Syntax Tree nodes and operations for logical expressions.
package ast

import (
	"errors"
	"fmt"
	"hash/fnv"
	"logicka/lib/lexer"
	"logicka/lib/utils"
	"slices"
	"sort"
	"strings"
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
	Equals(other ASTNode) bool
	String() string
	Hash() uint64
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
	return ok && g.Hash() == node.Hash()
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

func (g *GroupingNode) Hash() uint64 {
	h := fnv.New64a()
	h.Write([]byte("grouping"))
	h.Write(utils.Uint64ToBytes(g.Expr.Hash()))
	return h.Sum64()
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
	return ok && l.Hash() == node.Hash()
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

func (l *LiteralNode) Hash() uint64 {
	h := fnv.New64a()
	h.Write([]byte("literal"))
	if l.Value {
		h.Write([]byte("true"))
	} else {
		h.Write([]byte("false"))
	}
	return h.Sum64()
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
	return ok && v.Hash() == node.Hash()
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

func (v *VariableNode) Hash() uint64 {
	h := fnv.New64a()
	h.Write([]byte("variable"))
	h.Write([]byte(v.Name))
	return h.Sum64()
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
	return ok && b.Hash() == node.Hash()
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

func (b *BinaryNode) Hash() uint64 {
	h := fnv.New64a()
	h.Write([]byte("binary"))
	h.Write([]byte(b.Operator.String()))

	// For implication (IMPL), order matters
	if b.Operator == lexer.IMPL {
		h.Write(utils.Uint64ToBytes(b.Left.Hash()))
		h.Write(utils.Uint64ToBytes(b.Right.Hash()))
		return h.Sum64()
	}
	// For associative operations (CONJ, DISJ), use order-independent hash
	leftHash := b.Left.Hash()
	rightHash := b.Right.Hash()

	// Combine hashes in a commutative way
	if leftHash < rightHash {
		h.Write(utils.Uint64ToBytes(leftHash))
		h.Write(utils.Uint64ToBytes(rightHash))
	} else {
		h.Write(utils.Uint64ToBytes(rightHash))
		h.Write(utils.Uint64ToBytes(leftHash))
	}

	return h.Sum64()
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
	return ok && c.Hash() == node.Hash()
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

func (c *ChainNode) Hash() uint64 {
	h := fnv.New64a()
	h.Write([]byte("chain"))
	h.Write([]byte(c.Operator.String()))

	// For associative operations (CONJ, DISJ), sort hashes for order independence
	if c.Operator == lexer.CONJ || c.Operator == lexer.DISJ {
		hashes := make([]uint64, len(c.Operands))
		for i, operand := range c.Operands {
			hashes[i] = operand.Hash()
		}
		sort.Slice(hashes, func(i, j int) bool {
			return hashes[i] < hashes[j]
		})

		for _, hash := range hashes {
			h.Write(utils.Uint64ToBytes(hash))
		}
		return h.Sum64()
	}
	// For non-associative operations, preserve order
	for _, operand := range c.Operands {
		h.Write(utils.Uint64ToBytes(operand.Hash()))
	}

	return h.Sum64()
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
	return ok && u.Hash() == node.Hash()
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

func (u *UnaryNode) Hash() uint64 {
	h := fnv.New64a()
	h.Write([]byte("unary"))
	h.Write([]byte(u.Operator.String()))
	h.Write(utils.Uint64ToBytes(u.Operand.Hash()))
	return h.Sum64()
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
	return ok && p.Hash() == node.Hash()
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

func (p *PredicateNode) Hash() uint64 {
	panic("implement me")
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
	return ok && q.Hash() == node.Hash()
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

func (q *QuantifierNode) Hash() uint64 {
	panic("implement me")
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
