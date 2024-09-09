package evaluator

import (
	"monkey/object"
	"strings"
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
					"invalid argument count in call to `len`: found (%s) want (STRING)",
					strings.Join(types, ", "),
				)
			}

			switch arg := args[0].(type) {
			case *object.String:
				return &object.Integer{
					Value: int64(len(arg.Value)),
				}
			default:
				return toErrorObject(
					"invalid argument types in call to `len`: found (%s) want (STRING)", arg.Type(),
				)
			}
		},
	},
}
