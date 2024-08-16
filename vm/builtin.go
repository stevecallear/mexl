package vm

import (
	"fmt"
	"strings"

	"github.com/stevecallear/mexl/types"
)

var builtIns = map[string]types.Func{
	"len": func(args ...types.Object) (types.Object, error) {
		if err := expectArgsLen("len", args, 1); err != nil {
			return nil, err
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
		if err := expectArgsLen("lower", args, 1); err != nil {
			return nil, err
		}

		switch args[0].Type() {
		case types.TypeNull:
			return objNull, nil

		case types.TypeString:
			return &types.String{
				Value: strings.ToLower(args[0].(*types.String).Value),
			}, nil

		default:
			return nil, fmt.Errorf("lower: wrong arg type: %s, expected %s", args[0].Type(), types.TypeString)
		}
	},

	"upper": func(args ...types.Object) (types.Object, error) {
		if err := expectArgsLen("upper", args, 1); err != nil {
			return nil, err
		}

		switch args[0].Type() {
		case types.TypeNull:
			return objNull, nil

		case types.TypeString:
			return &types.String{
				Value: strings.ToUpper(args[0].(*types.String).Value),
			}, nil

		default:
			return nil, fmt.Errorf("upper: wrong arg type: %s, expected %s", args[0].Type(), types.TypeString)
		}
	},
}

func expectArgsLen(name string, args []types.Object, l int) error {
	if len(args) != l {
		return fmt.Errorf("%s: wrong number of arguments: %d, expected %d", name, len(args), l)
	}
	return nil
}
