package main

// stack based on slice
type scopes []map[string]bool

func (s scopes) peek() map[string]bool {
	return s[len(s)-1]
}

func (s scopes) isEmpty() bool {
	return len(s) == 0
}

type functionType int
type classType int

const (
	// NONE means we are not in a function statement body
	NONE functionType = iota
	// FUNCTION means we are in a function statement body
	FUNCTION
	// METHOD means we are in a class method statement body
	METHOD
	// INITIALIZER means we are in a class init method statement body
	INITIALIZER
)

const (
	// NoneClass means we are not in a class statement body
	NONECLASS classType = iota
	// CLASS means we are in a class statement body
	INCLASS
)

/*
	Resolver do a single walk(better O(n) ) over the ast tree to resolve all of the variables it contains
	This layer is often used for Semantic Analysis, for example:

	- Type check

	- Optimization

	- semantic error report

*/
type Resolver struct {
	lox         *Lox
	interpreter *Interpreter
	scopes      scopes
	// currentFunction shows if we are visiting a function statement
	currentFunction functionType
	// currentClass shows if we are visiting a class statement
	currentClass classType
}

// NewResolver create a Resolver instance
func NewResolver(l *Lox, interpreter *Interpreter) Resolver {
	return Resolver{
		l,
		interpreter,
		make(scopes, 0),
		NONE,
		NONECLASS,
	}
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
	r.resolveFunction(stmt, FUNCTION)
}

func (r Resolver) visitFunExpr(expr FunExpr) interface{} {
	parentFunctionType := r.currentFunction
	r.currentFunction = FUNCTION

	// like the blockStatement
	r.scopes = r.beginScope()
	for _, param := range expr.params {
		r.declare(param)
		r.define(param)
	}
	r.resolveBody(expr.body.statements)
	r.endScope()

	r.currentFunction = parentFunctionType
	return nil
}

func (r Resolver) resolveFunction(stmt FunStmt, ftype functionType) {
	parentFunctionType := r.currentFunction
	r.currentFunction = ftype

	// like the blockStatement
	r.scopes = r.beginScope()
	for _, param := range stmt.params {
		r.declare(param)
		r.define(param)
	}
	r.resolveBody(stmt.body.statements)
	r.endScope()

	// In case there is nest function,
	// after function has been resolved, we reset current function type to previous
	r.currentFunction = parentFunctionType
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

func (r Resolver) visitClassStmt(stmt ClassStmt) {
	parentClassType := r.currentClass
	r.currentClass = INCLASS

	r.declare(stmt.name)
	r.define(stmt.name)

	if stmt.super != nil {
		if stmt.super.name.literal == stmt.name.literal {
			r.lox.errorReporter.error(ParseError{
				stmt.super.name,
				"Cannot Access '" + stmt.super.name.literal + "' before initialization",
			})
		}
		r.resolveExpr(*stmt.super)
	}

	// No this or super in static methods
	for _, fun := range stmt.staticMethods {
		functionType := METHOD
		r.resolveFunction(fun, functionType)
	}

	// create a scope to store `super`
	if stmt.super != nil {
		r.scopes = r.beginScope()
		r.scopes.peek()["super"] = true
	}

	// create a scope to store `this`
	r.scopes = r.beginScope()
	r.scopes.peek()["this"] = true

	for _, fun := range stmt.methods {
		functionType := METHOD
		if fun.name.literal == "init" {
			functionType = INITIALIZER
		}
		r.resolveFunction(fun, functionType)
	}
	r.endScope()

	if stmt.super != nil {
		r.endScope()
	}

	r.currentClass = parentClassType
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

	scope := r.scopes.peek()

	if _, ok := scope[name.literal]; ok {
		r.lox.errorReporter.error(ParseError{
			name,
			"a redeclared in this block",
		})
	}

	scope[name.literal] = false
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

func (r Resolver) visitThisExpr(expr ThisExpr) interface{} {
	if r.currentClass == NONECLASS {
		r.lox.errorReporter.error(ParseError{
			expr.keyword,
			"Illegal 'this'",
		})
	}

	// resolve this like any other local variable.
	r.resolveLocal(expr, expr.keyword)
	return nil
}

func (r Resolver) visitSuperExpr(expr SuperExpr) interface{} {
	// give the distance about super to interpreter
	r.resolveLocal(expr, expr.keyword)
	return nil
}

func (r Resolver) visitIdentifierExpr(expr IdentifierExpr) interface{} {
	// lox Cannot read local variable in its own initializer.
	if !r.scopes.isEmpty() {
		if isDefined, ok := r.scopes.peek()[expr.name.literal]; ok && isDefined == false {
			r.lox.errorReporter.error(ParseError{
				expr.name,
				"Cannot read local variable in its own initializer.",
			})
		}
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
	// with `r.currentFunction`, we can know if we are in function body when we met a return statement
	if r.currentFunction == NONE {
		r.lox.errorReporter.error(ParseError{
			stmt.keyword,
			"Illegal return statement",
		})
	}

	if stmt.value != nil {
		if r.currentFunction == INITIALIZER {
			r.lox.errorReporter.error(ParseError{
				stmt.keyword,
				"Illegal return statement, Cannot return a value form a initializer",
			})
		}
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

func (r Resolver) visitSetExpr(expr SetExpr) interface{} {
	r.resolveExpr(expr.value)
	r.resolveExpr(expr.object)
	return nil
}

func (r Resolver) visitGetExpr(expr GetExpr) interface{} {
	r.resolveExpr(expr.object)
	return nil
}
