package visitor

import (
	"logicka/lib/ast"
	"logicka/lib/lexer"
)

type TruthTableEntry struct {
	Result    bool
	Variables []TruthTableVariable
}

type TruthTableVariable struct {
	Name  string
	Value bool
}
type BooleanSolver struct {
	context *EvaluationContext
}

func NewBooleanSolver(context *EvaluationContext) *BooleanSolver {
	return &BooleanSolver{context: context}
}

func (s *BooleanSolver) Visit(node ast.ASTNode) []TruthTableEntry {
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

func (s *BooleanSolver) visitGrouping(node *ast.GroupingNode) []TruthTableEntry {
	return s.Visit(node.Expr)
}

func (s *BooleanSolver) visitLiteral(node *ast.LiteralNode) []TruthTableEntry {
	return []TruthTableEntry{
		{Result: node.Value, Variables: []TruthTableVariable{}},
	}
}

func (s *BooleanSolver) visitVariable(node *ast.VariableNode) []TruthTableEntry {
	if val, ok := s.context.Variables[node.Name]; ok {
		return []TruthTableEntry{
			{Result: val, Variables: []TruthTableVariable{
				{Name: node.Name, Value: val},
			}},
		}
	}
	return []TruthTableEntry{
		{Result: true, Variables: []TruthTableVariable{
			{Name: node.Name, Value: true},
		}},
		{Result: false, Variables: []TruthTableVariable{
			{Name: node.Name, Value: false},
		}},
	}
}

func (s *BooleanSolver) visitBinary(node *ast.BinaryNode) []TruthTableEntry {
	left := s.Visit(node.Left)
	right := s.Visit(node.Right)

	res := make([]TruthTableEntry, 0)

	for _, l := range left {
		for _, r := range right {
			merged := mergeVariables(l.Variables, r.Variables)
			if merged == nil {
				continue
			}
			switch node.Operator {
			case lexer.IMPL:
				res = append(res, TruthTableEntry{
					Result:    !l.Result || r.Result,
					Variables: merged,
				})
			case lexer.EQUIV:
				res = append(res, TruthTableEntry{
					Result:    l.Result == r.Result,
					Variables: merged,
				})
			case lexer.CONJ:
				res = append(res, TruthTableEntry{
					Result:    l.Result && r.Result,
					Variables: merged,
				})
			case lexer.DISJ:
				res = append(res, TruthTableEntry{
					Result:    l.Result || r.Result,
					Variables: merged,
				})
			default:
				panic("unhandled default case")
			}
		}
	}

	return res
}

func (s *BooleanSolver) visitUnary(node *ast.UnaryNode) []TruthTableEntry {
	operands := s.Visit(node.Operand)

	res := make([]TruthTableEntry, 0)

	for _, o := range operands {
		switch node.Operator {
		case lexer.NEG:
			res = append(res, TruthTableEntry{
				Result:    !o.Result,
				Variables: o.Variables,
			})
		default:
			panic("unhandled default case")
		}
	}

	return res
}

func (s *BooleanSolver) visitPredicate(node *ast.PredicateNode) []TruthTableEntry {
	// TODO implement me
	panic("implement me")
}

func (s *BooleanSolver) visitQuantifier(node *ast.QuantifierNode) []TruthTableEntry {
	// TODO implement me
	panic("implement me")
}

func mergeVariables(left, right []TruthTableVariable) []TruthTableVariable {
	varMap := make(map[string]bool)

	for _, v := range left {
		varMap[v.Name] = v.Value
	}

	for _, v := range right {
		if val, ok := varMap[v.Name]; ok {
			if val != v.Value {
				return nil
			}
		} else {
			varMap[v.Name] = v.Value
		}
	}

	result := make([]TruthTableVariable, 0, len(varMap))
	for name, value := range varMap {
		result = append(result, TruthTableVariable{Name: name, Value: value})
	}
	return result
}
