package token

import (
	"fmt"
	"strings"
)

type TokenType string

const (
	NONE TokenType = "\\0"

	ILLEGAL TokenType = "ILLEGAL"

	EOL TokenType = "EOL"
	EOF TokenType = "EOF"

	INT    TokenType = "INT"
	CHAR   TokenType = "CHAR"
	STRING TokenType = "STRING"
	ARRAY  TokenType = "ARRAY"
	BOOL   TokenType = "BOOL"

	TRUE  TokenType = "TRUE"
	FALSE TokenType = "FALSE"

	DECLARE_ASSIGN TokenType = ":="

	PLUS     TokenType = "+"
	MINUS    TokenType = "-"
	ASTERISK TokenType = "*"
	SLASH    TokenType = "/"

	LPARENTHESIS TokenType = "("
	RPARENTHESIS TokenType = ")"

	EQ  TokenType = "="
	NEQ TokenType = "!="
	GT  TokenType = ">"
	GTE TokenType = ">="
	LT  TokenType = "<"
	LTE TokenType = "<="

	NEGATION TokenType = "!"

	NOT TokenType = "NOT"
	AND TokenType = "AND"
	OR  TokenType = "OR"
	XOR TokenType = "XOR"

	IF    TokenType = "IF"
	ELSE  TokenType = "ELSE"
	FOR   TokenType = "FOR"
	END   TokenType = "END"
	PRINT TokenType = "PRINT"

	IDENTIFIER TokenType = "IDENTIFIER"
)

func (e TokenType) String() string {
	switch e {
	case EOL:
		return "EOL"
	case EOF:
		return "EOF"
	case ILLEGAL:
		return "ILLEGAL"
	case INT:
		return "INT"
	case CHAR:
		return "CHAR"
	case BOOL:
		return "BOOL"
	case ARRAY:
		return "ARRAY"
	case DECLARE_ASSIGN:
		return "DECLASSIGN"
	case PLUS:
		return "PLUS"
	case MINUS:
		return "MINUS"
	case ASTERISK:
		return "ASTERISK"
	case SLASH:
		return "SLASH"
	case LPARENTHESIS:
		return "LPARENT"
	case RPARENTHESIS:
		return "RPARENT"
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
	case NEGATION:
		return "NEG"
	case NOT:
		return "NOT"
	case AND:
		return "AND"
	case OR:
		return "OR"
	case XOR:
		return "XOR"
	case IDENTIFIER:
		return "IDENTIFIER"
	case IF:
		return "IF"
	case ELSE:
		return "ELSE"
	case FOR:
		return "FOR"
	case END:
		return "END"
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

func NewToken(tokenType TokenType, position, line, positionInLine int, value string) *Token {
	t := new(Token)
	t.Type = tokenType
	t.Position = position
	t.Line = line
	t.PositionInLine = positionInLine
	t.Value = value
	return t
}

var keywords = map[string]TokenType{
	"true":  TRUE,
	"false": FALSE,
	"if":    IF,
	"else":  ELSE,
	"end":   END,
	"print": PRINT,
	"not":   NOT,
	"and":   AND,
	"or":    OR,
	"xor":   XOR,
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

var operators = map[TokenType]byte{
	PLUS:     0,
	MINUS:    0,
	ASTERISK: 0,
	SLASH:    0,

	EQ:  0,
	NEQ: 0,
	GT:  0,
	GTE: 0,
	LT:  0,
	LTE: 0,

	AND: 0,
	OR:  0,
	XOR: 0,
}

func (t TokenType) IsOperator() bool {
	_, ok := operators[t]
	return ok
}
