package main

type Stmt interface {
	accept(visitor Visitor) interface{}
}

type Expression struct {
	expression Expr
}

func (s Expression) accept(visitor Visitor) interface{} {
	return visitor.visitExpression(s)
}

type Print struct {
	expression Expr
}

func (s Print) accept(visitor Visitor) interface{} {
	return visitor.visitPrint(s)
}
