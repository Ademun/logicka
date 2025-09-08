package ast

import "logicka/lib/lexer"

func IsTrue(node ASTNode) bool {
	literal, ok := node.(*LiteralNode)
	return ok && literal.Value
}

func IsFalse(node ASTNode) bool {
	literal, ok := node.(*LiteralNode)
	return ok && !literal.Value
}

func IsNegation(node ASTNode) bool {
	neg, ok := node.(*UnaryNode)
	return ok && neg.Operator == lexer.NEG
}

func IsNegationOfSame(left, right ASTNode) bool {
	neg, ok := left.(*UnaryNode)
	return ok && neg.Operator == lexer.NEG && neg.Operand.Equals(right)
}

func IsConjunctionWith(a, b ASTNode) bool {
	bin, ok := a.(*BinaryNode)
	if !ok || bin.Operator != lexer.CONJ {
		return false
	}
	return b.Equals(bin.Left) || b.Equals(bin.Right)
}
