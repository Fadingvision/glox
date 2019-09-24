package main

type Token struct {
	tokentype TokenType
	lexeme    string
	literal   string
	line      int
	column    int
}

func (t Token) New(tokentype TokenType, line int, column int) Token {
	return Token{
		tokentype: tokentype,
		line:      line,
		column:    column,
	}
}

type Scanner struct {
	tokens []Token
}

func (t *Scanner) idenrifiyToken(token string, line int, column int) {
	switch token {
	case "(":
		t.addToken(LEFT_PAREN, line, column)
		break
	case ")":
		t.addToken(RIGHT_PAREN, line, column)
		break
	case "{":
		t.addToken(LEFT_BRACE, line, column)
		break
	case "}":
		t.addToken(RIGHT_BRACE, line, column)
		break
	case ",":
		t.addToken(COMMA, line, column)
		break
	case ".":
		t.addToken(DOT, line, column)
		break
	case "-":
		t.addToken(MINUS, line, column)
		break
	case "+":
		t.addToken(PLUS, line, column)
		break
	case ";":
		t.addToken(SEMICOLON, line, column)
		break
	case "*":
		t.addToken(STAR, line, column)
		break
	}
}

func (t *Scanner) addToken(tokentype TokenType, line int, column int) Token {
	return Token{
		tokentype: tokentype,
		line:      line,
		column:    column,
	}
}
