package main

// Parser does two things
// 1. tranform tokens to AST tree
// 2. report synatx error based on our language CFG
type Parser struct {
	tokens []Token
	lox    *Lox
}

func (p *Parser) parse() {}
