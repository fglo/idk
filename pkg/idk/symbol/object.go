package symbol

import (
	"bytes"
	"fmt"
	"hash/fnv"
	"strings"

	"github.com/fglo/idk/pkg/idk/ast"
)

type BuiltinFunction func(args ...Object) Object

type ObjectType string

const (
	NULL_OBJ  ObjectType = "NULL"
	ERROR_OBJ ObjectType = "ERROR"

	TYPE_OBJ ObjectType = "TYPE"

	INTEGER_OBJ        ObjectType = "INTEGER"
	FLOATING_POINT_OBJ ObjectType = "FLOAT"
	BOOLEAN_OBJ        ObjectType = "BOOLEAN"
	CHARACTER_OBJ      ObjectType = "CHARACTER"
	STRING_OBJ         ObjectType = "STRING"

	ARRAY_OBJ ObjectType = "ARRAY"
	HASH_OBJ  ObjectType = "HASH"

	RETURN_VALUE_OBJ ObjectType = "RETURN_VALUE"

	FUNCTION_OBJ ObjectType = "FUNCTION"
	BUILTIN_OBJ  ObjectType = "BUILTIN"
)

type HashKey struct {
	Type  ObjectType
	Value uint64
}

type Hashable interface {
	HashKey() HashKey
}

type Object interface {
	Type() ObjectType
	Inspect() string
}

func IsError(obj Object) bool {
	if obj != nil {
		return obj.Type() == ERROR_OBJ
	}
	return false
}

type Type struct {
	Value ObjectType
}

func (t *Type) Type() ObjectType { return TYPE_OBJ }
func (t *Type) Inspect() string  { return string(t.Value) }
func (t *Type) HashKey() HashKey {
	h := fnv.New64a()
	h.Write([]byte(t.Value))

	return HashKey{Type: t.Type(), Value: h.Sum64()}
}

type Integer struct {
	Value int64
}

func (i *Integer) Type() ObjectType { return INTEGER_OBJ }
func (i *Integer) Inspect() string  { return fmt.Sprintf("%d", i.Value) }
func (i *Integer) HashKey() HashKey {
	return HashKey{Type: i.Type(), Value: uint64(i.Value)}
}

type FloatingPoint struct {
	Value float64
}

func (fp *FloatingPoint) Type() ObjectType { return FLOATING_POINT_OBJ }
func (fp *FloatingPoint) Inspect() string  { return fmt.Sprintf("%f", fp.Value) }
func (fp *FloatingPoint) HashKey() HashKey {
	return HashKey{Type: fp.Type(), Value: uint64(fp.Value)}
}

type Boolean struct {
	Value bool
}

func (b *Boolean) Type() ObjectType { return BOOLEAN_OBJ }
func (b *Boolean) Inspect() string  { return fmt.Sprintf("%t", b.Value) }
func (b *Boolean) HashKey() HashKey {
	var value uint64

	if b.Value {
		value = 1
	} else {
		value = 0
	}

	return HashKey{Type: b.Type(), Value: value}
}

type Character struct {
	Value rune
}

func (c *Character) Type() ObjectType { return CHARACTER_OBJ }
func (c *Character) Inspect() string  { return fmt.Sprintf("%c", c.Value) }
func (c *Character) HashKey() HashKey {
	return HashKey{Type: c.Type(), Value: uint64(c.Value)}
}

type String struct {
	Value string
}

func (s *String) Type() ObjectType { return STRING_OBJ }
func (s *String) Inspect() string  { return s.Value }
func (s *String) HashKey() HashKey {
	h := fnv.New64a()
	h.Write([]byte(s.Value))

	return HashKey{Type: s.Type(), Value: h.Sum64()}
}

type Null struct{}

func (n *Null) Type() ObjectType { return NULL_OBJ }
func (n *Null) Inspect() string  { return "null" }

type ReturnValue struct {
	Value Object
}

func (rv *ReturnValue) Type() ObjectType { return RETURN_VALUE_OBJ }
func (rv *ReturnValue) Inspect() string  { return rv.Value.Inspect() }

type Error struct {
	File           string
	LineNumber     int
	PositionInLine int
	Message        string
}

func (e *Error) Type() ObjectType { return ERROR_OBJ }
func (e *Error) Inspect() string {
	switch {
	case e.LineNumber != 0 && e.PositionInLine != 0:
		return fmt.Sprintf("ERROR: Evaluator error in file %s on line %v, position %v: %s", e.File, e.LineNumber, e.PositionInLine, e.Message)
	case e.LineNumber != 0:
		return fmt.Sprintf("ERROR: Evaluator error in file %s on line %v: %s", e.File, e.LineNumber, e.Message)
	default:
		return fmt.Sprintf("ERROR: Evaluator error in file %s: %s", e.File, e.Message)
	}
}

type Function struct {
	Identifier string
	Parameters []*ast.DeclareStatement
	Body       *ast.BlockStatement
	Scope      *Scope
	ReturnType ObjectType
}

func (f *Function) Type() ObjectType { return FUNCTION_OBJ }
func (f *Function) Inspect() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range f.Parameters {
		params = append(params, p.String())
	}

	out.WriteString("func")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n")
	out.WriteString(f.Body.String())
	out.WriteString("\n}")

	return out.String()
}

type Builtin struct {
	Fn BuiltinFunction
}

func (b *Builtin) Type() ObjectType { return BUILTIN_OBJ }
func (b *Builtin) Inspect() string  { return "builtin function" }

type Array struct {
	Elements []Object
}

func (ao *Array) Type() ObjectType { return ARRAY_OBJ }
func (ao *Array) Inspect() string {
	var out bytes.Buffer

	elements := []string{}
	for _, e := range ao.Elements {
		elements = append(elements, e.Inspect())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")

	return out.String()
}

type HashPair struct {
	Key   Object
	Value Object
}

type Hash struct {
	Pairs map[HashKey]HashPair
}

func (h *Hash) Type() ObjectType { return HASH_OBJ }
func (h *Hash) Inspect() string {
	var out bytes.Buffer

	pairs := []string{}
	for _, pair := range h.Pairs {
		pairs = append(pairs, fmt.Sprintf("%s: %s",
			pair.Key.Inspect(), pair.Value.Inspect()))
	}

	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")

	return out.String()
}
