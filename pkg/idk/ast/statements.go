package ast

import (
	"bytes"

	"github.com/fglo/idk/pkg/idk/token"
)

type Statement interface {
	Node
	statementNode()
}

type ExpressionStatement struct {
	Expression Expression
}

func NewExpressionStatement(expression Expression) *ExpressionStatement {
	es := &ExpressionStatement{
		Expression: expression,
	}
	return es
}

func (es *ExpressionStatement) statementNode()                {}
func (es *ExpressionStatement) GetTokenValue() string         { return "" }
func (es *ExpressionStatement) GetTokenType() token.TokenType { return es.Expression.GetTokenType() }
func (es *ExpressionStatement) GetLineNumber() int            { return es.Expression.GetLineNumber() }
func (es *ExpressionStatement) GetPositionInLine() int        { return es.Expression.GetPositionInLine() }
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

func NewDeclareAssignStatement(identifier *Identifier, expression Expression) *DeclareAssignStatement {
	das := &DeclareAssignStatement{
		Identifier: identifier,
		Expression: expression,
	}
	return das
}

func (das *DeclareAssignStatement) statementNode()                {}
func (das *DeclareAssignStatement) GetTokenValue() string         { return "" }
func (das *DeclareAssignStatement) GetTokenType() token.TokenType { return token.DECLASSIGN }
func (das *DeclareAssignStatement) GetLineNumber() int            { return das.Identifier.GetLineNumber() }
func (das *DeclareAssignStatement) GetPositionInLine() int        { return das.Identifier.GetPositionInLine() }
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

func NewDeclareStatement(identifier *Identifier, assignment *AssignStatement) *DeclareStatement {
	ds := &DeclareStatement{
		Identifier: identifier,
		Assignment: assignment,
	}
	return ds
}

func (ds *DeclareStatement) statementNode()                {}
func (ds *DeclareStatement) GetTokenValue() string         { return "" }
func (ds *DeclareStatement) GetTokenType() token.TokenType { return token.DECLARE }
func (ds *DeclareStatement) GetLineNumber() int            { return ds.Identifier.GetLineNumber() }
func (ds *DeclareStatement) GetPositionInLine() int        { return ds.Identifier.GetPositionInLine() }
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
	out.WriteString(string(ds.Identifier.GetType()))

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

func NewAssignStatement(identifier *Identifier, expression Expression) *AssignStatement {
	as := &AssignStatement{
		Identifier: identifier,
		Expression: expression,
	}
	return as
}

func (as *AssignStatement) statementNode()                {}
func (as *AssignStatement) GetTokenValue() string         { return "" }
func (as *AssignStatement) GetTokenType() token.TokenType { return token.ASSIGN }
func (as *AssignStatement) GetLineNumber() int            { return as.Identifier.GetLineNumber() }
func (as *AssignStatement) GetPositionInLine() int        { return as.Identifier.GetPositionInLine() }
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

func NewIfStatement(condition Expression, consequence *BlockStatement, alternative *BlockStatement) *IfStatement {
	is := &IfStatement{
		Condition:   condition,
		Consequence: consequence,
		Alternative: alternative,
	}
	return is
}

func (is *IfStatement) statementNode()                {}
func (is *IfStatement) GetTokenValue() string         { return "" }
func (is *IfStatement) GetTokenType() token.TokenType { return token.IF }
func (is *IfStatement) GetLineNumber() int            { return is.Condition.GetLineNumber() }
func (is *IfStatement) GetPositionInLine() int        { return is.Condition.GetPositionInLine() }
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

func NewForLoopStatement(condition Expression, consequence *BlockStatement) *ForLoopStatement {
	fls := &ForLoopStatement{
		Condition:   condition,
		Consequence: consequence,
	}
	return fls
}

func (fls *ForLoopStatement) statementNode()                {}
func (fls *ForLoopStatement) GetTokenValue() string         { return "" }
func (fls *ForLoopStatement) GetTokenType() token.TokenType { return token.FOR }
func (fls *ForLoopStatement) GetLineNumber() int            { return fls.Condition.GetLineNumber() }
func (fls *ForLoopStatement) GetPositionInLine() int        { return fls.Condition.GetPositionInLine() }
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

func NewFunctionDefinitionStatement(identifier Identifier, parameters []*DeclareStatement, returnType token.Token, body *BlockStatement) *FunctionDefinitionStatement {
	fds := &FunctionDefinitionStatement{
		Identifier: identifier,
		Parameters: parameters,
		ReturnType: returnType,
		Body:       body,
	}
	return fds
}

func (fds *FunctionDefinitionStatement) statementNode()                {}
func (fds *FunctionDefinitionStatement) GetTokenValue() string         { return "" }
func (fds *FunctionDefinitionStatement) GetTokenType() token.TokenType { return token.FUNC }
func (fds *FunctionDefinitionStatement) GetLineNumber() int            { return fds.Identifier.GetLineNumber() }
func (fds *FunctionDefinitionStatement) GetPositionInLine() int {
	return fds.Identifier.GetPositionInLine()
}
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
		out.WriteString(parameter.String() + " : " + string(parameter.Identifier.GetType()) + ", ")
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

func NewReturnStatement(expression Expression) *ReturnStatement {
	rs := &ReturnStatement{
		Expression: expression,
	}
	return rs
}

func (rs *ReturnStatement) statementNode()                {}
func (rs *ReturnStatement) GetTokenValue() string         { return "" }
func (rs *ReturnStatement) GetTokenType() token.TokenType { return token.RETURN }
func (rs *ReturnStatement) GetLineNumber() int            { return rs.Expression.GetLineNumber() }
func (rs *ReturnStatement) GetPositionInLine() int        { return rs.Expression.GetPositionInLine() }
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

func NewBlockStatement(statements []Statement) *BlockStatement {
	bs := &BlockStatement{
		Statements: statements,
	}
	return bs
}

func (bs *BlockStatement) statementNode()                {}
func (bs *BlockStatement) GetTokenValue() string         { return "" }
func (bs *BlockStatement) GetTokenType() token.TokenType { return "─┐" }
func (bs *BlockStatement) GetLineNumber() int            { return bs.Statements[0].GetLineNumber() }
func (bs *BlockStatement) GetPositionInLine() int        { return bs.Statements[0].GetPositionInLine() }
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

type ImportStatement struct {
	Identifier *Identifier
}

func NewImportStatement(identifier token.Token) *ImportStatement {
	is := &ImportStatement{
		Identifier: NewIdentifier(identifier),
	}
	return is
}

func (is *ImportStatement) statementNode()                {}
func (is *ImportStatement) GetTokenValue() string         { return is.Identifier.GetTokenValue() }
func (is *ImportStatement) GetTokenType() token.TokenType { return token.IMPORT }
func (is *ImportStatement) GetLineNumber() int            { return is.Identifier.GetLineNumber() }
func (is *ImportStatement) GetPositionInLine() int        { return is.Identifier.GetPositionInLine() }
func (is *ImportStatement) GetChildren() []Node {
	return []Node{is.Identifier}
}
func (is *ImportStatement) String() string {
	var out bytes.Buffer

	out.WriteString(token.IMPORT.String())
	out.WriteString(" : ")
	out.WriteString(is.Identifier.GetTokenValue())

	return out.String()
}
