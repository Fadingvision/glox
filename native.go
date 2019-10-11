package main

import (
	"fmt"
	"time"
)

// Clock shows current time
type Clock struct{}

func (c Clock) call(interpreter Interpreter, args []interface{}) interface{} {
	return time.Now().Format("20060102150405")
}

func (c Clock) arity() int {
	return 0
}

func (c Clock) String() string {
	return "<native fn>"
}

// Print shows stuff
type Print struct{}

func (p Print) call(interpreter Interpreter, args []interface{}) interface{} {
	fmt.Println(args...)
	return nil
}

func (p Print) arity() int {
	return -1
}

func (p Print) String() string {
	return "<native fn>"
}
