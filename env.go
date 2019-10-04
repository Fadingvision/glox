package main

type env struct {
	values map[string]interface{}
	parent *env
}

func (e env) set(key string, value interface{}) {
	e.values[key] = value
}

func (e env) get(token Token) (interface{}, error) {
	value, ok := e.values[token.literal]
	if ok {
		return value, nil
	}

	// if it's not in global scope, lookup
	if e.parent != nil {
		return e.parent.get(token)
	}

	return nil, RuntimeError{token, "Undefined variable: " + token.literal}
}

func (e env) assign(token Token, value interface{}) (interface{}, error) {
	_, ok := e.values[token.literal]
	if ok {
		e.values[token.literal] = value
		return value, nil
	}

	// if it's not in global scope, lookup
	if e.parent != nil {
		return e.parent.assign(token, value)
	}

	return nil, RuntimeError{token, "Undefined variable: " + token.literal}
}
