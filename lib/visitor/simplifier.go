package visitor

import (
	"errors"
	"fmt"
	"logicka/lib/ast"
	"logicka/lib/lexer"
	"slices"
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
		fmt.Println("Old:", current, "New:", next)
		current = next
	}
}

func (s *Simplifier) VisitGrouping(node *ast.GroupingNode) (ast.ASTNode, error) {
	expr, err := Accept[ast.ASTNode](node.Expr, s)
	if err != nil {
		return nil, err
	}
	if len(expr.Children()) == 1 {
		return expr, nil
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
	switch op := node.Operator; op {
	case lexer.CONJ:
		return s.simplifyConjunction(node)
	case lexer.DISJ:
		return s.simplifyDisjunction(node)
	case lexer.IMPL:
		return s.simplifyImplication(node)
	case lexer.EQUIV:
		return s.simplifyEquivalence(node)
	default:
		return nil, fmt.Errorf("unknown binary operand: %s", op.String())
	}
}

func (s *Simplifier) simplifyConjunction(node *ast.BinaryNode) (ast.ASTNode, error) {
	left, errl := Accept[ast.ASTNode](node.Left, s)
	right, errr := Accept[ast.ASTNode](node.Right, s)
	if errl != nil || errr != nil {
		return nil, errors.Join(errl, errr)
	}

	// Conjunction with 1 A & 1 => A
	if ast.IsTrue(left) {
		return right, nil
	}
	if ast.IsTrue(right) {
		return left, nil
	}
	// Conjunction with 0 0 & A => 0
	if ast.IsFalse(left) || ast.IsFalse(right) {
		return &ast.LiteralNode{Value: false}, nil
	}
	// Tautology A & A => A
	if left.Equals(right) {
		return left, nil
	}
	// Contradiction A & !A => 0
	if ast.IsNegationOfSame(left, right) || ast.IsNegationOfSame(right, left) {
		return &ast.LiteralNode{Value: false}, nil
	}

	if absorbed, err, ok := s.applyAbsorptionConjunction(left, right); err == nil && ok {
		return absorbed, nil
	} else if err != nil {
		return nil, err
	}
	if absorbed, err, ok := s.applyAbsorptionConjunction(right, left); err == nil && ok {
		return absorbed, nil
	} else if err != nil {
		return nil, err
	}

	if ast.CanFlatten(node) {
		return ast.ChainFromBinary(node), nil
	}

	return &ast.BinaryNode{
		Operator: lexer.CONJ,
		Left:     left,
		Right:    right,
	}, nil
}

func (s *Simplifier) applyAbsorptionDisjunction(left, right ast.ASTNode) (ast.ASTNode, error, bool) {
	if group, ok := right.(*ast.GroupingNode); ok {
		var chainExpr *ast.ChainNode
		switch n := group.Expr.(type) {
		case *ast.BinaryNode:
			chainExpr = ast.ChainFromBinary(n)
		case *ast.ChainNode:
			chainExpr = n
		default:
			return nil, fmt.Errorf("unknown expression type: %T", n), false
		}
		if chainExpr.Operator != lexer.CONJ {
			return &ast.BinaryNode{
				Operator: lexer.DISJ,
				Left:     left,
				Right:    right,
			}, nil, false
		}
		if chainExpr.Contains(left) {
			return left, nil, true
		}
		if neg, ok := left.(*ast.UnaryNode); ok && neg.Operator == lexer.NEG {
			if chainExpr.Contains(neg.Operand) {
				chainExpr.Remove(neg.Operand)
				return &ast.BinaryNode{
					Operator: lexer.DISJ,
					Left:     left,
					Right:    &ast.GroupingNode{Expr: chainExpr},
				}, nil, true
			}
		}
		if chainExpr.Contains(&ast.UnaryNode{Operator: lexer.NEG, Operand: left}) {
			chainExpr.Remove(&ast.UnaryNode{Operator: lexer.NEG, Operand: left})
			return &ast.BinaryNode{
				Operator: lexer.DISJ,
				Left:     left,
				Right:    &ast.GroupingNode{Expr: chainExpr},
			}, nil, true
		}
	}
	return &ast.BinaryNode{
		Operator: lexer.DISJ,
		Left:     left,
		Right:    right,
	}, nil, false
}

func (s *Simplifier) applyAbsorptionConjunction(left, right ast.ASTNode) (ast.ASTNode, error, bool) {
	if group, ok := right.(*ast.GroupingNode); ok {
		var chainExpr *ast.ChainNode
		switch n := group.Expr.(type) {
		case *ast.BinaryNode:
			chainExpr = ast.ChainFromBinary(n)
		case *ast.ChainNode:
			chainExpr = n
		default:
			return nil, fmt.Errorf("unknown expression type: %T", n), false
		}
		if chainExpr.Operator != lexer.DISJ {
			return &ast.BinaryNode{
				Operator: lexer.CONJ,
				Left:     left,
				Right:    right,
			}, nil, false
		}
		if chainExpr.Contains(left) {
			return left, nil, true
		}
		if neg, ok := left.(*ast.UnaryNode); ok && neg.Operator == lexer.NEG {
			if chainExpr.Contains(neg.Operand) {
				chainExpr.Remove(neg.Operand)
				return &ast.BinaryNode{
					Operator: lexer.CONJ,
					Left:     left,
					Right:    &ast.GroupingNode{Expr: chainExpr},
				}, nil, true
			}
		}
		if chainExpr.Contains(&ast.UnaryNode{Operator: lexer.NEG, Operand: left}) {
			chainExpr.Remove(&ast.UnaryNode{Operator: lexer.NEG, Operand: left})
			return &ast.BinaryNode{
				Operator: lexer.CONJ,
				Left:     left,
				Right:    &ast.GroupingNode{Expr: chainExpr},
			}, nil, true
		}
	}
	return &ast.BinaryNode{
		Operator: lexer.CONJ,
		Left:     left,
		Right:    right,
	}, nil, false
}

func (s *Simplifier) simplifyDisjunction(node *ast.BinaryNode) (ast.ASTNode, error) {
	left, errl := Accept[ast.ASTNode](node.Left, s)
	right, errr := Accept[ast.ASTNode](node.Right, s)
	if errl != nil || errr != nil {
		return nil, errors.Join(errl, errr)
	}

	// Disjunction with 1 A \/ 1 => 1
	if ast.IsTrue(left) || ast.IsTrue(right) {
		return &ast.LiteralNode{Value: true}, nil
	}
	// Disjunction with 0 0 \/ A => A
	if ast.IsFalse(left) {
		return right, nil
	}
	if ast.IsFalse(right) {
		return left, nil
	}
	// Tautology A \/ A => A
	if left.Equals(right) {
		return left, nil
	}
	// Contradiction A \/ !A => 1
	if ast.IsNegationOfSame(left, right) || ast.IsNegationOfSame(right, left) {
		return &ast.LiteralNode{Value: true}, nil
	}

	if absorbed, err, ok := s.applyAbsorptionDisjunction(left, right); err == nil && ok {
		return absorbed, nil
	} else if err != nil {
		return nil, err
	}
	if absorbed, err, ok := s.applyAbsorptionDisjunction(right, left); err == nil && ok {
		return absorbed, nil
	} else if err != nil {
		return nil, err
	}

	if ast.CanFlatten(node) {
		return ast.ChainFromBinary(node), nil
	}

	return &ast.BinaryNode{
		Operator: lexer.DISJ,
		Left:     left,
		Right:    right,
	}, nil
}

func (s *Simplifier) simplifyImplication(node *ast.BinaryNode) (ast.ASTNode, error) {
	left, errl := Accept[ast.ASTNode](node.Left, s)
	right, errr := Accept[ast.ASTNode](node.Right, s)
	if errl != nil || errr != nil {
		return nil, errors.Join(errl, errr)
	}
	return &ast.BinaryNode{
		Operator: lexer.DISJ,
		Left: &ast.UnaryNode{
			Operator: lexer.NEG,
			Operand:  left,
		},
		Right: right,
	}, nil
}

func (s *Simplifier) simplifyEquivalence(node *ast.BinaryNode) (ast.ASTNode, error) {
	left, errl := Accept[ast.ASTNode](node.Left, s)
	right, errr := Accept[ast.ASTNode](node.Right, s)
	if errl != nil || errr != nil {
		return nil, errors.Join(errl, errr)
	}

	return &ast.BinaryNode{
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
	}, nil
}

func (s *Simplifier) VisitChain(node *ast.ChainNode) (ast.ASTNode, error) {
	simplifiedOperands := make([]ast.ASTNode, 0)
	for _, op := range node.Operands {
		if op == nil {
			continue
		}
		simplifiedNode, err := Accept[ast.ASTNode](op, s)
		if err != nil {
			return nil, err
		}
		simplifiedOperands = append(simplifiedOperands, simplifiedNode)
	}
	simplifiedChain := &ast.ChainNode{Operator: node.Operator, Operands: simplifiedOperands}
	simplifiedCChain, err := s.simplifyChain(simplifiedChain)
	if err != nil {
		return nil, err
	}

	return simplifiedCChain, nil
}

func (s *Simplifier) simplifyChain(node *ast.ChainNode) (ast.ASTNode, error) {
	simplifiedOperands := make([]ast.ASTNode, 0)
	for i := len(node.Operands) - 1; i >= 0; i-- {
		if i == 0 {
			simplifiedOperands = append(simplifiedOperands, node.Operands[i])
			break
		}
		simplified := false
		for j := i - 1; j >= 0; j-- {
			if ast.IsTrue(node.Operands[i]) || ast.IsTrue(node.Operands[j]) {
				if node.Operator == lexer.DISJ {
					return &ast.LiteralNode{Value: true}, nil
				}
			}
			if ast.IsFalse(node.Operands[i]) || ast.IsFalse(node.Operands[j]) {
				if node.Operator == lexer.CONJ {
					return &ast.LiteralNode{Value: false}, nil
				}
			}
			comb := &ast.BinaryNode{
				Operator: node.Operator,
				Left:     node.Operands[i],
				Right:    node.Operands[j],
			}
			simplifiedComb, err := Accept[ast.ASTNode](comb, s)
			if err != nil {
				return nil, err
			}
			if !comb.Equals(simplifiedComb) {
				simplified = true
				simplifiedOperands = append(simplifiedOperands, simplifiedComb.Children()...)
				node.Operands = append(node.Operands[:i], node.Operands[i+1:]...)
				node.Operands = append(node.Operands[:j], node.Operands[j+1:]...)
				i -= 1
				break
			}
		}
		if !simplified {
			simplifiedOperands = append(simplifiedOperands, node.Operands[i])
		}
	}
	slices.Reverse(simplifiedOperands)
	return &ast.ChainNode{Operator: node.Operator, Operands: simplifiedOperands}, nil
}

func (s *Simplifier) VisitUnary(node *ast.UnaryNode) (ast.ASTNode, error) {
	switch op := node.Operator; op {
	case lexer.NEG:
		return s.simplifyNegation(node)
	default:
		return nil, fmt.Errorf("unknown unary operand: %s", op.String())
	}
}

func (s *Simplifier) simplifyNegation(node *ast.UnaryNode) (ast.ASTNode, error) {
	operand, err := Accept[ast.ASTNode](node.Operand, s)
	if err != nil {
		return nil, err
	}
	// Double negation !!A => A
	if un, ok := operand.(*ast.UnaryNode); ok && ast.IsNegation(un) {
		return un.Operand, nil
	}

	// Negation of literal !1 => 0; !0 => 1
	if lit, ok := operand.(*ast.LiteralNode); ok {
		return &ast.LiteralNode{Value: !lit.Value}, nil
	}
	// De Morgan laws
	if group, ok := operand.(*ast.GroupingNode); ok {
		if bin, ok := group.Expr.(*ast.BinaryNode); ok {
			// !(A & B) => !A \/ !B
			if bin.Operator == lexer.CONJ {
				return &ast.BinaryNode{
					Operator: lexer.DISJ,
					Left: &ast.UnaryNode{
						Operator: lexer.NEG,
						Operand:  bin.Left,
					},
					Right: &ast.UnaryNode{
						Operator: lexer.NEG,
						Operand:  bin.Right,
					},
				}, nil
			}
			// !(A \/ B) => !A & !B
			if bin.Operator == lexer.DISJ {
				return &ast.BinaryNode{
					Operator: lexer.CONJ,
					Left: &ast.UnaryNode{
						Operator: lexer.NEG,
						Operand:  bin.Left,
					},
					Right: &ast.UnaryNode{
						Operator: lexer.NEG,
						Operand:  bin.Right,
					},
				}, nil
			}
		}
		if chain, ok := group.Expr.(*ast.ChainNode); ok {
			negated := make([]ast.ASTNode, len(chain.Operands))
			for i, op := range chain.Operands {
				negated[i] = &ast.UnaryNode{
					Operator: lexer.NEG,
					Operand:  op,
				}
			}
			var negatedOperator lexer.BooleanTokenType
			if chain.Operator == lexer.DISJ {
				negatedOperator = lexer.CONJ
			} else {
				negatedOperator = lexer.DISJ
			}
			return &ast.GroupingNode{Expr: &ast.ChainNode{Operator: negatedOperator, Operands: negated}}, nil
		}
	}
	return &ast.UnaryNode{
		Operator: lexer.NEG,
		Operand:  operand,
	}, nil
}

func (s *Simplifier) VisitPredicate(node *ast.PredicateNode) (ast.ASTNode, error) {
	// TODO implement me
	panic("implement me")
}

func (s *Simplifier) VisitQuantifier(node *ast.QuantifierNode) (ast.ASTNode, error) {
	// TODO implement me
	panic("implement me")
}
