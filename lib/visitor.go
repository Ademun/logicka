package lib

type EvaluationContext struct {
	Variables map[string]bool
}

type Visitor interface {
	VisitGrouping(node *GroupingNode, context *EvaluationContext) any
	VisitLiteral(node *LiteralNode, context *EvaluationContext) any
	VisitVariable(node *VariableNode, context *EvaluationContext) any
	VisitBinary(node *BinaryNode, context *EvaluationContext) any
	VisitUnary(node *UnaryNode, context *EvaluationContext) any
	VisitPredicate(node *PredicateNode, context *EvaluationContext) any
	VisitQuantifier(node *QuantifierNode, context *EvaluationContext) any
}
