package lib

type ASTNode interface {
	Accept(visitor Visitor, context *EvaluationContext) any
}

type GroupingNode struct {
	Expr ASTNode
}

func (g GroupingNode) Accept(visitor Visitor, context *EvaluationContext) any {
	return visitor.VisitGrouping(&g, context)
}

type VariableNode struct {
	Name string
}

func (v VariableNode) Accept(visitor Visitor, context *EvaluationContext) any {
	return visitor.VisitVariable(&v, context)
}

type BinaryNode struct {
	Operator    TokenType
	Left, Right ASTNode
}

func (b BinaryNode) Accept(visitor Visitor, context *EvaluationContext) any {
	return visitor.VisitBinary(&b, context)
}

type UnaryNode struct {
	Operator TokenType
	Operand  ASTNode
}

func (u UnaryNode) Accept(visitor Visitor, context *EvaluationContext) any {
	return visitor.VisitUnary(&u, context)
}

type PredicateNode struct {
	Name string
	Body any
}

func (p PredicateNode) Accept(visitor Visitor, context *EvaluationContext) any {
	return visitor.VisitPredicate(&p, context)
}

type QuantifierNode struct {
	Type     TokenType
	Variable string
	Domain   any
}

func (q QuantifierNode) Accept(visitor Visitor, context *EvaluationContext) any {
	return visitor.VisitQuantifier(&q, context)
}
