package vm_test

import (
	"testing"

	"github.com/stevecallear/mexl/types"
)

func TestBuiltIn(t *testing.T) {
	tests := []testCase{
		{
			name: "len args count error",
			prog: compile(`len()`),
			err:  true,
		},
		{
			name: "len args type error",
			prog: compile(`len(1)`),
			err:  true,
		},
		{
			name: "len null",
			prog: compile("len(null)"),
			exp:  0,
		},
		{
			name: "len string",
			prog: compile(`len("abc")`),
			exp:  3,
		},
		{
			name: "len array",
			prog: compile(`len([1, 2])`),
			exp:  2,
		},
		{
			name: "len map",
			prog: compile(`len(m)`),
			env: types.Map{
				"m": types.Map{
					"1": &types.Integer{Value: 1},
				},
			},
			exp: 1,
		},
		{
			name: "lower args count error",
			prog: compile(`lower("A", "b")`),
			err:  true,
		},
		{
			name: "lower arg type error",
			prog: compile(`lower(1)`),
			err:  true,
		},
		{
			name: "lower null",
			prog: compile("lower(null)"),
			exp:  nil,
		},
		{
			name: "lower",
			prog: compile(`lower("ABC")`),
			exp:  "abc",
		},
		{
			name: "upper args count error",
			prog: compile(`upper("A", "b")`),
			err:  true,
		},
		{
			name: "upper arg type error",
			prog: compile(`upper(1)`),
			err:  true,
		},
		{
			name: "upper",
			prog: compile(`upper("abc")`),
			exp:  "ABC",
		},
		{
			name: "upper null",
			prog: compile("upper(null)"),
			exp:  nil,
		},
		{
			name: "custom",
			prog: compile(`reverse("abc")`),
			env: map[string]types.Object{
				"reverse": types.Func(func(args ...types.Object) (types.Object, error) {
					runes := []rune(args[0].(*types.String).Value)
					for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
						runes[i], runes[j] = runes[j], runes[i]
					}
					return &types.String{Value: string(runes)}, nil
				}),
			},
			exp: "cba",
		},
	}

	testVM(t, tests)
}
