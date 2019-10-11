package main

type Object interface {
	get(name Token) (interface{}, error)
	set(name Token, value interface{}) error
}

type Class struct {
	name          string
	super         *Class
	staticMethods map[string]Function
	methods       map[string]Function
	fields        map[string]interface{}
}

// TODO: Add Class string representation
func (c Class) String() string {
	return "<class " + c.name + ">"
}

func (c Class) call(interpreter Interpreter, args []interface{}) interface{} {
	instance := ClassInstance{
		class:  c,
		fields: make(map[string]interface{}, 0),
	}

	// call constructor
	if init, ok := c.findMethod("init"); ok {
		init.bind(instance).call(interpreter, args)
	}

	return instance
}

func (c Class) get(name Token) (interface{}, error) {
	if value, ok := c.fields[name.literal]; ok {
		return value, nil
	}

	if method, ok := c.staticMethods[name.literal]; ok {
		return method, nil
	}

	return nil, RuntimeError{
		name,
		"Undefined property",
	}
}

func (c Class) set(name Token, value interface{}) error {
	c.fields[name.literal] = value
	return nil
}

func (c Class) arity() int {
	if init, ok := c.findMethod("init"); ok {
		return init.arity()
	}
	return 0
}

func (c Class) findMethod(name string) (Function, bool) {
	val, ok := c.methods[name]

	if !ok && c.super != nil {
		val, ok = c.super.findMethod(name)
	}

	return val, ok
}

type ClassInstance struct {
	class  Class
	fields map[string]interface{}
}

func (c ClassInstance) get(name Token) (interface{}, error) {
	if value, ok := c.fields[name.literal]; ok {
		return value, nil
	}

	if method, ok := c.class.findMethod(name.literal); ok {
		return method.bind(c), nil
	}

	return nil, RuntimeError{
		name,
		"Undefined property",
	}
}

func (c ClassInstance) set(name Token, value interface{}) error {
	c.fields[name.literal] = value
	return nil
}

// TODO: Add object string representation
func (c ClassInstance) String() string {
	return "<classInstance " + c.class.name + ">"
}
