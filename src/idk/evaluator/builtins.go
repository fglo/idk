package evaluator

import (
	"fmt"
	"idk/symbol"
)

var builtins = map[string]*symbol.Builtin{
	"len": &symbol.Builtin{Fn: func(args ...symbol.Object) symbol.Object {
		if len(args) != 1 {
			return newError("wrong number of arguments. got=%d, want=1",
				len(args))
		}

		switch arg := args[0].(type) {
		case *symbol.Array:
			return &symbol.Integer{Value: int64(len(arg.Elements))}
		case *symbol.String:
			return &symbol.Integer{Value: int64(len(arg.Value))}
		default:
			return newError("argument to `len` not supported, got %s",
				args[0].Type())
		}
	},
	},
	"puts": &symbol.Builtin{
		Fn: func(args ...symbol.Object) symbol.Object {
			for _, arg := range args {
				fmt.Println(arg.Inspect())
			}

			return NULL
		},
	},
	"first": &symbol.Builtin{
		Fn: func(args ...symbol.Object) symbol.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1",
					len(args))
			}
			if args[0].Type() != symbol.ARRAY_OBJ {
				return newError("argument to `first` must be ARRAY, got %s",
					args[0].Type())
			}

			arr := args[0].(*symbol.Array)
			if len(arr.Elements) > 0 {
				return arr.Elements[0]
			}

			return NULL
		},
	},
	"last": &symbol.Builtin{
		Fn: func(args ...symbol.Object) symbol.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1",
					len(args))
			}
			if args[0].Type() != symbol.ARRAY_OBJ {
				return newError("argument to `last` must be ARRAY, got %s",
					args[0].Type())
			}

			arr := args[0].(*symbol.Array)
			length := len(arr.Elements)
			if length > 0 {
				return arr.Elements[length-1]
			}

			return NULL
		},
	},
	"rest": &symbol.Builtin{
		Fn: func(args ...symbol.Object) symbol.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1",
					len(args))
			}
			if args[0].Type() != symbol.ARRAY_OBJ {
				return newError("argument to `rest` must be ARRAY, got %s",
					args[0].Type())
			}

			arr := args[0].(*symbol.Array)
			length := len(arr.Elements)
			if length > 0 {
				newElements := make([]symbol.Object, length-1, length-1)
				copy(newElements, arr.Elements[1:length])
				return &symbol.Array{Elements: newElements}
			}

			return NULL
		},
	},
	"push": &symbol.Builtin{
		Fn: func(args ...symbol.Object) symbol.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments. got=%d, want=2",
					len(args))
			}
			if args[0].Type() != symbol.ARRAY_OBJ {
				return newError("argument to `push` must be ARRAY, got %s",
					args[0].Type())
			}

			arr := args[0].(*symbol.Array)
			length := len(arr.Elements)

			newElements := make([]symbol.Object, length+1, length+1)
			copy(newElements, arr.Elements)
			newElements[length] = args[1]

			return &symbol.Array{Elements: newElements}
		},
	},
}
