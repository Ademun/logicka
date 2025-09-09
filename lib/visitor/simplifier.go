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

func (s *Simplifier) Simplify(node ast.ASTNode) (ast.ASTNode, error) {
	current := node
	for {
		next, err := Accept[ast.ASTNode](current, s)
		if err != nil {
			return nil, err
		}
		if next.Equals(current) {
			return next, nil
		}
		current = next
	}
}

func (s *Simplifier) VisitGrouping(node *ast.GroupingNode) (ast.ASTNode, error) {
	expr, err := Accept[ast.ASTNode](node.Expr, s)
	if err != nil {
		return nil, err
	}
	return &ast.GroupingNode{Expr: expr}, nil
}

func (s *Simplifier) VisitLiteral(node *ast.LiteralNode) (ast.ASTNode, error) {
	return node, nil
}

func (s *Simplifier) VisitVariable(node *ast.VariableNode) (ast.ASTNode, error) {
	return node, nil
}

func (s *Simplifier) VisitBinary(node *ast.BinaryNode) (ast.ASTNode, error) {
	left, err := Accept[ast.ASTNode](node.Left, s)
	if err != nil {
		return nil, err
	}
	right, err := Accept[ast.ASTNode](node.Right, s)
	if err != nil {
		return nil, err
	}

	switch node.Operator {
	case lexer.IMPL:
		return s.simplifyImplication(left, right)
	case lexer.EQUIV:
		return s.simplifyEquivalence(left, right)
	case lexer.CONJ:
		return s.simplifyConjunction(left, right)
	case lexer.DISJ:
		return s.simplifyDisjunction(left, right)
	default:
		panic("unhandled default case")
	}
}

func (s *Simplifier) simplifyConjunction(left, right ast.ASTNode) (ast.ASTNode, error) {
	if ast.IsFalse(left) || ast.IsFalse(right) {
		return &ast.LiteralNode{Value: false}, nil
	}
	if ast.IsTrue(left) {
		return right, nil
	}
	if ast.IsTrue(right) {
		return left, nil
	}
	if ast.IsNegationOfSame(left, right) || ast.IsNegationOfSame(right, left) {
		return &ast.LiteralNode{Value: false}, nil
	}
	if left.Equals(right) {
		return left, nil
	}
	return &ast.BinaryNode{Operator: lexer.CONJ, Left: left, Right: right}, nil
}

func (s *Simplifier) simplifyDisjunction(left, right ast.ASTNode) (ast.ASTNode, error) {
	if ast.IsTrue(left) || ast.IsTrue(right) {
		return &ast.LiteralNode{Value: true}, nil
	}
	if ast.IsFalse(left) {
		return right, nil
	}
	if ast.IsFalse(right) {
		return left, nil
	}
	if ast.IsNegationOfSame(left, right) || ast.IsNegationOfSame(right, left) {
		return &ast.LiteralNode{Value: true}, nil
	}
	if left.Equals(right) {
		return left, nil
	}
	if ast.IsConjunctionWith(left, right) {
		return right, nil
	}
	if ast.IsConjunctionWith(right, left) {
		return left, nil
	}
	return &ast.BinaryNode{Operator: lexer.DISJ, Left: left, Right: right}, nil
}

func (s *Simplifier) simplifyImplication(left, right ast.ASTNode) (ast.ASTNode, error) {
	negLeft := &ast.UnaryNode{Operator: lexer.NEG, Operand: left}
	return Accept[ast.ASTNode](&ast.BinaryNode{
		Operator: lexer.DISJ,
		Left:     negLeft,
		Right:    right,
	}, s)
}

func (s *Simplifier) simplifyEquivalence(left, right ast.ASTNode) (ast.ASTNode, error) {
	return Accept[ast.ASTNode](&ast.BinaryNode{
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
	}, s)
}

func (s *Simplifier) VisitUnary(node *ast.UnaryNode) (ast.ASTNode, error) {
	operand, err := Accept[ast.ASTNode](node.Operand, s)
	if err != nil {
		return nil, err
	}

	switch node.Operator {
	case lexer.NEG:
		return s.simplifyNegation(operand)
	default:
		panic("unhandled default case")
	}
}

func (s *Simplifier) simplifyNegation(operand ast.ASTNode) (ast.ASTNode, error) {
	if unary, ok := operand.(*ast.UnaryNode); ok && unary.Operator == lexer.NEG {
		return unary.Operand, nil
	}
	if literal, ok := operand.(*ast.LiteralNode); ok {
		return &ast.LiteralNode{Value: !literal.Value}, nil
	}
	return &ast.UnaryNode{Operator: lexer.NEG, Operand: operand}, nil
}

func (s *Simplifier) VisitPredicate(node *ast.PredicateNode) (ast.ASTNode, error) {
	// TODO implement me
	panic("implement me")
}

func (s *Simplifier) VisitQuantifier(node *ast.QuantifierNode) (ast.ASTNode, error) {
	// TODO implement me
	panic("implement me")
}
