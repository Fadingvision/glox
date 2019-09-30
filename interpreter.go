package main

// in lox, we treat everything as True except `nil` and `false`
func toBool(val interface{}) bool {
	if val == nil || val == false {
		return false
	}
	return true
}

type Interpreter struct{}

func (v Interpreter) evaluate(expr Expr) interface{} {
	return expr.accept(v)
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
		return left.(float64) / right.(float64)
	case MINUS:
		return left.(float64) - right.(float64)
	case STAR:
		return left.(float64) * right.(float64)
	case PLUS:
		if leftFloatOk && rightFloatOk {
			return leftFloat + rightFloat
		}
		if leftOk && rightOk {
			return leftString + rightString
		}
	// The ordering operators <, <=, >, and >= apply to operands that are ordered.
	// which in our case is string and number;
	case GREATER:
		if leftFloatOk && rightFloatOk {
			return leftFloat > rightFloat
		}
		if leftOk && rightOk {
			return leftString > rightString
		}
	case GREATER_EQUAL:
		if leftFloatOk && rightFloatOk {
			return leftFloat >= rightFloat
		}
		if leftOk && rightOk {
			return leftString >= rightString
		}
	case LESS:
		if leftFloatOk && rightFloatOk {
			return leftFloat < rightFloat
		}
		if leftOk && rightOk {
			return leftString < rightString
		}
	case LESS_EQUAL:
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
		return -right.(float64)
	}
	if expr.operator.tokentype == BANG {
		return !toBool(right)
	}
	return nil
}
func (v Interpreter) visitConditionExpr(expr ConditionExpr) interface{} {
	test := expr.test.accept(v)
	if toBool(test) {
		return expr.consequent.accept(v)
	} else {
		return expr.alternate.accept(v)
	}
}
func (v Interpreter) visitSequenceExpr(expr SequenceExpr) interface{} {
	// Take the last one of expressions to be the sequence result
	return expr.exprs[len(expr.exprs)-1].accept(v)
}
