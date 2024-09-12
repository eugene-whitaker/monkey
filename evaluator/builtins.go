package evaluator

import (
	"fmt"
	"strings"

	"github.com/eugene-whitaker/monkey/object"
)

var builtins = map[string]*object.Builtin{
	"len": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				types := []string{}
				for _, arg := range args {
					types = append(types, string(arg.Type()))
				}

				return toErrorObject(
					"invalid argument count in call to `len`: found (%s) want (STRING) or (ARRAY)",
					strings.Join(types, ", "),
				)
			}

			switch arg := args[0].(type) {
			case *object.String:
				return toIntegerObject(int64(len(arg.Value)))
			case *object.Array:
				return toIntegerObject(int64(len(arg.Elements)))
			default:
				return toErrorObject(
					"invalid argument types in call to `len`: found (%s) want (STRING) or (ARRAY)",
					arg.Type(),
				)
			}
		},
	},
	"first": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				types := []string{}
				for _, arg := range args {
					types = append(types, string(arg.Type()))
				}

				return toErrorObject(
					"invalid argument count in call to `first`: found (%s) want (ARRAY)",
					strings.Join(types, ", "),
				)
			}

			switch arg := args[0].(type) {
			case *object.Array:
				if len(arg.Elements) > 0 {
					return arg.Elements[0]
				}
				return NULL
			default:
				return toErrorObject(
					"invalid argument types in call to `first`: found (%s) want (ARRAY)",
					arg.Type(),
				)
			}
		},
	},
	"last": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				types := []string{}
				for _, arg := range args {
					types = append(types, string(arg.Type()))
				}

				return toErrorObject(
					"invalid argument count in call to `last`: found (%s) want (ARRAY)",
					strings.Join(types, ", "),
				)
			}

			switch arg := args[0].(type) {
			case *object.Array:
				if len(arg.Elements) > 0 {
					return arg.Elements[len(arg.Elements)-1]
				}
				return NULL
			default:
				return toErrorObject(
					"invalid argument types in call to `last`: found (%s) want (ARRAY)",
					arg.Type(),
				)
			}
		},
	},
	"rest": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				types := []string{}
				for _, arg := range args {
					types = append(types, string(arg.Type()))
				}

				return toErrorObject(
					"invalid argument count in call to `rest`: found (%s) want (ARRAY)",
					strings.Join(types, ", "),
				)
			}

			switch arg := args[0].(type) {
			case *object.Array:
				length := len(arg.Elements)
				if length > 0 {
					elements := make([]object.Object, length-1)
					copy(elements, arg.Elements[1:])
					return toArrayObject(elements)
				}
				return NULL
			default:
				return toErrorObject(
					"invalid argument types in call to `rest`: found (%s) want (ARRAY)",
					arg.Type(),
				)
			}
		},
	},
	"push": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				types := []string{}
				for _, arg := range args {
					types = append(types, string(arg.Type()))
				}

				return toErrorObject(
					"invalid argument count in call to `push`: found (%s) want (ARRAY, ANY)",
					strings.Join(types, ", "),
				)
			}

			switch arg := args[0].(type) {
			case *object.Array:
				length := len(arg.Elements)
				elements := make([]object.Object, length+1)
				copy(elements, arg.Elements)
				elements[length] = args[1]
				return toArrayObject(elements)
			default:
				types := []string{}
				for _, arg := range args {
					types = append(types, string(arg.Type()))
				}

				return toErrorObject(
					"invalid argument types in call to `push`: found (%s) want (ARRAY, ANY)",
					strings.Join(types, ", "),
				)
			}
		},
	},
	"puts": {
		Fn: func(args ...object.Object) object.Object {
			for _, arg := range args {
				fmt.Println(arg.Inspect())
			}

			return NULL
		},
	},
}
