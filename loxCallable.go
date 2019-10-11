package main

// Callable interface, function and methods should both implement
type Callable interface {
	call(interpreter Interpreter, args []interface{}) interface{}
	arity() int
}

// ReturnValue is function's return value, Distinguish between normal panic errors
type ReturnValue struct {
	value interface{}
}

type Function struct {
	stmt FunStmt
	// function declare environment, which is known as `closure`
	closure env
	isInit  bool
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

	// If the function is an initializer,
	// We ignore the nil return value or override the actual return value,
	// and forcibly return this.
	if f.isInit {
		return f.closure.values["this"]
	}

	return returnVal
}

func (f Function) bind(instance ClassInstance) Function {
	// Add a new env to store `this` in interpreter to sync the resolver
	environment := env{
		values: make(map[string]interface{}, 0),
		parent: &f.closure,
	}

	// this is dynamic, it refers to different instances
	environment.set("this", instance)

	return Function{
		stmt:    f.stmt,
		closure: environment,
		isInit:  f.isInit,
	}
}

func (f Function) arity() int {
	return len(f.stmt.params)
}
