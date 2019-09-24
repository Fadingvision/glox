package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"text/scanner"
)

type Lox struct {
	errorReporter *ErrorReporter
	hasError      bool
}

func main() {
	args := os.Args[1:]

	lox := &Lox{
		errorReporter: &ErrorReporter{},
	}

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
	var s scanner.Scanner
	s.Init(strings.NewReader(src))
	token := s.Scan()
	for ; token != scanner.EOF; token = s.Scan() {
		fmt.Println(s.Pos().Line, s.Pos().Column, token, s.TokenText())
	}
}

// TokenError implement the std err interface
type TokenError struct {
	code   string
	msg    string
	line   int
	Column int
}

func (e *TokenError) Error() string {
	// TODO: more elgant error display
	return fmt.Sprintf("[GLOX] Error: Line%d Column %d, %s", e.line, e.Column, e.msg)
}

// you will likely have multiple ways errors get displayed
// on stderr, in an IDE’s error window, logged to a file, etc.
type ErrorReporter struct{}

func (e *ErrorReporter) error(err TokenError) {
	// by default, we just print on the interface and then exit with 1
	log.Fatal(err)
}
