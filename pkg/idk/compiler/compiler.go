package compiler

import (
	"fmt"
	"strconv"

	"github.com/fglo/idk/pkg/idk/chunk"
	"github.com/fglo/idk/pkg/idk/lexer"
	"github.com/fglo/idk/pkg/idk/opcodes"
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
	prefixCompileFn func() ([]byte, opcodes.ValType)
	infixCompileFn  func(left []byte) ([]byte, opcodes.ValType)
)

type Compiler struct {
	input string
	lexer *lexer.Lexer

	previous token.Token
	current  token.Token
	next     token.Token

	prefixCompileFns map[token.TokenType]prefixCompileFn
	infixCompileFns  map[token.TokenType]infixCompileFn

	errors []string

	currentScope *scope

	chunk *chunk.Chunk
}

func NewCompiler(input string) *Compiler {
	compiler := &Compiler{
		input:            input,
		lexer:            lexer.NewLexer(input),
		prefixCompileFns: make(map[token.TokenType]prefixCompileFn),
		infixCompileFns:  make(map[token.TokenType]infixCompileFn),
		currentScope:     NewScope(),
		chunk:            chunk.NewChunk(),
	}

	compiler.registerPrefixes()
	compiler.registerInfixes()

	compiler.consume()

	return compiler
}

/// COMPILING

func (c *Compiler) CompileProgram() *chunk.Chunk {
	for !c.next.Is(token.EOF) {
		c.advance()
		if c.current.Is(token.EOL) {
			continue
		}
		c.compileStatement()
		c.skipEols()
	}

	fmt.Println(c.chunk.Disassemble())
	return c.chunk
}

/// statements

func (c *Compiler) compileStatement() {
	switch {
	case c.current.Is(token.LINE_COMMENT):
		c.skipCommentedLine()
	// case c.current.Is(token.IMPORT):
	// 	return c.compileImportStatement()
	case c.current.Is(token.IDENTIFIER) && c.next.Is(token.DECLASSIGN):
		c.compileDeclareAssignStatement()
	// case c.current.Is(token.IDENTIFIER) && c.current.Is(token.DECLARE):
	// 	return c.compileDeclareStatement()
	// case c.current.Is(token.IDENTIFIER) && c.current.Is(token.ASSIGN):
	// 	return c.compileAssignStatement()
	case c.current.Is(token.IDENTIFIER) && c.next.Is(token.LPARENTHESIS):
		c.compileFunctionCallStatement()
	// case c.current.Is(token.IF):
	// 	return c.compileIfStatement()
	// case c.current.Is(token.FOR):
	// 	return c.compileForStatement()
	// case c.current.Is(token.FUNC):
	// 	return c.compileFunctionDefinitionStatement()
	// case c.current.Is(token.RETURN):
	// 	return c.compileReturnStatement()
	default:
		c.reportUnexpectedFirstToken(c.current)
		// return c.compileExpressionStatement()
	}
}

func (c *Compiler) skipCommentedLine() {
	c.expectCurrent(token.LINE_COMMENT)
	for !c.current.Is(token.EOL) && !c.current.Is(token.EOF) {
		c.consumeTokenWithoutCheckingForIllegals()
	}
}

func (c *Compiler) compileDeclareAssignStatement() {
	identifier := c.current

	c.expectNext(token.DECLASSIGN)
	c.advance()

	c.advance()
	expr, varType := c.compileExpression(LOWEST)

	if expr == nil {
		return
	}

	addr := c.chunk.AddStringConstant(identifier.Value)
	c.currentScope.Insert(identifier.Value, addr, varType)

	c.chunk.WriteBytes(expr)
	c.chunk.Write(opcodes.VarBind(varType))
	c.chunk.Write(byte(addr))
}

// func (c *Compiler) compileDeclareStatement() *ast.DeclareStatement {
// 	identifier := ast.NewIdentifier(c.current)

// 	c.consumeToken() // declare operator
// 	vartype := c.consumeToken()

// 	identifier.SetType(token.LookupType(vartype.Value))

// 	var ass *ast.AssignStatement
// 	if c.next.Is(token.ASSIGN) {
// 		c.consumeToken() // skip type
// 		c.consumeToken() // skip the assign operator

// 		expr := c.compileExpression(LOWEST)
// 		if expr != nil {
// 			ass = ast.NewAssignStatement(identifier, expr)
// 		}
// 	}

// 	return ast.NewDeclareStatement(identifier, ass)
// }

// func (c *Compiler) compileAssignStatement() *ast.AssignStatement {
// 	identifier := ast.NewIdentifier(c.current)

// 	c.consumeToken() // assign operator
// 	c.consumeToken() // skip the assign operator

// 	expr := c.compileExpression(LOWEST)

// 	if expr == nil {
// 		return nil
// 	}

// 	return ast.NewAssignStatement(identifier, expr)
// }

// func (c *Compiler) compileIfStatement() *ast.IfStatement {
// 	innerIf := c.previousTokenWas(token.ELSE)

// 	if c.expectCurrentTokenType(token.IF) {
// 		c.consumeToken() // skip if keyword
// 	}

// 	condition := c.compileExpression(LOWEST)

// 	c.skipEols()

// 	consequence := c.compileBlockStatement()

// 	c.skipEols()

// 	var alternative *ast.BlockStatement
// 	if c.next.Is(token.ELSE) {
// 		c.consumeToken()
// 		alternative = c.compileBlockStatement()
// 	}
// 	if !innerIf && c.expectNextTokenType(token.END) {
// 		c.consumeToken() // skip end keyword
// 	}

// 	return ast.NewIfStatement(condition, consequence, alternative)
// }

// func (c *Compiler) compileForStatement() *ast.ForLoopStatement { // TODO: parsing for loop different from while
// 	if c.expectCurrentTokenType(token.FOR) {
// 		c.consumeToken() // skip for keyword
// 	}

// 	condition := c.compileExpression(LOWEST)

// 	c.skipEols()

// 	consequence := c.compileBlockStatement()

// 	c.skipEols()

// 	if c.expectNextTokenType(token.END) {
// 		c.consumeToken() // skip end keyword
// 	}

// 	return ast.NewForLoopStatement(condition, consequence)
// }

// func (c *Compiler) compileFunctionDefinitionStatement() *ast.FunctionDefinitionStatement {
// 	if c.expectCurrentTokenType(token.FUNC) {
// 		c.consumeToken() // skip func keyword
// 	}

// 	identifier := ast.NewIdentifier(c.current)
// 	identifier.SetType(token.FUNC)

// 	c.consumeToken() // skip identifier

// 	parameters := c.compileFunctionDefinitionParametersList()

// 	vartype := *token.NewTokenNotDefaultValue(token.TYPE, c.current.Position, c.current.Line, c.current.PositionInLine, string(token.VOID))
// 	if c.next.Is(token.RETURN_TYPE) {
// 		c.consumeToken()
// 		c.expectNextTokenType(token.TYPE)
// 		vartype = c.consumeToken()
// 	}

// 	c.expectNextTokenType(token.EOL)

// 	c.skipEols()

// 	body := c.compileBlockStatement()

// 	if c.expectNextTokenType(token.END) {
// 		c.consumeToken() // skip end keyword
// 	}

// 	return ast.NewFunctionDefinitionStatement(*identifier, parameters, vartype, body)
// }

// func (c *Compiler) compileReturnStatement() *ast.ReturnStatement {
// 	c.consumeToken() // return keyword

// 	expr := c.compileExpression(LOWEST)

// 	c.skipEols()

// 	if expr == nil {
// 		return nil
// 	}

// 	return ast.NewReturnStatement(expr)
// }

func (c *Compiler) compileFunctionCallStatement() {
	expr, _ := c.compileFunctionCallExpression()
	c.chunk.WriteBytes(expr)
}

// func (c *Compiler) compileImportStatement() *ast.ImportStatement {
// 	c.expectNextTokenType(token.IDENTIFIER)
// 	c.consumeToken()
// 	stmt := ast.NewImportStatement(c.current)
// 	c.consumeToken()
// 	return stmt
// }

// func (c *Compiler) compileBlockStatement() *ast.BlockStatement {
// 	statements := []ast.Statement{}

// 	c.skipEols()

// 	for !c.next.Is(token.END) && !c.next.Is(token.EOF) && !c.next.Is(token.ELSE) {
// 		c.consumeToken()
// 		s := c.compileStatement()
// 		if s != nil {
// 			statements = append(statements, s)
// 		}
// 		c.skipEols()
// 	}

// 	return ast.NewBlockStatement(statements)
// }

/// expressions

func (c *Compiler) compileExpression(precedence int) ([]byte, opcodes.ValType) {
	compilePrefix := c.prefixCompileFns[c.current.Type]
	if compilePrefix == nil {
		return nil, 0
	}
	expr, varType := compilePrefix()

	previousVarType := varType

	for !c.next.Is(token.EOL) && !c.next.Is(token.COMMA) && precedence < c.nextPrecedence() {
		operator := c.next

		compileInfix := c.infixCompileFns[c.peek().Type]
		if compileInfix == nil {
			return expr, varType
		}

		c.advance()
		expr, varType = compileInfix(expr)

		if previousVarType != varType {
			c.reportTypeMismatch(operator, previousVarType, varType)
		}

		c.expectOperatorOrEndOfExpression()
	}

	return expr, varType
}

func (c *Compiler) compilePrefixExpression() ([]byte, opcodes.ValType) {
	bytecode := make([]byte, 0)

	operator := c.current.Type
	c.consume() // skip the operator
	right, varType := c.compileExpression(PREFIX)

	bytecode = append(bytecode, right...)
	bytecode = append(bytecode, opcodes.GetPrefixOperator(operator, varType))

	return bytecode, varType
}

func (c *Compiler) compileInfixExpression(left []byte) ([]byte, opcodes.ValType) {
	bytecode := make([]byte, 0)

	bytecode = append(bytecode, left...)

	operator := c.current.Type
	precedence := c.currentPrecedence()
	c.advance() // skip the operator

	for c.current.Is(token.EOL) {
		c.advance()
	}

	right, varType := c.compileExpression(precedence)

	bytecode = append(bytecode, right...)
	bytecode = append(bytecode, opcodes.GetInfixOperator(operator, varType))

	return bytecode, varType
}

func (c *Compiler) compileIdentifier() ([]byte, opcodes.ValType) {
	if c.current.Is(token.IDENTIFIER) && c.next.Is(token.LPARENTHESIS) {
		return c.compileFunctionCallExpression()
	}
	return c.compileIdentifierLiteral()
}

func (c *Compiler) compileFunctionCallExpression() ([]byte, opcodes.ValType) {
	bytecode := make([]byte, 0)

	identifier := c.current.Value

	c.expectNext(token.LPARENTHESIS)
	c.advance()

	// exp.Parameters = c.compileFunctionCallParametersList()

	switch identifier {
	case "print":
		c.advance()
		expr, varType := c.compileExpression(LOWEST)
		bytecode = append(bytecode, expr...)
		// bytecode = append(bytecode, compileExpression(node.Parameters[0])...)
		bytecode = append(bytecode, opcodes.IPRINT)

		c.expectNext(token.RPARENTHESIS)
		c.advance()

		return bytecode, varType
	}

	return bytecode, 0
}

// func (c *Compiler) compileType() ast.Expression {
// 	typ := ast.NewType(c.current)
// 	return typ
// }

// func (c *Compiler) compileGroupedExpression() ast.Expression {
// 	c.consumeToken() //skip opening parenthesis
// 	exp := c.compileExpression(LOWEST)
// 	if !c.expectNextTokenType(token.RPARENTHESIS) {
// 		return nil
// 	}
// 	c.consumeToken() //skip closing parenthesis
// 	return exp
// }

/// literals

func (c *Compiler) compileIdentifierLiteral() ([]byte, opcodes.ValType) {
	bytecode := make([]byte, 0)

	varName := c.current.Value
	symbol, ok := c.currentScope.Lookup(varName)

	if !ok {
		// handle error
	}

	bytecode = append(bytecode, opcodes.VarLookup(symbol.valType))
	bytecode = append(bytecode, byte(symbol.cpAddr))

	return bytecode, symbol.valType
}

func (c *Compiler) compileIntegerLiteral() ([]byte, opcodes.ValType) {
	bytecode := make([]byte, 0)

	val, _ := strconv.Atoi(c.current.Value)
	addr := c.chunk.AddIntConstant(val)

	bytecode = append(bytecode, opcodes.IPUSH)
	bytecode = append(bytecode, byte(addr))

	return bytecode, opcodes.INT
}

func (c *Compiler) compileFloatingPointLiteral() ([]byte, opcodes.ValType) {
	bytecode := make([]byte, 0)

	val, _ := strconv.ParseFloat(c.current.Value, 64)
	addr := c.chunk.AddFloatConstant(val)

	bytecode = append(bytecode, opcodes.FPUSH)
	bytecode = append(bytecode, byte(addr))

	return bytecode, opcodes.FLOAT
}

func (c *Compiler) compileBooleanLiteral() ([]byte, opcodes.ValType) {
	bytecode := make([]byte, 0)

	val, _ := strconv.ParseBool(c.current.Value)
	addr := c.chunk.AddBoolConstant(val)

	bytecode = append(bytecode, opcodes.BPUSH)
	bytecode = append(bytecode, byte(addr))

	return bytecode, opcodes.BOOL
}

func (c *Compiler) compileCharacterLiteral() ([]byte, opcodes.ValType) {
	bytecode := make([]byte, 0)

	val := []rune(c.current.Value)[0]
	addr := c.chunk.AddCharConstant(val)

	bytecode = append(bytecode, opcodes.CPUSH)
	bytecode = append(bytecode, byte(addr))

	return bytecode, opcodes.CHAR
}

func (c *Compiler) compileStringLiteral() ([]byte, opcodes.ValType) {
	bytecode := make([]byte, 0)

	addr := c.chunk.AddStringConstant(c.current.Value)

	bytecode = append(bytecode, opcodes.SPUSH)
	bytecode = append(bytecode, byte(addr))

	return bytecode, opcodes.STRING
}

/// HELPERS

func (c *Compiler) registerPrefixes() {
	c.registerPrefix(token.IDENTIFIER, c.compileIdentifier)
	c.registerPrefix(token.INT, c.compileIntegerLiteral)
	c.registerPrefix(token.FLOAT, c.compileFloatingPointLiteral)
	c.registerPrefix(token.BOOL, c.compileBooleanLiteral)
	c.registerPrefix(token.CHAR, c.compileCharacterLiteral)
	c.registerPrefix(token.STRING, c.compileStringLiteral)
	// c.registerPrefix(token.LPARENTHESIS, c.compileGroupedExpression)
	c.registerPrefix(token.MINUS, c.compilePrefixExpression)
	c.registerPrefix(token.BANG, c.compilePrefixExpression)
	// c.registerPrefix(token.TYPE, c.compileType)
	// c.registerPrefix(token.FUNC, c.compileType)
}

func (c *Compiler) registerInfixes() {
	c.registerInfix(token.PLUS, c.compileInfixExpression)
	c.registerInfix(token.MINUS, c.compileInfixExpression)
	c.registerInfix(token.ASTERISK, c.compileInfixExpression)
	c.registerInfix(token.SLASH, c.compileInfixExpression)
	c.registerInfix(token.MODULO, c.compileInfixExpression)
	c.registerInfix(token.IN, c.compileInfixExpression)
	c.registerInfix(token.RANGE, c.compileInfixExpression)
	c.registerInfix(token.EQ, c.compileInfixExpression)
	c.registerInfix(token.NEQ, c.compileInfixExpression)
	c.registerInfix(token.GT, c.compileInfixExpression)
	c.registerInfix(token.GTE, c.compileInfixExpression)
	c.registerInfix(token.LT, c.compileInfixExpression)
	c.registerInfix(token.LTE, c.compileInfixExpression)
	c.registerInfix(token.AND, c.compileInfixExpression)
	c.registerInfix(token.OR, c.compileInfixExpression)
	c.registerInfix(token.XOR, c.compileInfixExpression)
	// c.registerInfix(token.DOT, c.compileProperty)
}

func (c *Compiler) registerPrefix(tokenType token.TokenType, fn prefixCompileFn) {
	c.prefixCompileFns[tokenType] = fn
}

func (c *Compiler) registerInfix(tokenType token.TokenType, fn infixCompileFn) {
	c.infixCompileFns[tokenType] = fn
}

func (c *Compiler) peek() token.Token {
	return c.next
}

func (c *Compiler) advance() {
	if c.current.Type == token.EOF {
		return
	}

	c.previous = c.current
	c.current = c.next
	c.next = c.lexer.ReadToken()

	if c.current.Type == token.ILLEGAL {
		c.reportIllegalToken()
	}
}

func (c *Compiler) consume() token.Token {
	c.advance()
	return c.current
}

func (c *Compiler) consumeTokenWithoutCheckingForIllegals() token.Token {
	c.previous = c.current
	c.current = c.next
	if c.current.Type == token.EOF {
		return c.current
	}
	c.next = c.lexer.ReadToken()
	return c.current
}

func (c *Compiler) expectCurrent(t token.TokenType) bool {
	if c.current.Is(t) {
		return true
	} else {
		c.reportUnexpectedToken(c.current, t)
		return false
	}
}

func (c *Compiler) expectNext(t token.TokenType) bool {
	if c.next.Is(t) {
		return true
	} else {
		c.reportUnexpectedToken(c.next, t)
		return false
	}
}

func (c *Compiler) expectOperatorOrEndOfExpression() bool {
	if c.next.Type.IsOperator() || c.next.Is(token.EOL) || c.next.Is(token.EOF) || c.next.Is(token.COMMA) || c.next.Is(token.RPARENTHESIS) {
		return true
	} else {
		c.reportExpectedOperatorOrEndOfExpression(c.next)
		return false
	}
}

func (c *Compiler) currentPrecedence() int {
	if p, ok := precedences[c.current.Type]; ok {
		return p
	}
	return LOWEST
}

func (c *Compiler) nextPrecedence() int {
	if p, ok := precedences[c.next.Type]; ok {
		return p
	}
	return LOWEST
}

func (c *Compiler) skipEols() {
	for c.next.Is(token.EOL) {
		c.advance()
	}
}

func (c *Compiler) Errors() []string {
	return c.errors
}

func (c *Compiler) reportUnexpectedToken(unexpected token.Token, expectedType token.TokenType) {
	msg := fmt.Sprintf("ERROR: Unexpected token <%v> on line %v, position %v. <%v> was expected.",
		unexpected.Type,
		unexpected.Line,
		unexpected.PositionInLine,
		expectedType)
	c.errors = append(c.errors, msg)
}

func (c *Compiler) reportUnexpectedFirstToken(unexpected token.Token) {
	msg := fmt.Sprintf("ERROR: Unexpected token <%v> on line %v, position %v. Expected declaration or a statement.",
		unexpected.Type,
		unexpected.Line,
		unexpected.PositionInLine)
	c.errors = append(c.errors, msg)
}

func (c *Compiler) reportExpectedOperatorOrEndOfExpression(unexpected token.Token) {
	msg := fmt.Sprintf("ERROR: Unexpected token <%v> on line %v, position %v. Expected operator, <EOL>, <EOF>, ',' or ')'.",
		unexpected.Type,
		unexpected.Line,
		unexpected.PositionInLine)
	c.errors = append(c.errors, msg)
}

func (c *Compiler) reportIllegalToken() {
	msg := fmt.Sprintf("ERROR: Illegal token '%v' on line %v, position %v.",
		c.next.Value,
		c.next.Line,
		c.next.PositionInLine)
	c.errors = append(c.errors, msg)
}

func (c *Compiler) reportTypeMismatch(operator token.Token, type1, type2 opcodes.ValType) {
	msg := fmt.Sprintf("ERROR: Type mismatch: '%v' %s '%v' on line %v, position %v.",
		type1,
		operator.Value,
		type2,
		operator.Line,
		operator.PositionInLine)
	c.errors = append(c.errors, msg)
}
