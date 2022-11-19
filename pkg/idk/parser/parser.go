package parser

import (
	"fmt"

	"github.com/fglo/idk/pkg/idk/ast"
	"github.com/fglo/idk/pkg/idk/lexer"
	"github.com/fglo/idk/pkg/idk/token"
)

const (
	_ int = iota
	LOWEST
	DECLARE_ASSIGN // :=
	DECLARE        // :
	ASSIGN         // =
	IN
	OR
	AND
	XOR
	NOT
	EQUALS      // ==
	LESSGREATER // > or <
	SUM         // +
	PRODUCT     // *
	PREFIX      // -X or !X
	RANGE       // .. or ..=
	CALL        // myFunction(X)
	INDEX       // array[index]
)

var precedences = map[token.TokenType]int{
	token.DECLARE_ASSIGN:  DECLARE_ASSIGN,
	token.DECLARE:         DECLARE,
	token.ASSIGN:          ASSIGN,
	token.IN:              IN,
	token.AND:             AND,
	token.OR:              OR,
	token.XOR:             XOR,
	token.NOT:             NOT,
	token.EQ:              EQUALS,
	token.NEQ:             EQUALS,
	token.LT:              LESSGREATER,
	token.GT:              LESSGREATER,
	token.PLUS:            SUM,
	token.MINUS:           SUM,
	token.SLASH:           PRODUCT,
	token.ASTERISK:        PRODUCT,
	token.RANGE:           RANGE,
	token.RANGE_INCLUSIVE: RANGE,
	token.LPARENTHESIS:    CALL,
}

type (
	prefixParseFn func() ast.Expression
	binaryParseFn func(ast.Expression) ast.Expression
)

type Parser struct {
	input string
	lexer *lexer.Lexer

	previous token.Token
	current  token.Token
	next     token.Token

	unaryParseFns  map[token.TokenType]prefixParseFn
	binaryParseFns map[token.TokenType]binaryParseFn

	errors []string
}

func NewParser(input string) *Parser {
	p := new(Parser)
	p.input = input
	p.lexer = lexer.NewLexer(input)

	p.unaryParseFns = make(map[token.TokenType]prefixParseFn)
	p.registerUnary(token.MINUS, p.parseUnaryExpression)
	p.registerUnary(token.IDENTIFIER, p.parseIdentifier)
	p.registerUnary(token.INT, p.parseIntegerLiteral)
	p.registerUnary(token.BOOL, p.parseBooleanLiteral)
	p.registerUnary(token.CHAR, p.parseCharacterLiteral)
	p.registerUnary(token.STRING, p.parseStringLiteral)
	p.registerUnary(token.LPARENTHESIS, p.parseGroupedExpression)
	p.registerUnary(token.NOT, p.parseUnaryExpression)

	p.binaryParseFns = make(map[token.TokenType]binaryParseFn)
	p.registerBinary(token.PLUS, p.parseBinaryExpression)
	p.registerBinary(token.MINUS, p.parseBinaryExpression)
	p.registerBinary(token.SLASH, p.parseBinaryExpression)
	p.registerBinary(token.ASTERISK, p.parseBinaryExpression)
	p.registerBinary(token.IN, p.parseBinaryExpression)
	p.registerBinary(token.RANGE, p.parseBinaryExpression)

	p.registerBinary(token.DECLARE_ASSIGN, p.parseBinaryExpression)

	p.registerBinary(token.EQ, p.parseBinaryExpression)
	p.registerBinary(token.NEQ, p.parseBinaryExpression)
	p.registerBinary(token.GT, p.parseBinaryExpression)
	p.registerBinary(token.GTE, p.parseBinaryExpression)
	p.registerBinary(token.LT, p.parseBinaryExpression)
	p.registerBinary(token.LTE, p.parseBinaryExpression)

	p.registerBinary(token.AND, p.parseBinaryExpression)
	p.registerBinary(token.OR, p.parseBinaryExpression)
	p.registerBinary(token.XOR, p.parseBinaryExpression)

	p.consumeToken()

	return p
}

func (p *Parser) registerUnary(tokenType token.TokenType, fn prefixParseFn) {
	p.unaryParseFns[tokenType] = fn
}

func (p *Parser) registerBinary(tokenType token.TokenType, fn binaryParseFn) {
	p.binaryParseFns[tokenType] = fn
}

func (p *Parser) peekNext() token.Token {
	return p.next
}

func (p *Parser) consumeToken() token.Token {
	p.previous = p.current
	p.current = p.next
	if p.current.Type == token.ILLEGAL {
		p.reportIllegalToken()
	}
	p.next = p.lexer.ReadToken()
	return p.current
}

func (p *Parser) previousTokenWas(t token.TokenType) bool {
	return p.previous.Type == t
}

func (p *Parser) currentTokenIs(t token.TokenType) bool {
	return p.current.Type == t
}

func (p *Parser) nextTokenIs(t token.TokenType) bool {
	return p.next.Type == t
}

func (p *Parser) expectCurrentTokenType(t token.TokenType) bool {
	if p.currentTokenIs(t) {
		return true
	} else {
		p.reportUnexpectedToken(p.next, t)
		return false
	}
}

func (p *Parser) expectNextTokenType(t token.TokenType) bool {
	if p.nextTokenIs(t) {
		return true
	} else {
		p.reportUnexpectedToken(p.next, t)
		return false
	}
}

func (p *Parser) expectOperatorOrEolOrEof() bool {
	if p.next.Type.IsOperator() || p.nextTokenIs(token.EOL) || p.nextTokenIs(token.EOF) || p.nextTokenIs(token.RPARENTHESIS) {
		return true
	} else {
		p.reportExpectedOperatorOrEolOrEof(p.next)
		return false
	}
}

func (p *Parser) currentPrecedence() int {
	if p, ok := precedences[p.current.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) nextPrecedence() int {
	if p, ok := precedences[p.next.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) reportUnexpectedToken(unexpected token.Token, expectedType token.TokenType) {
	msg := fmt.Sprintf("Unexpected token <%v> on line %v, position %v. <%v> was expected.",
		unexpected.Type,
		unexpected.Line,
		unexpected.PositionInLine,
		expectedType)
	p.errors = append(p.errors, msg)
}

func (p *Parser) reportUnexpectedFirstToken(unexpected token.Token) {
	msg := fmt.Sprintf("Unexpected token <%v> on line %v, position %v. Expected declaration or a statement.",
		unexpected.Type,
		unexpected.Line,
		unexpected.PositionInLine)
	p.errors = append(p.errors, msg)
}

func (p *Parser) reportExpectedOperatorOrEolOrEof(unexpected token.Token) {
	msg := fmt.Sprintf("Unexpected token <%v> on line %v, position %v. Expected operator, <EOL> or <EOF>.",
		unexpected.Type,
		unexpected.Line,
		unexpected.PositionInLine)
	p.errors = append(p.errors, msg)
}

func (p *Parser) reportIllegalToken() {
	msg := fmt.Sprintf("Illegal token '%v' on line %v, position %v.",
		p.next.Value,
		p.next.Line,
		p.next.PositionInLine)
	p.errors = append(p.errors, msg)
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) ifEolIsNextThenSkip() {
	for p.nextTokenIs(token.EOL) {
		p.consumeToken()
	}
}

/// PARSING

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for !p.nextTokenIs(token.EOF) {
		p.ifEolIsNextThenSkip()
		p.consumeToken()
		s := p.parseStatement()
		if s != nil {
			program.Statements = append(program.Statements, s)
		}
	}

	return program
}

func (p *Parser) parseStatement() ast.Statement {
	switch {
	case p.currentTokenIs(token.IDENTIFIER) && p.nextTokenIs(token.DECLARE_ASSIGN):
		return p.parseDeclareAssignStatement()
	case p.currentTokenIs(token.IDENTIFIER) && p.nextTokenIs(token.DECLARE):
		return p.parseDeclareStatement()
	case p.currentTokenIs(token.IDENTIFIER) && p.nextTokenIs(token.ASSIGN):
		return p.parseAssignStatement()
	case p.currentTokenIs(token.IDENTIFIER) && p.nextTokenIs(token.LPARENTHESIS):
		return p.parseFunctionCallStatement()
	case p.currentTokenIs(token.IF):
		return p.parseIfStatement()
	case p.currentTokenIs(token.FOR):
		return p.parseForStatement()
	case p.currentTokenIs(token.FUNC):
		return p.parseFunctionDefinitionStatement()
	case p.currentTokenIs(token.RETURN):
		return p.parseReturnStatement()
	default:
		p.reportUnexpectedFirstToken(p.current)
		return nil
		// return p.parseExpressionStatement()
	}
}

func (p *Parser) parseDeclareAssignStatement() *ast.DeclareAssignStatement {
	identifier := ast.NewIdentifier(p.current)

	p.consumeToken() // declare-assign operator
	p.consumeToken() // skip the declare-assign operator

	expr := p.parseExpression(LOWEST)

	p.ifEolIsNextThenSkip()

	if expr == nil {
		return nil
	}

	return ast.NewDeclareAssignStatement(identifier, expr)
}

func (p *Parser) parseDeclareStatement() *ast.DeclareStatement {
	identifier := ast.NewIdentifier(p.current)

	p.consumeToken() // declare operator
	vartype := p.consumeToken()

	identifier.Type = token.LookupType(vartype.Value)

	p.consumeToken() // skip type declaration

	var ass *ast.AssignStatement

	if p.currentTokenIs(token.ASSIGN) {
		p.consumeToken() // skip the assign operator

		expr := p.parseExpression(LOWEST)
		if expr != nil {
			ass = ast.NewAssignStatement(identifier, expr)
		}

	}

	p.ifEolIsNextThenSkip()

	return ast.NewDeclareStatement(identifier, ass)
}

func (p *Parser) parseAssignStatement() *ast.AssignStatement {
	identifier := ast.NewIdentifier(p.current)

	p.consumeToken() // assign operator
	p.consumeToken() // skip the assign operator

	expr := p.parseExpression(LOWEST)

	p.ifEolIsNextThenSkip()

	if expr == nil {
		return nil
	}

	return ast.NewAssignStatement(identifier, expr)
}

func (p *Parser) parseIfStatement() *ast.IfStatement {
	innerIf := false
	if p.previousTokenWas(token.ELSE) {
		innerIf = true
	}
	p.consumeToken() // skip the if keyword

	condition := p.parseExpression(LOWEST)

	p.ifEolIsNextThenSkip()

	consequence := p.parseBlockStatement()

	p.ifEolIsNextThenSkip()

	var alternative *ast.BlockStatement
	if p.nextTokenIs(token.ELSE) {
		p.consumeToken()
		alternative = p.parseBlockStatement()
	}
	if !innerIf && p.expectNextTokenType(token.END) {
		p.consumeToken()
	}

	p.ifEolIsNextThenSkip()

	return ast.NewIfStatement(condition, consequence, alternative)
}

func (p *Parser) parseForStatement() *ast.ForLoopStatement {
	p.consumeToken() // skip the for keyword

	condition := p.parseExpression(LOWEST)

	p.ifEolIsNextThenSkip()

	consequence := p.parseBlockStatement()

	p.ifEolIsNextThenSkip()

	if p.expectNextTokenType(token.END) {
		p.consumeToken()
	}

	return ast.NewForLoopStatement(condition, consequence)
}

func (p *Parser) parseFunctionDefinitionStatement() *ast.FunctionDefinitionStatement {
	p.consumeToken() // skip func keyword

	identifier := ast.NewIdentifier(p.current)

	p.consumeToken() // parameters

	parameters := p.parseFunctionDefinitionParametersList()

	vartype := p.consumeToken()

	identifier.Type = token.FUNC

	p.ifEolIsNextThenSkip()

	body := p.parseBlockStatement()

	if p.expectNextTokenType(token.END) {
		p.consumeToken()
	}

	p.ifEolIsNextThenSkip()

	return ast.NewFunctionDefinitionStatement(*identifier, parameters, vartype, body)
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	p.consumeToken() // return keyword

	expr := p.parseExpression(LOWEST)

	p.ifEolIsNextThenSkip()

	if expr == nil {
		return nil
	}

	return ast.NewReturnStatement(expr)
}

func (p *Parser) parseFunctionCallStatement() *ast.ExpressionStatement {
	stmt := new(ast.ExpressionStatement)
	stmt.Expression = p.parseFunctionCallExpression()
	p.ifEolIsNextThenSkip()
	return stmt
}

func (p *Parser) parseFunctionCallExpression() ast.Expression {
	exp := ast.NewFunctionCallExpression(p.current)
	p.consumeToken()
	exp.Parameters = p.parseFunctionCallParametersList()
	return exp
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	statements := []ast.Statement{}

	p.ifEolIsNextThenSkip()

	for !p.nextTokenIs(token.END) && !p.nextTokenIs(token.EOF) && !p.nextTokenIs(token.ELSE) {
		p.consumeToken()
		s := p.parseStatement()
		if s != nil {
			statements = append(statements, s)
		}
		p.ifEolIsNextThenSkip()
	}

	return ast.NewBlockStatement(statements)
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := new(ast.ExpressionStatement)
	stmt.Expression = p.parseExpression(LOWEST)
	p.ifEolIsNextThenSkip()
	return stmt
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	parsePrefix := p.unaryParseFns[p.current.Type]
	if parsePrefix == nil {
		return nil
	}
	expr := parsePrefix()

	for !p.nextTokenIs(token.EOL) && precedence < p.nextPrecedence() {
		parseInfix := p.binaryParseFns[p.peekNext().Type]
		if parseInfix == nil {
			return expr
		}
		p.consumeToken()
		expr = parseInfix(expr)

		p.expectOperatorOrEolOrEof()
	}

	return expr
}

func (p *Parser) parseFunctionDefinitionParametersList() []*ast.DeclareStatement {
	list := []*ast.DeclareStatement{}

	if p.nextTokenIs(token.RPARENTHESIS) {
		p.consumeToken()
		return list
	}

	p.consumeToken()
	list = append(list, p.parseDeclareStatement())

	for p.currentTokenIs(token.COMMA) {
		p.consumeToken()
		list = append(list, p.parseDeclareStatement())
	}

	if !p.expectCurrentTokenType(token.RPARENTHESIS) {
		return nil
	}

	return list
}

func (p *Parser) parseFunctionCallParametersList() []ast.Expression {
	list := []ast.Expression{}

	if p.nextTokenIs(token.RPARENTHESIS) {
		p.consumeToken()
		return list
	}

	p.consumeToken()
	list = append(list, p.parseExpression(LOWEST))

	for p.nextTokenIs(token.COMMA) {
		p.consumeToken()
		p.consumeToken()
		list = append(list, p.parseExpression(LOWEST))
	}

	if !p.expectNextTokenType(token.RPARENTHESIS) {
		return nil
	}
	p.consumeToken()

	return list
}

func (p *Parser) parseUnaryExpression() ast.Expression {
	operator := p.current
	p.consumeToken() // skip the operator
	right := p.parseExpression(PREFIX)
	expression := ast.NewUnaryExpression(operator, right)
	return expression
}

func (p *Parser) parseBinaryExpression(left ast.Expression) ast.Expression {
	operator := p.current
	precedence := p.currentPrecedence()
	p.consumeToken() // skip the operator
	right := p.parseExpression(precedence)
	expr := ast.NewBinaryExpression(left, operator, right)
	return expr
}

func (p *Parser) parseIdentifier() ast.Expression {
	if p.currentTokenIs(token.IDENTIFIER) && p.nextTokenIs(token.LPARENTHESIS) {
		return p.parseFunctionCallExpression()
	}
	return ast.NewIdentifier(p.current)
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit, _ := ast.NewIntegerLiteral(p.current)
	return lit
}

func (p *Parser) parseBooleanLiteral() ast.Expression {
	lit, _ := ast.NewBooleanLiteral(p.current)
	return lit
}

func (p *Parser) parseCharacterLiteral() ast.Expression {
	lit := ast.NewCharacterLiteral(p.current)
	return lit
}

func (p *Parser) parseStringLiteral() ast.Expression {
	lit := ast.NewStringLiteral(p.current)
	return lit
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.consumeToken() //skip opening parenthesis
	exp := p.parseExpression(LOWEST)
	if !p.expectNextTokenType(token.RPARENTHESIS) {
		return nil
	}
	p.consumeToken() //skip closing parenthesis
	return exp
}
