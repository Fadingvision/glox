package main

type Visitor interface {
	visitBinaryExpr(expr BinaryExpr)
	visitGroupingExpr(expr GroupingExpr)
	visitLiteralExpr(expr LiteralExpr)
	visitUnaryExpr(expr UnaryExpr)
}
