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
	hasError      bool
}

func main() {
	args := os.Args[1:]

	lox := &Lox{
		errorReporter: &ErrorReporter{},
		scanner: &Scanner{
			tokens: []Token{},
		},
	}
	lox.scanner.lox = lox

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
}

// TokenError implement the std err interface
type TokenError struct {
	code   string
	msg    string
	line   int
	column int
}

func (e *TokenError) Error() string {
	// TODO: more elgant error display
	return fmt.Sprintf("[GLOX] Error: Line%d c %d, %s", e.line, e.column, e.msg)
}

// you will likely have multiple ways errors get displayed
// on stderr, in an IDE’s error window, logged to a file, etc.
type ErrorReporter struct{}

func (e *ErrorReporter) error(err error) {
	// by default, we just print on the interface and then exit with 1
	log.Fatal(err)
}
