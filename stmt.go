
package main

type Stmt interface {
	accept(visitor StmtVisitor)
}


type ExpressionStmt struct {
	
	expression Expr
	
}
func (s ExpressionStmt) accept(visitor StmtVisitor){
	visitor.visitExpressionStmt(s)
}

type PrintStmt struct {
	
	expression Expr
	
}
func (s PrintStmt) accept(visitor StmtVisitor){
	visitor.visitPrintStmt(s)
}

type BlockStmt struct {
	
	statements []Stmt
	
}
func (s BlockStmt) accept(visitor StmtVisitor){
	visitor.visitBlockStmt(s)
}

type VarStmt struct {
	
	name Token
	
	 init Expr
	
}
func (s VarStmt) accept(visitor StmtVisitor){
	visitor.visitVarStmt(s)
}

