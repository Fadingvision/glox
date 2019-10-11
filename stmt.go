
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

type VarStmt    	 struct {
	
	name Token
	
	 init Expr
	
}
func (s VarStmt    	) accept(visitor StmtVisitor){
	visitor.visitVarStmt    	(s)
}

type ClassStmt struct {
	
	name Token
	
	 super *IdentifierExpr
	
	 methods []FunStmt
	
	 staticMethods []FunStmt
	
}
func (s ClassStmt) accept(visitor StmtVisitor){
	visitor.visitClassStmt(s)
}

type ReturnStmt struct {
	
	keyword Token
	
	 value Expr
	
}
func (s ReturnStmt) accept(visitor StmtVisitor){
	visitor.visitReturnStmt(s)
}

type FunStmt    	 struct {
	
	name Token
	
	 params []Token
	
	 body BlockStmt
	
}
func (s FunStmt    	) accept(visitor StmtVisitor){
	visitor.visitFunStmt    	(s)
}

type IfStmt    		 struct {
	
	condition Expr
	
	 consequent Stmt
	
	 alternate Stmt
	
}
func (s IfStmt    		) accept(visitor StmtVisitor){
	visitor.visitIfStmt    		(s)
}

type WhileStmt struct {
	
	condition Expr
	
	 body Stmt
	
}
func (s WhileStmt) accept(visitor StmtVisitor){
	visitor.visitWhileStmt(s)
}

