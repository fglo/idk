package symbol

type Symbol struct {
	Object Object
	Type   ObjectType
}

type Scope struct {
	Name        string
	symbolTable map[string]Symbol
	outer       *Scope
	namedScopes map[string]*Scope
}

func NewScope() *Scope {
	return &Scope{
		symbolTable: make(map[string]Symbol),
		namedScopes: make(map[string]*Scope),
		outer:       nil,
	}
}

func NewInnerScope(outer *Scope) *Scope {
	env := NewScope()
	env.outer = outer
	return env
}

func (s *Scope) GetNamedScope(name string) *Scope {
	if scope, ok := s.namedScopes[name]; ok {
		return scope
	}

	env := NewInnerScope(s)
	env.Name = name

	s.namedScopes[name] = env

	return env
}

func (s *Scope) Lookup(name string) (Symbol, bool) {
	obj, ok := s.symbolTable[name]
	if !ok && s.outer != nil {
		obj, ok = s.outer.Lookup(name)
	}
	return obj, ok
}

func (s *Scope) LookupInCurrentScope(name string) (Symbol, bool) {
	obj, ok := s.symbolTable[name]
	return obj, ok
}

func (s *Scope) Insert(name string, val Object, typ ObjectType) Symbol {
	symbol := Symbol{val, typ}
	s.symbolTable[name] = symbol
	return symbol
}

func (s *Scope) TryToAssign(name string, val Object, typ ObjectType) bool {
	if _, ok := s.LookupInCurrentScope(name); ok {
		s.symbolTable[name] = Symbol{val, typ}
		return true
	} else if s.outer != nil {
		return s.outer.TryToAssign(name, val, typ)
	}
	return false
}
