package object

// Environment is HashMap
type Environment struct {
	store map[string]Object
	outer *Environment
}

func NewEnclosedEnvironment(e *Environment) *Environment {
	env := NewEnvironment()
	env.outer = e
	return env
}

func NewEnvironment() *Environment {
	s := make(map[string]Object)
	return &Environment{store: s}
}

func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(name)
	}
	return obj, ok
}

func (e *Environment) Set(name string, obj Object) Object {
	e.store[name] = obj
	return obj
}
