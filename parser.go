package main

import "fmt"

// Parser does two things
// 1. tranform tokens to AST tree
// 2. report synatx error based on our language CFG
type Parser struct {
	tokens  []Token
	current int
	lox     *Lox
}

func (p *Parser) parse() Expr {
	return p.expression()
}

// TokenError implement the std err interface
type ParseError struct {
	token Token
	msg   string
}

func (e ParseError) Error() string {
	// TODO: more elgant error display
	return fmt.Sprintf("[GLOX] SyntaxError: Line %d, Cloumn %d, %s", e.token.line, e.token.column, e.msg)
}

/* utils, kind of self-explained */

// match checks if the next token is matching the specific type
func (p *Parser) match(types ...TokenType) bool {
	for _, tokentype := range types {
		if p.checkType(tokentype) {
			// since we already know the next token
			// we need to mark it as consumed
			p.advance()
			return true
		}
	}

	return false
}

func (p *Parser) checkType(tokentype TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().tokentype == tokentype
}

func (p *Parser) isAtEnd() bool {
	return p.peek().tokentype == EOF
}

// previous returns the most recently consumed token
func (p *Parser) previous() Token {
	return p.tokens[p.current-1]
}

// peek returns the current token we have yet to consume
func (p *Parser) peek() Token {
	return p.tokens[p.current]
}

// advance comsume A token and returns it
func (p *Parser) advance() Token {
	if !p.isAtEnd() {
		p.current += 1
	}
	return p.previous()
}

/*
	Each method for parsing a grammar rule produces a syntax tree
	for that rule and ruturns it to the caller.

	When the body of the rule contains a nonterminal --
	a reference to another rule -- we call that rule's method
*/

// expression     → equality
// equality       → comparison ( ( "!=" | "==" ) comparison )* ;
// comparison     → addition ( ( ">" | ">=" | "<" | "<=" ) addition )* ;
// addition       → multiplication ( ( "-" | "+" ) multiplication )* ;
// multiplication → unary ( ( "/" | "*" ) unary )* ;
// unary          → ( "!" | "-" ) unary
//                | primary ;
// primary        → NUMBER | STRING | "false" | "true" | "nil"
//                | "(" expression ")" ;

// expression rule
func (p *Parser) expression() Expr {
	return p.equality()
}

// equality rule
func (p *Parser) equality() Expr {
	expr := p.comparison()

	// use `for` to implement the `*` wildcard
	for p.match(BANG_EQUAL, EQUAL_EQUAL) {
		operator := p.previous()
		// use comparison to check whether the right opreater is another comparsion
		right := p.comparison()
		// constitute the current expression
		expr = BinaryExpr{expr, operator, right}
	}

	return expr
}

// comparison rule
func (p *Parser) comparison() Expr {
	expr := p.addition()

	for p.match(LESS, LESS_EQUAL, GREATER, GREATER_EQUAL) {
		operator := p.previous()
		right := p.addition()
		expr = BinaryExpr{expr, operator, right}
	}

	return expr
}

// addition rule
func (p *Parser) addition() Expr {
	expr := p.multiplication()

	for p.match(MINUS, PLUS) {
		operator := p.previous()
		right := p.multiplication()
		expr = BinaryExpr{expr, operator, right}
	}

	return expr
}

// multiplication rule
func (p *Parser) multiplication() Expr {
	expr := p.unary()

	for p.match(SLASH, STAR) {
		operator := p.previous()
		right := p.unary()
		expr = BinaryExpr{expr, operator, right}
	}

	return expr
}

// unary rule
func (p *Parser) unary() Expr {
	if p.match(BANG, MINUS) {
		operator := p.previous()
		right := p.unary()
		return UnaryExpr{operator, right}
	}

	return p.primary()
}

// unary rule
func (p *Parser) primary() Expr {
	if p.match(FALSE) {
		return LiteralExpr{false}
	}
	if p.match(TRUE) {
		return LiteralExpr{true}
	}
	if p.match(NIL) {
		return LiteralExpr{nil}
	}
	if p.match(NUMBER, STRING) {
		return LiteralExpr{p.previous().lexeme}
	}
	if p.match(LEFT_PAREN) {
		expr := p.expression()
		if !p.match(RIGHT_PAREN) {
			p.lox.errorReporter.errorWithoutExit(ParseError{
				p.peek(),
				"Unexpected end of input",
			})
		}
		return GroupingExpr{expr}
	}

	p.lox.errorReporter.error(ParseError{
		p.peek(),
		"Invalid or unexpected token",
	})

	return nil
}
