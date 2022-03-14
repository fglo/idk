package ast

import (
	"bytes"
	"fmt"
	"idk/token"
	"strconv"
)

type Node interface {
	GetValue() string
	GetTokenType() token.TokenType
	GetChildren() []Node
	String() string
}

type Expression interface {
	Node
	expressionNode()
}

type Statement interface {
	Node
	statementNode()
}

type Program struct {
	Statements []Statement
}

func PrettyPrintProgram(program *Program) {
	for _, s := range program.Statements {
		fmt.Println(s)
		PrettyPrint(s, "", true)
	}
}

func PrettyPrint(node Node, indent string, isLast bool) {
	marker := "├──"
	if isLast {
		marker = "└──"
	}

	fmt.Print(indent)
	fmt.Print(marker)
	fmt.Print(node.GetTokenType())
	fmt.Print(" ")
	fmt.Println(node.GetValue())

	if !isLast {
		indent += "│   "
	} else {
		indent += "    "
	}

	children := node.GetChildren()
	for i, child := range children {
		PrettyPrint(child, indent, i == len(children)-1)
	}
}

/// expressions

type UnaryExpression struct {
	Token token.Token
	Right Expression
}

func NewUnaryExpression(Operator token.Token, Right Expression) *UnaryExpression {
	e := new(UnaryExpression)
	e.Token = Operator
	e.Right = Right
	return e
}

func (e *UnaryExpression) expressionNode()               {}
func (e *UnaryExpression) GetValue() string              { return "" }
func (e *UnaryExpression) GetTokenType() token.TokenType { return e.Token.Type }
func (e *UnaryExpression) GetChildren() []Node           { return []Node{e.Right} }
func (e *UnaryExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(e.Token.Value)
	out.WriteString(e.Right.String())
	out.WriteString(")")

	return out.String()
}

type BinaryExpression struct {
	Token token.Token
	Left  Expression
	Right Expression
}

func NewBinaryExpression(Left Expression, Operator token.Token, Right Expression) *BinaryExpression {
	e := new(BinaryExpression)
	e.Token = Operator
	e.Left = Left
	e.Right = Right
	return e
}

func (e *BinaryExpression) expressionNode()               {}
func (e *BinaryExpression) GetValue() string              { return "" }
func (e *BinaryExpression) GetTokenType() token.TokenType { return e.Token.Type }
func (e *BinaryExpression) GetChildren() []Node           { return []Node{e.Left, e.Right} }
func (e *BinaryExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(e.Left.String())
	out.WriteString(" " + e.Token.Value + " ")
	out.WriteString(e.Right.String())
	out.WriteString(")")

	return out.String()
}

type Identifier struct {
	Token token.Token
	Value string
}

func NewIdentifier(Token token.Token) *Identifier {
	l := new(Identifier)
	l.Token = Token
	l.Value = Token.Value
	return l
}

func (e *Identifier) expressionNode()               {}
func (e *Identifier) GetValue() string              { return e.Token.Value }
func (e *Identifier) GetTokenType() token.TokenType { return token.IDENTIFIER }
func (e *Identifier) GetChildren() []Node           { return []Node{} }
func (e *Identifier) String() string                { return e.Value }

type IntegerLiteral struct {
	Token token.Token
	Value int
}

func NewIntegerLiteral(Token token.Token) (*IntegerLiteral, error) {
	l := new(IntegerLiteral)
	l.Token = Token
	val, err := strconv.Atoi(Token.Value)
	l.Value = val
	return l, err
}

func (e *IntegerLiteral) expressionNode()               {}
func (e *IntegerLiteral) GetValue() string              { return e.Token.Value }
func (e *IntegerLiteral) GetTokenType() token.TokenType { return token.INT }
func (e *IntegerLiteral) GetChildren() []Node           { return []Node{} }
func (e *IntegerLiteral) String() string                { return e.Token.Value }

/// statements

type ExpressionStatement struct {
	Token      token.Token
	Expression Expression
}

func (s *ExpressionStatement) statementNode()                {}
func (s *ExpressionStatement) GetValue() string              { return "" }
func (s *ExpressionStatement) GetTokenType() token.TokenType { return s.Expression.GetTokenType() }
func (s *ExpressionStatement) GetChildren() []Node           { return s.Expression.GetChildren() }
func (s *ExpressionStatement) String() string {
	if s.Expression != nil {
		return s.Expression.String()
	}
	return ""
}

type DeclareAssignStatement struct {
	Token      token.Token
	Identifier *Identifier
	Expression Expression
}

func NewDeclareAssignStatement(Token token.Token, Identifier *Identifier, Expression Expression) *DeclareAssignStatement {
	s := new(DeclareAssignStatement)
	s.Token = Token
	s.Identifier = Identifier
	s.Expression = Expression
	return s
}

func (s *DeclareAssignStatement) statementNode()                {}
func (s *DeclareAssignStatement) GetValue() string              { return "" }
func (s *DeclareAssignStatement) GetTokenType() token.TokenType { return s.Token.Type }
func (s *DeclareAssignStatement) GetChildren() []Node           { return []Node{s.Identifier, s.Expression} }
func (s *DeclareAssignStatement) String() string {
	var out bytes.Buffer

	out.WriteString(s.Identifier.String())
	out.WriteString(" := ")

	if s.Expression != nil {
		out.WriteString(s.Expression.String())
	}

	return out.String()
}
