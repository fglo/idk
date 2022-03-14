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
	ue := new(UnaryExpression)
	ue.Token = Operator
	ue.Right = Right
	return ue
}

func (ue *UnaryExpression) expressionNode()               {}
func (ue *UnaryExpression) GetValue() string              { return "" }
func (ue *UnaryExpression) GetTokenType() token.TokenType { return ue.Token.Type }
func (ue *UnaryExpression) GetChildren() []Node           { return []Node{ue.Right} }
func (ue *UnaryExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ue.Token.Value)
	out.WriteString(ue.Right.String())
	out.WriteString(")")

	return out.String()
}

type BinaryExpression struct {
	Token token.Token
	Left  Expression
	Right Expression
}

func NewBinaryExpression(Left Expression, Operator token.Token, Right Expression) *BinaryExpression {
	be := new(BinaryExpression)
	be.Token = Operator
	be.Left = Left
	be.Right = Right
	return be
}

func (be *BinaryExpression) expressionNode()               {}
func (be *BinaryExpression) GetValue() string              { return "" }
func (be *BinaryExpression) GetTokenType() token.TokenType { return be.Token.Type }
func (be *BinaryExpression) GetChildren() []Node           { return []Node{be.Left, be.Right} }
func (be *BinaryExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(be.Left.String())
	out.WriteString(" " + be.Token.Value + " ")
	out.WriteString(be.Right.String())
	out.WriteString(")")

	return out.String()
}

type IfExpression struct {
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func NewIfExpression(Condition Expression, Consequence *BlockStatement, Alternative *BlockStatement) *IfExpression {
	ie := new(IfExpression)
	ie.Condition = Condition
	ie.Consequence = Consequence
	ie.Alternative = Alternative
	return ie
}

func (ie *IfExpression) expressionNode()               {}
func (ie *IfExpression) GetValue() string              { return "" }
func (ie *IfExpression) GetTokenType() token.TokenType { return token.IF }
func (ie *IfExpression) GetChildren() []Node {
	if ie.Alternative != nil {
		return []Node{ie.Condition, ie.Consequence, ie.Alternative}
	}
	return []Node{ie.Condition, ie.Consequence}
}
func (ie *IfExpression) String() string {
	var out bytes.Buffer

	out.WriteString("if")
	out.WriteString(ie.Condition.String())
	out.WriteString(" ")
	out.WriteString(ie.Consequence.String())

	if ie.Alternative != nil {
		out.WriteString("else ")
		out.WriteString(ie.Alternative.String())
	}

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
	das := new(DeclareAssignStatement)
	das.Token = Token
	das.Identifier = Identifier
	das.Expression = Expression
	return das
}

func (das *DeclareAssignStatement) statementNode()                {}
func (das *DeclareAssignStatement) GetValue() string              { return "" }
func (das *DeclareAssignStatement) GetTokenType() token.TokenType { return das.Token.Type }
func (das *DeclareAssignStatement) GetChildren() []Node {
	return []Node{das.Identifier, das.Expression}
}
func (das *DeclareAssignStatement) String() string {
	var out bytes.Buffer

	out.WriteString(das.Identifier.String())
	out.WriteString(" := ")

	if das.Expression != nil {
		out.WriteString(das.Expression.String())
	}

	return out.String()
}

type BlockStatement struct {
	Statements []Statement
}

func NewBlockStatement(Statements []Statement) *BlockStatement {
	bs := new(BlockStatement)
	bs.Statements = Statements
	return bs
}

func (bs *BlockStatement) statementNode()                {}
func (bs *BlockStatement) GetValue() string              { return "" }
func (bs *BlockStatement) GetTokenType() token.TokenType { return "─┐" }
func (bs *BlockStatement) GetChildren() []Node {
	var nodes []Node
	for _, s := range bs.Statements {
		nodes = append(nodes, s)
	}
	return nodes
}
func (bs *BlockStatement) String() string {
	var out bytes.Buffer

	for _, s := range bs.Statements {
		out.WriteString(s.String())
		out.WriteString(" ")
	}

	return out.String()
}
