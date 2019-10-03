package main

/*
	GO generics workarounds:

	1. Find a well-fitting interface
	2. Use multiple functions
	3. Use the empty interface
*/

// Visitor is the interface expression visitor should implement
type Visitor interface {
	visitBinaryExpr(expr BinaryExpr) interface{}
	visitGroupingExpr(expr GroupingExpr) interface{}
	visitLiteralExpr(expr LiteralExpr) interface{}
	visitUnaryExpr(expr UnaryExpr) interface{}
	visitSequenceExpr(expr SequenceExpr) interface{}
	visitConditionExpr(expr ConditionExpr) interface{}
}

// StmtVisitor is the interface statements visitor should implement
type StmtVisitor interface {
	// unlike expressions, statements produce no values,
	// so the return type of the visit methods is void, not Object
	visitExpressionStmt(stmt ExpressionStmt)
	visitPrintStmt(stmt PrintStmt)
}