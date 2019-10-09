package main

type scopes []map[string]bool

func (s scopes) peek() map[string]bool {
	return s[len(s)-1]
}

func (s scopes) isEmpty() bool {
	return len(s) == 0
}

// type RuntimeError struct {
// 	token Token
// 	msg   string
// }

// func (e RuntimeError) Error() string {
// 	return fmt.Sprintf("[GLOX] RuntimeError: Line %d, Cloumn %d, %s", e.token.line, e.token.column, e.msg)
// }

// Resolver is for evaluating codes
type Resolver struct {
	lox         *Lox
	interpreter *Interpreter
	scopes      scopes
}

func (r Resolver) resolveStmt(stmt Stmt) {
	stmt.accept(r)
}

func (r Resolver) resolveBody(stmts []Stmt) {
	for _, statement := range stmts {
		r.resolveStmt(statement)
	}
}

func (r Resolver) resolveExpr(expr Expr) {
	expr.accept(r)
}

func (r Resolver) visitFunStmt(stmt FunStmt) {
	r.declare(stmt.name)
	r.define(stmt.name)
	r.resolveFunction(stmt)
}

func (r Resolver) visitFunExpr(expr FunExpr) interface{} {
	// like the blockStatement
	r.scopes = r.beginScope()
	for _, param := range expr.params {
		r.declare(param)
		r.define(param)
	}
	r.resolveBody(expr.body.statements)
	r.endScope()
	return nil
}

func (r Resolver) resolveFunction(stmt FunStmt) {
	// like the blockStatement
	r.scopes = r.beginScope()
	for _, param := range stmt.params {
		r.declare(param)
		r.define(param)
	}
	r.resolveBody(stmt.body.statements)
	r.endScope()
}

func (r Resolver) beginScope() scopes {
	// add A new scope for block
	return append(r.scopes, make(map[string]bool, 0))
}

func (r Resolver) endScope() {
	// pop
	r.scopes = r.scopes[:len(r.scopes)-1]
}

func (r Resolver) visitBlockStmt(stmt BlockStmt) {
	r.scopes = r.beginScope()
	r.resolveBody(stmt.statements)
	r.endScope()
}

func (r Resolver) visitVarStmt(stmt VarStmt) {
	r.declare(stmt.name)
	if stmt.init != nil {
		r.resolveExpr(stmt.init)
	}
	r.define(stmt.name)
}

func (r Resolver) declare(name Token) {
	// skip the global vars
	if len(r.scopes) == 0 {
		return
	}
	r.scopes.peek()[name.literal] = false
}

func (r Resolver) define(name Token) {
	// skip the global vars
	if len(r.scopes) == 0 {
		return
	}
	r.scopes.peek()[name.literal] = true
}

func (r Resolver) visitAssignExpr(expr AssignExpr) interface{} {
	r.resolveExpr(expr.right)
	r.resolveLocal(expr, expr.left)
	return nil
}

func (r Resolver) visitIdentifierExpr(expr IdentifierExpr) interface{} {
	// lox Cannot read local variable in its own initializer.
	if !r.scopes.isEmpty() && r.scopes.peek()[expr.name.literal] == false {
		r.lox.errorReporter.error(ParseError{
			expr.name,
			"Cannot read local variable in its own initializer.",
		})
	}

	r.resolveLocal(expr, expr.name)
	return nil
}

func (r Resolver) resolveLocal(expr Expr, name Token) {
	// When accessing a local variable
	// We calculate the distance from the scope the var is accessed to the scope the var is declared.
	// Then give the constant distance to interpreter
	for i := len(r.scopes) - 1; i >= 0; i-- {
		if _, ok := r.scopes[i][name.literal]; ok {
			distance := len(r.scopes) - i - 1
			r.interpreter.resolve(expr, distance)
			return
		}
	}

	// Not found. it's global var or a undefined error
}

/*
	The rest of statements and expressions won't directly get us any scope and variables.
	So we just pass it down.
*/
func (r Resolver) visitExpressionStmt(stmt ExpressionStmt) {
	r.resolveExpr(stmt.expression)
}

func (r Resolver) visitIfStmt(stmt IfStmt) {
	// unlike interpreter, we need record every branch of ifstatement,
	// since every branch may produces Var and function Declaration.
	// we need record all their distance.
	r.resolveExpr(stmt.condition)
	r.resolveStmt(stmt.consequent)
	if stmt.alternate != nil {
		r.resolveStmt(stmt.alternate)
	}
}

func (r Resolver) visitWhileStmt(stmt WhileStmt) {
	r.resolveExpr(stmt.condition)
	r.resolveStmt(stmt.body)
}

func (r Resolver) visitPrintStmt(stmt PrintStmt) {
	r.resolveExpr(stmt.expression)
}

func (r Resolver) visitReturnStmt(stmt ReturnStmt) {
	if stmt.value != nil {
		r.resolveExpr(stmt.value)
	}
}

func (r Resolver) visitGroupingExpr(expr GroupingExpr) interface{} {
	r.resolveExpr(expr.expression)
	return nil
}
func (r Resolver) visitLiteralExpr(expr LiteralExpr) interface{} {
	return nil
}
func (r Resolver) visitLogicalExpr(expr LogicalExpr) interface{} {
	r.resolveExpr(expr.left)
	r.resolveExpr(expr.right)
	return nil
}

func (r Resolver) visitBinaryExpr(expr BinaryExpr) interface{} {
	r.resolveExpr(expr.left)
	r.resolveExpr(expr.right)
	return nil
}

func (r Resolver) visitCallExpr(expr CallExpr) interface{} {
	r.resolveExpr(expr.callee)
	for _, arg := range expr.arguments {
		r.resolveExpr(arg)
	}
	return nil
}

func (r Resolver) visitUnaryExpr(expr UnaryExpr) interface{} {
	r.resolveExpr(expr.right)
	return nil
}

func (r Resolver) visitConditionExpr(expr ConditionExpr) interface{} {
	r.resolveExpr(expr.test)
	r.resolveExpr(expr.consequent)
	r.resolveExpr(expr.alternate)
	return nil
}

func (r Resolver) visitSequenceExpr(expr SequenceExpr) interface{} {
	for _, expression := range expr.exprs {
		r.resolveExpr(expression)
	}
	return nil
}
