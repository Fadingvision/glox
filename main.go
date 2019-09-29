package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

type Lox struct {
	errorReporter *ErrorReporter
	scanner       *Scanner
	parser        *Parser
	hasError      bool
}

func main() {
	args := os.Args[1:]

	lox := &Lox{
		errorReporter: &ErrorReporter{},
		parser:        &Parser{},
		scanner: &Scanner{
			tokens: []Token{},
		},
	}
	lox.scanner.lox = lox
	lox.errorReporter.lox = lox
	lox.parser.lox = lox

	if len(args) > 1 {
		fmt.Println("GLOX]: Usage: glox [script]")
		os.Exit(1)
	} else if len(args) == 1 {
		lox.runFile(args[0])
	} else {
		lox.runREPL()
	}
}

func (l *Lox) runFile(path string) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println("[GLOX]: unvalid filepath")
		os.Exit(1)
	}

	l.run(string(content))
}

func (l *Lox) runREPL() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		l.run(scanner.Text())
	}

	if scanner.Err() != nil {
		// handle error.
	}
}

func (l *Lox) run(src string) {
	l.scanner.source = src
	l.scanner.scanTokens()
	l.parser.tokens = l.scanner.tokens
	expr := l.parser.parse()
	// check if our expr works as we expect
	AstPrinter{}.print(expr, os.Stdout)
}

// you will likely have multiple ways errors get displayed
// on stderr, in an IDE’s error window, logged to a file, etc.
type ErrorReporter struct {
	lox *Lox
}

func (e *ErrorReporter) error(err error) {
	e.lox.hasError = true
	// by default, we just print on the interface and then exit with 1
	log.Fatal(err)
}

func (e *ErrorReporter) errorWithoutExit(err error) {
	e.lox.hasError = true
	fmt.Println(err)
}
