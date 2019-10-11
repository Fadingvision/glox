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

func (e env) getAt(token Token, distance int) (interface{}, error) {
	environment := e.ancestor(distance)
	// If our resolver did not go wrong, this must be valid
	value, ok := environment.values[token.literal]
	if ok {
		return value, nil
	}

	return nil, RuntimeError{token, "Undefined variable: " + token.literal}
}

func (e env) getAtByName(name string, distance int) (interface{}, bool) {
	environment := e.ancestor(distance)
	// If our resolver did not go wrong, this must be valid
	value, ok := environment.values[name]
	return value, ok
}

func (e env) ancestor(distance int) *env {
	var environment = &e
	for i := 0; i < distance; i++ {
		environment = environment.parent
	}
	return environment
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

func (e env) assignAt(token Token, value interface{}, distance int) (interface{}, error) {
	environment := e.ancestor(distance)
	_, ok := environment.values[token.literal]
	if ok {
		environment.values[token.literal] = value
		return value, nil
	}

	return nil, RuntimeError{token, "Undefined variable: " + token.literal}
}
