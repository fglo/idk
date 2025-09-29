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
	typeCompileFn   func() ([]byte, opcodes.ValType)
	returnCompileFn func() ([]byte, opcodes.ValType)
	prefixCompileFn func() ([]byte, opcodes.ValType)
	infixCompileFn  func(left []byte, leftValType opcodes.ValType) ([]byte, opcodes.ValType)
)

type Compiler struct {
	input string
	lexer *lexer.Lexer

	previous token.Token
	current  token.Token
	next     token.Token

	typeCompileFns   map[opcodes.ValType]typeCompileFn
	returnCompileFns map[opcodes.ValType]returnCompileFn
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
		typeCompileFns:   make(map[opcodes.ValType]typeCompileFn),
		returnCompileFns: make(map[opcodes.ValType]returnCompileFn),
		prefixCompileFns: make(map[token.TokenType]prefixCompileFn),
		infixCompileFns:  make(map[token.TokenType]infixCompileFn),
		currentScope:     NewScope(),
		chunk:            chunk.NewChunk(),
	}

	compiler.registerTypes()
	compiler.registerPrefixes()
	compiler.registerInfixes()
	compiler.registerReturns()

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
		c.chunk.WriteBytes(c.compileStatement())
		c.skipEols()
	}

	fmt.Println(c.chunk.Disassemble())
	return c.chunk
}

/// statements

func (c *Compiler) compileStatement() []byte {
	bytecode := make([]byte, 0)

	switch {
	case c.current.Is(token.LINE_COMMENT):
		c.skipCommentedLine()
	// case c.current.Is(token.IMPORT):
	// 	return c.compileImportStatement()
	case c.current.Is(token.IDENTIFIER) && c.next.Is(token.DECLASSIGN):
		bytecode = append(bytecode, c.compileDeclareAssignStatement(c.current)...)
	case c.current.Is(token.IDENTIFIER) && c.next.Is(token.DECLARE):
		bytecode = append(bytecode, c.compileDeclareStatement(c.current)...)
	case c.current.Is(token.IDENTIFIER) && c.next.Is(token.ASSIGN):
		bytecode = append(bytecode, c.compileAssignStatement(c.current)...)
	case c.current.Is(token.IDENTIFIER) && c.next.Is(token.LPARENTHESIS):
		bytecode = append(bytecode, c.compileFunctionCallStatement()...)
	case c.current.Is(token.IF):
		bytecode = append(bytecode, c.compileIfStatement()...)
	// case c.current.Is(token.FOR):
	// 	c.compileForStatement()
	case c.current.Is(token.FUNC):
		bytecode = append(bytecode, c.compileFunctionDefinitionStatement()...)
	case c.current.Is(token.RETURN):
		bytecode = append(bytecode, c.compileReturnStatement()...)
	case c.current.Is(token.END):
		c.advance()
	default:
		c.reportUnexpectedFirstToken(c.current)
		// c.compileExpressionStatement()
	}

	return bytecode
}

func (c *Compiler) skipCommentedLine() {
	c.expectCurrent(token.LINE_COMMENT)
	for !c.current.Is(token.EOL) && !c.current.Is(token.EOF) {
		c.consumeTokenWithoutCheckingForIllegals()
	}
}

func (c *Compiler) compileDeclareAssignStatement(identifier token.Token) []byte {
	bytecode := make([]byte, 0)

	c.expectNext(token.DECLASSIGN)
	c.advance()

	c.advance()

	expr, valType := c.compileExpression(LOWEST)
	if expr == nil {
		return bytecode
	}

	addr := c.chunk.AddStringConstant(identifier.Value)
	c.currentScope.Insert(identifier.Value, addr, valType)

	bytecode = append(bytecode, expr...)
	bytecode = append(bytecode, opcodes.VarBind(valType))
	bytecode = append(bytecode, byte(addr))

	return bytecode
}

func (c *Compiler) compileDeclareStatement(identifier token.Token) []byte {
	bytecode := make([]byte, 0)

	c.expectNext(token.DECLARE)
	c.advance()

	c.expectNext(token.TYPE)
	c.advance()

	bytes, valType := c.compileType()
	if bytes == nil {
		return bytecode
	}

	addr := c.chunk.AddStringConstant(identifier.Value)
	c.currentScope.Insert(identifier.Value, addr, valType)

	bytecode = append(bytecode, bytes...)
	bytecode = append(bytecode, opcodes.VarBind(valType))
	bytecode = append(bytecode, byte(addr))

	if c.next.Type == token.ASSIGN {
		bytecode = append(bytecode, c.compileAssignStatement(identifier)...)
	}

	return bytecode
}

func (c *Compiler) compileAssignStatement(identifier token.Token) []byte {
	bytecode := make([]byte, 0)

	c.expectNext(token.ASSIGN)
	c.advance()

	c.advance()
	expr, valType := c.compileExpression(LOWEST)
	if expr == nil {
		return bytecode
	}

	addr := c.chunk.AddStringConstant(identifier.Value)
	c.currentScope.Insert(identifier.Value, addr, valType)

	bytecode = append(bytecode, expr...)
	bytecode = append(bytecode, opcodes.VarBind(valType))
	bytecode = append(bytecode, byte(addr))

	return bytecode
}

func (c *Compiler) compileIfStatement() []byte {
	bytecode := make([]byte, 0)

	c.expectCurrent(token.IF)

	token := c.current

	c.advance()

	expr, valType := c.compileExpression(LOWEST)
	if expr == nil {
		return bytecode
	}

	if valType != opcodes.BOOL {
		c.reportWrongType(token, opcodes.BOOL, valType)
	}

	blockStatementByteCode := c.compileBlockStatement()

	bytecode = append(bytecode, expr...)
	bytecode = append(bytecode, opcodes.IF)
	bytecode = append(bytecode, byte(len(blockStatementByteCode)))
	bytecode = append(bytecode, blockStatementByteCode...)

	return bytecode
}

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

func (c *Compiler) compileFunctionDefinitionStatement() []byte {
	bytecode := make([]byte, 0)

	c.expectCurrent(token.FUNC)

	c.expectNext(token.IDENTIFIER)
	c.advance()

	identifier := c.current

	c.expectNext(token.LPARENTHESIS)
	c.advance()

	noParams := 0

	paramsBytecode := make([]byte, 0)

	if c.next.Is(token.IDENTIFIER) {
		for !c.current.Is(token.RPARENTHESIS) {
			c.expectNext(token.IDENTIFIER)
			c.advance()

			paramIdentifier := c.current

			c.expectNext(token.DECLARE)
			c.advance()

			c.expectNext(token.TYPE)
			c.advance()

			bytes, valType := c.compileType()
			if bytes == nil {
				return bytecode // TODO: better error handling
			}

			addr := c.chunk.AddStringConstant(paramIdentifier.Value)
			c.currentScope.Insert(paramIdentifier.Value, addr, valType)

			paramsBytecode = append(paramsBytecode, bytes...)
			paramsBytecode = append(paramsBytecode, opcodes.VarBind(valType))
			paramsBytecode = append(paramsBytecode, byte(addr))

			c.expectNext(token.COMMA, token.RPARENTHESIS)
			c.advance()

			noParams += 1
		}
	} else {
		c.expectNext(token.RPARENTHESIS)
		c.advance()
	}

	c.expectNext(token.RETURN_TYPE)
	c.advance()

	c.expectNext(token.TYPE)
	c.advance()

	valType := opcodes.ValTypeFromString(c.current.Value)
	if valType == opcodes.UNKNOWN_TYPE {
		return bytecode
	}

	addr := c.chunk.AddStringConstant(identifier.Value)
	c.currentScope.Insert(identifier.Value, addr, valType)

	bytecode = append(bytecode, opcodes.FuncCreate(valType))
	bytecode = append(bytecode, byte(addr))
	bytecode = append(bytecode, byte(noParams)) // TODO: actual number of args

	bytecode = append(bytecode, paramsBytecode...)

	bytecode = append(bytecode, c.compileBlockStatement()...)

	return bytecode
}

func (c *Compiler) compileReturnStatement() []byte {
	bytecode := make([]byte, 0)

	c.expectCurrent(token.RETURN)
	c.advance()

	expr, valType := c.compileExpression(LOWEST)
	if expr == nil {
		return bytecode
	}

	bytecode = append(bytecode, expr...)
	bytecode = append(bytecode, opcodes.FuncReturn(valType))

	return bytecode
}

func (c *Compiler) compileFunctionCallStatement() []byte {
	bytecode := make([]byte, 0)
	expr, _ := c.compileFunctionCallExpression()
	bytecode = append(bytecode, expr...)
	return bytecode
}

// func (c *Compiler) compileImportStatement() *ast.ImportStatement {
// 	c.expectNextTokenType(token.IDENTIFIER)
// 	c.consumeToken()
// 	stmt := ast.NewImportStatement(c.current)
// 	c.consumeToken()
// 	return stmt
// }

func (c *Compiler) compileBlockStatement() []byte {
	bytecode := make([]byte, 0)

	for !c.next.Is(token.END) {
		c.advance()
		if c.current.Is(token.EOL) {
			continue
		}
		bytecode = append(bytecode, c.compileStatement()...)
		c.skipEols()
	}

	c.expectNext(token.END)
	c.advance()

	return bytecode
}

/// expressions

func (c *Compiler) compileExpression(precedence int) ([]byte, opcodes.ValType) {
	compilePrefix := c.prefixCompileFns[c.current.Type]
	if compilePrefix == nil {
		return nil, 0
	}
	expr, valType := compilePrefix()

	for !c.next.Is(token.EOL) && !c.next.Is(token.COMMA) && precedence < c.nextPrecedence() {
		compileInfix := c.infixCompileFns[c.peek().Type]
		if compileInfix == nil {
			return expr, valType
		}

		c.advance()
		expr, valType = compileInfix(expr, valType)

		c.expectOperatorOrEndOfExpression()
	}

	return expr, valType
}

func (c *Compiler) compilePrefixExpression() ([]byte, opcodes.ValType) {
	bytecode := make([]byte, 0)

	operator := c.current.Type
	c.consume() // skip the operator
	right, valType := c.compileExpression(PREFIX)

	bytecode = append(bytecode, right...)
	bytecode = append(bytecode, opcodes.PrefixOperator(operator, valType))

	return bytecode, valType
}

func (c *Compiler) compileGroupedExpression() ([]byte, opcodes.ValType) {
	bytecode := make([]byte, 0)

	c.expectCurrent(token.LPARENTHESIS)
	c.advance()

	expr, valType := c.compileExpression(LOWEST)
	if expr == nil {
		return nil, opcodes.UNKNOWN_TYPE
	}

	bytecode = append(bytecode, expr...)

	c.expectNext(token.RPARENTHESIS)
	c.advance()

	return expr, valType
}

func (c *Compiler) compileInfixExpression(left []byte, leftValType opcodes.ValType) ([]byte, opcodes.ValType) {
	bytecode := make([]byte, 0)

	bytecode = append(bytecode, left...)

	operator := c.current
	precedence := c.currentPrecedence()
	c.advance() // skip the operator

	for c.current.Is(token.EOL) {
		c.advance()
	}

	right, rightValType := c.compileExpression(precedence)

	if leftValType != rightValType {
		c.reportTypeMismatch(operator, leftValType, rightValType)
	}

	bytecode = append(bytecode, right...)
	bytecode = append(bytecode, opcodes.InfixOperator(operator.Type, rightValType))

	resultValType := opcodes.OperatorResultType(operator.Type, rightValType)

	return bytecode, resultValType
}

func (c *Compiler) compileIdentifier() ([]byte, opcodes.ValType) {
	if c.current.Is(token.IDENTIFIER) && c.next.Is(token.LPARENTHESIS) {
		return c.compileFunctionCallExpression()
	}
	return c.compileIdentifierLiteral()
}

func (c *Compiler) compileFunctionCallExpression() ([]byte, opcodes.ValType) {
	bytecode := make([]byte, 0)
	var expr []byte
	valType := opcodes.UNKNOWN_TYPE

	identifier := c.current.Value

	c.expectNext(token.LPARENTHESIS)
	c.advance()

	// TODO: exp.Parameters = c.compileFunctionCallParametersList()

	switch identifier {
	case "print":
		c.advance()
		expr, valType = c.compileExpression(LOWEST)
		bytecode = append(bytecode, expr...)
		bytecode = append(bytecode, opcodes.ValPrint(valType))
		c.expectNext(token.RPARENTHESIS)
		c.advance()
	default:
		symbol, ok := c.currentScope.Lookup(identifier)
		if !ok {
			return nil, opcodes.UNKNOWN_TYPE
		}

		noParams := 0
		if !c.next.Is(token.RPARENTHESIS) {
			for !c.current.Is(token.RPARENTHESIS) {
				expr, valType = c.compileExpression(LOWEST)
				if expr == nil {
					return bytecode, valType // TODO: better error handling
				}

				bytecode = append(bytecode, expr...)

				c.expectCurrent(token.COMMA, token.RPARENTHESIS)

				noParams += 1
			}
		} else {
			c.expectNext(token.RPARENTHESIS)
			c.advance()
		}

		bytecode = append(bytecode, opcodes.FuncCall(symbol.valType))
		bytecode = append(bytecode, byte(symbol.cpAddr))
		bytecode = append(bytecode, byte(noParams))

		valType = symbol.valType
	}

	return bytecode, valType
}

/// types

func (c *Compiler) compileType() ([]byte, opcodes.ValType) {
	valType := opcodes.ValTypeFromString(c.current.Value)

	compileType := c.typeCompileFns[valType]
	if compileType == nil {
		return nil, 0
	}

	return compileType()
}

func (c *Compiler) compileIntegerType() ([]byte, opcodes.ValType) {
	bytecode := make([]byte, 0)

	addr := c.chunk.AddIntConstant(0)

	bytecode = append(bytecode, opcodes.IPUSH)
	bytecode = append(bytecode, byte(addr))

	return bytecode, opcodes.INT
}

func (c *Compiler) compileFloatingPointType() ([]byte, opcodes.ValType) {
	bytecode := make([]byte, 0)

	addr := c.chunk.AddFloatConstant(0)

	bytecode = append(bytecode, opcodes.FPUSH)
	bytecode = append(bytecode, byte(addr))

	return bytecode, opcodes.FLOAT
}

func (c *Compiler) compileBooleanType() ([]byte, opcodes.ValType) {
	bytecode := make([]byte, 0)

	bytecode = append(bytecode, opcodes.BPUSH)
	bytecode = append(bytecode, 0)

	return bytecode, opcodes.BOOL
}

func (c *Compiler) compileCharacterType() ([]byte, opcodes.ValType) {
	bytecode := make([]byte, 0)

	addr := c.chunk.AddCharConstant(0)

	bytecode = append(bytecode, opcodes.CPUSH)
	bytecode = append(bytecode, byte(addr))

	return bytecode, opcodes.CHAR
}

func (c *Compiler) compileStringType() ([]byte, opcodes.ValType) {
	bytecode := make([]byte, 0)

	addr := c.chunk.AddStringConstant("")

	bytecode = append(bytecode, opcodes.SPUSH)
	bytecode = append(bytecode, byte(addr))

	return bytecode, opcodes.STRING
}

/// literals

func (c *Compiler) compileIdentifierLiteral() ([]byte, opcodes.ValType) {
	bytecode := make([]byte, 0)

	symbol, ok := c.currentScope.Lookup(c.current.Value)
	if !ok {
		return nil, opcodes.UNKNOWN_TYPE
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

	byteVal := byte(0)
	if val {
		byteVal = 1
	}

	bytecode = append(bytecode, opcodes.BPUSH)
	bytecode = append(bytecode, byteVal)

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

/// returns

func (c *Compiler) compileReturn(valType opcodes.ValType) ([]byte, opcodes.ValType) {
	compileReturn := c.returnCompileFns[valType]
	if compileReturn == nil {
		return nil, 0
	}

	return compileReturn()
}

func (c *Compiler) compileIntegerReturn() ([]byte, opcodes.ValType) {
	bytecode := make([]byte, 0)

	addr := c.chunk.AddIntConstant(0)

	bytecode = append(bytecode, opcodes.IFUNC_RETURN)
	bytecode = append(bytecode, byte(addr))

	return bytecode, opcodes.INT
}

func (c *Compiler) compileFloatingPointReturn() ([]byte, opcodes.ValType) {
	bytecode := make([]byte, 0)

	addr := c.chunk.AddFloatConstant(0)

	bytecode = append(bytecode, opcodes.FFUNC_RETURN)
	bytecode = append(bytecode, byte(addr))

	return bytecode, opcodes.FLOAT
}

func (c *Compiler) compileBooleanReturn() ([]byte, opcodes.ValType) {
	bytecode := make([]byte, 0)

	bytecode = append(bytecode, opcodes.BPUSH)
	bytecode = append(bytecode, 0)

	return bytecode, opcodes.BOOL
}

func (c *Compiler) compileCharacterReturn() ([]byte, opcodes.ValType) {
	bytecode := make([]byte, 0)

	addr := c.chunk.AddCharConstant(0)

	bytecode = append(bytecode, opcodes.CFUNC_RETURN)
	bytecode = append(bytecode, byte(addr))

	return bytecode, opcodes.CHAR
}

func (c *Compiler) compileStringReturn() ([]byte, opcodes.ValType) {
	bytecode := make([]byte, 0)

	addr := c.chunk.AddStringConstant("")

	bytecode = append(bytecode, opcodes.SFUNC_RETURN)
	bytecode = append(bytecode, byte(addr))

	return bytecode, opcodes.STRING
}

/// HELPERS

func (c *Compiler) registerTypes() {
	c.registerType(opcodes.INT, c.compileIntegerType)
	c.registerType(opcodes.FLOAT, c.compileFloatingPointType)
	c.registerType(opcodes.BOOL, c.compileBooleanType)
	c.registerType(opcodes.CHAR, c.compileCharacterType)
	c.registerType(opcodes.STRING, c.compileStringType)
}

func (c *Compiler) registerPrefixes() {
	c.registerPrefix(token.IDENTIFIER, c.compileIdentifier)
	c.registerPrefix(token.INT, c.compileIntegerLiteral)
	c.registerPrefix(token.FLOAT, c.compileFloatingPointLiteral)
	c.registerPrefix(token.BOOL, c.compileBooleanLiteral)
	c.registerPrefix(token.CHAR, c.compileCharacterLiteral)
	c.registerPrefix(token.STRING, c.compileStringLiteral)
	c.registerPrefix(token.LPARENTHESIS, c.compileGroupedExpression)
	c.registerPrefix(token.MINUS, c.compilePrefixExpression)
	c.registerPrefix(token.BANG, c.compilePrefixExpression)
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

func (c *Compiler) registerReturns() {
	c.registerReturn(opcodes.INT, c.compileIntegerReturn)
	c.registerReturn(opcodes.FLOAT, c.compileFloatingPointReturn)
	c.registerReturn(opcodes.BOOL, c.compileBooleanReturn)
	c.registerReturn(opcodes.CHAR, c.compileCharacterReturn)
	c.registerReturn(opcodes.STRING, c.compileStringReturn)
}

func (c *Compiler) registerType(valType opcodes.ValType, fn typeCompileFn) {
	c.typeCompileFns[valType] = fn
}

func (c *Compiler) registerPrefix(tokenType token.TokenType, fn prefixCompileFn) {
	c.prefixCompileFns[tokenType] = fn
}

func (c *Compiler) registerInfix(tokenType token.TokenType, fn infixCompileFn) {
	c.infixCompileFns[tokenType] = fn
}

func (c *Compiler) registerReturn(valType opcodes.ValType, fn returnCompileFn) {
	c.returnCompileFns[valType] = fn
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

func (c *Compiler) expectCurrent(tokens ...token.TokenType) bool {
	found := false
	for _, t := range tokens {
		found = found || c.current.Is(t)
	}

	if found {
		return true
	} else {
		c.reportUnexpectedToken(c.current, tokens...)
		return false
	}
}

func (c *Compiler) expectNext(tokens ...token.TokenType) bool {
	found := false
	for _, t := range tokens {
		found = found || c.next.Is(t)
	}

	if found {
		return true
	} else {
		c.reportUnexpectedToken(c.next, tokens...)
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

func (c *Compiler) reportUnexpectedToken(unexpected token.Token, expectedTypes ...token.TokenType) {
	expectedStr := ""
	if len(expectedTypes) == 1 {
		expectedStr = fmt.Sprintf("<%v>", expectedTypes[0])
	} else {
		for i, expectedType := range expectedTypes {
			if i == 0 {
				expectedStr = fmt.Sprintf("<%v>", expectedType)
			} else {
				expectedStr = fmt.Sprintf("%s or <%v>", expectedStr, expectedType)
			}
		}
	}

	msg := fmt.Sprintf("ERROR: Unexpected token <%v> on line %v, position %v. %s was expected.",
		unexpected.Type,
		unexpected.Line,
		unexpected.PositionInLine,
		expectedStr,
	)
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

func (c *Compiler) reportWrongType(token token.Token, expected, actual opcodes.ValType) {
	msg := fmt.Sprintf("ERROR: wrong type. Expected: '%v', got '%v' on line %v, position %v.",
		expected,
		actual,
		token.Line,
		token.PositionInLine)
	c.errors = append(c.errors, msg)
}

//// EXAMPLE 1

// ICONST 2
// ICONST 3
// VAR_BIND "a"
// ICONST 4
// ICONST 5
// VAR_BIND "b"
// VAR_LOOKUP "a"
// VAR_LOOKUP "b"
// IADD
// IPRINT
// HALT

//// EXAMPLE 2

// ; function definition: add two numbers
// FCREATE "add" 2 [ICONST 2 ICONST 3 IADD FRETURN]

// ; function call: add(1, 2)
// ICONST 1
// ICONST 2
// FCALL "add" 2
// IPRINT

// ; function call: add(2, 3)
// ICONST 2
// ICONST 3
// FCALL "add" 2
// IPRINT

// HALT
