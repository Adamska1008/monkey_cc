package evaluator

import "monkey_cc/object"

var builtins = map[string]*object.BuiltIn{
	"len": {func(args ...object.Object) object.Object {
		if len(args) != 1 {
			return newError("wrong number of arguments, expect: %d, found: %d.", 1, len(args))
		}
		switch arg := args[0].(type) {
		case *object.String:
			return &object.Integer{Value: int64(len(arg.Value))}
		default:
			return newError("argument type %s to `len` is not supported", arg.Type())
		}
	}},
}
