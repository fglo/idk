package common

import (
	"github.com/fglo/idk/pkg/idk/symbol"
	"github.com/fglo/idk/pkg/idk/token"
)

func ToObjectType(tt token.TokenType) symbol.ObjectType {
	switch tt {
	case token.INT:
		return symbol.INTEGER_OBJ
	case token.FLOAT:
		return symbol.FLOATING_POINT_OBJ
	case token.BOOL:
		return symbol.BOOLEAN_OBJ
	case token.CHAR:
		return symbol.CHARACTER_OBJ
	case token.STRING:
		return symbol.STRING_OBJ
	case token.ARRAY:
		return symbol.ARRAY_OBJ
	case token.FUNC:
		return symbol.FUNCTION_OBJ
	default:
		return symbol.NULL_OBJ
	}
}

func ToTokenType(ot symbol.ObjectType) token.TokenType {
	switch ot {
	case symbol.INTEGER_OBJ:
		return token.INT
	case symbol.FLOATING_POINT_OBJ:
		return token.FLOAT
	case symbol.BOOLEAN_OBJ:
		return token.BOOL
	case symbol.CHARACTER_OBJ:
		return token.CHAR
	case symbol.STRING_OBJ:
		return token.STRING
	case symbol.ARRAY_OBJ:
		return token.ARRAY
	case symbol.FUNCTION_OBJ:
		return token.FUNC
	default:
		return token.NONE
	}
}
