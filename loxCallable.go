package main

import "time"

// Callable interface, function and methods should both implement
type Callable interface {
	call(interpreter Interpreter, args []interface{}) interface{}
	arity() int
}

// ReturnValue is function's return value, Distinguish between normal panic errors
type ReturnValue struct {
	value interface{}
}

// Clock is first built-in function, to show current time
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

type Function struct {
	stmt FunStmt
	// function declare environment, which is known as `closure`
	closure env
}

// TODO: Add function string representation
func (f Function) String() string {
	return "<fn " + f.stmt.name.literal + ">"
}

func (f Function) call(interpreter Interpreter, args []interface{}) interface{} {
	var returnVal interface{}

	environment := env{
		values: make(map[string]interface{}, 0),
		parent: &f.closure,
	}

	// build a local variable for each one param
	for index, param := range f.stmt.params {
		environment.set(param.lexeme.(string), args[index])
	}

	func() {
		// recover from return statement
		defer func() {
			err := recover()
			if val, ok := err.(ReturnValue); ok {
				returnVal = val.value
			} else if err != nil {
				panic(err)
			}
		}()
		// execute this function body in this environment
		interpreter.executeBlockStmt(f.stmt.body, environment)
	}()

	return returnVal
}

func (f Function) arity() int {
	return len(f.stmt.params)
}
