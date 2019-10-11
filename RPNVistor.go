package main

import (
	"fmt"
	"io"
)

// RPNVisitor is the Reverse Polish notation visitor
// which turns normal binary operation to RPN strings
type RPNVisitor struct{}

func (v RPNVisitor) generate(expr Expr, out io.Writer) {
	str, _ := expr.accept(v).(string)
	out.Write([]byte(str))
}
func (v RPNVisitor) visitBinaryExpr(expr BinaryExpr) interface{} {
	return expr.left.accept(v).(string) + " " + expr.right.accept(v).(string) + " " + expr.operator.literal
}
func (v RPNVisitor) visitSequenceExpr(expr SequenceExpr) interface{} {
	return "1 1 +"
}
func (v RPNVisitor) visitConditionExpr(expr ConditionExpr) interface{} {
	return "1 1 +"
}
func (v RPNVisitor) visitGroupingExpr(expr GroupingExpr) interface{} {
	return "1 1 +"
}
func (v RPNVisitor) visitAssignExpr(expr AssignExpr) interface{} {
	return "1 1 +"
}
func (v RPNVisitor) visitIdentifierExpr(expr IdentifierExpr) interface{} {
	return "1 1 +"
}
func (v RPNVisitor) visitLiteralExpr(expr LiteralExpr) interface{} {
	return fmt.Sprint(expr.value)
}
func (v RPNVisitor) visitUnaryExpr(expr UnaryExpr) interface{} {
	return "1 1 +"
}
func (v RPNVisitor) visitLogicalExpr(expr LogicalExpr) interface{} {
	return "1 1 +"
}
func (v RPNVisitor) visitCallExpr(expr CallExpr) interface{} {
	return "1 1 +"
}
func (v RPNVisitor) visitFunExpr(expr FunExpr) interface{} {
	return "1 1 +"
}
func (v RPNVisitor) visitSetExpr(expr SetExpr) interface{} {
	return "1 1 +"
}
func (v RPNVisitor) visitGetExpr(expr GetExpr) interface{} {
	return "1 1 +"
}
func (v RPNVisitor) visitThisExpr(expr ThisExpr) interface{} {
	return "1 1 +"
}
func (v RPNVisitor) visitSuperExpr(expr SuperExpr) interface{} {
	return "1 1 +"
}

// func init() {
// 	expression := BinaryExpr{
// 		left: LiteralExpr{123},
// 		operator: Token{
// 			literal:   "+",
// 			tokentype: PLUS,
// 		},
// 		right: LiteralExpr{456},
// 	}
// 	RPNVisitor{}.generate(expression, os.Stdout)
// }
