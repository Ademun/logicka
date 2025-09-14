package visitor

import (
	"errors"
	"fmt"
	"logicka/lib/ast"
	"logicka/lib/lexer"
	"slices"
	"strings"
)

// ErrMaxIterationsExceeded indicates simplification took too many iterations
var ErrMaxIterationsExceeded = errors.New("maximum simplification iterations exceeded")

// TraceStep represents a single step in the simplification trace
type TraceStep struct {
	Rule   string      // Russian name of the applied rule
	Before ast.ASTNode // Expression before the rule application
	After  ast.ASTNode // Expression after the rule application
	Depth  int         // Nesting depth of the operation
}

// SimplificationTrace contains the complete trace of simplification steps
type SimplificationTrace struct {
	Steps []TraceStep
}

func (t SimplificationTrace) String() string {
	if len(t.Steps) == 0 {
		return "Трассировка пуста"
	}

	var result strings.Builder
	result.WriteString("Трассировка упрощения:\n")
	result.WriteString(strings.Repeat("=", 50))
	result.WriteString("\n")

	for i, step := range t.Steps {
		result.WriteString(fmt.Sprintf("%d. %s\n", i+1, step.Rule))
		result.WriteString(fmt.Sprintf("До: %s\n", step.Before.String()))
		result.WriteString(fmt.Sprintf("После: %s\n", step.After.String()))

		if i < len(t.Steps)-1 {
			result.WriteString("\n")
		}
	}

	return result.String()
}

// SimplificationOptions configures the simplification process.
type SimplificationOptions struct {
	MaxIterations int  // Maximum number of simplification iterations
	Debug         bool // Enable debug output
	Trace         bool // Enable Russian trace collection
}

var DefaultSimplificationOptions = SimplificationOptions{
	MaxIterations: 100,
	Debug:         true,
	Trace:         true,
}

// SimplificationResult contains both the simplified expression and trace
type SimplificationResult struct {
	Result ast.ASTNode
	Trace  SimplificationTrace
}

// Simplifier implements logical expression simplification using algebraic rules.
type Simplifier struct {
	options SimplificationOptions
	trace   SimplificationTrace
	depth   int
}

func NewSimplifier(options ...SimplificationOptions) *Simplifier {
	opts := DefaultSimplificationOptions
	if len(options) > 0 {
		opts = options[0]
	}

	return &Simplifier{
		options: opts,
		trace:   SimplificationTrace{Steps: make([]TraceStep, 0)},
		depth:   0,
	}
}

// addTraceStep adds a step to the trace if tracing is enabled
func (s *Simplifier) addTraceStep(rule string, before, after ast.ASTNode) {
	if !s.options.Trace {
		return
	}

	// Only add if there was an actual change
	if !before.Equals(after) {
		s.trace.Steps = append(s.trace.Steps, TraceStep{
			Rule:   rule,
			Before: before,
			After:  after,
			Depth:  s.depth,
		})
	}
}

// Simplify applies algebraic simplification rules until a fixed point is reached.
func (s *Simplifier) Simplify(node ast.ASTNode) (*SimplificationResult, error) {
	if node == nil {
		return nil, errors.New("cannot simplify nil node")
	}

	// Reset trace for new simplification
	if s.options.Trace {
		s.trace = SimplificationTrace{Steps: make([]TraceStep, 0)}
	}

	current := node

	for i := range s.options.MaxIterations {
		s.depth = 0
		next, err := Accept[ast.ASTNode](current, s)
		if err != nil {
			return nil, fmt.Errorf("simplification error at iteration %d: %w", i, err)
		}

		if next.Equals(current) {
			if s.options.Debug {
				fmt.Printf("Simplification converged after %d iterations\n", i)
			}
			return &SimplificationResult{
				Result: next,
				Trace:  s.trace,
			}, nil
		}

		if s.options.Debug {
			fmt.Printf("Iteration %d: %s -> %s\n", i, current.String(), next.String())
		}

		current = next
	}

	return nil, fmt.Errorf("%w: reached %d iterations", ErrMaxIterationsExceeded, s.options.MaxIterations)
}

// Visitor implementation for Simplifier

func (s *Simplifier) VisitGrouping(node *ast.GroupingNode) (ast.ASTNode, error) {
	expr, err := Accept[ast.ASTNode](node.Expr, s)
	if err != nil {
		return nil, err
	}

	// Remove unnecessary grouping if the inner expression is atomic
	if isAtomicExpression(expr) {
		return expr, nil
	}

	result := ast.NewGroupingNode(expr)
	return result, nil
}

func (s *Simplifier) VisitLiteral(node *ast.LiteralNode) (ast.ASTNode, error) {
	return node, nil
}

func (s *Simplifier) VisitVariable(node *ast.VariableNode) (ast.ASTNode, error) {
	return node, nil
}

func (s *Simplifier) VisitBinary(node *ast.BinaryNode) (ast.ASTNode, error) {
	s.depth++
	defer func() { s.depth-- }()

	switch node.Operator {
	case lexer.CONJ:
		return s.simplifyConjunction(node)
	case lexer.DISJ:
		return s.simplifyDisjunction(node)
	case lexer.IMPL:
		return s.simplifyImplication(node)
	case lexer.EQUIV:
		return s.simplifyEquivalence(node)
	default:
		return nil, OperatorError{Operator: node.Operator.String()}
	}
}

// Simplification methods for specific operations

func (s *Simplifier) simplifyConjunction(node *ast.BinaryNode) (ast.ASTNode, error) {
	left, errL := Accept[ast.ASTNode](node.Left, s)
	right, errR := Accept[ast.ASTNode](node.Right, s)
	if err := errors.Join(errL, errR); err != nil {
		return nil, err
	}

	// Identity: A ∧ 1 = A
	if ast.IsTrue(left) {
		s.addTraceStep("Операция с константой (A ∧ 1 = A)", node, right)
		return right, nil
	}
	if ast.IsTrue(right) {
		s.addTraceStep("Операция с константой (1 ∧ A = A)", node, left)
		return left, nil
	}

	// Annihilator: A ∧ 0 = 0
	if ast.IsFalse(left) || ast.IsFalse(right) {
		result := ast.NewLiteralNode(false)
		s.addTraceStep("Операция с константой (A ∧ 0 = 0)", node, result)
		return result, nil
	}

	// Idempotency: A ∧ A = A
	if left.Equals(right) {
		s.addTraceStep("Закон идемпотентности (A ∧ A = A)", node, left)
		return left, nil
	}

	// Contradiction: A ∧ ¬A = 0
	if ast.IsNegationOf(left, right) || ast.IsNegationOf(right, left) {
		result := ast.NewLiteralNode(false)
		s.addTraceStep("Закон исключённого третьего (A ∧ ¬A = 0)", node, result)
		return result, nil
	}

	// Apply absorption laws
	if absorbed, applied, err := s.applyAbsorption(left, right, lexer.CONJ); err != nil {
		return nil, err
	} else if applied {
		s.addTraceStep("Итог распределительного закона", node, absorbed)
		return absorbed, nil
	}

	if absorbed, applied, err := s.applyAbsorption(right, left, lexer.CONJ); err != nil {
		return nil, err
	} else if applied {
		s.addTraceStep("Итог распределительного закона", node, absorbed)
		return absorbed, nil
	}

	// Try to flatten to chain if beneficial
	newNode := ast.NewBinaryNode(lexer.CONJ, left, right)
	if ast.CanFlatten(newNode) {
		result := ast.ChainFromBinary(newNode)
		return result, nil
	}

	s.addTraceStep("Упрощённая конъюнкция", node, newNode)
	return newNode, nil
}

func (s *Simplifier) simplifyDisjunction(node *ast.BinaryNode) (ast.ASTNode, error) {
	left, errL := Accept[ast.ASTNode](node.Left, s)
	right, errR := Accept[ast.ASTNode](node.Right, s)
	if err := errors.Join(errL, errR); err != nil {
		return nil, err
	}

	// Annihilator: A ∨ 1 = 1
	if ast.IsTrue(left) || ast.IsTrue(right) {
		result := ast.NewLiteralNode(true)
		s.addTraceStep("Операция с константой (A ∨ 1 = 1)", node, result)
		return result, nil
	}

	// Identity: A ∨ 0 = A
	if ast.IsFalse(left) {
		s.addTraceStep("Операция с константой (0 ∨ A = A)", node, right)
		return right, nil
	}
	if ast.IsFalse(right) {
		s.addTraceStep("Операция с константой (A ∨ 0 = A)", node, left)
		return left, nil
	}

	// Idempotency: A ∨ A = A
	if left.Equals(right) {
		s.addTraceStep("Закон идемпотентности (A ∨ A = A)", node, left)
		return left, nil
	}

	// Tautology: A ∨ ¬A = 1
	if ast.IsNegationOf(left, right) || ast.IsNegationOf(right, left) {
		result := ast.NewLiteralNode(true)
		s.addTraceStep("Закон исключённого третьего (A ∨ ¬A = 1)", node, result)
		return result, nil
	}

	// Apply absorption laws
	if absorbed, applied, err := s.applyAbsorption(left, right, lexer.DISJ); err != nil {
		return nil, err
	} else if applied {
		s.addTraceStep("Итог распределительного закона", node, absorbed)
		return absorbed, nil
	}

	if absorbed, applied, err := s.applyAbsorption(right, left, lexer.DISJ); err != nil {
		return nil, err
	} else if applied {
		s.addTraceStep("Итог распределительного закона", node, absorbed)
		return absorbed, nil
	}

	// Try to flatten to chain if beneficial
	newNode := ast.NewBinaryNode(lexer.DISJ, left, right)
	if ast.CanFlatten(newNode) {
		result := ast.ChainFromBinary(newNode)
		return result, nil
	}

	s.addTraceStep("Упрощённая дизъюнкция", node, newNode)
	return newNode, nil
}

func (s *Simplifier) simplifyImplication(node *ast.BinaryNode) (ast.ASTNode, error) {
	left, errL := Accept[ast.ASTNode](node.Left, s)
	right, errR := Accept[ast.ASTNode](node.Right, s)
	if err := errors.Join(errL, errR); err != nil {
		return nil, err
	}

	// Convert A → B to ¬A ∨ B
	result := ast.NewBinaryNode(lexer.DISJ,
		ast.NewUnaryNode(lexer.NEG, left),
		right)
	s.addTraceStep("Упрощение импликации (A → B = ¬A ∨ B)", node, result)
	return result, nil
}

func (s *Simplifier) simplifyEquivalence(node *ast.BinaryNode) (ast.ASTNode, error) {
	left, errL := Accept[ast.ASTNode](node.Left, s)
	right, errR := Accept[ast.ASTNode](node.Right, s)
	if err := errors.Join(errL, errR); err != nil {
		return nil, err
	}

	// Convert A ↔ B to (A → B) ∧ (B → A)
	result := ast.NewBinaryNode(lexer.CONJ,
		ast.NewGroupingNode(ast.NewBinaryNode(lexer.IMPL, left, right)),
		ast.NewGroupingNode(ast.NewBinaryNode(lexer.IMPL, right, left)))
	s.addTraceStep("Упрощение эквивалентности (A ↔ B = (A → B) ∧ (B → A))", node, result)
	return result, nil
}

// applyAbsorption applies absorption laws: A ∧ (A ∨ B) = A, A ∨ (A ∧ B) = A
func (s *Simplifier) applyAbsorption(left, right ast.ASTNode, operator lexer.BooleanTokenType) (ast.ASTNode, bool, error) {
	grouping, ok := right.(*ast.GroupingNode)
	if !ok {
		return nil, false, nil
	}

	var targetOp lexer.BooleanTokenType
	if operator == lexer.CONJ {
		targetOp = lexer.DISJ // For conjunction, look for disjunction inside
	} else {
		targetOp = lexer.CONJ // For disjunction, look for conjunction inside
	}

	// Handle both binary and chain expressions inside grouping
	if binary, ok := grouping.Expr.(*ast.BinaryNode); ok && binary.Operator == targetOp {
		if binary.Left.Equals(left) || binary.Right.Equals(left) {
			return left, true, nil
		}
		if binary.Left.Equals(ast.NewUnaryNode(lexer.NEG, left)) {
			return ast.NewBinaryNode(operator, left, ast.NewGroupingNode(binary.Right)), true, nil
		}
		if binary.Right.Equals(ast.NewUnaryNode(lexer.NEG, left)) {
			return ast.NewBinaryNode(operator, left, ast.NewGroupingNode(binary.Left)), true, nil
		}
		if neg, ok := left.(*ast.UnaryNode); ok && neg.Operator == lexer.NEG {
			if binary.Left.Equals(neg.Operand) {
				return ast.NewBinaryNode(operator, left, ast.NewGroupingNode(binary.Right)), true, nil
			}
			if binary.Right.Equals(neg.Operand) {
				return ast.NewBinaryNode(operator, left, ast.NewGroupingNode(binary.Left)), true, nil
			}
		}
	}

	if chain, ok := grouping.Expr.(*ast.ChainNode); ok && chain.Operator == targetOp {
		if chain.Contains(left) {
			return left, true, nil
		}
		if neg := ast.NewUnaryNode(lexer.NEG, left); chain.Contains(neg) {
			chain.Remove(neg)
			return ast.NewBinaryNode(operator, left, ast.NewGroupingNode(chain)), true, nil
		}
		if neg, ok := left.(*ast.UnaryNode); ok && neg.Operator == lexer.NEG {
			if chain.Contains(neg.Operand) {
				chain.Remove(neg.Operand)
				return ast.NewBinaryNode(operator, left, ast.NewGroupingNode(chain)), true, nil
			}
		}
	}

	return nil, false, nil
}

func (s *Simplifier) VisitChain(node *ast.ChainNode) (ast.ASTNode, error) {
	s.depth++
	defer func() { s.depth-- }()

	// First, simplify all operands
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

	newChain := &ast.ChainNode{
		Operator: node.Operator,
		Operands: simplified,
	}

	result, err := s.simplifyChain(newChain)
	if err != nil {
		return nil, err
	}
	s.addTraceStep("Упрощённая цепочка операций", node, result)
	return result, nil
}

func (s *Simplifier) simplifyChain(node *ast.ChainNode) (ast.ASTNode, error) {
	originalNode := &ast.ChainNode{
		Operator: node.Operator,
		Operands: slices.Clone(node.Operands),
	}

	operands := slices.Clone(node.Operands)
	newOperands := make([]ast.ASTNode, 0, len(operands))
Outer:
	for i := len(operands) - 1; i >= 0; i-- {
		one := operands[i]
		if collapsed := s.tryCollapseChain(one, node.Operator); collapsed != nil {
			s.addTraceStep("Свертывание цепочки по константе", originalNode, collapsed)
			return collapsed, nil
		}
		for j := i - 1; j >= 0; j-- {
			other := operands[j]
			if collapsed := s.tryCollapseChain(other, node.Operator); collapsed != nil {
				s.addTraceStep("Свертывание цепочки по константе", originalNode, collapsed)
				return collapsed, nil
			}
			combination := ast.NewBinaryNode(node.Operator, one, other)
			simplifiedCombination, err := Accept[ast.ASTNode](combination, s)
			if err != nil {
				return nil, err
			}
			if !combination.Equals(simplifiedCombination) {
				if collapsed := s.tryCollapseChain(simplifiedCombination, node.Operator); collapsed != nil {
					s.addTraceStep("Свертывание цепочки по константе", originalNode, collapsed)
					return collapsed, nil
				}
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

// tryCollapseChain attempt to collapse conjunction or disjunction chains based on node value
func (s *Simplifier) tryCollapseChain(node ast.ASTNode, operator lexer.BooleanTokenType) ast.ASTNode {
	if ast.IsTrue(node) {
		if operator == lexer.DISJ {
			return ast.NewLiteralNode(true)
		}
		return nil
	}
	if ast.IsFalse(node) {
		if operator == lexer.CONJ {
			return ast.NewLiteralNode(false)
		}
	}
	return nil
}

func (s *Simplifier) VisitUnary(node *ast.UnaryNode) (ast.ASTNode, error) {
	s.depth++
	defer func() { s.depth-- }()

	switch node.Operator {
	case lexer.NEG:
		return s.simplifyNegation(node)
	default:
		return nil, OperatorError{Operator: node.Operator.String()}
	}
}

func (s *Simplifier) simplifyNegation(node *ast.UnaryNode) (ast.ASTNode, error) {
	operand, err := Accept[ast.ASTNode](node.Operand, s)
	if err != nil {
		return nil, err
	}

	// Double negation: ¬¬A = A
	if unary, ok := operand.(*ast.UnaryNode); ok && unary.Operator == lexer.NEG {
		s.addTraceStep("Закон двойного отрицания (¬¬A = A)", node, unary.Operand)
		return unary.Operand, nil
	}

	// Literal negation: ¬1 = 0, ¬0 = 1
	if literal, ok := operand.(*ast.LiteralNode); ok {
		result := ast.NewLiteralNode(!literal.Value)
		s.addTraceStep("Отрицание константы", node, result)
		return result, nil
	}

	// De Morgan's laws
	if grouping, ok := operand.(*ast.GroupingNode); ok {
		if binary, ok := grouping.Expr.(*ast.BinaryNode); ok {
			result, err := s.applyDeMorgan(binary)
			if err != nil {
				return nil, err
			}
			s.addTraceStep("Закон де Моргана", node, result)
			return result, nil
		}

		if chain, ok := grouping.Expr.(*ast.ChainNode); ok {
			result, err := s.applyDeMorganChain(chain)
			if err != nil {
				return nil, err
			}
			s.addTraceStep("Закон де Моргана", node, result)
			return result, nil
		}
	}

	result := ast.NewUnaryNode(lexer.NEG, operand)
	s.addTraceStep("Упрощённое отрицание", node, result)
	return result, nil
}

// applyDeMorgan applies De Morgan's laws: ¬(A ∧ B) = ¬A ∨ ¬B, ¬(A ∨ B) = ¬A ∧ ¬B
func (s *Simplifier) applyDeMorgan(node *ast.BinaryNode) (ast.ASTNode, error) {
	var newOperator lexer.BooleanTokenType

	switch node.Operator {
	case lexer.CONJ:
		newOperator = lexer.DISJ
	case lexer.DISJ:
		newOperator = lexer.CONJ
	default:
		return nil, OperatorError{Operator: node.Operator.String()}
	}

	return ast.NewGroupingNode(ast.NewBinaryNode(newOperator,
		ast.NewUnaryNode(lexer.NEG, node.Left),
		ast.NewUnaryNode(lexer.NEG, node.Right))), nil
}

// applyDeMorganChain applies De Morgan's laws to chain nodes
func (s *Simplifier) applyDeMorganChain(node *ast.ChainNode) (ast.ASTNode, error) {
	var newOperator lexer.BooleanTokenType

	switch node.Operator {
	case lexer.CONJ:
		newOperator = lexer.DISJ
	case lexer.DISJ:
		newOperator = lexer.CONJ
	default:
		return nil, OperatorError{Operator: node.Operator.String()}
	}

	negatedOperands := make([]ast.ASTNode, len(node.Operands))
	for i, operand := range node.Operands {
		negatedOperands[i] = ast.NewUnaryNode(lexer.NEG, operand)
	}

	newChain, err := ast.NewChainNode(newOperator, negatedOperands...)
	if err != nil {
		return nil, err
	}

	return ast.NewGroupingNode(newChain), nil
}

func (s *Simplifier) VisitPredicate(node *ast.PredicateNode) (ast.ASTNode, error) {
	// For now, predicates are not simplified
	// This could be extended with domain-specific rules
	return node, nil
}

func (s *Simplifier) VisitQuantifier(node *ast.QuantifierNode) (ast.ASTNode, error) {
	s.depth++
	defer func() { s.depth-- }()

	// Simplify the domain and body
	domain, errD := Accept[ast.ASTNode](node.Domain, s)
	body, errB := Accept[ast.ASTNode](node.Body, s)

	if err := errors.Join(errD, errB); err != nil {
		return nil, err
	}

	// Basic quantifier simplifications could be added here
	// For example, ∀x: false . P(x) = true, ∃x: false . P(x) = false

	result := ast.NewQuantifierNode(node.Type, node.Variable, domain, body)
	s.addTraceStep("Упрощение квантора", node, result)
	return result, nil
}

// Helper functions

// isAtomicExpression returns true if the expression doesn't need grouping
func isAtomicExpression(node ast.ASTNode) bool {
	switch node.(type) {
	case *ast.LiteralNode, *ast.VariableNode, *ast.UnaryNode, *ast.GroupingNode:
		return true
	default:
		return false
	}
}
