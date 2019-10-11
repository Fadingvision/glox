package main

type TokenType int

const (
	// Single-character tokens.
	LEFT_PAREN TokenType = iota + 1
	RIGHT_PAREN
	LEFT_BRACE
	RIGHT_BRACE
	COMMA
	DOT
	MINUS
	PLUS
	SEMICOLON
	SLASH
	STAR

	QUESTION
	COLON

	// One or two character tokens.
	BANG
	BANG_EQUAL

	EQUAL
	EQUAL_EQUAL

	GREATER
	GREATER_EQUAL

	LESS
	LESS_EQUAL

	// Literals.
	IDENTIFIER
	STRING
	NUMBER

	// Keywords.
	AND
	CLASS
	STATIC
	ELSE
	FALSE
	FUN
	FOR
	IF
	NIL
	OR

	// PRINT

	RETURN
	SUPER
	THIS
	TRUE
	VAR
	WHILE

	EOF
)

var keywords = map[string]TokenType{
	"and":   AND,
	"class": CLASS,
	"else":  ELSE,
	"false": FALSE,
	"for":   FOR,
	"fun":   FUN,
	"if":    IF,
	"nil":   NIL,
	"or":    OR,
	// "print":  PRINT,
	"return": RETURN,
	"super":  SUPER,
	"this":   THIS,
	"true":   TRUE,
	"var":    VAR,
	"while":  WHILE,
	"static": STATIC,
}
