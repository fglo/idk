package compiler

// import (
// 	"github.com/fglo/idk/pkg/idk/ast"
// 	"github.com/fglo/idk/pkg/idk/opcodes"
// 	"github.com/fglo/idk/pkg/idk/symbol"
// 	"github.com/fglo/idk/pkg/idk/token"
// )

// var filepath string
// var scope *symbol.Scope

// func Compile(file string, program *ast.Program) []byte {
// 	scope = symbol.NewScope()

// 	filepath = file

// 	bytecode := make([]byte, 0)

// 	for _, statement := range program.Statements {
// 		bytecode = append(bytecode, compileStatement(statement)...)
// 	}

// 	return bytecode
// }

// func compileStatement(node ast.Statement) []byte {
// 	bytecode := make([]byte, 0)

// 	switch node := node.(type) {
// 	case *ast.DeclareAssignStatement:
// 		return compileDeclareAssignStatement(node)
// 	case *ast.DeclareStatement:
// 		return compileDeclareStatement(node)
// 	case *ast.AssignStatement:
// 		return compileAssignStatement(node)
// 	case *ast.ExpressionStatement:
// 		return compileExpression(node.Expression)
// 	}

// 	return bytecode
// }

// func compileExpression(node ast.Node) []byte {
// 	bytecode := make([]byte, 0)

// 	switch node := node.(type) {
// 	case *ast.Identifier:
// 		return compileIdentifierLiteral(node)
// 	case *ast.IntegerLiteral:
// 		return compileIntegerLiteral(node)
// 	case *ast.FloatingPointLiteral:
// 		return compileFloatingPointLiteral(node)
// 	case *ast.BooleanLiteral:
// 		return compileBooleanLiteral(node)
// 	case *ast.CharacterLiteral:
// 		return compileCharacterLiteral(node)
// 	case *ast.StringLiteral:
// 		return compileStringLiteral(node)
// 	case *ast.PrefixExpression:
// 		return compilePrefixExpression(node)
// 	case *ast.InfixExpression:
// 		return compileInfixExpression(node)
// 	case *ast.FunctionCallExpression:
// 		return compileFunctionCallExpression(node)
// 	}

// 	return bytecode
// }

// /// statements

// func compileDeclareAssignStatement(node *ast.DeclareAssignStatement) []byte {
// 	bytecode := make([]byte, 0)

// 	bytecode = append(bytecode, compileExpression(node.Expression)...)
// 	varName := node.Identifier.GetValue()
// 	bytecode = append(bytecode, opcodes.IVAR_BIND)
// 	bytecode = append(bytecode, byte(len(varName)))
// 	bytecode = append(bytecode, []byte(varName)...)

// 	return bytecode
// }

// func compileDeclareStatement(node *ast.DeclareStatement) []byte {
// 	bytecode := make([]byte, 0)

// 	bytecode = append(bytecode, opcodes.IPUSH)
// 	bytecode = append(bytecode, 0)
// 	varName := node.Identifier.GetValue()
// 	bytecode = append(bytecode, opcodes.IVAR_BIND)
// 	bytecode = append(bytecode, byte(len(varName)))
// 	bytecode = append(bytecode, []byte(varName)...)

// 	if node.Assignment != nil {
// 		bytecode = append(bytecode, []byte(compileAssignStatement(node.Assignment))...)
// 	}

// 	return bytecode
// }

// func compileAssignStatement(node *ast.AssignStatement) []byte {
// 	bytecode := make([]byte, 0)

// 	bytecode = append(bytecode, compileExpression(node.Expression)...)
// 	varName := node.Identifier.GetValue()
// 	bytecode = append(bytecode, opcodes.IVAR_BIND)
// 	bytecode = append(bytecode, byte(len(varName)))
// 	bytecode = append(bytecode, []byte(varName)...)

// 	return bytecode
// }

// // func compileIfStatement(node *ast.IfStatement) []byte {
// // 	bytecode := make([]byte, 0)

// // 	bytecode = append(bytecode, compileExpression(node.Condition)...)
// // 	varName := node.Identifier.GetValue()
// // 	bytecode = append(bytecode, opcodes.IVAR_BIND)
// // 	bytecode = append(bytecode, byte(len(varName)))
// // 	bytecode = append(bytecode, []byte(varName)...)

// // 	return bytecode
// // }

// /// literals

// func compileIdentifierLiteral(node *ast.Identifier) []byte {
// 	bytecode := make([]byte, 0)

// 	varName := node.GetValue()
// 	// switch node.GetTokenType() {

// 	// }
// 	bytecode = append(bytecode, opcodes.IVAR_LOOKUP)
// 	bytecode = append(bytecode, byte(len(varName)))
// 	bytecode = append(bytecode, []byte(varName)...)

// 	return bytecode
// }

// func compileIntegerLiteral(node *ast.IntegerLiteral) []byte {
// 	bytecode := make([]byte, 0)

// 	bytecode = append(bytecode, opcodes.IPUSH)
// 	bytecode = append(bytecode, byte(node.GetValue()))

// 	return bytecode
// }

// func compileFloatingPointLiteral(node *ast.FloatingPointLiteral) []byte {
// 	bytecode := make([]byte, 0)

// 	bytecode = append(bytecode, opcodes.FPUSH)
// 	bytecode = append(bytecode, byte(node.GetValue()))

// 	return bytecode
// }

// func compileBooleanLiteral(node *ast.BooleanLiteral) []byte {
// 	bytecode := make([]byte, 0)

// 	bytecode = append(bytecode, opcodes.BPUSH)
// 	if node.GetValue() {
// 		bytecode = append(bytecode, 1)
// 	} else {
// 		bytecode = append(bytecode, 0)
// 	}

// 	return bytecode
// }

// func compileCharacterLiteral(node *ast.CharacterLiteral) []byte {
// 	bytecode := make([]byte, 0)

// 	bytecode = append(bytecode, opcodes.FPUSH)
// 	bytecode = append(bytecode, byte(node.GetValue()))

// 	return bytecode
// }

// func compileStringLiteral(node *ast.StringLiteral) []byte {
// 	bytecode := make([]byte, 0)

// 	val := node.GetValue()
// 	bytecode = append(bytecode, opcodes.FPUSH)
// 	bytecode = append(bytecode, byte(len(val)))
// 	bytecode = append(bytecode, []byte(val)...)

// 	return bytecode
// }

// /// expresions

// func compilePrefixExpression(node *ast.PrefixExpression) []byte {
// 	bytecode := make([]byte, 0)

// 	switch node.GetTokenType() {
// 	case token.MINUS:
// 		bytecode = append(bytecode, opcodes.IPUSH)
// 		bytecode = append(bytecode, 0)
// 		bytecode = append(bytecode, compileExpression(node.Right)...)
// 		bytecode = append(bytecode, opcodes.ISUB)
// 	}

// 	return bytecode
// }

// func compileInfixExpression(node *ast.InfixExpression) []byte {
// 	bytecode := make([]byte, 0)

// 	bytecode = append(bytecode, compileExpression(node.Left)...)
// 	bytecode = append(bytecode, compileExpression(node.Right)...)

// 	switch node.GetTokenType() {
// 	case token.PLUS:
// 		bytecode = append(bytecode, opcodes.IADD)
// 	case token.MINUS:
// 		bytecode = append(bytecode, opcodes.ISUB)
// 	case token.ASTERISK:
// 		bytecode = append(bytecode, opcodes.IMUL)
// 	case token.SLASH:
// 		bytecode = append(bytecode, opcodes.IDIV)
// 	}

// 	return bytecode
// }

// func compileFunctionCallExpression(node *ast.FunctionCallExpression) []byte {
// 	bytecode := make([]byte, 0)

// 	switch node.Identifier.GetValue() {
// 	case "print":
// 		bytecode = append(bytecode, compileExpression(node.Parameters[0])...)
// 		bytecode = append(bytecode, opcodes.IPRINT)
// 	}

// 	return bytecode
// }

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
