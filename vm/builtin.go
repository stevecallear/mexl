package vm

import (
	"fmt"
	"strings"

	"github.com/stevecallear/mexl/types"
)

var builtIns = map[string]types.Func{
	"len": func(args ...types.Object) (types.Object, error) {
		if len(args) != 1 {
			return nil, fmt.Errorf("wrong number of arguments: %d, expected 1", len(args))
		}

		var l int
		switch args[0].Type() {
		case types.TypeNull:
			l = 0
		case types.TypeString:
			l = len(args[0].(*types.String).Value)
		case types.TypeArray:
			l = len(args[0].(types.Array))
		case types.TypeMap:
			l = len(args[0].(types.Map))
		default:
			return nil, fmt.Errorf("invalid argument: %T", args[0])
		}

		return &types.Integer{Value: int64(l)}, nil
	},

	"lower": func(args ...types.Object) (types.Object, error) {
		if err := expectArgs(args, types.TypeString); err != nil {
			return nil, err
		}

		return &types.String{
			Value: strings.ToLower(args[0].(*types.String).Value),
		}, nil
	},

	"upper": func(args ...types.Object) (types.Object, error) {
		if err := expectArgs(args, types.TypeString); err != nil {
			return nil, err
		}

		return &types.String{
			Value: strings.ToUpper(args[0].(*types.String).Value),
		}, nil
	},
}

func expectArgs(args []types.Object, types ...types.Type) error {
	if len(args) != len(types) {
		return fmt.Errorf("wrong number of arguments: %d, expected %d", len(args), len(types))
	}
	for i, v := range types {
		if args[i].Type() != v {
			return fmt.Errorf("invalid type at %d: %s, expected %s", i, args[i].Type(), v)
		}
	}
	return nil
}
