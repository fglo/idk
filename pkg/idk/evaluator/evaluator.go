package evaluator

import (
	"fmt"

	"github.com/fglo/idk/pkg/idk/ast"
	"github.com/fglo/idk/pkg/idk/common"
	"github.com/fglo/idk/pkg/idk/symbol"
	"github.com/fglo/idk/pkg/idk/token"
)

var (
	NULL  = &symbol.Null{}
	TRUE  = &symbol.Boolean{Value: true}
	FALSE = &symbol.Boolean{Value: false}
)

func GetDefaultValue(identifier ast.Identifier) symbol.Object {
	switch identifier.Type {
	case token.INT:
		return &symbol.Integer{Value: int64(0)}
	case token.FLOAT:
		return &symbol.FloatingPoint{Value: float64(0)}
	case token.CHAR:
		return &symbol.Character{Value: 0}
	case token.STRING:
		return &symbol.String{Value: ""}
	// case token.ARRAY:
	// 	return &symbol.Array{Value: ""}
	case token.BOOL:
		return &symbol.Boolean{Value: false}
	case token.FUNC:
		return &symbol.Function{}
	}
	return &symbol.Null{}
}

func Eval(node ast.Node, scope *symbol.Scope) symbol.Object {
	switch node := node.(type) {

	// Statements
	case *ast.Program:
		return evalProgram(node, scope)

	case *ast.BlockStatement:
		return evalBlockStatement(node, scope)

	case *ast.ExpressionStatement:
		return Eval(node.Expression, scope)

	case *ast.ReturnStatement:
		val := Eval(node.Expression, scope)
		if symbol.IsError(val) {
			return val
		}

		return &symbol.ReturnValue{Value: val}

	case *ast.DeclareAssignStatement:
		result := evalDeclareAssignStetment(node, scope)
		if symbol.IsError(result) {
			return newEvaluatorError(node.Identifier.Token.Line, result.Inspect())
		}

		return result

	case *ast.DeclareStatement:
		result := evalDeclareStatement(node, scope)
		if symbol.IsError(result) {
			return newEvaluatorError(node.Identifier.Token.Line, result.Inspect())
		}

		return result

	case *ast.AssignStatement:
		result := evalAssignStatement(node, scope)
		if symbol.IsError(result) {
			return newEvaluatorError(node.Identifier.Token.Line, result.Inspect())
		}

		return result

	// Expressions
	case *ast.Type:
		tokenType := token.LookupType(node.Token.Value)
		objType := common.ToObjectType(tokenType)
		return &symbol.Type{Value: objType}

	case *ast.IntegerLiteral:
		return &symbol.Integer{Value: int64(node.Value)}

	case *ast.FloatingPointLiteral:
		return &symbol.FloatingPoint{Value: float64(node.Value)}

	case *ast.BooleanLiteral:
		return nativeBoolToBooleanObject(node.Value)

	case *ast.CharacterLiteral:
		return &symbol.Character{Value: node.Value}

	case *ast.StringLiteral:
		return &symbol.String{Value: node.Value}

	// case *ast.ArrayLiteral:
	// 	elements := evalExpressions(node.Elements, env)
	// 	if len(elements) == 1 && symbol.IsError(elements[0]) {
	// 		return elements[0]
	// 	}
	// 	return &symbol.Array{Elements: elements}

	// case *ast.IndexExpression:
	// 	left := Eval(node.Left, env)
	// 	if symbol.IsError(left) {
	// 		return left
	// 	}
	// 	index := Eval(node.Index, env)
	// 	if symbol.IsError(index) {
	// 		return index
	// 	}
	// 	return evalIndexExpression(left, index)

	// case *ast.HashLiteral:
	// 	return evalHashLiteral(node, env)

	// case *ast.FunctionLiteral:
	// 	params := node.Parameters
	// 	body := node.Body
	// 	return &symbol.Function{Parameters: params, Env: env, Body: body}

	case *ast.PrefixExpression:
		right := Eval(node.Right, scope)
		if symbol.IsError(right) {
			return right
		}

		result := evalPrefixExpression(node.Token.Value, right)
		if symbol.IsError(result) {
			return newEvaluatorError(node.Token.Line, result.Inspect())
		}

		return result

	case *ast.InfixExpression:
		left := Eval(node.Left, scope)
		if symbol.IsError(left) {
			return left
		}

		right := Eval(node.Right, scope)
		if symbol.IsError(right) {
			return right
		}

		result := evalInfixExpression(node.Token.Value, left, right)
		if symbol.IsError(result) {
			return newEvaluatorError(node.Token.Line, result.Inspect())
		}

		return result

	case *ast.IfStatement:
		return evalIfStatement(node, scope)

	case *ast.ForLoopStatement:
		return evalForLoopStatement(node, scope)

	case *ast.IfExpression:
		return evalIfExpression(node, scope)

	case *ast.Identifier:
		return evalIdentifier(node, scope)

	case *ast.FunctionDefinitionStatement:
		function := evalIdentifier(&node.Identifier, scope)
		if !symbol.IsError(function) {
			return newError("Evaluator error on line %v, position %v: identifier %s is already taken", node.Identifier.Token.Line, node.Identifier.Token.PositionInLine, node.Identifier.Value)
		}

		returnType := token.LookupType(node.ReturnType.Value)
		function = &symbol.Function{
			Identifier: node.Identifier.Value,
			Parameters: node.Parameters,
			Scope:      scope,
			Body:       node.Body,
			ReturnType: common.ToObjectType(returnType),
		}

		scope.Insert(node.Identifier.Value, function, symbol.FUNCTION_OBJ)

	case *ast.FunctionCallExpression:
		return evalFunctionCallExpression(node, scope)
	}

	return nil
}

func evalProgram(program *ast.Program, env *symbol.Scope) symbol.Object {
	var result symbol.Object

	for _, statement := range program.Statements {
		result = Eval(statement, env)

		switch result := result.(type) {
		case *symbol.ReturnValue:
			return result.Value
		case *symbol.Error:
			return result
		}
	}

	return result
}

func evalBlockStatement(
	block *ast.BlockStatement,
	env *symbol.Scope,
) symbol.Object {
	var result symbol.Object

	for _, statement := range block.Statements {
		result = Eval(statement, env)

		if result != nil {
			rt := result.Type()
			if rt == symbol.RETURN_VALUE_OBJ || rt == symbol.ERROR_OBJ {
				return result
			}
		}
	}

	return result
}

func nativeBoolToBooleanObject(input bool) *symbol.Boolean {
	if input {
		return TRUE
	}
	return FALSE
}

func evalPrefixExpression(operator string, right symbol.Object) symbol.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	default:
		return newError("unknown operator: %s%s", operator, right.Type())
	}
}

// TODO: adding bool and chars to ints
// TODO: operations on ints and floats
// TODO: flaot to int and int to float conversion
func evalInfixExpression(
	operator string,
	left, right symbol.Object,
) symbol.Object {
	switch {
	case left.Type() == symbol.TYPE_OBJ && right.Type() == symbol.TYPE_OBJ:
		return evalTypeInfixExpression(operator, left, right)
	case left.Type() == symbol.INTEGER_OBJ && right.Type() == symbol.INTEGER_OBJ:
		return evalIntegerInfixExpression(operator, left, right)
	case left.Type() == symbol.FLOATING_POINT_OBJ && right.Type() == symbol.FLOATING_POINT_OBJ:
		return evalFloatingPointInfixExpression(operator, left, right)
	case left.Type() == symbol.BOOLEAN_OBJ && right.Type() == symbol.BOOLEAN_OBJ:
		return evalBooleanInfixExpression(operator, left, right)
	case left.Type() == symbol.CHARACTER_OBJ && right.Type() == symbol.CHARACTER_OBJ:
		return evalCharacterInfixExpression(operator, left, right)
	case left.Type() == symbol.STRING_OBJ && right.Type() == symbol.STRING_OBJ:
		return evalStringInfixExpression(operator, left, right)
	// case operator == "==":
	// 	return nativeBoolToBooleanObject(left == right)
	// case operator == "!=":
	// 	return nativeBoolToBooleanObject(left != right)
	case left.Type() != right.Type():
		return newError("type mismatch: %s %s %s",
			left.Type(), operator, right.Type())
	default:
		return newError("unknown operator: %s %s %s",
			left.Type(), operator, right.Type())
	}
}

func evalBangOperatorExpression(right symbol.Object) symbol.Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		return FALSE
	}
}

func evalMinusPrefixOperatorExpression(right symbol.Object) symbol.Object {
	if right.Type() != symbol.INTEGER_OBJ && right.Type() != symbol.FLOATING_POINT_OBJ {
		return newError("unknown operator: -%s", right.Type())
	}

	if right.Type() == symbol.INTEGER_OBJ {
		value := right.(*symbol.Integer).Value
		return &symbol.Integer{Value: -value}
	}

	value := right.(*symbol.FloatingPoint).Value
	return &symbol.FloatingPoint{Value: -value}
}

func evalTypeInfixExpression(
	operator string,
	left, right symbol.Object,
) symbol.Object {
	leftVal := left.(*symbol.Type).Value
	rightVal := right.(*symbol.Type).Value

	switch operator {
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	default:
		return newError("unknown operator: %s %s %s",
			left.Type(), operator, right.Type())
	}
}

func evalIntegerInfixExpression(
	operator string,
	left, right symbol.Object,
) symbol.Object {
	leftVal := left.(*symbol.Integer).Value
	rightVal := right.(*symbol.Integer).Value

	switch operator {
	case "+":
		return &symbol.Integer{Value: leftVal + rightVal}
	case "-":
		return &symbol.Integer{Value: leftVal - rightVal}
	case "*":
		return &symbol.Integer{Value: leftVal * rightVal}
	case "/":
		return &symbol.Integer{Value: leftVal / rightVal}
	case "%":
		return &symbol.Integer{Value: leftVal % rightVal}
	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case "<=":
		return nativeBoolToBooleanObject(leftVal <= rightVal)
	case ">=":
		return nativeBoolToBooleanObject(leftVal >= rightVal)
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	default:
		return newError("unknown operator: %s %s %s",
			left.Type(), operator, right.Type())
	}
}

func evalFloatingPointInfixExpression(
	operator string,
	left, right symbol.Object,
) symbol.Object {
	leftVal := left.(*symbol.FloatingPoint).Value
	rightVal := right.(*symbol.FloatingPoint).Value

	switch operator {
	case "+":
		return &symbol.FloatingPoint{Value: leftVal + rightVal}
	case "-":
		return &symbol.FloatingPoint{Value: leftVal - rightVal}
	case "*":
		return &symbol.FloatingPoint{Value: leftVal * rightVal}
	case "/":
		return &symbol.FloatingPoint{Value: leftVal / rightVal}
	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case "<=":
		return nativeBoolToBooleanObject(leftVal <= rightVal)
	case ">=":
		return nativeBoolToBooleanObject(leftVal >= rightVal)
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	default:
		return newError("unknown operator: %s %s %s",
			left.Type(), operator, right.Type())
	}
}

func evalBooleanInfixExpression(
	operator string,
	left, right symbol.Object,
) symbol.Object {
	leftVal := left.(*symbol.Boolean).Value
	rightVal := right.(*symbol.Boolean).Value

	switch operator {
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	case "and":
		return nativeBoolToBooleanObject(leftVal && rightVal)
	case "or":
		return nativeBoolToBooleanObject(leftVal || rightVal)
	default:
		return newError("unknown operator: %s %s %s",
			left.Type(), operator, right.Type())
	}
}

func evalCharacterInfixExpression(
	operator string,
	left, right symbol.Object,
) symbol.Object {
	leftVal := left.(*symbol.Character).Value
	rightVal := right.(*symbol.Character).Value

	switch operator {
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	default:
		return newError("unknown operator: %s %s %s",
			left.Type(), operator, right.Type())
	}
}

func evalStringInfixExpression(
	operator string,
	left, right symbol.Object,
) symbol.Object {
	leftVal := left.(*symbol.String).Value
	rightVal := right.(*symbol.String).Value

	switch operator {
	case "+":
		return &symbol.String{Value: leftVal + rightVal}
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	default:
		return newError("unknown operator: %s %s %s",
			left.Type(), operator, right.Type())
	}
}

func evalIfStatement(
	ie *ast.IfStatement,
	scope *symbol.Scope,
) symbol.Object {
	condition := Eval(ie.Condition, scope)
	if symbol.IsError(condition) {
		return condition
	}

	extendedScope := symbol.NewInnerScope(scope)
	if isTruthy(condition) {
		return Eval(ie.Consequence, extendedScope)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative, extendedScope)
	} else {
		return NULL
	}
}

func evalForLoopStatement(
	ie *ast.ForLoopStatement,
	scope *symbol.Scope,
) symbol.Object {
	condition := Eval(ie.Condition, scope)
	if symbol.IsError(condition) {
		return condition
	}

	extendedScope := symbol.NewInnerScope(scope)
	for isTruthy(condition) {
		result := Eval(ie.Consequence, extendedScope)

		if result != nil {
			rt := result.Type()
			if rt == symbol.RETURN_VALUE_OBJ || rt == symbol.ERROR_OBJ {
				return result
			}
		}

		condition = Eval(ie.Condition, scope)
		if symbol.IsError(condition) {
			return condition
		}
	}

	return NULL
}

func evalIfExpression(
	ie *ast.IfExpression,
	scope *symbol.Scope,
) symbol.Object {
	condition := Eval(ie.Condition, scope)
	if symbol.IsError(condition) {
		return condition
	}

	if isTruthy(condition) {
		return Eval(*ie.Consequence, scope)
	} else if ie.Alternative != nil {
		return Eval(*ie.Alternative, scope)
	} else {
		return NULL
	}
}

func evalDeclareAssignStetment(
	node *ast.DeclareAssignStatement,
	scope *symbol.Scope,
) symbol.Object {
	variable := evalIdentifierInCurrentScope(node.Identifier, scope)
	if !symbol.IsError(variable) {
		return newError("identifier already taken: %s", node.Identifier.Value)
	}

	val := Eval(node.Expression, scope)
	if symbol.IsError(val) {
		return val
	}

	scope.Insert(node.Identifier.Value, val, val.Type())
	return nil
}

func evalDeclareStatement(
	node *ast.DeclareStatement,
	scope *symbol.Scope,
) symbol.Object {
	variable := evalIdentifierInCurrentScope(node.Identifier, scope)
	if !symbol.IsError(variable) {
		return newError("identifier already taken: %s", node.Identifier.Value)
	}

	scope.Insert(node.Identifier.Value, GetDefaultValue(*node.Identifier), common.ToObjectType(node.Identifier.Type))

	if node.Assignment != nil {
		val := evalAssignStatement(node.Assignment, scope)
		if symbol.IsError(val) {
			return val
		}
	}
	return nil
}

func evalAssignStatement(
	node *ast.AssignStatement,
	scope *symbol.Scope,
) symbol.Object {
	variable := evalIdentifier(node.Identifier, scope)
	if symbol.IsError(variable) {
		return variable
	}

	sym, ok := scope.Lookup(node.Identifier.Value)
	identifierType := symbol.NULL_OBJ
	if ok {
		identifierType = sym.Type
	}

	val := Eval(node.Expression, scope)
	if symbol.IsError(val) {
		return val
	}

	if val.Type() != identifierType {
		return newError("type mismatch: %s = %s", identifierType, val.Type())
	}

	scope.TryToAssign(node.Identifier.Value, val, val.Type())

	return nil
}

func evalIdentifier(
	node *ast.Identifier,
	scope *symbol.Scope,
) symbol.Object {
	if val, ok := scope.Lookup(node.Value); ok {
		return val.Object
	}

	if builtin, ok := builtins[node.Value]; ok {
		return builtin
	}

	return newError("identifier not found: %s", node.Value)
}

func evalIdentifierInCurrentScope(
	node *ast.Identifier,
	scope *symbol.Scope,
) symbol.Object {
	if val, ok := scope.LookupInCurrentScope(node.Value); ok {
		return val.Object
	}

	if builtin, ok := builtins[node.Value]; ok {
		return builtin
	}

	return newError("identifier not found: %s", node.Value)
}

func evalFunctionCallExpression(
	node *ast.FunctionCallExpression,
	scope *symbol.Scope,
) symbol.Object {
	function := evalIdentifier(&node.Identifier, scope)
	if symbol.IsError(function) {
		return function
	}

	args := evalExpressions(node.Parameters, scope)
	if len(args) == 1 && symbol.IsError(args[0]) {
		return args[0]
	}

	return applyFunctionOrBuiltin(function, args)
}

func isTruthy(obj symbol.Object) bool {
	switch obj {
	case NULL:
		return false
	case TRUE:
		return true
	case FALSE:
		return false
	default:
		return true
	}
}

func newEvaluatorError(line int, message string) *symbol.Error {
	return &symbol.Error{Message: fmt.Sprintf("Evaluator error on line %v: %s", line, message)}
}

func newError(format string, a ...interface{}) *symbol.Error {
	return &symbol.Error{Message: fmt.Sprintf(format, a...)}
}

func evalExpressions(
	exps []ast.Expression,
	env *symbol.Scope,
) []symbol.Object {
	var result []symbol.Object

	for _, e := range exps {
		evaluated := Eval(e, env)
		if symbol.IsError(evaluated) {
			return []symbol.Object{evaluated}
		}
		result = append(result, evaluated)
	}

	return result
}

func applyFunctionOrBuiltin(fn symbol.Object, args []symbol.Object) symbol.Object {
	switch fn := fn.(type) {
	case *symbol.Function:
		return applyFunction(fn, args)
	case *symbol.Builtin:
		return fn.Fn(args...)
	default:
		return newError("not a function: %s", fn.Type())
	}
}

func applyFunction(fn *symbol.Function, args []symbol.Object) symbol.Object {
	for i := 0; i < len(args); i++ {
		parameterType := fn.Parameters[i].Identifier.Type
		x := args[i].Type()
		_ = x
		argType := common.ToTokenType(args[i].Type())
		if parameterType != argType {
			return newError("function parameter type mismatch: %s, wanted: %s, got: %s", fn.Identifier, parameterType, argType)
		}
	}

	extendedScope := extendFunctionScope(fn, args)
	for i, param := range fn.Parameters {
		arg := args[i]
		extendedScope.Insert(param.Identifier.Value, arg, symbol.ObjectType(param.Identifier.Type))
	}
	evaluated := Eval(fn.Body, extendedScope)
	result := unwrapReturnValue(evaluated)

	if result.Type() != fn.ReturnType {
		return newError("cannot use %s as %s in return statement", result.Type(), fn.ReturnType)
	}

	return result
}

func extendFunctionScope(
	fn *symbol.Function,
	args []symbol.Object,
) *symbol.Scope {
	scope := symbol.NewInnerScope(fn.Scope)

	for _, param := range fn.Parameters {
		evalDeclareStatement(param, scope)
	}

	return scope
}

func unwrapReturnValue(obj symbol.Object) symbol.Object {
	if returnValue, ok := obj.(*symbol.ReturnValue); ok {
		return returnValue.Value
	}

	return obj
}

// func evalIndexExpression(left, index symbol.Object) symbol.Object {
// 	switch {
// 	case left.Type() == symbol.ARRAY_OBJ && index.Type() == symbol.INTEGER_OBJ:
// 		return evalArrayIndexExpression(left, index)
// 	case left.Type() == symbol.HASH_OBJ:
// 		return evalHashIndexExpression(left, index)
// 	default:
// 		return newError("index operator not supported: %s", left.Type())
// 	}
// }

// func evalArrayIndexExpression(array, index symbol.Object) symbol.Object {
// 	arrayObject := array.(*symbol.Array)
// 	idx := index.(*symbol.Integer).Value
// 	max := int64(len(arrayObject.Elements) - 1)

// 	if idx < 0 || idx > max {
// 		return NULL
// 	}

// 	return arrayObject.Elements[idx]
// }

// func evalHashLiteral(
// 	node *ast.HashLiteral,
// 	env *symbol.Scope,
// ) symbol.Object {
// 	pairs := make(map[symbol.HashKey]symbol.HashPair)

// 	for keyNode, valueNode := range node.Pairs {
// 		key := Eval(keyNode, env)
// 		if symbol.IsError(key) {
// 			return key
// 		}

// 		hashKey, ok := key.(symbol.Hashable)
// 		if !ok {
// 			return newError("unusable as hash key: %s", key.Type())
// 		}

// 		value := Eval(valueNode, env)
// 		if symbol.IsError(value) {
// 			return value
// 		}

// 		hashed := hashKey.HashKey()
// 		pairs[hashed] = symbol.HashPair{Key: key, Value: value}
// 	}

// 	return &symbol.Hash{Pairs: pairs}
// }

// func evalHashIndexExpression(hash, index symbol.Object) symbol.Object {
// 	hashObject := hash.(*symbol.Hash)

// 	key, ok := index.(symbol.Hashable)
// 	if !ok {
// 		return newError("unusable as hash key: %s", index.Type())
// 	}

// 	pair, ok := hashsymbol.Pairs[key.HashKey()]
// 	if !ok {
// 		return NULL
// 	}

// 	return pair.Value
// }
