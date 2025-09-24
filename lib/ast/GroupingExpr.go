package ast

type GroupingExpr struct {
	Body Expr
}

func NewGroupingExpr(body Expr) *GroupingExpr {
	return &GroupingExpr{Body: body}
}

func (n *GroupingExpr) String() string {
	return "(" + n.Body.String() + ")"
}

func (n *GroupingExpr) Children() []Expr {
	return []Expr{n.Body}
}

func (n *GroupingExpr) expr() {}
