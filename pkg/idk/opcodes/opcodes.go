package opcodes

import (
	"strings"

	"github.com/fglo/idk/pkg/idk/token"
)

const (
	UNKNOWN_OP byte = iota

	// integer operations
	IPUSH
	IADD
	ISUB
	IMUL
	IDIV
	IMOD
	INEG
	IEQ
	INEQ
	ILT
	ILE
	IGT
	IGE

	// float operations
	FPUSH
	FADD
	FSUB
	FMUL
	FDIV
	FMOD
	FNEG
	FEQ
	FNEQ
	FLT
	FLE
	FGT
	FGE

	// bool operations
	BPUSH
	BNEG
	BEQ
	BNEQ
	BAND
	BOR
	BXOR

	// char operations
	CPUSH
	CEQ
	CNEQ

	// string operations
	SPUSH
	SCONCAT
	SEQ
	SNEQ

	// branches
	BT
	BF

	// jumps
	JMP

	// variables
	IVAR_BIND
	IVAR_LOOKUP

	FVAR_BIND
	FVAR_LOOKUP

	BVAR_BIND
	BVAR_LOOKUP

	CVAR_BIND
	CVAR_LOOKUP

	SVAR_BIND
	SVAR_LOOKUP

	FUNC_VAR_BIND
	FUNC_VAR_LOOKUP

	// functions
	IFUNC_CREATE
	IFUNC_CALL
	IFUNC_RETURN

	FFUNC_CREATE
	FFUNC_CALL
	FFUNC_RETURN

	BFUNC_CREATE
	BFUNC_CALL
	BFUNC_RETURN

	CFUNC_CREATE
	CFUNC_CALL
	CFUNC_RETURN

	SFUNC_CREATE
	SFUNC_CALL
	SFUNC_RETURN

	// built-ins
	IPRINT
	FPRINT
	BPRINT
	CPRINT
	SPRINT

	IF
	ELSE
	ENDIF
	FOR
	NEXT
	BREAK

	HALT
)

func ToString(opc byte) string {
	switch opc {
	case IPUSH:
		return "IPUSH"
	case IADD:
		return "IADD"
	case ISUB:
		return "ISUB"
	case IMUL:
		return "IMUL"
	case IDIV:
		return "IDIV"
	case IMOD:
		return "IMOD"
	case INEG:
		return "INEG"
	case IGT:
		return "IGT"
	case ILT:
		return "ILT"
	case IGE:
		return "IGE"
	case ILE:
		return "ILE"
	case IEQ:
		return "IEQ"
	case INEQ:
		return "INEQ"
	case IPRINT:
		return "IPRINT"
	case IVAR_BIND:
		return "IVAR_BIND"
	case IVAR_LOOKUP:
		return "IVAR_LOOKUP"
	case IFUNC_CREATE:
		return "IFUNC_CREATE"
	case IFUNC_CALL:
		return "IFUNC_CALL"
	case IFUNC_RETURN:
		return "IFUNC_RETURN"
	case FPUSH:
		return "FPUSH"
	case FADD:
		return "FADD"
	case FSUB:
		return "FSUB"
	case FMUL:
		return "FMUL"
	case FDIV:
		return "FDIV"
	case FMOD:
		return "FMOD"
	case FNEG:
		return "FNEG"
	case FGT:
		return "FGT"
	case FLT:
		return "FLT"
	case FGE:
		return "FGE"
	case FLE:
		return "FLE"
	case FEQ:
		return "FEQ"
	case FNEQ:
		return "FNEQ"
	case FPRINT:
		return "FPRINT"
	case FVAR_BIND:
		return "FVAR_BIND"
	case FVAR_LOOKUP:
		return "FVAR_LOOKUP"
	case BPUSH:
		return "BPUSH"
	case BNEG:
		return "BNEG"
	case BEQ:
		return "BEQ"
	case BNEQ:
		return "BNEQ"
	case BAND:
		return "BAND"
	case BOR:
		return "BOR"
	case BXOR:
		return "BXOR"
	case BPRINT:
		return "BPRINT"
	case BVAR_BIND:
		return "BVAR_BIND"
	case BVAR_LOOKUP:
		return "BVAR_LOOKUP"
	case CPUSH:
		return "CPUSH"
	case CEQ:
		return "CEQ"
	case CNEQ:
		return "CNEQ"
	case CPRINT:
		return "CPRINT"
	case CVAR_BIND:
		return "CVAR_BIND"
	case CVAR_LOOKUP:
		return "CVAR_LOOKUP"
	case SPUSH:
		return "SPUSH"
	case SEQ:
		return "SEQ"
	case SNEQ:
		return "SNEQ"
	case SPRINT:
		return "SPRINT"
	case SVAR_BIND:
		return "SVAR_BIND"
	case SVAR_LOOKUP:
		return "SVAR_LOOKUP"
	default:
		return string(opc)
	}
}

type ValType int

const (
	UNKNOWN_TYPE ValType = iota
	INT
	FLOAT
	BOOL
	CHAR
	STRING
	FUNC
)

func ValTypeFromString(vtStr string) ValType {
	switch strings.ToUpper(vtStr) {
	case "INT":
		return INT
	case "FLOAT":
		return FLOAT
	case "BOOL":
		return BOOL
	case "CHAR":
		return CHAR
	case "STRING":
		return STRING
	case "FUNC":
		return FUNC
	default:
		return UNKNOWN_TYPE
	}
}

func (vt ValType) String() string {
	switch vt {
	case INT:
		return "INT"
	case FLOAT:
		return "FLOAT"
	case BOOL:
		return "BOOL"
	case CHAR:
		return "CHAR"
	case STRING:
		return "STRING"
	case FUNC:
		return "FUNC"
	default:
		return "UNKNOWN"
	}
}

func InfixOperator(operator token.TokenType, valType ValType) byte {
	switch operator {
	case token.PLUS:
		switch valType {
		case INT:
			return IADD
		case FLOAT:
			return FADD
		case STRING:
			return SCONCAT
		}
	case token.MINUS:
		switch valType {
		case INT:
			return ISUB
		case FLOAT:
			return FSUB
		}
	case token.ASTERISK:
		switch valType {
		case INT:
			return IMUL
		case FLOAT:
			return FMUL
		}
	case token.SLASH:
		switch valType {
		case INT:
			return IDIV
		case FLOAT:
			return FDIV
		}
	case token.MODULO:
		switch valType {
		case INT:
			return IMOD
		case FLOAT:
			return FMOD
		}
	case token.GT:
		switch valType {
		case INT:
			return IGT
		case FLOAT:
			return FGT
		}
	case token.LT:
		switch valType {
		case INT:
			return ILT
		case FLOAT:
			return FLT
		}
	case token.GTE:
		switch valType {
		case INT:
			return IGE
		case FLOAT:
			return FGE
		}
	case token.LTE:
		switch valType {
		case INT:
			return ILE
		case FLOAT:
			return FLE
		}
	case token.EQ:
		switch valType {
		case INT:
			return IEQ
		case FLOAT:
			return FEQ
		case BOOL:
			return BEQ
		case CHAR:
			return CEQ
		case STRING:
			return SEQ
		}
	case token.NEQ:
		switch valType {
		case INT:
			return INEQ
		case FLOAT:
			return FNEQ
		case BOOL:
			return BNEQ
		case CHAR:
			return CNEQ
		case STRING:
			return SNEQ
		}
	case token.AND:
		switch valType {
		case BOOL:
			return BAND
		}
	case token.OR:
		switch valType {
		case BOOL:
			return BOR
		}
	case token.XOR:
		switch valType {
		case BOOL:
			return BXOR
		}
	}

	return UNKNOWN_OP
}

func PrefixOperator(operator token.TokenType, valType ValType) byte {
	switch operator {
	case token.MINUS:
		switch valType {
		case INT:
			return INEG
		case FLOAT:
			return FNEG
		}
	case token.BANG:
		switch valType {
		case BOOL:
			return BNEG
		}
	}

	return UNKNOWN_OP
}

func OperatorResultType(operator token.TokenType, inputType ValType) ValType {
	switch operator {
	case token.PLUS:
		return inputType
	case token.MINUS:
		return inputType
	case token.ASTERISK:
		return inputType
	case token.SLASH:
		return inputType
	case token.MODULO:
		return inputType
	case token.GT:
		return BOOL
	case token.LT:
		return BOOL
	case token.GTE:
		return BOOL
	case token.LTE:
		return BOOL
	case token.EQ:
		return BOOL
	case token.NEQ:
		return BOOL
	}

	return inputType
}

func VarBind(valType ValType) byte {
	switch valType {
	case INT:
		return IVAR_BIND
	case FLOAT:
		return FVAR_BIND
	case BOOL:
		return BVAR_BIND
	case CHAR:
		return CVAR_BIND
	case STRING:
		return SVAR_BIND
	case FUNC:
		return FUNC_VAR_BIND
	}

	return UNKNOWN_OP
}

func VarLookup(valType ValType) byte {
	switch valType {
	case INT:
		return IVAR_LOOKUP
	case FLOAT:
		return FVAR_LOOKUP
	case BOOL:
		return BVAR_LOOKUP
	case CHAR:
		return CVAR_LOOKUP
	case STRING:
		return SVAR_LOOKUP
	case FUNC:
		return FUNC_VAR_LOOKUP
	}

	return UNKNOWN_OP
}

func FuncCreate(valType ValType) byte {
	switch valType {
	case INT:
		return IFUNC_CREATE
	case FLOAT:
		return FFUNC_CREATE
	case BOOL:
		return BFUNC_CREATE
	case CHAR:
		return CFUNC_CREATE
	case STRING:
		return SFUNC_CREATE
	}

	return UNKNOWN_OP
}

func FuncCall(valType ValType) byte {
	switch valType {
	case INT:
		return IFUNC_CALL
	case FLOAT:
		return FFUNC_CALL
	case BOOL:
		return BFUNC_CALL
	case CHAR:
		return CFUNC_CALL
	case STRING:
		return SFUNC_CALL
	}

	return UNKNOWN_OP
}

func FuncReturn(valType ValType) byte {
	switch valType {
	case INT:
		return IFUNC_RETURN
	case FLOAT:
		return FFUNC_RETURN
	case BOOL:
		return BFUNC_RETURN
	case CHAR:
		return CFUNC_RETURN
	case STRING:
		return SFUNC_RETURN
	}

	return UNKNOWN_OP
}

func ValPrint(valType ValType) byte {
	switch valType {
	case INT:
		return IPRINT
	case FLOAT:
		return FPRINT
	case BOOL:
		return BPRINT
	case CHAR:
		return CPRINT
	case STRING:
		return SPRINT
	}

	return UNKNOWN_OP
}
