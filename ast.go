package main

type Expr interface {
	accept(visitor Visitor) interface{}
}

type BinaryExpr struct {
	left Expr

	operator Token

	right Expr
}

func (s BinaryExpr) accept(visitor Visitor) interface{} {
	return visitor.visitBinaryExpr(s)
}

type SequenceExpr struct {
	exprs []Expr
}

func (s SequenceExpr) accept(visitor Visitor) interface{} {
	return visitor.visitSequenceExpr(s)
}

type ConditionExpr struct {
	test       Expr
	consequent Expr
	alternate  Expr
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
