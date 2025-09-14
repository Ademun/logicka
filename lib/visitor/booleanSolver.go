package visitor

import (
	"fmt"
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

func (s *BooleanSolver) Solve(node ast.ASTNode) ([]TruthTableEntry, error) {
	return Accept[[]TruthTableEntry](node, s)
}

func (s *BooleanSolver) VisitGrouping(node *ast.GroupingNode) ([]TruthTableEntry, error) {
	return Accept[[]TruthTableEntry](node.Expr, s)
}

func (s *BooleanSolver) VisitLiteral(node *ast.LiteralNode) ([]TruthTableEntry, error) {
	return []TruthTableEntry{
		{Result: node.Value, Variables: []TruthTableVariable{}},
	}, nil
}

func (s *BooleanSolver) VisitVariable(node *ast.VariableNode) ([]TruthTableEntry, error) {
	if val, ok := s.context.Variables[node.Name]; ok {
		return []TruthTableEntry{
			{Result: val, Variables: []TruthTableVariable{
				{Name: node.Name, Value: val},
			}},
		}, nil
	}
	return []TruthTableEntry{
		{Result: true, Variables: []TruthTableVariable{
			{Name: node.Name, Value: true},
		}},
		{Result: false, Variables: []TruthTableVariable{
			{Name: node.Name, Value: false},
		}},
	}, nil
}

func (s *BooleanSolver) VisitBinary(node *ast.BinaryNode) ([]TruthTableEntry, error) {
	left, err := Accept[[]TruthTableEntry](node.Left, s)
	if err != nil {
		return nil, err
	}
	right, err := Accept[[]TruthTableEntry](node.Right, s)
	if err != nil {
		return nil, err
	}

	res := make([]TruthTableEntry, 0)

	for _, l := range left {
		for _, r := range right {
			merged := mergeVariables(l.Variables, r.Variables)
			if merged == nil {
				continue
			}
			switch op := node.Operator; op {
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
				return nil, fmt.Errorf("unkown operator: %s", op)
			}
		}
	}

	return res, nil
}

func (s *BooleanSolver) VisitChain(node *ast.ChainNode) ([]TruthTableEntry, error) {
	if len(node.Operands) < 2 {
		return nil, fmt.Errorf("chain must have at least 2 operands, got %d", len(node.Operands))
	}

	result, err := Accept[[]TruthTableEntry](node.Operands[0], s)
	if err != nil {
		return nil, err
	}

	for _, operand := range node.Operands[1:] {
		operandEntries, err := Accept[[]TruthTableEntry](operand, s)
		if err != nil {
			return nil, err
		}

		newResult := make([]TruthTableEntry, 0)

		for _, leftEntry := range result {
			for _, rightEntry := range operandEntries {
				merged := mergeVariables(leftEntry.Variables, rightEntry.Variables)
				if merged == nil {
					continue
				}

				var combinedResult bool
				switch node.Operator {
				case lexer.CONJ:
					combinedResult = leftEntry.Result && rightEntry.Result
				case lexer.DISJ:
					combinedResult = leftEntry.Result || rightEntry.Result
				default:
					return nil, fmt.Errorf("unsupported chain operator: %s", node.Operator)
				}

				newResult = append(newResult, TruthTableEntry{
					Result:    combinedResult,
					Variables: merged,
				})
			}
		}

		result = newResult
	}

	return result, nil
}

func (s *BooleanSolver) VisitUnary(node *ast.UnaryNode) ([]TruthTableEntry, error) {
	operands, err := Accept[[]TruthTableEntry](node.Operand, s)
	if err != nil {
		return nil, err
	}

	res := make([]TruthTableEntry, 0)

	for _, o := range operands {
		switch op := node.Operator; op {
		case lexer.NEG:
			res = append(res, TruthTableEntry{
				Result:    !o.Result,
				Variables: o.Variables,
			})
		default:
			return nil, fmt.Errorf("unkown operator: %s", op)
		}
	}

	return res, nil
}

func (s *BooleanSolver) VisitPredicate(node *ast.PredicateNode) ([]TruthTableEntry, error) {
	// TODO implement me
	panic("implement me")
}

func (s *BooleanSolver) VisitQuantifier(node *ast.QuantifierNode) ([]TruthTableEntry, error) {
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
