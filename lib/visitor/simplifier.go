package visitor

import (
	"logicka/lib/ast"
	"logicka/lib/lexer"
)

type Simplifier struct {
}

func NewSimplifier() *Simplifier {
	return &Simplifier{}
}

func (s *Simplifier) Visit(node ast.ASTNode) ast.ASTNode {
	current := node
	for {
		next := s.visitSwitch(current)
		if next.Equals(current) {
			return next
		}
		current = next
	}
}

func (s *Simplifier) visitSwitch(node ast.ASTNode) ast.ASTNode {
	switch n := node.(type) {
	case *ast.GroupingNode:
		return s.visitGrouping(n)
	case *ast.LiteralNode:
		return s.visitLiteral(n)
	case *ast.VariableNode:
		return s.visitVariable(n)
	case *ast.BinaryNode:
		return s.visitBinary(n)
	case *ast.UnaryNode:
		return s.visitUnary(n)
	case *ast.PredicateNode:
		return s.visitPredicate(n)
	case *ast.QuantifierNode:
		return s.visitQuantifier(n)
	default:
		// TODO: implement me
		panic("implement me")
	}
}

func (s *Simplifier) visitGrouping(node *ast.GroupingNode) ast.ASTNode {
	return &ast.GroupingNode{Expr: s.Visit(node.Expr)}
}

func (s *Simplifier) visitLiteral(node *ast.LiteralNode) ast.ASTNode {
	return node
}

func (s *Simplifier) visitVariable(node *ast.VariableNode) ast.ASTNode {
	return node
}

func (s *Simplifier) visitBinary(node *ast.BinaryNode) ast.ASTNode {
	switch node.Operator {
	case lexer.IMPL:
		return s.simplifyImplication(node)
	case lexer.EQUIV:
		return s.simplifyEquivalence(node)
	case lexer.CONJ:
		return s.simplifyConjunction(node)
	case lexer.DISJ:
		return s.simplifyDisjunction(node)
	default:
		panic("unhandled default case")
	}
}

func (s *Simplifier) simplifyConjunction(node *ast.BinaryNode) ast.ASTNode {
	left, right := s.Visit(node.Left), s.Visit(node.Right)
	if ast.IsFalse(left) || ast.IsFalse(right) {
		return &ast.LiteralNode{Value: false}
	}
	if ast.IsTrue(left) {
		return right
	}
	if ast.IsTrue(right) {
		return left
	}
	if ast.IsNegationOfSame(left, right) || ast.IsNegationOfSame(right, left) {
		return &ast.LiteralNode{Value: false}
	}
	if left.Equals(right) {
		return left
	}
	return &ast.BinaryNode{Operator: lexer.CONJ, Left: left, Right: right}
}

func (s *Simplifier) simplifyDisjunction(node *ast.BinaryNode) ast.ASTNode {
	left, right := s.Visit(node.Left), s.Visit(node.Right)
	if ast.IsTrue(left) || ast.IsTrue(right) {
		return &ast.LiteralNode{Value: true}
	}
	if ast.IsFalse(left) {
		return right
	}
	if ast.IsFalse(right) {
		return left
	}
	if ast.IsNegationOfSame(left, right) || ast.IsNegationOfSame(right, left) {
		return &ast.LiteralNode{Value: true}
	}
	if left.Equals(right) {
		return left
	}
	if ast.IsConjunctionWith(left, right) {
		return right
	}
	if ast.IsConjunctionWith(right, left) {
		return left
	}
	return &ast.BinaryNode{Operator: lexer.DISJ, Left: left, Right: right}
}

func (s *Simplifier) simplifyImplication(node *ast.BinaryNode) ast.ASTNode {
	left, right := s.Visit(node.Left), s.Visit(node.Right)
	return s.Visit(&ast.BinaryNode{
		Operator: lexer.DISJ,
		Left:     &ast.UnaryNode{Operator: lexer.NEG, Operand: left},
		Right:    right,
	})
}

func (s *Simplifier) simplifyEquivalence(node *ast.BinaryNode) ast.ASTNode {
	left, right := s.Visit(node.Left), s.Visit(node.Right)
	return s.Visit(&ast.BinaryNode{
		Operator: lexer.CONJ,
		Left: &ast.GroupingNode{
			Expr: &ast.BinaryNode{
				Operator: lexer.IMPL,
				Left:     left,
				Right:    right,
			},
		},
		Right: &ast.GroupingNode{
			Expr: &ast.BinaryNode{
				Operator: lexer.IMPL,
				Left:     right,
				Right:    left,
			},
		},
	})
}

func (s *Simplifier) visitUnary(node *ast.UnaryNode) ast.ASTNode {
	switch node.Operator {
	case lexer.NEG:
		return s.simplifyNegation(node)
	default:
		panic("unhandled default case")
	}
}

func (s *Simplifier) simplifyNegation(node *ast.UnaryNode) ast.ASTNode {
	operand := s.Visit(node.Operand)
	if ast.IsNegation(operand) {
		return operand.(*ast.UnaryNode).Operand
	}
	if literal, ok := node.Operand.(*ast.LiteralNode); ok {
		return &ast.LiteralNode{Value: !literal.Value}
	}
	return node
}

func (s *Simplifier) visitPredicate(node *ast.PredicateNode) ast.ASTNode {
	// TODO implement me
	panic("implement me")
}

func (s *Simplifier) visitQuantifier(node *ast.QuantifierNode) ast.ASTNode {
	// TODO implement me
	panic("implement me")
}
