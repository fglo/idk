package symbol

type Scope struct {
	symbolTable map[string]Object
	outer       *Scope
}

func NewScope() *Scope {
	st := make(map[string]Object)
	return &Scope{symbolTable: st, outer: nil}
}

func NewInnerScope(outer *Scope) *Scope {
	env := NewScope()
	env.outer = outer
	return env
}

func (e *Scope) Lookup(name string) (Object, bool) {
	obj, ok := e.symbolTable[name]
	if !ok && e.outer != nil {
		obj, ok = e.outer.Lookup(name)
	}
	return obj, ok
}

func (e *Scope) LookupInCurrentScope(name string) (Object, bool) {
	obj, ok := e.symbolTable[name]
	return obj, ok
}

func (e *Scope) Insert(name string, val Object) Object {
	e.symbolTable[name] = val
	return val
}
