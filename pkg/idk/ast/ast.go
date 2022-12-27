package ast

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/fglo/idk/pkg/idk/token"
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

func (p *Program) GetValue() string              { return "" }
func (p *Program) GetTokenType() token.TokenType { return "" }
func (p *Program) GetChildren() []Node           { return []Node{} }
func (p *Program) String() string {
	var out bytes.Buffer

	for _, s := range p.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

func PrettyPrintProgram(program *Program) {
	for i, s := range program.Statements {
		fmt.Println(s)
		PrettyPrint(s, "", i == len(program.Statements)-1)
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

type PrefixExpression struct {
	Token token.Token
	Right Expression
}

func NewPrefixExpression(Operator token.Token, Right Expression) *PrefixExpression {
	pe := new(PrefixExpression)
	pe.Token = Operator
	pe.Right = Right
	return pe
}

func (pe *PrefixExpression) expressionNode()               {}
func (pe *PrefixExpression) GetValue() string              { return "" }
func (pe *PrefixExpression) GetTokenType() token.TokenType { return pe.Token.Type }
func (pe *PrefixExpression) GetChildren() []Node           { return []Node{pe.Right} }
func (pe *PrefixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(pe.Token.Value)
	out.WriteString(pe.Right.String())
	out.WriteString(")")

	return out.String()
}

type InfixExpression struct {
	Token token.Token
	Left  Expression
	Right Expression
}

func NewInfixExpression(Left Expression, Operator token.Token, Right Expression) *InfixExpression {
	ie := new(InfixExpression)
	ie.Token = Operator
	ie.Left = Left
	ie.Right = Right
	return ie
}

func (ie *InfixExpression) expressionNode()               {}
func (ie *InfixExpression) GetValue() string              { return "" }
func (ie *InfixExpression) GetTokenType() token.TokenType { return ie.Token.Type }
func (ie *InfixExpression) GetChildren() []Node           { return []Node{ie.Left, ie.Right} }
func (ie *InfixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString(" " + ie.Token.Value + " ")
	out.WriteString(ie.Right.String())
	out.WriteString(")")

	return out.String()
}

type IfExpression struct {
	Condition   Expression
	Consequence *Expression
	Alternative *Expression
}

// TODO: implement if expression
func NewIfExpression(Condition Expression, Consequence *Expression, Alternative *Expression) *IfExpression {
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

func NewFunctionCallExpression(Identifier token.Token) *FunctionCallExpression {
	fce := new(FunctionCallExpression)
	fce.Identifier = *NewIdentifier(Identifier)
	return fce
}

func (fce *FunctionCallExpression) expressionNode()               {}
func (fce *FunctionCallExpression) GetValue() string              { return fce.Identifier.Value }
func (fce *FunctionCallExpression) GetTokenType() token.TokenType { return fce.Identifier.Type }
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
	Token token.Token
	Type  token.TokenType
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
func (e *Identifier) String() string                { return e.Token.Value }

type Type struct {
	Token token.Token
}

func NewType(Token token.Token) *Type {
	t := new(Type)
	t.Token = Token
	return t
}

func (t *Type) expressionNode()               {}
func (t *Type) GetValue() string              { return t.Token.Value }
func (t *Type) GetTokenType() token.TokenType { return token.TYPE }
func (t *Type) GetChildren() []Node           { return []Node{} }
func (t *Type) String() string                { return t.Token.Value }

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

type FloatingPointLiteral struct {
	Token token.Token
	Value float64
}

func NewFloatingPointLiteral(Token token.Token) (*FloatingPointLiteral, error) {
	l := new(FloatingPointLiteral)
	l.Token = Token
	val, err := strconv.ParseFloat(Token.Value, 64)
	l.Value = val
	return l, err
}

func (e *FloatingPointLiteral) expressionNode()               {}
func (e *FloatingPointLiteral) GetValue() string              { return e.Token.Value }
func (e *FloatingPointLiteral) GetTokenType() token.TokenType { return token.FLOAT }
func (e *FloatingPointLiteral) GetChildren() []Node           { return []Node{} }
func (e *FloatingPointLiteral) String() string                { return e.Token.Value }

type BooleanLiteral struct {
	Token token.Token
	Value bool
}

func NewBooleanLiteral(Token token.Token) (*BooleanLiteral, error) {
	l := new(BooleanLiteral)
	l.Token = Token
	val, err := strconv.ParseBool(Token.Value)
	l.Value = val
	return l, err
}

func (e *BooleanLiteral) expressionNode()               {}
func (e *BooleanLiteral) GetValue() string              { return e.Token.Value }
func (e *BooleanLiteral) GetTokenType() token.TokenType { return token.BOOL }
func (e *BooleanLiteral) GetChildren() []Node           { return []Node{} }
func (e *BooleanLiteral) String() string                { return e.Token.Value }

type CharacterLiteral struct {
	Token token.Token
	Value rune
}

func NewCharacterLiteral(Token token.Token) *CharacterLiteral {
	l := new(CharacterLiteral)
	l.Token = Token
	val := []rune(Token.Value)[0]
	l.Value = val
	return l
}

func (e *CharacterLiteral) expressionNode()               {}
func (e *CharacterLiteral) GetValue() string              { return e.Token.Value }
func (e *CharacterLiteral) GetTokenType() token.TokenType { return token.CHAR }
func (e *CharacterLiteral) GetChildren() []Node           { return []Node{} }
func (e *CharacterLiteral) String() string                { return e.Token.Value }

type StringLiteral struct {
	Token token.Token
	Value string
}

func NewStringLiteral(Token token.Token) *StringLiteral {
	l := new(StringLiteral)
	l.Token = Token
	l.Value = Token.Value
	return l
}

func (e *StringLiteral) expressionNode()               {}
func (e *StringLiteral) GetValue() string              { return e.Token.Value }
func (e *StringLiteral) GetTokenType() token.TokenType { return token.STRING }
func (e *StringLiteral) GetChildren() []Node           { return []Node{} }
func (e *StringLiteral) String() string                { return e.Token.Value }

/// statements

type ExpressionStatement struct {
	Expression Expression
}

func NewExpressionStatement(Expression Expression) *ExpressionStatement {
	es := new(ExpressionStatement)
	es.Expression = Expression
	return es
}

func (es *ExpressionStatement) statementNode()                {}
func (es *ExpressionStatement) GetValue() string              { return "" }
func (es *ExpressionStatement) GetTokenType() token.TokenType { return es.Expression.GetTokenType() }
func (es *ExpressionStatement) GetChildren() []Node           { return es.Expression.GetChildren() }
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

type DeclareAssignStatement struct {
	Identifier *Identifier
	Expression Expression
}

func NewDeclareAssignStatement(Identifier *Identifier, Expression Expression) *DeclareAssignStatement {
	das := new(DeclareAssignStatement)
	das.Identifier = Identifier
	das.Expression = Expression
	return das
}

func (das *DeclareAssignStatement) statementNode()                {}
func (das *DeclareAssignStatement) GetValue() string              { return "" }
func (das *DeclareAssignStatement) GetTokenType() token.TokenType { return token.DECLASSIGN }
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

type DeclareStatement struct {
	Identifier *Identifier
	Assignment *AssignStatement
}

func NewDeclareStatement(Identifier *Identifier, Assignment *AssignStatement) *DeclareStatement {
	ds := new(DeclareStatement)
	ds.Identifier = Identifier
	ds.Assignment = Assignment
	return ds
}

func (ds *DeclareStatement) statementNode()                {}
func (ds *DeclareStatement) GetValue() string              { return "" }
func (ds *DeclareStatement) GetTokenType() token.TokenType { return token.DECLARE }
func (ds *DeclareStatement) GetChildren() []Node {
	if ds.Assignment != nil {
		return []Node{ds.Identifier, ds.Assignment}
	}
	return []Node{ds.Identifier}
}
func (ds *DeclareStatement) String() string {
	var out bytes.Buffer

	out.WriteString(ds.Identifier.String())
	out.WriteString(" : ")
	out.WriteString(string(ds.Identifier.Type))

	if ds.Assignment != nil {
		out.WriteString(" = ")
		if ds.Assignment.Expression != nil {
			out.WriteString(ds.Assignment.Expression.String())
		}
	}

	return out.String()
}

type AssignStatement struct {
	Identifier *Identifier
	Expression Expression
}

func NewAssignStatement(Identifier *Identifier, Expression Expression) *AssignStatement {
	as := new(AssignStatement)
	as.Identifier = Identifier
	as.Expression = Expression
	return as
}

func (as *AssignStatement) statementNode()                {}
func (as *AssignStatement) GetValue() string              { return "" }
func (as *AssignStatement) GetTokenType() token.TokenType { return token.ASSIGN }
func (as *AssignStatement) GetChildren() []Node {
	return []Node{as.Identifier, as.Expression}
}
func (as *AssignStatement) String() string {
	var out bytes.Buffer

	out.WriteString(as.Identifier.String())
	out.WriteString(" = ")

	if as.Expression != nil {
		out.WriteString(as.Expression.String())
	}

	return out.String()
}

type IfStatement struct {
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func NewIfStatement(Condition Expression, Consequence *BlockStatement, Alternative *BlockStatement) *IfStatement {
	is := new(IfStatement)
	is.Condition = Condition
	is.Consequence = Consequence
	is.Alternative = Alternative
	return is
}

func (is *IfStatement) statementNode()                {}
func (is *IfStatement) GetValue() string              { return "" }
func (is *IfStatement) GetTokenType() token.TokenType { return token.IF }
func (is *IfStatement) GetChildren() []Node {
	if is.Alternative != nil {
		return []Node{is.Condition, is.Consequence, is.Alternative}
	}
	return []Node{is.Condition, is.Consequence}
}
func (is *IfStatement) String() string {
	var out bytes.Buffer

	out.WriteString("{")
	out.WriteString("if")
	out.WriteString(is.Condition.String())
	out.WriteString(" ")
	out.WriteString(is.Consequence.String())

	if is.Alternative != nil {
		out.WriteString("else ")
		out.WriteString(is.Alternative.String())
	}
	out.WriteString("}")

	return out.String()
}

type ForLoopStatement struct {
	Condition   Expression
	Consequence *BlockStatement
}

func NewForLoopStatement(Condition Expression, Consequence *BlockStatement) *ForLoopStatement {
	fls := new(ForLoopStatement)
	fls.Condition = Condition
	fls.Consequence = Consequence
	return fls
}

func (fls *ForLoopStatement) statementNode()                {}
func (fls *ForLoopStatement) GetValue() string              { return "" }
func (fls *ForLoopStatement) GetTokenType() token.TokenType { return token.FOR }
func (fls *ForLoopStatement) GetChildren() []Node {
	return []Node{fls.Condition, fls.Consequence}
}
func (fls *ForLoopStatement) String() string {
	var out bytes.Buffer

	out.WriteString("{")
	out.WriteString("for")
	out.WriteString(fls.Condition.String())
	out.WriteString(" ")
	out.WriteString(fls.Consequence.String())
	out.WriteString("}")

	return out.String()
}

type FunctionDefinitionStatement struct {
	Identifier Identifier
	Parameters []*DeclareStatement
	ReturnType token.Token
	Body       *BlockStatement
}

func NewFunctionDefinitionStatement(Identifier Identifier, Parameters []*DeclareStatement, ReturnType token.Token, Body *BlockStatement) *FunctionDefinitionStatement {
	fds := new(FunctionDefinitionStatement)
	fds.Identifier = Identifier
	fds.Parameters = Parameters
	fds.ReturnType = ReturnType
	fds.Body = Body
	return fds
}

func (fds *FunctionDefinitionStatement) statementNode()                {}
func (fds *FunctionDefinitionStatement) GetValue() string              { return "" }
func (fds *FunctionDefinitionStatement) GetTokenType() token.TokenType { return token.FUNC }
func (fds *FunctionDefinitionStatement) GetChildren() []Node {
	return []Node{fds.Body}
}
func (fds *FunctionDefinitionStatement) String() string {
	var out bytes.Buffer

	out.WriteString("{")
	out.WriteString("func ")
	out.WriteString(fds.Identifier.String())
	out.WriteString(" ")
	out.WriteString("(")
	for _, parameter := range fds.Parameters {
		out.WriteString(parameter.String() + " : " + string(parameter.Identifier.Type) + ", ")
	}
	out.WriteString(")")
	out.WriteString(" ")
	out.WriteString(fds.Body.String())
	out.WriteString("}")

	return out.String()
}

type ReturnStatement struct {
	Expression Expression
}

func NewReturnStatement(Expression Expression) *ReturnStatement {
	rs := new(ReturnStatement)
	rs.Expression = Expression
	return rs
}

func (rs *ReturnStatement) statementNode()                {}
func (rs *ReturnStatement) GetValue() string              { return "" }
func (rs *ReturnStatement) GetTokenType() token.TokenType { return token.RETURN }
func (rs *ReturnStatement) GetChildren() []Node           { return rs.Expression.GetChildren() }
func (rs *ReturnStatement) String() string {
	if rs.Expression != nil {
		return rs.Expression.String()
	}
	return ""
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
