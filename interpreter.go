package main

import "fmt"

type RuntimeError struct {
	token Token
	msg   string
}

func (e RuntimeError) Error() string {
	return fmt.Sprintf("[GLOX] RuntimeError: Line %d, Cloumn %d, %s", e.token.line, e.token.column, e.msg)
}

// in lox, we treat everything as True except `nil` and `false`
func toBool(val interface{}) bool {
	if val == nil || val == false {
		return false
	}
	return true
}

// Interpreter is for evaluating codes
type Interpreter struct {
	lox *Lox
	env env
}

func (v Interpreter) checkNumberOperands(token Token, exprs ...interface{}) {
	for _, expr := range exprs {
		_, ok := expr.(float64)
		if !ok {
			v.lox.errorReporter.error(RuntimeError{
				token,
				"invalid operation, mismatched types",
			})
			return
		}
	}
}

func (v Interpreter) checkNumberOrStringOperands(token Token, exprs ...interface{}) {
	for _, expr := range exprs {
		_, ok := expr.(float64)
		_, okString := expr.(string)
		if !ok && !okString {
			v.lox.errorReporter.error(RuntimeError{
				token,
				"invalid operation, mismatched types",
			})
			return
		}
	}
}

func (v Interpreter) execute(stmt Stmt) {
	stmt.accept(v)
}

func (v Interpreter) evaluate(expr Expr) interface{} {
	return expr.accept(v)
}

func (v Interpreter) visitExpressionStmt(stmt ExpressionStmt) {
	v.evaluate(stmt.expression)
}

func (v Interpreter) visitIfStmt(stmt IfStmt) {
	if toBool(v.evaluate(stmt.condition)) {
		v.execute(stmt.consequent)
	} else if stmt.alternate != nil {
		v.execute(stmt.alternate)
	}
}

func (v Interpreter) visitWhileStmt(stmt WhileStmt) {
	for toBool(v.evaluate(stmt.condition)) {
		v.execute(stmt.body)
	}
}

func (v Interpreter) visitPrintStmt(stmt PrintStmt) {
	value := v.evaluate(stmt.expression)
	fmt.Println(value)
}

func (v Interpreter) visitBlockStmt(stmt BlockStmt) {
	/*
		This is how we fully support local scope.
		When block statement is called, store the parent scope,
		build a new env, so the block will executed in this new one.
		After these, restore the previous env.
	*/
	parent := v.env
	blockEnv := env{
		values: make(map[string]interface{}, 0),
		parent: &parent,
	}
	v.env = blockEnv
	defer func() {
		v.env = parent
	}()

	for _, statement := range stmt.statements {
		v.execute(statement)
	}
}

func (v Interpreter) visitVarStmt(stmt VarStmt) {
	var value interface{}
	if stmt.init != nil {
		value = v.evaluate(stmt.init)
	}
	// NOTE: Instead of implicitly initializing variables to nil,
	// we also can make it a runtime error to access a variable that has not been initialized or assigned to
	v.env.set(stmt.name.literal, value)
}

func (v Interpreter) visitAssignExpr(stmt AssignExpr) interface{} {
	var value interface{}
	if stmt.right != nil {
		value = v.evaluate(stmt.right)
	}

	_, err := v.env.assign(stmt.left, value)
	if err != nil {
		v.lox.errorReporter.errorWithoutExit(err)
	}
	return value
}
func (v Interpreter) visitLogicalExpr(expr LogicalExpr) interface{} {
	// logic operator is short-circuit
	left := v.evaluate(expr.left)
	if expr.operator.tokentype == OR {
		if toBool(left) {
			return left
		}
	} else if expr.operator.tokentype == AND {
		if !toBool(left) {
			return left
		}
	}

	return v.evaluate(expr.right)
}

func (v Interpreter) visitBinaryExpr(expr BinaryExpr) interface{} {
	left := expr.left.accept(v)
	right := expr.right.accept(v)
	leftFloat, leftFloatOk := left.(float64)
	rightFloat, rightFloatOk := right.(float64)
	leftString, leftOk := left.(string)
	rightString, rightOk := right.(string)
	switch expr.operator.tokentype {
	// The four standard arithmetic operators (+, -, *, /) apply to numbers;
	// + also applies to strings.
	case SLASH:
		v.checkNumberOperands(expr.operator, left, right)
		// in sutiation like dividing a number by zero, we preserve it as go did, which results Infinity
		return leftFloat / rightFloat
	case MINUS:
		v.checkNumberOperands(expr.operator, left, right)
		return leftFloat - rightFloat
	case STAR:
		v.checkNumberOperands(expr.operator, left, right)
		return leftFloat * rightFloat
	case PLUS:
		v.checkNumberOrStringOperands(expr.operator, left, right)
		if leftFloatOk && rightFloatOk {
			return leftFloat + rightFloat
		}
		if leftOk && rightOk {
			return leftString + rightString
		}
		if leftOk && rightFloatOk {
			return leftString + fmt.Sprintf("%g", rightFloat)
		}
		if leftFloatOk && rightOk {
			return fmt.Sprintf("%g", leftFloat) + rightString
		}
	// The ordering operators <, <=, >, and >= apply to operands that are ordered.
	// which in our case is string and number;
	case GREATER:
		v.checkNumberOrStringOperands(expr.operator, left, right)
		if leftFloatOk && rightFloatOk {
			return leftFloat > rightFloat
		}
		if leftOk && rightOk {
			return leftString > rightString
		}
	case GREATER_EQUAL:
		v.checkNumberOrStringOperands(expr.operator, left, right)
		if leftFloatOk && rightFloatOk {
			return leftFloat >= rightFloat
		}
		if leftOk && rightOk {
			return leftString >= rightString
		}
	case LESS:
		v.checkNumberOrStringOperands(expr.operator, left, right)
		if leftFloatOk && rightFloatOk {
			return leftFloat < rightFloat
		}
		if leftOk && rightOk {
			return leftString < rightString
		}
	case LESS_EQUAL:
		v.checkNumberOrStringOperands(expr.operator, left, right)
		if leftFloatOk && rightFloatOk {
			return leftFloat <= rightFloat
		}
		if leftOk && rightOk {
			return leftString <= rightString
		}
	// accroding to `Go Programming Language Specification`:
	// Two interface values are equal
	// if they have identical dynamic types and equal dynamic values
	// or if both have value nil.
	// so we don't need to assert their types
	case BANG_EQUAL:
		return left != right
	case EQUAL_EQUAL:
		return left == right
	}
	return nil
}

func (v Interpreter) visitGroupingExpr(expr GroupingExpr) interface{} {
	return expr.expression.accept(v)
}

func (v Interpreter) visitLiteralExpr(expr LiteralExpr) interface{} {
	return expr.value
}

func (v Interpreter) visitUnaryExpr(expr UnaryExpr) interface{} {
	right := expr.right.accept(v)
	if expr.operator.tokentype == MINUS {
		v.checkNumberOperands(expr.operator, right)
		return -right.(float64)
	}
	if expr.operator.tokentype == BANG {
		return !toBool(right)
	}
	return nil
}

func (v Interpreter) visitIdentifierExpr(expr IdentifierExpr) interface{} {
	value, err := v.env.get(expr.name)
	if err != nil {
		v.lox.errorReporter.errorWithoutExit(err)
	}
	return value
}

func (v Interpreter) visitConditionExpr(expr ConditionExpr) interface{} {
	test := expr.test.accept(v)
	if toBool(test) {
		return expr.consequent.accept(v)
	}
	return expr.alternate.accept(v)
}

func (v Interpreter) visitSequenceExpr(expr SequenceExpr) interface{} {
	// Take the last one of expressions to be the sequence result
	return expr.exprs[len(expr.exprs)-1].accept(v)
}
