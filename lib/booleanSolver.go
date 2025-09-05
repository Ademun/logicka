package lib

type TruthTableEntry struct {
	Result    bool
	Variables []TruthTableVariable
}

type TruthTableVariable struct {
	Name  string
	Value bool
}
type BooleanSolver struct {
}

func (b *BooleanSolver) VisitGrouping(node *GroupingNode, context *EvaluationContext) any {
	return node.Expr.Accept(b, context)
}

func (b *BooleanSolver) VisitVariable(node *VariableNode, context *EvaluationContext) any {
	if val, ok := context.Variables[node.Name]; ok {
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

func (b *BooleanSolver) VisitBinary(node *BinaryNode, context *EvaluationContext) any {
	left := node.Left.Accept(b, context).([]TruthTableEntry)
	right := node.Right.Accept(b, context).([]TruthTableEntry)

	res := make([]TruthTableEntry, 0)

	for _, l := range left {
		for _, r := range right {
			merged := mergeVariables(l.Variables, r.Variables)
			if merged == nil {
				continue
			}
			switch node.Operator {
			case IMPL:
				res = append(res, TruthTableEntry{
					Result:    !l.Result || r.Result,
					Variables: merged,
				})
			case EQUIV:
				res = append(res, TruthTableEntry{
					Result:    l.Result == r.Result,
					Variables: merged,
				})
			case CONJ:
				res = append(res, TruthTableEntry{
					Result:    l.Result && r.Result,
					Variables: merged,
				})
			case DISJ:
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

func (b *BooleanSolver) VisitUnary(node *UnaryNode, context *EvaluationContext) any {
	operands := node.Operand.Accept(b, context).([]TruthTableEntry)

	res := make([]TruthTableEntry, 0)

	for _, o := range operands {
		switch node.Operator {
		case NEG:
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

func (b *BooleanSolver) VisitPredicate(node *PredicateNode, context *EvaluationContext) any {
	// TODO implement me
	panic("implement me")
}

func (b *BooleanSolver) VisitQuantifier(node *QuantifierNode, context *EvaluationContext) any {
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
