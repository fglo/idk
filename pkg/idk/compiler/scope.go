package compiler

import "github.com/fglo/idk/pkg/idk/opcodes"

type scope struct {
	symbolTable map[string]symbol
	outer       *scope
}

type symbol struct {
	valType opcodes.ValType
	cpAddr  int
}

func NewScope() *scope {
	return &scope{
		symbolTable: make(map[string]symbol),
		outer:       nil,
	}
}

func NewInnerScope(outer *scope) *scope {
	env := NewScope()
	env.outer = outer
	return env
}

func (s *scope) Insert(name string, addr int, varType opcodes.ValType) {
	s.symbolTable[name] = symbol{
		valType: varType,
		cpAddr:  addr,
	}
}

func (s *scope) Lookup(name string) (symbol, bool) {
	symbol, ok := s.symbolTable[name]
	if !ok && s.outer != nil {
		symbol, ok = s.outer.Lookup(name)
	}
	return symbol, ok
}

// func (s *Scope) LookupInCurrentScope(name string) (opcodes.VarType, bool) {
// 	obj, ok := s.symbolTable[name]
// 	return obj, ok
// }
