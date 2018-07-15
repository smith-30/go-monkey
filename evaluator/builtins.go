package evaluator

import (
	"fmt"

	"github.com/smith-30/go-monkey/object"
)

//
// monkey's array is static so push and rest do copy original array and return
//
var builtins = map[string]*object.Builtin{
	"len": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}

			switch arg := args[0].(type) {
			case *object.Array:
				return &object.Integer{Value: int64(len(arg.Elements))}
			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value))}
			default:
				return newError("argument to `len` not supported, got %s", args[0].Type())
			}
		},
	},
	"first": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}

			if args[0].Type() != object.ARRAY_OBJ {
				return newError("argument to `first` must be Array, got=%s", args[0].Type())
			}

			arr := args[0].(*object.Array)
			if len(arr.Elements) > 0 {
				return arr.Elements[0]
			}

			return NULL
		},
	},
	"last": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}

			if args[0].Type() != object.ARRAY_OBJ {
				return newError("argument to `last` must be Array, got=%s", args[0].Type())
			}

			arr := args[0].(*object.Array)
			l := len(arr.Elements)
			if l > 0 {
				return arr.Elements[l-1]
			}

			return NULL
		},
	},
	"rest": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}

			if args[0].Type() != object.ARRAY_OBJ {
				return newError("argument to `rest` must be Array, got=%s", args[0].Type())
			}

			arr := args[0].(*object.Array)
			l := len(arr.Elements)
			if l > 0 {
				// deep copy original elements
				// for nested call rest
				// Elements may be used after rest operation
				newElems := make([]object.Object, l-1, l-1)
				copy(newElems, arr.Elements[1:l])
				return &object.Array{Elements: newElems}
			}

			return NULL
		},
	},
	"push": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments. got=%d, want=2", len(args))
			}

			if args[0].Type() != object.ARRAY_OBJ {
				return newError("argument to `rest` must be Array, got=%s", args[0].Type())
			}

			arr := args[0].(*object.Array)
			l := len(arr.Elements)
			newElems := make([]object.Object, l+1, l+1)
			copy(newElems, arr.Elements)
			newElems[l] = args[1]

			return &object.Array{Elements: newElems}
		},
	},
	"puts": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			for _, v := range args {
				fmt.Println(v.Inspect())
			}
			return NULL
		},
	},
}
