package ast

import (
	"bytes"
	"strconv"

	"github.com/fglo/idk/pkg/idk/token"
)

type Expression interface {
	Node
	expressionNode()
}

type PrefixExpression struct {
	token token.Token
	Right Expression
}

func NewPrefixExpression(operator token.Token, expression Expression) *PrefixExpression {
	pe := &PrefixExpression{
		token: operator,
		Right: expression,
	}
	return pe
}

func (pe *PrefixExpression) expressionNode()               {}
func (pe *PrefixExpression) GetTokenValue() string         { return pe.token.Value }
func (pe *PrefixExpression) GetTokenType() token.TokenType { return pe.token.Type }
func (pe *PrefixExpression) GetLineNumber() int            { return pe.token.Line }
func (pe *PrefixExpression) GetPositionInLine() int        { return pe.token.PositionInLine }
func (pe *PrefixExpression) GetChildren() []Node           { return []Node{pe.Right} }
func (pe *PrefixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(pe.GetTokenValue())
	out.WriteString(pe.Right.String())
	out.WriteString(")")

	return out.String()
}

type InfixExpression struct {
	token token.Token
	Left  Expression
	Right Expression
}

func NewInfixExpression(left Expression, operator token.Token, right Expression) *InfixExpression {
	ie := &InfixExpression{
		token: operator,
		Left:  left,
		Right: right,
	}
	return ie
}

func (ie *InfixExpression) expressionNode()               {}
func (ie *InfixExpression) GetTokenValue() string         { return ie.token.Value }
func (ie *InfixExpression) GetTokenType() token.TokenType { return ie.token.Type }
func (ie *InfixExpression) GetLineNumber() int            { return ie.token.Line }
func (ie *InfixExpression) GetPositionInLine() int        { return ie.token.PositionInLine }
func (ie *InfixExpression) GetChildren() []Node           { return []Node{ie.Left, ie.Right} }
func (ie *InfixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString(" " + ie.GetTokenValue() + " ")
	out.WriteString(ie.Right.String())
	out.WriteString(")")

	return out.String()
}

type PropertyExpression struct {
	Parent   Expression
	Property Expression
}

func NewPropertyExpression(parent Expression, property Expression) *PropertyExpression {
	pe := &PropertyExpression{
		Parent:   parent,
		Property: property,
	}
	return pe
}

func (pe *PropertyExpression) expressionNode()               {}
func (pe *PropertyExpression) GetTokenValue() string         { return pe.Parent.GetTokenValue() }
func (pe *PropertyExpression) GetTokenType() token.TokenType { return token.DOT }
func (pe *PropertyExpression) GetLineNumber() int            { return pe.Parent.GetLineNumber() }
func (pe *PropertyExpression) GetPositionInLine() int        { return pe.Parent.GetPositionInLine() }
func (pe *PropertyExpression) GetChildren() []Node           { return []Node{pe.Parent, pe.Property} }
func (pe *PropertyExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(pe.Parent.String())
	// out.WriteString(" " + token.DOT.String() + " ")
	out.WriteString(".")
	out.WriteString(pe.Property.String())
	out.WriteString(")")

	return out.String()
}

type IfExpression struct {
	Condition   Expression
	Consequence *Expression
	Alternative *Expression
}

// TODO: implement if expression
func NewIfExpression(condition Expression, consequence *Expression, alternative *Expression) *IfExpression {
	ie := &IfExpression{
		Condition:   condition,
		Consequence: consequence,
		Alternative: alternative,
	}
	return ie
}

func (ie *IfExpression) expressionNode()               {}
func (ie *IfExpression) GetTokenValue() string         { return ie.Condition.GetTokenValue() }
func (ie *IfExpression) GetTokenType() token.TokenType { return token.IF }
func (ie *IfExpression) GetLineNumber() int            { return ie.Condition.GetLineNumber() }
func (ie *IfExpression) GetPositionInLine() int        { return ie.Condition.GetPositionInLine() }
func (ie *IfExpression) GetChildren() []Node {
	if ie.Alternative != nil {
		return []Node{ie.Condition, *ie.Consequence, *ie.Alternative}
	}
	return []Node{ie.Condition, *ie.Consequence}
}
func (ie *IfExpression) String() string {
	var out bytes.Buffer

	out.WriteString("if")
	out.WriteString(ie.Condition.String())
	out.WriteString(" ")
	out.WriteString((*ie.Consequence).String())

	if ie.Alternative != nil {
		out.WriteString("else ")
		out.WriteString((*ie.Alternative).String())
	}

	return out.String()
}

type FunctionCallExpression struct {
	Identifier Identifier
	Parameters []Expression
}

func NewFunctionCallExpression(tok token.Token) *FunctionCallExpression {
	fce := &FunctionCallExpression{
		Identifier: *NewIdentifier(tok),
	}
	return fce
}

func (fce *FunctionCallExpression) expressionNode()               {}
func (fce *FunctionCallExpression) GetTokenValue() string         { return fce.Identifier.GetTokenValue() }
func (fce *FunctionCallExpression) GetTokenType() token.TokenType { return fce.Identifier.GetType() }
func (fce *FunctionCallExpression) GetLineNumber() int            { return fce.Identifier.GetLineNumber() }
func (fce *FunctionCallExpression) GetPositionInLine() int        { return fce.Identifier.GetPositionInLine() }
func (fce *FunctionCallExpression) GetChildren() []Node           { return []Node{} }
func (fce *FunctionCallExpression) String() string {
	var out bytes.Buffer

	out.WriteString(fce.Identifier.String())
	out.WriteString("(")
	for i, parameter := range fce.Parameters {
		if i < len(fce.Parameters)-1 {
			out.WriteString(parameter.String() + ", ")
		} else {
			out.WriteString(parameter.String())
		}
	}
	out.WriteString(")")

	return out.String()
}

type Identifier struct {
	token token.Token
	_type token.TokenType
	value string
}

func NewIdentifier(identifier token.Token) *Identifier {
	i := &Identifier{
		token: identifier,
		value: identifier.Value,
	}
	return i
}

func (i *Identifier) expressionNode()               {}
func (i *Identifier) GetValue() string              { return i.value }
func (i *Identifier) GetType() token.TokenType      { return i._type }
func (i *Identifier) SetType(_type token.TokenType) { i._type = _type }

func (i *Identifier) GetTokenValue() string         { return i.token.Value }
func (i *Identifier) GetTokenType() token.TokenType { return token.IDENTIFIER }
func (i *Identifier) GetLineNumber() int            { return i.token.Line }
func (i *Identifier) GetPositionInLine() int        { return i.token.PositionInLine }
func (i *Identifier) GetChildren() []Node           { return []Node{} }
func (i *Identifier) String() string                { return i.value }

type Type struct {
	token token.Token
}

func NewType(tok token.Token) *Type {
	t := &Type{
		token: tok,
	}
	return t
}

func (t *Type) expressionNode()               {}
func (t *Type) GetTokenValue() string         { return t.token.Value }
func (t *Type) GetTokenType() token.TokenType { return token.TYPE }
func (t *Type) GetLineNumber() int            { return t.token.Line }
func (t *Type) GetPositionInLine() int        { return t.token.PositionInLine }
func (t *Type) GetChildren() []Node           { return []Node{} }
func (t *Type) String() string                { return t.token.Value }

type IntegerLiteral struct {
	token token.Token
	value int
}

func NewIntegerLiteral(tok token.Token) (*IntegerLiteral, error) {
	val, err := strconv.Atoi(tok.Value)
	l := &IntegerLiteral{
		token: tok,
		value: val,
	}
	return l, err
}

func (e *IntegerLiteral) expressionNode()               {}
func (e *IntegerLiteral) GetValue() int                 { return e.value }
func (e *IntegerLiteral) GetTokenValue() string         { return e.token.Value }
func (e *IntegerLiteral) GetTokenType() token.TokenType { return token.INT }
func (e *IntegerLiteral) GetLineNumber() int            { return e.token.Line }
func (e *IntegerLiteral) GetPositionInLine() int        { return e.token.PositionInLine }
func (e *IntegerLiteral) GetChildren() []Node           { return []Node{} }
func (e *IntegerLiteral) String() string                { return e.token.Value }

type FloatingPointLiteral struct {
	token token.Token
	value float64
}

func NewFloatingPointLiteral(tok token.Token) (*FloatingPointLiteral, error) {
	val, err := strconv.ParseFloat(tok.Value, 64)
	l := &FloatingPointLiteral{
		token: tok,
		value: val,
	}
	return l, err
}

func (e *FloatingPointLiteral) expressionNode()               {}
func (e *FloatingPointLiteral) GetValue() float64             { return e.value }
func (e *FloatingPointLiteral) GetTokenValue() string         { return e.token.Value }
func (e *FloatingPointLiteral) GetTokenType() token.TokenType { return token.FLOAT }
func (e *FloatingPointLiteral) GetLineNumber() int            { return e.token.Line }
func (e *FloatingPointLiteral) GetPositionInLine() int        { return e.token.PositionInLine }
func (e *FloatingPointLiteral) GetChildren() []Node           { return []Node{} }
func (e *FloatingPointLiteral) String() string                { return e.token.Value }

type BooleanLiteral struct {
	token token.Token
	value bool
}

func NewBooleanLiteral(tok token.Token) (*BooleanLiteral, error) {
	val, err := strconv.ParseBool(tok.Value)
	l := &BooleanLiteral{
		token: tok,
		value: val,
	}
	return l, err
}

func (e *BooleanLiteral) expressionNode()               {}
func (e *BooleanLiteral) GetValue() bool                { return e.value }
func (e *BooleanLiteral) GetTokenValue() string         { return e.token.Value }
func (e *BooleanLiteral) GetTokenType() token.TokenType { return token.BOOL }
func (e *BooleanLiteral) GetLineNumber() int            { return e.token.Line }
func (e *BooleanLiteral) GetPositionInLine() int        { return e.token.PositionInLine }
func (e *BooleanLiteral) GetChildren() []Node           { return []Node{} }
func (e *BooleanLiteral) String() string                { return e.token.Value }

type CharacterLiteral struct {
	token token.Token
	value rune
}

func NewCharacterLiteral(tok token.Token) *CharacterLiteral {
	val := []rune(tok.Value)[0]
	l := &CharacterLiteral{
		token: tok,
		value: val,
	}
	return l
}

func (e *CharacterLiteral) expressionNode()               {}
func (e *CharacterLiteral) GetValue() rune                { return e.value }
func (e *CharacterLiteral) GetTokenValue() string         { return e.token.Value }
func (e *CharacterLiteral) GetTokenType() token.TokenType { return token.CHAR }
func (e *CharacterLiteral) GetLineNumber() int            { return e.token.Line }
func (e *CharacterLiteral) GetPositionInLine() int        { return e.token.PositionInLine }
func (e *CharacterLiteral) GetChildren() []Node           { return []Node{} }
func (e *CharacterLiteral) String() string                { return e.token.Value }

type StringLiteral struct {
	token token.Token
	value string
}

func NewStringLiteral(tok token.Token) *StringLiteral {
	l := &StringLiteral{
		token: tok,
		value: tok.Value,
	}
	return l
}

func (e *StringLiteral) expressionNode()               {}
func (e *StringLiteral) GetValue() string              { return e.value }
func (e *StringLiteral) GetTokenValue() string         { return e.token.Value }
func (e *StringLiteral) GetTokenType() token.TokenType { return token.STRING }
func (e *StringLiteral) GetLineNumber() int            { return e.token.Line }
func (e *StringLiteral) GetPositionInLine() int        { return e.token.PositionInLine }
func (e *StringLiteral) GetChildren() []Node           { return []Node{} }
func (e *StringLiteral) String() string                { return e.token.Value }
