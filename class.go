package main

type Class struct {
	name        string
	methods     map[string]Function
	interpreter Interpreter
}

// TODO: Add Class string representation
func (c Class) String() string {
	return "<class " + c.name + ">"
}

func (c Class) call(interpreter Interpreter, args []interface{}) interface{} {
	instance := ClassInstance{
		class: c,
	}
	return instance
}

func (c Class) arity() int {
	return 0
}

func (c Class) findMethod(name string) Function {
	return c.methods[name]
}

type ClassInstance struct {
	class  Class
	fields map[string]interface{}
}

func (c ClassInstance) get(name Token) (interface{}, error) {
	if value, ok := c.fields[name.literal]; ok {
		return value, nil
	} else {
		if method := c.class.findMethod(name.literal); method != nil {
			return method, nil
		}

		return nil, RuntimeError{
			name,
			"Undefined property",
		}
	}
}

func (c ClassInstance) set(name Token, value interface{}) error {
	c.fields[name.literal] = value
	return nil
}

func (c ClassInstance) String() string {
	return "<classInstance " + c.class.name + ">"
}
