package token

import (
	"fmt"
	"strings"
)

type TokenType string

const (
	NONE TokenType = "\\0"

	LINE_COMMENT TokenType = "//"

	ILLEGAL TokenType = "ILLEGAL"

	EOL TokenType = "EOL"
	EOF TokenType = "EOF"

	TYPE TokenType = "TYPE"

	INT    TokenType = "INT"
	FLOAT  TokenType = "FLOAT"
	CHAR   TokenType = "CHAR"
	STRING TokenType = "STRING"
	ARRAY  TokenType = "ARRAY"
	BOOL   TokenType = "BOOL"
	VOID   TokenType = "VOID"

	TRUE  TokenType = "TRUE"
	FALSE TokenType = "FALSE"

	DECLASSIGN TokenType = "DECLASSIGN"
	DECLARE    TokenType = ":"
	ASSIGN     TokenType = "="

	RANGE           TokenType = ".."  // TODO: range
	RANGE_INCLUSIVE TokenType = "..=" // TODO: range inclusive

	TO TokenType = "TO"

	PLUS     TokenType = "+"
	MINUS    TokenType = "-"
	ASTERISK TokenType = "*"
	SLASH    TokenType = "/"
	MODULO   TokenType = "%"

	LPARENTHESIS TokenType = "("
	RPARENTHESIS TokenType = ")"

	EQ  TokenType = "=="
	NEQ TokenType = "!="
	GT  TokenType = ">"
	GTE TokenType = ">="
	LT  TokenType = "<"
	LTE TokenType = "<="

	COMMA TokenType = ","
	DOT   TokenType = "."

	NOT TokenType = "!"
	AND TokenType = "AND"
	OR  TokenType = "OR"
	XOR TokenType = "XOR"

	IF   TokenType = "IF"
	ELSE TokenType = "ELSE"
	FOR  TokenType = "FOR"
	END  TokenType = "END"
	IN   TokenType = "IN"

	FUNC   TokenType = "FUNC"
	RETURN TokenType = "RETURN"

	IDENTIFIER TokenType = "IDENTIFIER"
)

func (e TokenType) String() string {
	switch e {
	case EOL:
		return "EOL"
	case EOF:
		return "EOF"
	case LINE_COMMENT:
		return "LINE_COMMENT"
	case ARRAY:
		return "ARRAY"
	case DECLASSIGN:
		return "DECLASSIGN"
	case DECLARE:
		return "DECLARE"
	case ASSIGN:
		return "ASSIGN"
	case RANGE:
		return "RANGE"
	case RANGE_INCLUSIVE:
		return "RANGE_INCLUSIVE"
	case PLUS:
		return "PLUS"
	case MINUS:
		return "MINUS"
	case ASTERISK:
		return "ASTERISK"
	case SLASH:
		return "SLASH"
	case MODULO:
		return "MODULO"
	case LPARENTHESIS:
		return "LPARENTHESIS"
	case RPARENTHESIS:
		return "RPARENTHESIS"
	case EQ:
		return "EQ"
	case NEQ:
		return "NEQ"
	case GT:
		return "GT"
	case GTE:
		return "GTE"
	case LT:
		return "LT"
	case LTE:
		return "LTE"
	case COMMA:
		return "COMMA"
	case DOT:
		return "DOT"
	case NOT:
		return "NOT"
	default:
		return string(e)
	}
}

type Token struct {
	Type           TokenType
	Position       int
	Line           int
	PositionInLine int
	Value          string
}

func NewToken(tokenType TokenType, position, line, positionInLine int) *Token {
	t := new(Token)
	t.Type = tokenType
	t.Position = position
	t.Line = line
	t.PositionInLine = positionInLine
	t.Value = string(tokenType)
	return t
}

func NewTokenNotDefaultValue(tokenType TokenType, position, line, positionInLine int, value string) *Token {
	t := new(Token)
	t.Type = tokenType
	t.Position = position
	t.Line = line
	t.PositionInLine = positionInLine
	t.Value = value
	return t
}

func (t Token) Is(tokenType TokenType) bool {
	return t.Type == tokenType
}

func (t Token) Not(tokenType TokenType) bool {
	return t.Type != tokenType
}

var keywords = map[string]TokenType{
	"int":    TYPE,
	"float":  TYPE,
	"char":   TYPE,
	"string": TYPE,
	"bool":   TYPE,
	"void":   TYPE,
	"true":   BOOL,
	"false":  BOOL,
	"if":     IF,
	"else":   ELSE,
	"for":    FOR,
	"end":    END,
	"not":    NOT,
	"and":    AND,
	"or":     OR,
	"xor":    XOR,
	"in":     IN,
	"func":   FUNC,
	"return": RETURN,
}

var types = map[string]TokenType{
	"int":    INT,
	"float":  FLOAT,
	"char":   CHAR,
	"string": STRING,
	"bool":   BOOL,
	"void":   VOID,
	"func":   FUNC,
}

func (t Token) String() string {
	val := strings.Replace(t.Value, "\n", "\\n", -1)
	val = strings.Replace(val, "\r", "\\r", -1)
	val = strings.Replace(val, "\t", "\\t", -1)
	return fmt.Sprintf("type=%v, value='%v', line=%v, position=%v", t.Type, val, t.Line, t.PositionInLine)
}

func LookupKeyword(word string) TokenType {
	if tok, ok := keywords[word]; ok {
		return tok
	}
	return IDENTIFIER
}

func LookupType(typeword string) TokenType {
	if tok, ok := types[typeword]; ok {
		return tok
	}
	return ""
}

var operators = map[TokenType]byte{
	PLUS:     0,
	MINUS:    0,
	ASTERISK: 0,
	SLASH:    0,
	MODULO:   0,

	EQ:  0,
	NEQ: 0,
	GT:  0,
	GTE: 0,
	LT:  0,
	LTE: 0,

	NOT: 0,
	AND: 0,
	OR:  0,
	XOR: 0,
}

func (t TokenType) IsOperator() bool {
	_, ok := operators[t]
	return ok
}
