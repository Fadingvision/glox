package main

import (
	"html/template"
	"os"
	"strings"
)

func init() {
	generateAst("Expr", []string{
		"BinaryExpr   : left Expr,operator Token,right Expr",
		"GroupingExpr : expression Expr",
		"LiteralExpr  : value interface{}",
		"UnaryExpr    : operator Token,right Expr",
	})
}

const structTemplate = `
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

type class struct {
	Name   string
	Fields []string
}

func generateAst(super string, sub []string) {
	var subClass []class
	for _, item := range sub {
		className := strings.Trim(strings.Split(item, ":")[0], " ")
		fields := strings.Trim(strings.Split(item, ":")[1], " ")
		subClass = append(subClass, class{
			Name:   className,
			Fields: strings.Split(fields, ","),
		})
	}
	t, err := template.New("struct").Parse(structTemplate)
	newFile, err := os.Create("ast.go")
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
