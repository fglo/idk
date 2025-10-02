package virtualmachine

import (
	"fmt"

	"github.com/fglo/idk/pkg/idk/opcodes"
)

type symbolTable struct {
	types   map[string]opcodes.ValType
	symbols map[string]int

	intMemory    []int
	floatMemory  []float64
	boolMemory   []bool
	charMemory   []rune
	stringMemory []string
}

func newSymbolTable() *symbolTable {
	return &symbolTable{
		types:   make(map[string]opcodes.ValType),
		symbols: make(map[string]int),

		intMemory:    make([]int, 0),
		floatMemory:  make([]float64, 0),
		boolMemory:   make([]bool, 0),
		charMemory:   make([]rune, 0),
		stringMemory: make([]string, 0),
	}
}

func (s *symbolTable) lookupVarType(name string) opcodes.ValType {
	valType, exists := s.types[name]
	if !exists {
		return opcodes.UNKNOWN_TYPE
	}
	return valType
}

func (s *symbolTable) bindInt(name string, value int) {
	s.symbols[name] = len(s.intMemory)
	s.types[name] = opcodes.INT
	s.intMemory = append(s.intMemory, value)
}

func (s *symbolTable) bindFloat(name string, value float64) {
	s.symbols[name] = len(s.floatMemory)
	s.types[name] = opcodes.FLOAT
	s.floatMemory = append(s.floatMemory, value)
}

func (s *symbolTable) bindBool(name string, value bool) {
	s.symbols[name] = len(s.boolMemory)
	s.types[name] = opcodes.BOOL
	s.boolMemory = append(s.boolMemory, value)
}

func (s *symbolTable) bindChar(name string, value rune) {
	s.symbols[name] = len(s.charMemory)
	s.types[name] = opcodes.CHAR
	s.charMemory = append(s.charMemory, value)
}

func (s *symbolTable) bindString(name string, value string) {
	s.symbols[name] = len(s.stringMemory)
	s.types[name] = opcodes.STRING
	s.stringMemory = append(s.stringMemory, value)
}

func (s *symbolTable) lookupInt(name string) (int, error) {
	address, exists := s.symbols[name]
	if !exists {
		return 0, fmt.Errorf("Couldn't find symbol '%s'", name)
	}
	return s.intMemory[address], nil
}

func (s *symbolTable) lookupFloat(name string) (float64, error) {
	address, exists := s.symbols[name]
	if !exists {
		return 0, fmt.Errorf("Couldn't find symbol '%s'", name)
	}
	return s.floatMemory[address], nil
}

func (s *symbolTable) lookupBool(name string) (bool, error) {
	address, exists := s.symbols[name]
	if !exists {
		return false, fmt.Errorf("Couldn't find symbol '%s'", name)
	}
	return s.boolMemory[address], nil
}

func (s *symbolTable) lookupChar(name string) (rune, error) {
	address, exists := s.symbols[name]
	if !exists {
		return 0, fmt.Errorf("Couldn't find symbol '%s'", name)
	}
	return s.charMemory[address], nil
}

func (s *symbolTable) lookupString(name string) (string, error) {
	address, exists := s.symbols[name]
	if !exists {
		return "", fmt.Errorf("Couldn't find symbol '%s'", name)
	}
	return s.stringMemory[address], nil
}
