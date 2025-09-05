package lib

type ASTNode interface {
	Accept(visitor Visitor, context *EvaluationContext) any
	Equals(node ASTNode) bool
}

type GroupingNode struct {
	Expr ASTNode
}

func (g GroupingNode) Accept(visitor Visitor, context *EvaluationContext) any {
	return visitor.VisitGrouping(&g, context)
}

func (g GroupingNode) Equals(node ASTNode) bool {
	same, ok := node.(GroupingNode)
	return ok && same.Expr.Equals(g.Expr)
}

type LiteralNode struct {
	Value bool
}

func (l LiteralNode) Accept(visitor Visitor, context *EvaluationContext) any {
	return visitor.VisitLiteral(&l, context)
}

func (l LiteralNode) Equals(node ASTNode) bool {
	same, ok := node.(LiteralNode)
	return ok && same.Value == l.Value
}

type VariableNode struct {
	Name string
}

func (v VariableNode) Accept(visitor Visitor, context *EvaluationContext) any {
	return visitor.VisitVariable(&v, context)
}

func (v VariableNode) Equals(node ASTNode) bool {
	same, ok := node.(VariableNode)
	return ok && same.Name == v.Name
}

type BinaryNode struct {
	Operator    TokenType
	Left, Right ASTNode
}

func (b BinaryNode) Accept(visitor Visitor, context *EvaluationContext) any {
	return visitor.VisitBinary(&b, context)
}

func (b BinaryNode) Equals(node ASTNode) bool {
	same, ok := node.(BinaryNode)
	return ok && b.Operator == same.Operator && b.Left.Equals(same.Left) && b.Right.Equals(same.Right)
}

type UnaryNode struct {
	Operator TokenType
	Operand  ASTNode
}

func (u UnaryNode) Accept(visitor Visitor, context *EvaluationContext) any {
	return visitor.VisitUnary(&u, context)
}

func (u UnaryNode) Equals(node ASTNode) bool {
	same, ok := node.(UnaryNode)
	return ok && u.Operator == same.Operator && u.Operand.Equals(same.Operand)
}

type PredicateNode struct {
	Name string
	Body any
}

func (p PredicateNode) Accept(visitor Visitor, context *EvaluationContext) any {
	return visitor.VisitPredicate(&p, context)
}

func (p PredicateNode) Equals(node ASTNode) bool {
	return true
}

type QuantifierNode struct {
	Type     TokenType
	Variable string
	Domain   any
}

func (q QuantifierNode) Accept(visitor Visitor, context *EvaluationContext) any {
	return visitor.VisitQuantifier(&q, context)
}

func (q QuantifierNode) Equals(node ASTNode) bool {
	return true
}
