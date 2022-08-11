package evaluator

import (
	"fmt"
	"github.com/fglo/idk/pkg/idk/ast"
	"github.com/fglo/idk/pkg/idk/symbol"
)

var (
	NULL  = &symbol.Null{}
	TRUE  = &symbol.Boolean{Value: true}
	FALSE = &symbol.Boolean{Value: false}
)

func Eval(node ast.Node, scope *symbol.Scope) symbol.Object {
	switch node := node.(type) {

	// Statements
	case *ast.Program:
		return evalProgram(node, scope)

	case *ast.BlockStatement:
		return evalBlockStatement(node, scope)

	case *ast.ExpressionStatement:
		return Eval(node.Expression, scope)

	// case *ast.ReturnStatement:
	// 	val := Eval(node.ReturnValue, env)
	// 	if isError(val) {
	// 		return val
	// 	}
	// 	return &symbol.ReturnValue{Value: val}

	case *ast.DeclareAssignStatement:
		val := Eval(node.Expression, scope)
		if isError(val) {
			return val
		}
		scope.Insert(node.Identifier.Value, val)

	// Expressions
	case *ast.IntegerLiteral:
		return &symbol.Integer{Value: int64(node.Value)}

	// case *ast.StringLiteral:
	// 	return &symbol.String{Value: node.Value}

	// case *ast.Boolean:
	// 	return nativeBoolToBooleanObject(node.Value)

	// case *ast.PrefixExpression:
	// 	right := Eval(node.Right, env)
	// 	if isError(right) {
	// 		return right
	// 	}
	// 	return evalPrefixExpression(node.Operator, right)

	case *ast.BinaryExpression:
		left := Eval(node.Left, scope)
		if isError(left) {
			return left
		}

		right := Eval(node.Right, scope)
		if isError(right) {
			return right
		}

		return evalInfixExpression(node.Token.Value, left, right)

	case *ast.IfStatement:
		return evalIfStatement(node, scope)

	case *ast.IfExpression:
		return evalIfExpression(node, scope)

	case *ast.Identifier:
		return evalIdentifier(node, scope)

	// case *ast.FunctionLiteral:
	// 	params := node.Parameters
	// 	body := node.Body
	// 	return &symbol.Function{Parameters: params, Env: env, Body: body}

	case *ast.FunctionCallExpression:
		function := evalIdentifier(&node.Identifier, scope)
		if isError(function) {
			return function
		}

		args := evalExpressions(node.Parameters, scope)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}

		return applyFunction(function, args)

		// case *ast.ArrayLiteral:
		// 	elements := evalExpressions(node.Elements, env)
		// 	if len(elements) == 1 && isError(elements[0]) {
		// 		return elements[0]
		// 	}
		// 	return &symbol.Array{Elements: elements}

		// case *ast.IndexExpression:
		// 	left := Eval(node.Left, env)
		// 	if isError(left) {
		// 		return left
		// 	}
		// 	index := Eval(node.Index, env)
		// 	if isError(index) {
		// 		return index
		// 	}
		// 	return evalIndexExpression(left, index)

		// case *ast.HashLiteral:
		// 	return evalHashLiteral(node, env)

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

func evalInfixExpression(
	operator string,
	left, right symbol.Object,
) symbol.Object {
	switch {
	case left.Type() == symbol.INTEGER_OBJ && right.Type() == symbol.INTEGER_OBJ:
		return evalIntegerInfixExpression(operator, left, right)
	case left.Type() == symbol.STRING_OBJ && right.Type() == symbol.STRING_OBJ:
		return evalStringInfixExpression(operator, left, right)
	case operator == "==":
		return nativeBoolToBooleanObject(left == right)
	case operator == "!=":
		return nativeBoolToBooleanObject(left != right)
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
	if right.Type() != symbol.INTEGER_OBJ {
		return newError("unknown operator: -%s", right.Type())
	}

	value := right.(*symbol.Integer).Value
	return &symbol.Integer{Value: -value}
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
	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
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
	if operator != "+" {
		return newError("unknown operator: %s %s %s",
			left.Type(), operator, right.Type())
	}

	leftVal := left.(*symbol.String).Value
	rightVal := right.(*symbol.String).Value
	return &symbol.String{Value: leftVal + rightVal}
}

func evalIfStatement(
	ie *ast.IfStatement,
	env *symbol.Scope,
) symbol.Object {
	condition := Eval(ie.Condition, env)
	if isError(condition) {
		return condition
	}

	if isTruthy(condition) {
		return Eval(ie.Consequence, env)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative, env)
	} else {
		return NULL
	}
}

func evalIfExpression(
	ie *ast.IfExpression,
	env *symbol.Scope,
) symbol.Object {
	condition := Eval(ie.Condition, env)
	if isError(condition) {
		return condition
	}

	if isTruthy(condition) {
		return Eval(*ie.Consequence, env)
	} else if ie.Alternative != nil {
		return Eval(*ie.Alternative, env)
	} else {
		return NULL
	}
}

func evalIdentifier(
	node *ast.Identifier,
	env *symbol.Scope,
) symbol.Object {
	if val, ok := env.Lookup(node.Value); ok {
		return val
	}

	if builtin, ok := builtins[node.Value]; ok {
		return builtin
	}

	return newError("identifier not found: " + node.Value)
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

func newError(format string, a ...interface{}) *symbol.Error {
	return &symbol.Error{Message: fmt.Sprintf(format, a...)}
}

func isError(obj symbol.Object) bool {
	if obj != nil {
		return obj.Type() == symbol.ERROR_OBJ
	}
	return false
}

func evalExpressions(
	exps []ast.Expression,
	env *symbol.Scope,
) []symbol.Object {
	var result []symbol.Object

	for _, e := range exps {
		evaluated := Eval(e, env)
		if isError(evaluated) {
			return []symbol.Object{evaluated}
		}
		result = append(result, evaluated)
	}

	return result
}

func applyFunction(fn symbol.Object, args []symbol.Object) symbol.Object {
	switch fn := fn.(type) {

	case *symbol.Function:
		extendedEnv := extendFunctionEnv(fn, args)
		evaluated := Eval(fn.Body, extendedEnv)
		return unwrapReturnValue(evaluated)

	case *symbol.Builtin:
		return fn.Fn(args...)

	default:
		return newError("not a function: %s", fn.Type())
	}
}

func extendFunctionEnv(
	fn *symbol.Function,
	args []symbol.Object,
) *symbol.Scope {
	env := symbol.NewInnerScope(fn.Env)

	for paramIdx, param := range fn.Parameters {
		env.Insert(param.Value, args[paramIdx])
	}

	return env
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
// 		if isError(key) {
// 			return key
// 		}

// 		hashKey, ok := key.(symbol.Hashable)
// 		if !ok {
// 			return newError("unusable as hash key: %s", key.Type())
// 		}

// 		value := Eval(valueNode, env)
// 		if isError(value) {
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
