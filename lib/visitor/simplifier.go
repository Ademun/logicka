package visitor

import (
	"errors"
	"fmt"
	"logicka/lib/ast"
	"logicka/lib/lexer"
	"logicka/lib/simplification/rules/base"
	"slices"
)

type Simplifier struct {
	ruleSets []*base.RuleSet
	log      base.ApplicationLogger
}

func NewSimplifier() *Simplifier {
	return &Simplifier{ruleSets: make([]*base.RuleSet, 0), log: base.NewBasicApplicationLogger()}
}

func (s *Simplifier) AddRuleSet(ruleSet *base.RuleSet) {
	s.ruleSets = append(s.ruleSets, ruleSet)
}

func (s *Simplifier) Simplify(node ast.Node) (ast.Node, error) {
	if node == nil {
		return nil, fmt.Errorf("empty node")
	}

	s.log.Clear()
	current := node

	for range 100 {
		next, err := Accept[ast.Node](current, s)
		if err != nil {
			return nil, err
		}

		if current.Equals(next) {
			s.log.LogApplication("Итоговое выражение", "Описание", current.String(), next.String())
			fmt.Println(s.log.String())
			return current, nil
		}

		s.log.LogApplication("Выражение после перобразований", "Описание", current.String(), next.String())
		current = next
	}

	return current, nil
}

func (s *Simplifier) VisitGroupingExpr(node *ast.GroupingExpr) (ast.Node, error) {
	expr, err := Accept[ast.Node](node.Body, s)
	if err != nil {
		return nil, err
	}
	current := ast.NewGroupingExpr(expr.(ast.Expr))

	simplified, err := s.applyAllRuleSets(current)
	if err != nil {
		return nil, err
	}

	return s.tryUnwrap(simplified), nil
}

func (s *Simplifier) VisitBooleanExpr(node *ast.BooleanExpr) (ast.Node, error) {
	return node, nil
}

func (s *Simplifier) VisitVariable(node *ast.IdentifierExpr) (ast.Node, error) {
	return node, nil
}

func (s *Simplifier) VisitBinaryExpr(node *ast.BinaryExpr) (ast.Node, error) {
	left, errL := Accept[ast.Node](node.Left, s)
	right, errR := Accept[ast.Node](node.Right, s)
	if err := errors.Join(errL, errR); err != nil {
		return nil, err
	}

	current := ast.NewBinaryExpr(node.Operator, s.tryWrap(left.(ast.Expr)), s.tryWrap(right.(ast.Expr)))

	return s.applyAllRuleSets(current)
}

func (s *Simplifier) VisitChainExpr(node *ast.ChainExpr) (ast.Node, error) {
	simplified := make([]ast.Expr, 0, len(node.Elements))
	for _, operand := range node.Elements {
		if operand == nil {
			continue
		}

		simplifiedOperand, err := Accept[ast.Node](operand, s)
		if err != nil {
			return nil, err
		}
		simplified = append(simplified, s.tryWrap(simplifiedOperand.(ast.Expr)))
	}

	current, err := s.applyAllRuleSets(ast.NewChainExpr(node.Operator, simplified))
	if err != nil {
		return nil, err
	}

	if binary, ok := current.(*ast.BinaryExpr); ok {
		return s.applyAllRuleSets(binary)
	}

	simplifiedChain, err := s.simplifyChain(current.(*ast.ChainExpr))

	if err != nil {
		return nil, err
	}
	return s.applyAllRuleSets(simplifiedChain)
}

func (s *Simplifier) simplifyChain(node *ast.ChainExpr) (ast.Expr, error) {
	operands := slices.Clone(node.Elements)
	newOperands := make([]ast.Expr, 0, len(operands))
Outer:
	for i := len(operands) - 1; i >= 0; i-- {
		one := operands[i]
		for j := i - 1; j >= 0; j-- {
			other := operands[j]
			combination := ast.NewBinaryExpr(node.Operator, one, other)
			simplifiedCombination, err := Accept[ast.Node](combination, s)
			if err != nil {
				return nil, err
			}
			if !combination.Equals(simplifiedCombination) {
				if t, ok := simplifiedCombination.(ast.Traversable); ok && len(t.Children()) > 1 {
					newOperands = append(newOperands, t.Children()...)
				} else {
					newOperands = append(newOperands, simplifiedCombination.(ast.Expr))
				}
				operands = append(operands[:j], operands[j+1:]...)
				i--
				j--
				continue Outer
			}
		}
		newOperands = append(newOperands, one)
	}
	slices.Reverse(newOperands)

	var result ast.Expr
	switch len(newOperands) {
	case 0:
		result = ast.NewBooleanExpr(node.Operator == lexer.BlDisjunction)
		return result, nil
	case 1:
		result = newOperands[0]
		return result, nil
	case 2:
		result = ast.NewBinaryExpr(node.Operator, newOperands[0], newOperands[1])
		return result, nil
	default:
		result = ast.NewChainExpr(node.Operator, newOperands)
		return result, nil
	}
}

func (s *Simplifier) VisitUnaryExpr(node *ast.UnaryExpr) (ast.Node, error) {
	operand, err := Accept[ast.Node](node.Operand, s)
	if err != nil {
		return nil, err
	}

	current := ast.NewUnaryExpr(node.Operator, operand.(ast.Expr))
	return s.applyAllRuleSets(current)
}

func (s *Simplifier) VisitAssignmentExpr(node *ast.AssignmentExpr) (ast.Node, error) {
	return node, nil
}

func (s *Simplifier) VisitBlockStmt(node *ast.BlockStmt) (ast.Node, error) {
	simplified := make([]ast.Stmt, len(node.Body))
	for i, stmt := range node.Body {
		simplifiedStmt, err := Accept[ast.Node](stmt, s)
		if err != nil {
			return nil, err
		}
		simplified[i] = simplifiedStmt.(ast.Stmt)
	}
	return ast.NewBlockStmt(simplified), nil
}

func (s *Simplifier) VisitBracedExpr(node *ast.BracedExpr) (ast.Node, error) {
	return node, nil
}

func (s *Simplifier) VisitExprStmt(node *ast.ExprStmt) (ast.Node, error) {
	simplified, err := Accept[ast.Node](node.Expr, s)
	if err != nil {
		return nil, err
	}
	return ast.NewExprStmt(simplified.(ast.Expr)), nil
}

func (s *Simplifier) VisitFunctionCallExpr(node *ast.FunctionCallExpr) (ast.Node, error) {
	return node, nil
}

func (s *Simplifier) VisitFunctionDeclExpr(node *ast.FunctionDeclExpr) (ast.Node, error) {
	return node, nil
}

func (s *Simplifier) VisitIdentifierExpr(node *ast.IdentifierExpr) (ast.Node, error) {
	return node, nil
}

func (s *Simplifier) VisitNumberExpr(node *ast.NumberExpr) (ast.Node, error) {
	return node, nil
}

func (s *Simplifier) VisitQuantifierExpr(node *ast.QuantifierExpr) (ast.Node, error) {
	return node, nil
}

func (s *Simplifier) VisitStringExpr(node *ast.StringExpr) (ast.Node, error) {
	return node, nil
}

func (s *Simplifier) applyAllRuleSets(node ast.Expr) (ast.Expr, error) {
	current := node

	for _, ruleSet := range s.ruleSets {
		simplified, err := ruleSet.Apply(current, s.log)
		if err != nil {
			return nil, fmt.Errorf("error in set")
		}

		current = simplified
	}

	return current, nil
}

func (s *Simplifier) tryWrap(node ast.Expr) ast.Expr {
	switch v := node.(type) {
	case *ast.BinaryExpr, *ast.ChainExpr:
		return ast.NewGroupingExpr(v)
	default:
		return v
	}
}

func (s *Simplifier) tryUnwrap(node ast.Expr) ast.Expr {
	grouping, ok := node.(*ast.GroupingExpr)
	if !ok {
		return node
	}

	switch v := grouping.Body.(type) {
	case *ast.BinaryExpr, *ast.ChainExpr:
		return node
	default:
		return v
	}
}
