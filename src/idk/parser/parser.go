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
	OR
	AND
	XOR
	NOT
	EQUALS      // ==
	LESSGREATER // > or <
	SUM         // +
	PRODUCT     // *
	PREFIX      // -X or !X
	CALL        // myFunction(X)
	INDEX       // array[index]
)

var precedences = map[token.TokenType]int{
	token.DECLARE_ASSIGN: DECLARE_ASSIGN,
	token.AND:            AND,
	token.OR:             OR,
	token.XOR:            XOR,
	token.NOT:            NOT,
	token.EQ:             EQUALS,
	token.NEQ:            EQUALS,
	token.LT:             LESSGREATER,
	token.GT:             LESSGREATER,
	token.PLUS:           SUM,
	token.MINUS:          SUM,
	token.SLASH:          PRODUCT,
	token.ASTERISK:       PRODUCT,
	token.LPARENTHESIS:   CALL,
}

type (
	prefixParseFn func() ast.Expression
	binaryParseFn func(ast.Expression) ast.Expression
)

type Parser struct {
	input string
	lexer *lexer.Lexer

	current token.Token
	next    token.Token

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

	// p.registerPrefix(token.NEGATION, p.parseGroupedExpression)
	// p.registerPrefix(token.NOT, p.parseGroupedExpression)

	p.registerPrefix(token.IF, p.parseIfExpression)
	// p.registerPrefix(token.FOR, p.parseGroupedExpression)

	p.infixParseFns = make(map[token.TokenType]binaryParseFn)
	p.registerBinary(token.PLUS, p.parseBinaryExpression)
	p.registerBinary(token.MINUS, p.parseBinaryExpression)
	p.registerBinary(token.SLASH, p.parseBinaryExpression)
	p.registerBinary(token.ASTERISK, p.parseBinaryExpression)

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

	p.nextToken()

	return p
}

func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerBinary(tokenType token.TokenType, fn binaryParseFn) {
	p.infixParseFns[tokenType] = fn
}

func (p *Parser) nextToken() token.Token {
	p.current = p.next
	if p.current.Type == token.ILLEGAL {
		p.illegalToken()
	}
	p.next = p.lexer.ReadToken()
	return p.current
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
		p.unexpectedToken(p.current, t)
		return false
	}
}

func (p *Parser) expectNextTokenType(t token.TokenType) bool {
	if p.nextTokenIs(t) {
		return true
	} else {
		p.unexpectedToken(p.next, t)
		return false
	}
}

func (p *Parser) unexpectedToken(unexpected token.Token, expectedType token.TokenType) {
	msg := fmt.Sprintf("Unexpected token <%v> on line %v, position %v. <%v> was expected.",
		unexpected.Type,
		unexpected.Line,
		unexpected.PositionInLine,
		expectedType)
	p.errors = append(p.errors, msg)
}

func (p *Parser) illegalToken() {
	msg := fmt.Sprintf("Illegal token '%v' on line %v, position %v.",
		p.next.Value,
		p.next.Line,
		p.next.PositionInLine)
	p.errors = append(p.errors, msg)
}

func (p *Parser) Errors() []string {
	return p.errors
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

func (p *Parser) skipEol() {
	for p.nextTokenIs(token.EOL) {
		p.nextToken()
	}
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for {
		p.nextToken()
		if p.currentTokenIs(token.EOF) {
			break
		}
		s := p.parseStatement()
		if s != nil {
			program.Statements = append(program.Statements, s)
		}
		p.skipEol()
	}

	return program
}

func (p *Parser) parseStatement() ast.Statement {
	switch {
	case p.currentTokenIs(token.IDENTIFIER) && p.nextTokenIs(token.DECLARE_ASSIGN):
		return p.parseDeclareAssignStatement()
	// case p.current.Type == token.RETURN:
	// 	return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseDeclareAssignStatement() *ast.DeclareAssignStatement {
	identifier := ast.NewIdentifier(p.current)
	operator := p.nextToken()
	p.nextToken()
	expr := p.parseExpression(LOWEST)
	s := ast.NewDeclareAssignStatement(operator, identifier, expr)

	if expr == nil {
		return nil
	}

	return s
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.current}
	stmt.Expression = p.parseExpression(LOWEST)
	return stmt
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	parsePrefix := p.prefixParseFns[p.current.Type]
	if parsePrefix == nil {
		// p.noPrefixParseFnError(p.current.Type)
		return nil
	}
	expr := parsePrefix()

	for !p.nextTokenIs(token.EOL) && precedence < p.nextPrecedence() {
		parseInfix := p.infixParseFns[p.next.Type]
		if parseInfix == nil {
			return expr
		}
		expr = parseInfix(expr)
	}

	return expr
}

func (p *Parser) parseUnaryExpression() ast.Expression {
	operator := p.current
	right := p.parseExpression(PREFIX)
	expression := ast.NewUnaryExpression(operator, right)
	return expression
}

func (p *Parser) parseBinaryExpression(left ast.Expression) ast.Expression {
	operator := p.nextToken()
	precedence := p.currentPrecedence()
	p.nextToken()
	right := p.parseExpression(precedence)
	expr := ast.NewBinaryExpression(left, operator, right)
	return expr
}

func (p *Parser) parseIfExpression() ast.Expression {

	p.nextToken()
	condition := p.parseExpression(LOWEST)

	consequence := p.parseBlockStatement()

	var alternative *ast.BlockStatement
	innerIf := false
	if p.currentTokenIs(token.ELSE) {
		innerIf = p.nextTokenIs(token.IF)
		alternative = p.parseBlockStatement()
	}

	if !innerIf {
		if p.expectCurrentTokenType(token.END) {
			p.skipEol()
		}
	}

	return ast.NewIfExpression(condition, consequence, alternative)
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	statements := []ast.Statement{}
	p.skipEol()
	p.nextToken()

	for !p.currentTokenIs(token.END) && !p.currentTokenIs(token.EOF) && !p.currentTokenIs(token.ELSE) {
		s := p.parseStatement()
		if s != nil {
			statements = append(statements, s)
		}
		p.skipEol()
		p.nextToken()
	}

	return ast.NewBlockStatement(statements)
}

func (p *Parser) parseIdentifier() ast.Expression {
	return ast.NewIdentifier(p.current)
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit, _ := ast.NewIntegerLiteral(p.current)
	return lit
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()
	exp := p.parseExpression(LOWEST)
	if !p.expectNextTokenType(token.RPARENTHESIS) {
		return nil
	}
	p.nextToken()
	return exp
}
