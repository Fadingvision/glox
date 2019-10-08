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
	visitLogicalExpr(expr LogicalExpr) interface{}
	visitGroupingExpr(expr GroupingExpr) interface{}
	visitLiteralExpr(expr LiteralExpr) interface{}
	visitUnaryExpr(expr UnaryExpr) interface{}
	visitSequenceExpr(expr SequenceExpr) interface{}
	visitConditionExpr(expr ConditionExpr) interface{}
	visitAssignExpr(expr AssignExpr) interface{}
	visitIdentifierExpr(expr IdentifierExpr) interface{}
	visitCallExpr(expr CallExpr) interface{}
	visitFunExpr(expr FunExpr) interface{}
}

// StmtVisitor is the interface statements visitor should implement
type StmtVisitor interface {
	// unlike expressions, statements produce no values,
	// so the return type of the visit methods is void, not Object
	visitExpressionStmt(stmt ExpressionStmt)
	visitPrintStmt(stmt PrintStmt)
	visitFunStmt(stmt FunStmt)
	visitReturnStmt(stmt ReturnStmt)
	visitVarStmt(stmt VarStmt)
	visitBlockStmt(stmt BlockStmt)
	visitIfStmt(stmt IfStmt)
	visitWhileStmt(stmt WhileStmt)
}
