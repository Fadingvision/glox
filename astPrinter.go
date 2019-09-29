package main

import (
	"fmt"
	"io"
	"strings"
)

// AstPrinter prints expression AST to string
type AstPrinter struct{}

func (v AstPrinter) print(expr Expr, out io.Writer) {
	str, _ := expr.accept(v).(string)
	out.Write([]byte(str))
}
func (v AstPrinter) visitBinaryExpr(expr BinaryExpr) interface{} {
	return expr.left.accept(v).(string) + " " + expr.operator.literal + " " + expr.right.accept(v).(string)
}
func (v AstPrinter) visitGroupingExpr(expr GroupingExpr) interface{} {
	var b strings.Builder
	b.WriteString("(")
	b.WriteString(expr.expression.accept(v).(string))
	b.WriteString(")")
	return b.String()
}
func (v AstPrinter) visitLiteralExpr(expr LiteralExpr) interface{} {
	return fmt.Sprint(expr.value)
}
func (v AstPrinter) visitUnaryExpr(expr UnaryExpr) interface{} {
	return expr.operator.literal + expr.right.accept(v).(string)
}
func (v AstPrinter) visitConditionExpr(expr ConditionExpr) interface{} {
	var b strings.Builder
	b.WriteString("( ")
	b.WriteString(expr.test.accept(v).(string))
	b.WriteString(" ? ")
	b.WriteString(expr.consequent.accept(v).(string))
	b.WriteString(" : ")
	b.WriteString(expr.alternate.accept(v).(string))
	b.WriteString(" )")
	return b.String()
}
func (v AstPrinter) visitSequenceExpr(expr SequenceExpr) interface{} {
	var b strings.Builder
	for index, expression := range expr.exprs {
		b.WriteString(expression.accept(v).(string))
		if index < len(expr.exprs)-1 {
			b.WriteString(", ")
		}
	}
	return b.String()
}
