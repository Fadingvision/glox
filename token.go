package main

import (
	"fmt"
	"strings"
	"text/scanner"
)

// Token is A struct contains lexeme or token info
type Token struct {
	tokentype TokenType
	lexeme    string
	literal   string
	line      int
	column    int
}

// Scanner is Lexeme anglizer
type Scanner struct {
	lox         *Lox
	textScanner scanner.Scanner
	tokens      []Token
	source      string
}

func (t *Scanner) scanTokens() {
	t.textScanner = scanner.Scanner{}
	t.textScanner.Init(strings.NewReader(t.source))
	token := t.textScanner.Scan()
	for token != scanner.EOF {
		line, column := t.textScanner.Pos().Line, t.textScanner.Pos().Column
		tokenText := t.textScanner.TokenText()
		fmt.Println(line, column, tokenText, token)

		switch token {
		case scanner.Float:
		case scanner.Int:
			t.addToken(NUMBER, tokenText)
			break
		case scanner.String:
			t.addToken(STRING, tokenText[1:len(tokenText)-1])
			break
		default:
			t.identifiyToken(tokenText)
			break
		}
		token = t.textScanner.Scan()
	}
	fmt.Println(t.tokens)
}

func (t *Scanner) peek() string {
	fmt.Println("TokenType:", t.textScanner.Peek())
	if t.textScanner.Peek() == scanner.EOF {
		return ""
	}

	// FIXME: printbale string vs normal string ???
	printableString := scanner.TokenString(t.textScanner.Peek())
	return printableString[1 : len(printableString)-1]
}

func (t *Scanner) identifiyToken(tokenText string) {
	switch tokenText {
	case "(":
		t.addToken(LEFT_PAREN, tokenText)
		break
	case ")":
		t.addToken(RIGHT_PAREN, tokenText)
		break
	case "{":
		t.addToken(LEFT_BRACE, tokenText)
		break
	case "}":
		t.addToken(RIGHT_BRACE, tokenText)
		break
	case ",":
		t.addToken(COMMA, tokenText)
		break
	case ".":
		t.addToken(DOT, tokenText)
		break
	case "-":
		t.addToken(MINUS, tokenText)
		break
	case "+":
		t.addToken(PLUS, tokenText)
		break
	case ";":
		t.addToken(SEMICOLON, tokenText)
		break
	case "*":
		t.addToken(STAR, tokenText)
		break
	case "!":
		if ok, val := t.match("="); ok {
			t.addToken(BANG_EQUAL, tokenText+val)
		} else {
			t.addToken(BANG, tokenText)
		}
		break
	case "=":
		if ok, val := t.match("="); ok {
			t.addToken(EQUAL_EQUAL, tokenText+val)
		} else {
			t.addToken(EQUAL, tokenText)
		}
		break
	case "<":
		if ok, val := t.match("="); ok {
			t.addToken(LESS_EQUAL, tokenText+val)
		} else {
			t.addToken(LESS, tokenText)
		}
		break
	case ">":
		if ok, val := t.match("="); ok {
			t.addToken(GREATER_EQUAL, tokenText+val)
		} else {
			t.addToken(GREATER, tokenText)
		}
		break
	default:
		// FIXME: 英文开头的只含有英文数字的字符才能算作有效的标识符
		if true {
			fmt.Println(tokenText)
			if tokenType, ok := keywords[tokenText]; ok {
				t.addToken(tokenType, tokenText)
			} else {
				t.addToken(IDENTIFIER, tokenText)
			}
		} else {
			t.lox.hasError = true
			t.lox.errorReporter.error(&TokenError{
				msg:    "Unexpected token",
				line:   t.textScanner.Pos().Line,
				column: t.textScanner.Pos().Column,
			})
		}
	}
}

func (t *Scanner) addToken(tokentype TokenType, value string) {
	t.tokens = append(t.tokens, Token{
		tokentype: tokentype,
		lexeme:    value,
		line:      t.textScanner.Pos().Line,
		column:    t.textScanner.Pos().Column,
	})
}

func (t *Scanner) match(expected string) (bool, string) {
	next := t.peek()
	if next != expected {
		return false, ""
	}
	t.textScanner.Scan()
	return true, next
}
