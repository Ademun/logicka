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

func (s *Simplifier) Simplify(node ast.ASTNode) (ast.ASTNode, error) {
	if node == nil {
		return nil, fmt.Errorf("empty node")
	}

	s.log.Clear()
	current := node

	for range 100 {
		next, err := Accept[ast.ASTNode](current, s)
		if err != nil {
			return nil, err
		}
		next = s.tryUnwrap(next)
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

func (s *Simplifier) VisitGrouping(node *ast.GroupingNode) (ast.ASTNode, error) {
	expr, err := Accept[ast.ASTNode](node.Expr, s)
	if err != nil {
		return nil, err
	}
	current := ast.NewGroupingNode(expr)

	simplified, err := s.applyAllRuleSets(current)
	if err != nil {
		return nil, err
	}

	return s.tryUnwrap(simplified), nil
}

func (s *Simplifier) VisitLiteral(node *ast.LiteralNode) (ast.ASTNode, error) {
	return node, nil
}

func (s *Simplifier) VisitVariable(node *ast.VariableNode) (ast.ASTNode, error) {
	return node, nil
}

func (s *Simplifier) VisitBinary(node *ast.BinaryNode) (ast.ASTNode, error) {
	left, errL := Accept[ast.ASTNode](node.Left, s)
	right, errR := Accept[ast.ASTNode](node.Right, s)
	if err := errors.Join(errL, errR); err != nil {
		return nil, err
	}

	current := ast.NewBinaryNode(node.Operator, s.tryWrap(left), s.tryWrap(right))

	return s.applyAllRuleSets(current)
}

func (s *Simplifier) VisitChain(node *ast.ChainNode) (ast.ASTNode, error) {
	simplified := make([]ast.ASTNode, 0, len(node.Operands))
	for _, operand := range node.Operands {
		if operand == nil {
			continue
		}

		simplifiedOperand, err := Accept[ast.ASTNode](operand, s)
		if err != nil {
			return nil, err
		}
		simplified = append(simplified, s.tryWrap(simplifiedOperand))
	}

	current, err := s.applyAllRuleSets(ast.NewChainNode(node.Operator, simplified...))
	if err != nil {
		return nil, err
	}

	if binary, ok := current.(*ast.BinaryNode); ok {
		return s.applyAllRuleSets(binary)
	}

	simplifiedChain, err := s.simplifyChain(current.(*ast.ChainNode))

	if err != nil {
		return nil, err
	}
	return s.applyAllRuleSets(simplifiedChain)
}

func (s *Simplifier) simplifyChain(node *ast.ChainNode) (ast.ASTNode, error) {
	operands := slices.Clone(node.Operands)
	newOperands := make([]ast.ASTNode, 0, len(operands))
Outer:
	for i := len(operands) - 1; i >= 0; i-- {
		one := operands[i]
		for j := i - 1; j >= 0; j-- {
			other := operands[j]
			combination := ast.NewBinaryNode(node.Operator, one, other)
			simplifiedCombination, err := Accept[ast.ASTNode](combination, s)
			if err != nil {
				return nil, err
			}
			if !combination.Equals(simplifiedCombination) {
				if t, ok := simplifiedCombination.(ast.Traversable); ok && len(t.Children()) > 1 {
					newOperands = append(newOperands, t.Children()...)
				} else {
					newOperands = append(newOperands, simplifiedCombination)
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

	var result ast.ASTNode
	switch len(newOperands) {
	case 0:
		result = ast.NewLiteralNode(node.Operator == lexer.DISJ)
		return result, nil
	case 1:
		result = newOperands[0]
		return result, nil
	case 2:
		result = ast.NewBinaryNode(node.Operator, newOperands[0], newOperands[1])
		return result, nil
	default:
		result = &ast.ChainNode{Operator: node.Operator, Operands: newOperands}
		return result, nil
	}
}

func (s *Simplifier) VisitUnary(node *ast.UnaryNode) (ast.ASTNode, error) {
	operand, err := Accept[ast.ASTNode](node.Operand, s)
	if err != nil {
		return nil, err
	}

	current := ast.NewUnaryNode(node.Operator, operand)
	return s.applyAllRuleSets(current)
}

func (s *Simplifier) VisitPredicate(node *ast.PredicateNode) (ast.ASTNode, error) {
	return s.applyAllRuleSets(node)
}

func (s *Simplifier) VisitQuantifier(node *ast.QuantifierNode) (ast.ASTNode, error) {
	return s.applyAllRuleSets(node)
}

func (s *Simplifier) applyAllRuleSets(node ast.ASTNode) (ast.ASTNode, error) {
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

func (s *Simplifier) tryWrap(node ast.ASTNode) ast.ASTNode {
	switch v := node.(type) {
	case *ast.BinaryNode, *ast.ChainNode:
		return ast.NewGroupingNode(v)
	default:
		return v
	}
}

func (s *Simplifier) tryUnwrap(node ast.ASTNode) ast.ASTNode {
	grouping, ok := node.(*ast.GroupingNode)
	if !ok {
		return node
	}

	switch v := grouping.Expr.(type) {
	case *ast.BinaryNode, *ast.ChainNode:
		return node
	default:
		return v
	}
}
