package main

type Stmt interface {
	accept(visitor StmtVisitor)
}

type ExpressionStmt struct {
	expression Expr
}

func (s ExpressionStmt) accept(visitor StmtVisitor) {
	visitor.visitExpressionStmt(s)
}

type PrintStmt struct {
	expression Expr
}

func (s PrintStmt) accept(visitor StmtVisitor) {
	visitor.visitPrintStmt(s)
}
