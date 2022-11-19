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

func (s *Scope) Lookup(name string) (Object, bool) {
	obj, ok := s.symbolTable[name]
	if !ok && s.outer != nil {
		obj, ok = s.outer.Lookup(name)
	}
	return obj, ok
}

func (s *Scope) LookupInCurrentScope(name string) (Object, bool) {
	obj, ok := s.symbolTable[name]
	return obj, ok
}

func (s *Scope) Insert(name string, val Object) Object {
	s.symbolTable[name] = val
	return val
}

func (s *Scope) TryToAssign(name string, val Object) bool {
	if _, ok := s.LookupInCurrentScope(name); ok {
		s.symbolTable[name] = val
		return true
	} else if s.outer != nil {
		return s.outer.TryToAssign(name, val)
	}
	return false
}
