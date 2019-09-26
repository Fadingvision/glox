package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"text/scanner"
)

// Token is A struct contains lexeme or token info
type Token struct {
	tokentype TokenType
	lexeme    interface{}
	literal   string
	line      int
	column    int
}

// Scanner is Lexeme anglizer
type Scanner struct {
	lox         *Lox
	textScanner *scanner.Scanner
	tokens      []Token
	source      string
}

func (t *Scanner) scanTokens() {
	t.textScanner = &scanner.Scanner{}
	// FIXME: custom Error do not working
	t.textScanner.Error = func(s *scanner.Scanner, msg string) {
		t.lox.hasError = true
		t.lox.errorReporter.error(&TokenError{
			msg:    msg,
			line:   s.Pos().Line,
			column: s.Pos().Column,
		})
	}
	t.textScanner.Init(strings.NewReader(t.source))
	// don't skip comments
	// NOTE: consider Add support for comments liken `//` and '/* ... */'
	// consider allowing them to nest
	// t.textScanner.Mode ^= scanner.SkipComments
	token := t.textScanner.Scan()
	for token != scanner.EOF {
		line, column := t.textScanner.Pos().Line, t.textScanner.Pos().Column
		tokenText := t.textScanner.TokenText()
		// fmt.Println(line, column, tokenText, token)

		switch token {
		case scanner.Float, scanner.Int:
			// turn all numeric thing to float64, like javascript
			if num, err := strconv.ParseFloat(tokenText, 64); err != nil {
				t.lox.hasError = true
				t.lox.errorReporter.error(&TokenError{
					msg:    "Unexpected number",
					line:   line,
					column: column,
				})
			} else {
				t.addToken(NUMBER, num)
			}
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
	if t.textScanner.Peek() == scanner.EOF {
		return ""
	}

	// FIXME: printable string vs normal string ???
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
		// only begin with `_`, alpha and contains only alpha-numeric charactors
		// we treat it as a identifier
		if regexp.MustCompile(`^[_a-zA-Z]([_a-zA-Z0-9]*)$`).MatchString(tokenText) {
			if tokenType, ok := keywords[tokenText]; ok {
				t.addToken(tokenType, tokenText)
			} else {
				t.addToken(IDENTIFIER, tokenText)
			}
		} else {
			t.lox.hasError = true
			t.lox.errorReporter.error(&TokenError{
				msg:    "Invalid or unexpected token: " + tokenText,
				line:   t.textScanner.Pos().Line,
				column: t.textScanner.Pos().Column,
			})
		}
	}
}

func (t *Scanner) addToken(tokentype TokenType, value interface{}) {
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
