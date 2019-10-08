package main

type Expr interface {
	accept(visitor Visitor) interface{}
}

type AssignExpr struct {
	left Token

	right Expr
}

func (s AssignExpr) accept(visitor Visitor) interface{} {
	return visitor.visitAssignExpr(s)
}

type BinaryExpr struct {
	left Expr

	operator Token

	right Expr
}

func (s BinaryExpr) accept(visitor Visitor) interface{} {
	return visitor.visitBinaryExpr(s)
}

type LogicalExpr struct {
	left Expr

	operator Token

	right Expr
}

func (s LogicalExpr) accept(visitor Visitor) interface{} {
	return visitor.visitLogicalExpr(s)
}

type SequenceExpr struct {
	exprs []Expr
}

func (s SequenceExpr) accept(visitor Visitor) interface{} {
	return visitor.visitSequenceExpr(s)
}

type ConditionExpr struct {
	test Expr

	consequent Expr

	alternate Expr
}

func (s ConditionExpr) accept(visitor Visitor) interface{} {
	return visitor.visitConditionExpr(s)
}

type GroupingExpr struct {
	expression Expr
}

func (s GroupingExpr) accept(visitor Visitor) interface{} {
	return visitor.visitGroupingExpr(s)
}

type LiteralExpr struct {
	value interface{}
}

func (s LiteralExpr) accept(visitor Visitor) interface{} {
	return visitor.visitLiteralExpr(s)
}

type UnaryExpr struct {
	operator Token

	right Expr
}

func (s UnaryExpr) accept(visitor Visitor) interface{} {
	return visitor.visitUnaryExpr(s)
}

type CallExpr struct {
	callee Expr

	paren Token

	arguments []Expr
}

func (s CallExpr) accept(visitor Visitor) interface{} {
	return visitor.visitCallExpr(s)
}

type IdentifierExpr struct {
	name Token
}

func (s IdentifierExpr) accept(visitor Visitor) interface{} {
	return visitor.visitIdentifierExpr(s)
}
