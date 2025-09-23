package ast

type GroupingExpr struct {
	Body Expr
}

func NewGroupingExpr(body Expr) *GroupingExpr {
	return &GroupingExpr{Body: body}
}

func (n *GroupingExpr) expr() {}
