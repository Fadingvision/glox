package main

import (
	"fmt"
	"io"
	"os"
)

// Visitor is the interface all visitor should implement
type Visitor interface {
	visitBinaryExpr(expr BinaryExpr) interface{}
	visitGroupingExpr(expr GroupingExpr) interface{}
	visitLiteralExpr(expr LiteralExpr) interface{}
	visitUnaryExpr(expr UnaryExpr) interface{}
}

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
func (v RPNVisitor) visitGroupingExpr(expr GroupingExpr) interface{} {
	return "1 1 +"
}
func (v RPNVisitor) visitLiteralExpr(expr LiteralExpr) interface{} {
	return fmt.Sprint(expr.value)
}
func (v RPNVisitor) visitUnaryExpr(expr UnaryExpr) interface{} {
	return "1 1 +"
}

func init() {
	expression := BinaryExpr{
		left: LiteralExpr{123},
		operator: Token{
			literal:   "+",
			tokentype: PLUS,
		},
		right: LiteralExpr{456},
	}
	RPNVisitor{}.generate(expression, os.Stdout)
}
