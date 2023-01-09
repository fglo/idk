package ast

import (
	"github.com/fglo/idk/pkg/idk/token"
)

type Node interface {
	GetTokenValue() string
	GetTokenType() token.TokenType
	GetLineNumber() int
	GetPositionInLine() int
	GetChildren() []Node
	String() string
}
