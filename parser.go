package main

import "fmt"

// Parser does two things
// 1. tranform tokens to AST tree
// 2. report syntax error based on our language CFG
type Parser struct {
	tokens  []Token
	current int
	lox     *Lox
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

func (p *Parser) consume(tokentype TokenType, msg string) {
	if p.checkType(tokentype) {
		p.advance()
		return
	}

	p.lox.errorReporter.errorWithoutExit(ParseError{
		p.peek(),
		msg,
	})
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
		p.current++
	}
	return p.previous()
}

func (p *Parser) parse() []Stmt {
	statements := make([]Stmt, 0)
	for !p.isAtEnd() {
		statements = append(statements, p.declaration())
	}

	return statements
}

func (p *Parser) reset() {
	p.current = 0
	p.tokens = nil
}

/*
	Form top to down, each method for parsing a grammar rule produces a syntax tree
	for that rule and ruturns it to the caller.

	When the body of the rule contains a nonterminal --
	a reference to another rule -- we call that rule's method,

	A rule may refer to itself, so this's why it's called recursive descent parser
*/

/* STATEMENTS */

func (p *Parser) declaration() Stmt {
	if p.match(VAR) {
		return p.varDeclaration()
	}

	return p.statement()
}

func (p *Parser) varDeclaration() Stmt {
	p.consume(IDENTIFIER, "Unexpected token")
	name := p.previous()
	var init Expr
	if p.match(EQUAL) {
		init = p.expression()
	}
	p.consume(SEMICOLON, "Unexpected end of input, Expect ';' after value")
	return VarStmt{name, init}
}

func (p *Parser) statement() Stmt {
	if p.match(PRINT) {
		return p.printStatement()
	}

	if p.match(LEFT_BRACE) {
		return p.blockStatement()
	}
	if p.match(IF) {
		return p.ifstatement()
	}
	if p.match(WHILE) {
		return p.whileStatement()
	}
	if p.match(FOR) {
		return p.forStatement()
	}

	return p.expressionStatement()
}

func (p *Parser) printStatement() Stmt {
	expr := p.expression()
	p.consume(SEMICOLON, "Unexpected end of input, Expect ';' after value")
	return PrintStmt{expr}
}

func (p *Parser) ifstatement() Stmt {
	p.consume(LEFT_PAREN, "Unexpected token")
	condition := p.expression()
	p.consume(RIGHT_PAREN, "Unexpected token")

	consequent := p.statement()
	var alternate Stmt

	if p.match(ELSE) {
		alternate = p.statement()
	}

	return IfStmt{condition, consequent, alternate}
}

func (p *Parser) whileStatement() Stmt {
	p.consume(LEFT_PAREN, "Unexpected token")
	condition := p.expression()
	p.consume(RIGHT_PAREN, "Unexpected token")
	body := p.statement()
	return WhileStmt{condition, body}
}

/*
	the for loop can be transformed like this
	{
		var i = 0; 				// init
		while (i < 10) { 	// condition
			print i; 				// original body
			i = i + 1; 			// increment
		}
	}
*/
func (p *Parser) forStatement() Stmt {
	// forStmt is syntactic sugar for while loop
	p.consume(LEFT_PAREN, "Expected ')' before for clauses")

	var init Stmt
	var condition Expr
	var increment Expr
	var body Stmt

	if p.match(SEMICOLON) {
		init = nil
	} else if p.match(VAR) {
		init = p.varDeclaration()
	} else {
		init = p.expressionStatement()
	}

	if !p.checkType(SEMICOLON) {
		condition = p.expression()
	}

	// if condition is omited, we treat it as 'true'
	if condition == nil {
		condition = LiteralExpr{true}
	}

	p.consume(SEMICOLON, "Expected ';' after loop condition")

	if !p.checkType(RIGHT_PAREN) {
		increment = p.expression()
	}
	p.consume(RIGHT_PAREN, "Expected ')' after for clauses")

	body = p.statement()

	if increment != nil {
		// add increment expr to the end of for original boby
		// which will be executed after original body is executed
		body = BlockStmt{
			[]Stmt{body, ExpressionStmt{increment}},
		}
	}

	body = WhileStmt{condition, body}

	if init != nil {
		// add initializer stmt to the beginning of the while boby
		body = BlockStmt{
			[]Stmt{init, body},
		}
	}

	return body
}

func (p *Parser) expressionStatement() Stmt {
	expr := p.expression()
	p.consume(SEMICOLON, "Unexpected end of input, Expect ';' after value")
	return ExpressionStmt{expr}
}

func (p *Parser) blockStatement() Stmt {
	stmts := make([]Stmt, 0)
	for !p.checkType(RIGHT_BRACE) && !p.isAtEnd() {
		stmts = append(stmts, p.declaration())
	}
	p.consume(RIGHT_BRACE, "Unexpected end of input, Expect '}' after block")
	return BlockStmt{stmts}
}

/*
	EXPRESSIONS
*/

// expression rule
func (p *Parser) expression() Expr {
	return p.sequence()
}

func (p *Parser) sequence() Expr {
	expr := p.assignment()

	for p.match(COMMA) {
		next := p.assignment()
		if sequenceExpr, ok := expr.(SequenceExpr); ok {
			sequenceExpr.exprs = append(sequenceExpr.exprs, next)
			expr = sequenceExpr
		} else {
			expr = SequenceExpr{
				[]Expr{expr, next},
			}
		}
	}
	return expr
}

// sequence rule has the lowest piority just like c-like languages
func (p *Parser) assignment() Expr {
	expr := p.condition()

	for p.match(EQUAL) {
		equal := p.previous()
		value := p.assignment()
		if expr, ok := expr.(IdentifierExpr); ok {
			name := expr.name
			return AssignExpr{name, value}
		}
		p.lox.errorReporter.errorWithoutExit(RuntimeError{
			equal,
			"Invalid left-hand assignment target.",
		})
	}
	return expr
}

// condition → equality ("?" condition ":" condition)?
func (p *Parser) condition() Expr {
	expr := p.loginOr()

	if p.match(QUESTION) {
		consequent := p.condition()
		if p.match(COLON) {
			alternate := p.condition()
			expr = ConditionExpr{
				expr,
				consequent,
				alternate,
			}
		} else {
			p.lox.errorReporter.errorWithoutExit(ParseError{
				p.peek(),
				"Unexpected end of input",
			})
		}
	}
	return expr
}

// logic_or   → logic_and ( "or" logic_and )* ;
func (p *Parser) loginOr() Expr {
	expr := p.loginAnd()

	for p.match(OR) {
		token := p.previous()
		right := p.loginAnd()
		expr = LogicalExpr{expr, token, right}
	}
	return expr
}

// logic_and  → equality ( "and" equality )* ;
func (p *Parser) loginAnd() Expr {
	expr := p.equality()

	for p.match(AND) {
		token := p.previous()
		right := p.equality()
		expr = LogicalExpr{expr, token, right}
	}

	return expr
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
	if p.match(IDENTIFIER, STRING) {
		return IdentifierExpr{p.previous()}
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
		"Invalid or unexpected token:" + p.peek().literal,
	})

	return nil
}
