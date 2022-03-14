package token

import (
	"fmt"
	"strings"
)

type TokenType string

const (
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

	OPEN_PARENTHESIS  TokenType = "("
	CLOSE_PARENTHESIS TokenType = ")"

	EQ  TokenType = "="
	NEQ TokenType = "!="
	GT  TokenType = ">"
	GTE TokenType = ">="
	LT  TokenType = "<"
	LTE TokenType = "<="

	NEGATION TokenType = "!"

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
		return "DECLARE_ASSIGN"
	case PLUS:
		return "PLUS"
	case MINUS:
		return "MINUS"
	case ASTERISK:
		return "ASTERISK"
	case SLASH:
		return "SLASH"
	case OPEN_PARENTHESIS:
		return "OPEN_PARENTHESIS"
	case CLOSE_PARENTHESIS:
		return "CLOSE_PARENTHESIS"
	case EQ:
		return "EQUAL"
	case NEQ:
		return "NOT_EQUAL"
	case GT:
		return "GT"
	case GTE:
		return "GTE"
	case LT:
		return "LT"
	case LTE:
		return "LTE"
	case NEGATION:
		return "NEGATION"
	case IDENTIFIER:
		return "IDENTIFIER"
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
}

func (t Token) String() string {
	val := strings.Replace(t.Value, "\n", "\\n", -1)
	val = strings.Replace(val, "\r", "\\r", -1)
	val = strings.Replace(val, "\t", "\\t", -1)
	return fmt.Sprintf("type=%v, value='%v', position=%v", t.Type, val, t.Position) // OK: note conversion.
}

func LookupKeyword(word string) TokenType {
	if tok, ok := keywords[word]; ok {
		return tok
	}
	return IDENTIFIER
}
