package advanced

import (
	"logicka/lib/ast"
	"logicka/lib/lexer"
	"logicka/lib/simplification/rules/base"
)

type ResolutionRule struct {
	base.BaseRule
}

func NewResolutionRule() *ResolutionRule {
	return &ResolutionRule{
		BaseRule: *base.NewBaseRule("Resolution law"),
	}
}

func (r *ResolutionRule) CanApply(node ast.ASTNode) bool {
	binary, ok := node.(*ast.BinaryNode)
	if !ok {
		return false
	}

	if binary.Operator != lexer.CONJ {
		return false
	}

	_, lok := binary.Left.(*ast.GroupingNode)
	_, rok := binary.Right.(*ast.GroupingNode)

	return lok && rok
}

func (r *ResolutionRule) Apply(node ast.ASTNode) (ast.ASTNode, error) {
	binary := node.(*ast.BinaryNode)
	lGroup := binary.Left.(*ast.GroupingNode)
	rGroup := binary.Right.(*ast.GroupingNode)

	if leftBinary, lok := lGroup.Expr.(*ast.BinaryNode); lok && leftBinary.Operator == lexer.DISJ {
		if rightBinary, rok := rGroup.Expr.(*ast.BinaryNode); rok && rightBinary.Operator == lexer.DISJ {
			if leftBinary.Operator == rightBinary.Operator {
				return node, nil
			}
			if result, resolved := r.applyBinaryResolution(leftBinary, rightBinary); resolved {
				return result, nil
			}
		}
	}

	if leftChain, lok := lGroup.Expr.(*ast.ChainNode); lok && leftChain.Operator == lexer.DISJ {
		if rightChain, rok := rGroup.Expr.(*ast.ChainNode); rok && rightChain.Operator == lexer.DISJ {
			if result, resolved := r.applyChainResolution(leftChain, rightChain); resolved {
				return result, nil
			}
		}
	}

	if leftBinary, lok := lGroup.Expr.(*ast.BinaryNode); lok && leftBinary.Operator == lexer.DISJ {
		if rightChain, rok := rGroup.Expr.(*ast.ChainNode); rok && rightChain.Operator == lexer.DISJ {
			if result, resolved := r.applyMixedResolution(leftBinary, rightChain); resolved {
				return result, nil
			}
		}
	}

	if leftChain, lok := lGroup.Expr.(*ast.ChainNode); lok && leftChain.Operator == lexer.DISJ {
		if rightBinary, rok := rGroup.Expr.(*ast.BinaryNode); rok && rightBinary.Operator == lexer.DISJ {
			if result, resolved := r.applyMixedResolution(rightBinary, leftChain); resolved {
				return result, nil
			}
		}
	}

	return node, nil
}

func (r *ResolutionRule) applyBinaryResolution(left, right *ast.BinaryNode) (ast.ASTNode, bool) {
	if ast.IsNegationOf(left.Left, right.Left) || ast.IsNegationOf(right.Left, left.Left) {
		return r.createResult(left.Right, right.Right), true
	}
	if ast.IsNegationOf(left.Left, right.Right) || ast.IsNegationOf(right.Right, left.Left) {
		return r.createResult(left.Right, right.Left), true
	}
	if ast.IsNegationOf(left.Right, right.Left) || ast.IsNegationOf(right.Left, left.Right) {
		return r.createResult(left.Left, right.Right), true
	}
	if ast.IsNegationOf(left.Right, right.Right) || ast.IsNegationOf(right.Right, left.Right) {
		return r.createResult(left.Left, right.Left), true
	}

	return nil, false
}

func (r *ResolutionRule) applyChainResolution(left, right *ast.ChainNode) (ast.ASTNode, bool) {
	for i, leftOp := range left.Operands {
		for j, rightOp := range right.Operands {
			if ast.IsNegationOf(leftOp, rightOp) || ast.IsNegationOf(rightOp, leftOp) {
				leftOps := make([]ast.ASTNode, 0, len(left.Operands)-1)
				leftOps = append(leftOps, left.Operands[:i]...)
				leftOps = append(leftOps, left.Operands[i+1:]...)

				rightOps := make([]ast.ASTNode, 0, len(right.Operands)-1)
				rightOps = append(rightOps, right.Operands[:j]...)
				rightOps = append(rightOps, right.Operands[j+1:]...)

				allOps := append(leftOps, rightOps...)

				return ast.NewGroupingNode(&ast.ChainNode{
					Operator: lexer.DISJ,
					Operands: allOps,
				}), true
			}
		}
	}

	return nil, false
}

func (r *ResolutionRule) applyMixedResolution(binary *ast.BinaryNode, chain *ast.ChainNode) (ast.ASTNode, bool) {
	for i, chainOp := range chain.Operands {
		if ast.IsNegationOf(binary.Left, chainOp) || ast.IsNegationOf(chainOp, binary.Left) {
			chainOps := make([]ast.ASTNode, 0, len(chain.Operands))
			chainOps = append(chainOps, chain.Operands[:i]...)
			chainOps = append(chainOps, chain.Operands[i+1:]...)
			chainOps = append(chainOps, binary.Right)

			return ast.NewGroupingNode(&ast.ChainNode{
				Operator: lexer.DISJ,
				Operands: chainOps,
			}), true
		}

		if ast.IsNegationOf(binary.Right, chainOp) || ast.IsNegationOf(chainOp, binary.Right) {
			chainOps := make([]ast.ASTNode, 0, len(chain.Operands))
			chainOps = append(chainOps, chain.Operands[:i]...)
			chainOps = append(chainOps, chain.Operands[i+1:]...)
			chainOps = append(chainOps, binary.Left)

			if len(chainOps) == 1 {
				return chainOps[0], true
			}

			return ast.NewGroupingNode(&ast.ChainNode{
				Operator: lexer.DISJ,
				Operands: chainOps,
			}), true
		}
	}

	return nil, false
}

func (r *ResolutionRule) createResult(a, b ast.ASTNode) ast.ASTNode {
	return ast.NewGroupingNode(
		ast.NewBinaryNode(lexer.DISJ, a, b),
	)
}
