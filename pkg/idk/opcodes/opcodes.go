package opcodes

import (
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

	// string operations
	SPUSH
	SCONCAT
	SEQ

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
	case IPRINT:
		return "IPRINT"
	case IVAR_BIND:
		return "IVAR_BIND"
	case IVAR_LOOKUP:
		return "IVAR_LOOKUP"
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
)

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
	default:
		return "UNKNOWN"
	}
}

func GetInfixOperator(operator token.TokenType, varType ValType) byte {
	switch operator {
	case token.PLUS:
		switch varType {
		case INT:
			return IADD
		case FLOAT:
			return FADD
		case STRING:
			return SCONCAT
		}
	case token.MINUS:
		switch varType {
		case INT:
			return ISUB
		case FLOAT:
			return FSUB
		}
	case token.ASTERISK:
		switch varType {
		case INT:
			return IMUL
		case FLOAT:
			return FMUL
		}
	case token.SLASH:
		switch varType {
		case INT:
			return IDIV
		case FLOAT:
			return FDIV
		}
	case token.MODULO:
		switch varType {
		case INT:
			return IMOD
		case FLOAT:
			return FMOD
		}
	}

	return UNKNOWN_OP
}

func GetPrefixOperator(operator token.TokenType, varType ValType) byte {
	switch operator {
	case token.MINUS:
		switch varType {
		case INT:
			return INEG
		case FLOAT:
			return FNEG
		}
	case token.BANG:
		switch varType {
		case BOOL:
			return BNEG
		}
	}

	return UNKNOWN_OP
}

func VarBind(varType ValType) byte {
	switch varType {
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
	}

	return UNKNOWN_OP
}

func VarLookup(varType ValType) byte {
	switch varType {
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
	}

	return UNKNOWN_OP
}
