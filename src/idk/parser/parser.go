package parser

import (
	"fmt"
	"idk/ast"
	"idk/lexer"
	"idk/token"
)

const (
	_ int = iota
	LOWEST
	DECLARE_ASSIGN // :=
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

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]binaryParseFn

	errors []string
}

func NewParser(input string) *Parser {
	p := new(Parser)
	p.input = input
	p.lexer = lexer.NewLexer(input)

	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.registerPrefix(token.MINUS, p.parseUnaryExpression)
	p.registerPrefix(token.IDENTIFIER, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.LPARENTHESIS, p.parseGroupedExpression)

	// p.registerPrefix(token.NEGATION, p.parseUnaryExpression)
	// p.registerPrefix(token.NOT, p.parseUnaryExpression)

	p.infixParseFns = make(map[token.TokenType]binaryParseFn)
	p.registerInfix(token.PLUS, p.parseBinaryExpression)
	p.registerInfix(token.MINUS, p.parseBinaryExpression)
	p.registerInfix(token.SLASH, p.parseBinaryExpression)
	p.registerInfix(token.ASTERISK, p.parseBinaryExpression)
	p.registerInfix(token.IN, p.parseBinaryExpression)
	p.registerInfix(token.RANGE, p.parseBinaryExpression)

	p.registerInfix(token.DECLARE_ASSIGN, p.parseBinaryExpression)

	p.registerInfix(token.EQ, p.parseBinaryExpression)
	p.registerInfix(token.NEQ, p.parseBinaryExpression)
	p.registerInfix(token.GT, p.parseBinaryExpression)
	p.registerInfix(token.GTE, p.parseBinaryExpression)
	p.registerInfix(token.LT, p.parseBinaryExpression)
	p.registerInfix(token.LTE, p.parseBinaryExpression)

	p.registerInfix(token.AND, p.parseBinaryExpression)
	p.registerInfix(token.OR, p.parseBinaryExpression)
	p.registerInfix(token.XOR, p.parseBinaryExpression)

	p.consumeToken()

	return p
}

func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.TokenType, fn binaryParseFn) {
	p.infixParseFns[tokenType] = fn
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
	// case p.currentTokenIs(token.IDENTIFIER) && p.nextTokenIs(token.EQ):
	// 	return p.parseAssignmentStatement()
	case p.currentTokenIs(token.IDENTIFIER) && p.nextTokenIs(token.LPARENTHESIS):
		return p.parseFunctionCallStatement()
	case p.currentTokenIs(token.IF):
		return p.parseIfStatement()
	case p.currentTokenIs(token.FOR):
		return p.parseForStatement()
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
		// innerIf = p.nextTokenIs(token.IF)
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

	return ast.NewForLoopStatement(condition, consequence)
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
	exp.Parameters = p.parseFunctionParametersList()
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
	parsePrefix := p.prefixParseFns[p.current.Type]
	if parsePrefix == nil {
		return nil
	}
	expr := parsePrefix()

	p.expectOperatorOrEolOrEof()

	for !p.nextTokenIs(token.EOL) && precedence < p.nextPrecedence() {
		parseInfix := p.infixParseFns[p.peekNext().Type]
		if parseInfix == nil {
			return expr
		}
		p.consumeToken()
		expr = parseInfix(expr)

		p.expectOperatorOrEolOrEof()
	}

	return expr
}

func (p *Parser) parseFunctionParametersList() []ast.Expression {
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
	return ast.NewIdentifier(p.current)
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit, _ := ast.NewIntegerLiteral(p.current)
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
