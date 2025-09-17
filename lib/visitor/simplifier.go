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
}

func NewSimplifier() *Simplifier {
	return &Simplifier{ruleSets: make([]*base.RuleSet, 0)}
}

func (s *Simplifier) AddRuleSet(ruleSet *base.RuleSet) {
	s.ruleSets = append(s.ruleSets, ruleSet)
}
func (s *Simplifier) Simplify(node ast.ASTNode) (ast.ASTNode, error) {
	if node == nil {
		return nil, fmt.Errorf("empty node")
	}

	current := node

	for i := range 100 {
		fmt.Println("Iteration", i)
		next, err := Accept[ast.ASTNode](current, s)
		if err != nil {
			return nil, err
		}
		fmt.Println("Было:\n", current, "\n", "Стало:\n", next)
		if current.Equals(next) {
			fmt.Println("Converged on iteration", i)
			return current, nil
		}
		for _, ruleSet := range s.ruleSets {
			records := ruleSet.String(false, true)
			fmt.Println(records)
		}
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

	if s.canRemoveGrouping(simplified) {
		if grouping, ok := simplified.(*ast.GroupingNode); ok {
			return grouping.Expr, nil
		}
	}

	return simplified, nil
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

	if lch, ok := left.(*ast.ChainNode); ok {
		left = ast.NewGroupingNode(lch)
	}
	if rch, ok := right.(*ast.ChainNode); ok {
		right = ast.NewGroupingNode(rch)
	}

	current := ast.NewBinaryNode(node.Operator, left, right)

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
		simplified = append(simplified, simplifiedOperand)
	}

	current, err := s.applyAllRuleSets(ast.NewChainNode(node.Operator, simplified...))
	fmt.Println(current, node.Operands)
	if err != nil {
		return nil, err
	}

	if binary, ok := current.(*ast.BinaryNode); ok {
		fmt.Println("Binary:", binary)
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
		simplified, err := ruleSet.Apply(current)
		if err != nil {
			return nil, fmt.Errorf("error in set")
		}

		current = simplified
	}

	return current, nil
}

func (s *Simplifier) canRemoveGrouping(node ast.ASTNode) bool {
	grouping, ok := node.(*ast.GroupingNode)
	if !ok {
		return false
	}

	switch grouping.Expr.(type) {
	case *ast.LiteralNode, *ast.VariableNode:
		return true
	case *ast.UnaryNode:
		return true
	case *ast.GroupingNode:
		return true
	default:
		return false
	}
}
