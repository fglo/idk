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
	DECLARE_ASSIGN
	DECLARE
	ASSIGN
	IN
	OR
	AND
	XOR
	NOT
	EQUALS
	LESSGREATER
	SUM
	PRODUCT
	PREFIX
	RANGE
	CALL
	INDEX
	PROPERTY
)

var precedences = map[token.TokenType]int{
	token.DECLASSIGN:      DECLARE_ASSIGN,
	token.DECLARE:         DECLARE,
	token.ASSIGN:          ASSIGN,
	token.IN:              IN,
	token.AND:             AND,
	token.OR:              OR,
	token.XOR:             XOR,
	token.BANG:            NOT,
	token.EQ:              EQUALS,
	token.NEQ:             EQUALS,
	token.LT:              LESSGREATER,
	token.GT:              LESSGREATER,
	token.LTE:             LESSGREATER,
	token.GTE:             LESSGREATER,
	token.PLUS:            SUM,
	token.MINUS:           SUM,
	token.MODULO:          PRODUCT,
	token.SLASH:           PRODUCT,
	token.ASTERISK:        PRODUCT,
	token.RANGE:           RANGE,
	token.RANGE_INCLUSIVE: RANGE,
	token.LPARENTHESIS:    CALL,
	token.DOT:             PROPERTY,
}

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

type Parser struct {
	input string
	lexer *lexer.Lexer

	previous token.Token
	current  token.Token
	next     token.Token

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn

	errors []string
}

func NewParser(input string) *Parser {
	parser := &Parser{
		input:          input,
		lexer:          lexer.NewLexer(input),
		prefixParseFns: make(map[token.TokenType]prefixParseFn),
		infixParseFns:  make(map[token.TokenType]infixParseFn),
	}

	parser.registerPrefixes()
	parser.registerInfixes()

	parser.consumeToken()

	return parser
}

func (p *Parser) registerPrefixes() {
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	p.registerPrefix(token.IDENTIFIER, p.parseIdentifier)
	p.registerPrefix(token.TYPE, p.parseType)
	p.registerPrefix(token.FUNC, p.parseType)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.FLOAT, p.parseFloatingPointLiteral)
	p.registerPrefix(token.BOOL, p.parseBooleanLiteral)
	p.registerPrefix(token.CHAR, p.parseCharacterLiteral)
	p.registerPrefix(token.STRING, p.parseStringLiteral)
	p.registerPrefix(token.LPARENTHESIS, p.parseGroupedExpression)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
}

func (p *Parser) registerInfixes() {
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.ASTERISK, p.parseInfixExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)
	p.registerInfix(token.MODULO, p.parseInfixExpression)
	p.registerInfix(token.IN, p.parseInfixExpression)
	p.registerInfix(token.RANGE, p.parseInfixExpression)
	p.registerInfix(token.EQ, p.parseInfixExpression)
	p.registerInfix(token.NEQ, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)
	p.registerInfix(token.GTE, p.parseInfixExpression)
	p.registerInfix(token.LT, p.parseInfixExpression)
	p.registerInfix(token.LTE, p.parseInfixExpression)
	p.registerInfix(token.AND, p.parseInfixExpression)
	p.registerInfix(token.OR, p.parseInfixExpression)
	p.registerInfix(token.XOR, p.parseInfixExpression)
	p.registerInfix(token.DOT, p.parseProperty)
}

func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

func (p *Parser) peekNext() token.Token {
	return p.next
}

func (p *Parser) consumeToken() token.Token { // TODO: take token type as argument
	p.previous = p.current
	p.current = p.next
	if p.current.Type == token.ILLEGAL {
		p.reportIllegalToken()
	}
	if p.current.Type == token.EOF {
		return p.current
	}
	p.next = p.lexer.ReadToken()
	return p.current
}

func (p *Parser) consumeTokenWithoutCheckingForIllegals() token.Token {
	p.previous = p.current
	p.current = p.next
	if p.current.Type == token.EOF {
		return p.current
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

func (p *Parser) expectOperatorOrEndOfExpression() bool {
	if p.next.Type.IsOperator() || p.nextTokenIs(token.EOL) || p.nextTokenIs(token.EOF) || p.nextTokenIs(token.COMMA) || p.nextTokenIs(token.RPARENTHESIS) {
		return true
	} else {
		p.reportExpectedOperatorOrEndOfExpression(p.next)
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
	msg := fmt.Sprintf("ERROR: Unexpected token <%v> on line %v, position %v. <%v> was expected.",
		unexpected.Type,
		unexpected.Line,
		unexpected.PositionInLine,
		expectedType)
	p.errors = append(p.errors, msg)
}

func (p *Parser) reportUnexpectedFirstToken(unexpected token.Token) {
	msg := fmt.Sprintf("ERROR: Unexpected token <%v> on line %v, position %v. Expected declaration or a statement.",
		unexpected.Type,
		unexpected.Line,
		unexpected.PositionInLine)
	p.errors = append(p.errors, msg)
}

func (p *Parser) reportExpectedOperatorOrEndOfExpression(unexpected token.Token) {
	msg := fmt.Sprintf("ERROR: Unexpected token <%v> on line %v, position %v. Expected operator, <EOL>, <EOF>, ',' or ')'.",
		unexpected.Type,
		unexpected.Line,
		unexpected.PositionInLine)
	p.errors = append(p.errors, msg)
}

func (p *Parser) reportIllegalToken() {
	msg := fmt.Sprintf("ERROR: Illegal token '%v' on line %v, position %v.",
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

func (p *Parser) skipEols() {
	for p.currentTokenIs(token.EOL) {
		p.consumeToken()
	}
}

/// PARSING

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for !p.nextTokenIs(token.EOF) {
		if p.consumeToken().Not(token.EOL) {
			if s := p.parseStatement(); s != nil {
				program.Statements = append(program.Statements, s)
			}
		}
	}

	return program
}

func (p *Parser) parseStatement() ast.Statement {
	switch {
	case p.currentTokenIs(token.LINE_COMMENT):
		p.skipCommentedLine()
	case p.currentTokenIs(token.IMPORT):
		return p.parseImportStatement()
	case p.currentTokenIs(token.IDENTIFIER) && p.nextTokenIs(token.DECLASSIGN):
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
	return nil
}

func (p *Parser) skipCommentedLine() {
	p.expectCurrentTokenType(token.LINE_COMMENT)
	for !p.currentTokenIs(token.EOL) && !p.currentTokenIs(token.EOF) {
		p.consumeTokenWithoutCheckingForIllegals()
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

	identifier.SetType(token.LookupType(vartype.Value))

	var ass *ast.AssignStatement
	if p.nextTokenIs(token.ASSIGN) {
		p.consumeToken() // skip type
		p.consumeToken() // skip the assign operator

		expr := p.parseExpression(LOWEST)
		if expr != nil {
			ass = ast.NewAssignStatement(identifier, expr)
		}
	}

	return ast.NewDeclareStatement(identifier, ass)
}

func (p *Parser) parseAssignStatement() *ast.AssignStatement {
	identifier := ast.NewIdentifier(p.current)

	p.consumeToken() // assign operator
	p.consumeToken() // skip the assign operator

	expr := p.parseExpression(LOWEST)

	if expr == nil {
		return nil
	}

	return ast.NewAssignStatement(identifier, expr)
}

func (p *Parser) parseIfStatement() *ast.IfStatement {
	innerIf := p.previousTokenWas(token.ELSE)

	if p.expectCurrentTokenType(token.IF) {
		p.consumeToken() // skip if keyword
	}

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
		p.consumeToken() // skip end keyword
	}

	return ast.NewIfStatement(condition, consequence, alternative)
}

func (p *Parser) parseForStatement() *ast.ForLoopStatement { // TODO: parsing for loop different from while
	if p.expectCurrentTokenType(token.FOR) {
		p.consumeToken() // skip for keyword
	}

	condition := p.parseExpression(LOWEST)

	p.ifEolIsNextThenSkip()

	consequence := p.parseBlockStatement()

	p.ifEolIsNextThenSkip()

	if p.expectNextTokenType(token.END) {
		p.consumeToken() // skip end keyword
	}

	return ast.NewForLoopStatement(condition, consequence)
}

func (p *Parser) parseFunctionDefinitionStatement() *ast.FunctionDefinitionStatement {
	if p.expectCurrentTokenType(token.FUNC) {
		p.consumeToken() // skip func keyword
	}

	identifier := ast.NewIdentifier(p.current)
	identifier.SetType(token.FUNC)

	p.consumeToken() // skip identifier

	parameters := p.parseFunctionDefinitionParametersList()

	vartype := *token.NewTokenNotDefaultValue(token.TYPE, p.current.Position, p.current.Line, p.current.PositionInLine, string(token.VOID))
	if p.nextTokenIs(token.RETURN_TYPE) {
		p.consumeToken()
		p.expectNextTokenType(token.TYPE)
		vartype = p.consumeToken()
	}

	p.expectNextTokenType(token.EOL)

	p.ifEolIsNextThenSkip()

	body := p.parseBlockStatement()

	if p.expectNextTokenType(token.END) {
		p.consumeToken() // skip end keyword
	}

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
	stmt := &ast.ExpressionStatement{
		Expression: p.parseFunctionCallExpression(),
	}
	p.ifEolIsNextThenSkip()
	return stmt
}

func (p *Parser) parseFunctionCallExpression() *ast.FunctionCallExpression {
	exp := ast.NewFunctionCallExpression(p.current)
	p.consumeToken()
	exp.Parameters = p.parseFunctionCallParametersList()
	return exp
}

func (p *Parser) parseImportStatement() *ast.ImportStatement {
	p.expectNextTokenType(token.IDENTIFIER)
	p.consumeToken()
	stmt := ast.NewImportStatement(p.current)
	p.consumeToken()
	return stmt
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

// TODO: expression statements (if and how?)
// func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
// 	stmt := new(ast.ExpressionStatement)
// 	stmt.Expression = p.parseExpression(LOWEST)
// 	return stmt
// }

func (p *Parser) parseExpression(precedence int) ast.Expression {
	parsePrefix := p.prefixParseFns[p.current.Type]
	if parsePrefix == nil {
		return nil
	}
	expr := parsePrefix()

	for !p.nextTokenIs(token.EOL) && !p.nextTokenIs(token.COMMA) && precedence < p.nextPrecedence() {
		parseInfix := p.infixParseFns[p.peekNext().Type]
		if parseInfix == nil {
			return expr
		}
		p.consumeToken()
		expr = parseInfix(expr)

		p.expectOperatorOrEndOfExpression()
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

	for p.nextTokenIs(token.COMMA) {
		p.consumeToken()
		p.consumeToken()
		list = append(list, p.parseDeclareStatement())
	}

	if p.expectNextTokenType(token.RPARENTHESIS) {
		p.consumeToken()
		return list
	}

	return nil
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

func (p *Parser) parsePrefixExpression() ast.Expression {
	operator := p.current
	p.consumeToken() // skip the operator
	right := p.parseExpression(PREFIX)
	expression := ast.NewPrefixExpression(operator, right)
	return expression
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	operator := p.current
	precedence := p.currentPrecedence()
	p.consumeToken() // skip the operator
	p.skipEols()     // skip EOLs
	right := p.parseExpression(precedence)
	expr := ast.NewInfixExpression(left, operator, right)
	return expr
}

func (p *Parser) parseProperty(parent ast.Expression) ast.Expression {
	p.expectCurrentTokenType(token.DOT)
	precedence := p.currentPrecedence()
	p.consumeToken() // skip the operator
	p.expectCurrentTokenType(token.IDENTIFIER)
	property := p.parseExpression(precedence)
	expr := ast.NewPropertyExpression(parent.(*ast.Identifier), property)
	return expr
}

func (p *Parser) parseIdentifier() ast.Expression {
	if p.currentTokenIs(token.IDENTIFIER) && p.nextTokenIs(token.LPARENTHESIS) {
		return p.parseFunctionCallExpression()
	}
	return ast.NewIdentifier(p.current)
}

func (p *Parser) parseType() ast.Expression {
	typ := ast.NewType(p.current)
	return typ
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit, _ := ast.NewIntegerLiteral(p.current)
	return lit
}

func (p *Parser) parseFloatingPointLiteral() ast.Expression {
	lit, _ := ast.NewFloatingPointLiteral(p.current)
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
