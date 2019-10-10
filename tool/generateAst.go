package main

import (
	"os"
	"strings"
	"text/template"
)

func main() {
	generateAst("Expr", []string{
		"AssignExpr   : left Token,right Expr",
		"BinaryExpr   : left Expr,operator Token,right Expr",
		"LogicalExpr   : left Expr,operator Token,right Expr",
		"SequenceExpr   : exprs []Expr",
		"ConditionExpr   : test Expr,consequent Expr,alternate Expr",
		"GroupingExpr : expression Expr",
		"LiteralExpr  : value interface{}",
		"UnaryExpr    : operator Token,right Expr",
		"CallExpr    : callee Expr,paren Token,arguments []Expr",
		"GetExpr    : object Expr,name Token",
		"SetExpr    : object Expr,name Token,value Expr",
		"ThisExpr    : keyword Token",
		"IdentifierExpr    : name Token",
		"FunExpr    : params []Token, body BlockStmt",
	}, "expr.go", exprTemplate)

	generateAst("Stmt", []string{
		"ExpressionStmt   : expression Expr",
		"PrintStmt    : expression Expr",
		"BlockStmt    : statements []Stmt",
		"VarStmt    	: name Token, init Expr",
		"ClassStmt    : name Token, methods []FunStmt",
		"ReturnStmt   : keyword Token, value Expr",
		"FunStmt    	: name Token, params []Token, body BlockStmt",
		"IfStmt    		: condition Expr, consequent Stmt, alternate Stmt",
		"WhileStmt    : condition Expr, body Stmt",
	}, "stmt.go", stmtTemplate)
}

const exprTemplate = `
package main

type {{.Super}} interface {
	accept(visitor Visitor) interface{}
}

{{ range $i, $v := .Sub }}
type {{ $v.Name }} struct {
	{{ range $i, $v2 := $v.Fields }}
	{{ $v2 }}
	{{ end }}
}
func (s {{$v.Name}}) accept(visitor Visitor) interface{} {
	return visitor.visit{{ $v.Name }}(s)
}
{{ end }}
`
const stmtTemplate = `
package main

type {{.Super}} interface {
	accept(visitor StmtVisitor)
}

{{ range $i, $v := .Sub }}
type {{ $v.Name }} struct {
	{{ range $i, $v2 := $v.Fields }}
	{{ $v2 }}
	{{ end }}
}
func (s {{$v.Name}}) accept(visitor StmtVisitor){
	visitor.visit{{ $v.Name }}(s)
}
{{ end }}
`

type class struct {
	Name   string
	Fields []string
}

func generateAst(super string, sub []string, filename string, templates string) {
	var subClass []class
	for _, item := range sub {
		className := strings.Trim(strings.Split(item, ":")[0], " ")
		fields := strings.Trim(strings.Split(item, ":")[1], " ")
		subClass = append(subClass, class{
			Name:   className,
			Fields: strings.Split(fields, ","),
		})
	}
	t, err := template.New("struct").Parse(templates)
	newFile, err := os.Create(filename)
	if err == nil {
		t.Execute(newFile, struct {
			Super string
			Sub   []class
		}{
			Super: super,
			Sub:   subClass,
		})
	}
}
