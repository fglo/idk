package evaluator

import (
	"fmt"

	"github.com/fglo/idk/pkg/idk/symbol"
)

var builtins = map[string]*symbol.Builtin{
	"print": {
		Fn: func(args ...symbol.Object) symbol.Object {
			for _, arg := range args {
				fmt.Print(arg.Inspect())
				fmt.Print(" ")
			}
			fmt.Println()

			return NULL
		},
	},
	"typeof": {
		Fn: func(args ...symbol.Object) symbol.Object {
			if len(args) != 1 {
				return newError("typeof: wrong number of arguments. got=%d, want=1",
					len(args))
			}

			return &symbol.Type{Value: args[0].Type()}
		},
	},
	"int": {
		Fn: func(args ...symbol.Object) symbol.Object {
			if len(args) != 1 {
				return newError("typeof: wrong number of arguments. got=%d, want=1",
					len(args))
			}

			arg := args[0]
			if arg.Type() != symbol.FLOATING_POINT_OBJ {
				return newError("typeof: wrong argument type. got=%s, want=FLOATING_POINT_OBJ",
					arg.Type())
			}

			return &symbol.Integer{Value: int64(arg.(*symbol.FloatingPoint).Value)}
		},
	},
	"float": {
		Fn: func(args ...symbol.Object) symbol.Object {
			if len(args) != 1 {
				return newError("typeof: wrong number of arguments. got=%d, want=1",
					len(args))
			}

			arg := args[0]
			if arg.Type() != symbol.INTEGER_OBJ {
				return newError("typeof: wrong argument type. got=%s, want=INTEGER_OBJ",
					arg.Type())
			}

			return &symbol.FloatingPoint{Value: float64(arg.(*symbol.Integer).Value)}
		},
	},

	// TODO: arrays

	// "len": {Fn: func(args ...symbol.Object) symbol.Object {
	// 	if len(args) != 1 {
	// 		return newError("len: wrong number of arguments. got=%d, want=1",
	// 			len(args))
	// 	}

	// 	switch arg := args[0].(type) {
	// 	case *symbol.Array:
	// 		return &symbol.Integer{Value: int64(len(arg.Elements))}
	// 	case *symbol.String:
	// 		return &symbol.Integer{Value: int64(len(arg.Value))}
	// 	default:
	// 		return newError("argument to `len` not supported, got %s",
	// 			args[0].Type())
	// 	}
	// },
	// },
	// "first": {
	// 	Fn: func(args ...symbol.Object) symbol.Object {
	// 		if len(args) != 1 {
	// 			return newError("first: wrong number of arguments. got=%d, want=1",
	// 				len(args))
	// 		}
	// 		if args[0].Type() != symbol.ARRAY_OBJ {
	// 			return newError("argument to `first` must be ARRAY, got %s",
	// 				args[0].Type())
	// 		}

	// 		arr := args[0].(*symbol.Array)
	// 		if len(arr.Elements) > 0 {
	// 			return arr.Elements[0]
	// 		}

	// 		return NULL
	// 	},
	// },
	// "last": {
	// 	Fn: func(args ...symbol.Object) symbol.Object {
	// 		if len(args) != 1 {
	// 			return newError("last: wrong number of arguments. got=%d, want=1",
	// 				len(args))
	// 		}
	// 		if args[0].Type() != symbol.ARRAY_OBJ {
	// 			return newError("argument to `last` must be ARRAY, got %s",
	// 				args[0].Type())
	// 		}

	// 		arr := args[0].(*symbol.Array)
	// 		length := len(arr.Elements)
	// 		if length > 0 {
	// 			return arr.Elements[length-1]
	// 		}

	// 		return NULL
	// 	},
	// },
	// "rest": {
	// 	Fn: func(args ...symbol.Object) symbol.Object {
	// 		if len(args) != 1 {
	// 			return newError("rest: wrong number of arguments. got=%d, want=1",
	// 				len(args))
	// 		}
	// 		if args[0].Type() != symbol.ARRAY_OBJ {
	// 			return newError("argument to `rest` must be ARRAY, got %s",
	// 				args[0].Type())
	// 		}

	// 		arr := args[0].(*symbol.Array)
	// 		length := len(arr.Elements)
	// 		if length > 0 {
	// 			newElements := make([]symbol.Object, length-1, length-1)
	// 			copy(newElements, arr.Elements[1:length])
	// 			return &symbol.Array{Elements: newElements}
	// 		}

	// 		return NULL
	// 	},
	// },
	// "push": {
	// 	Fn: func(args ...symbol.Object) symbol.Object {
	// 		if len(args) != 2 {
	// 			return newError("push: wrong number of arguments. got=%d, want=2",
	// 				len(args))
	// 		}
	// 		if args[0].Type() != symbol.ARRAY_OBJ {
	// 			return newError("argument to `push` must be ARRAY, got %s",
	// 				args[0].Type())
	// 		}

	// 		arr := args[0].(*symbol.Array)
	// 		length := len(arr.Elements)

	// 		newElements := make([]symbol.Object, length+1, length+1)
	// 		copy(newElements, arr.Elements)
	// 		newElements[length] = args[1]

	// 		return &symbol.Array{Elements: newElements}
	// 	},
	// },
}
