package main

import "time"

// Clock to show current time
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

// TODO: make print as the native function instead of expression
