// Package ast provides Abstract Syntax Tree nodes and operations for logical expressions.
package ast

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
	Add(ast ASTNode)
	Remove(node ASTNode)
}
