package vm_test

import (
	"testing"

	"github.com/stevecallear/mexl/types"
)

func TestBuiltIn(t *testing.T) {
	tests := []testCase{
		{
			name:  "len args count error",
			input: `len()`,
			err:   true,
		},
		{
			name:  "len args type error",
			input: `len(1)`,
			err:   true,
		},
		{
			name:  "len null",
			input: "len(null)",
			exp:   0,
		},
		{
			name:  "len string",
			input: `len("abc")`,
			exp:   3,
		},
		{
			name:  "len array",
			input: `len([1, 2])`,
			exp:   2,
		},
		{
			name:  "len map",
			input: `len(m)`,
			env: types.Map{
				"m": types.Map{
					"1": &types.Integer{Value: 1},
				},
			},
			exp: 1,
		},
		{
			name:  "lower args count error",
			input: `lower("A", "b")`,
			err:   true,
		},
		{
			name:  "lower arg type error",
			input: `lower(1)`,
			err:   true,
		},
		{
			name:  "lower",
			input: `lower("ABC")`,
			exp:   "abc",
		},
		{
			name:  "lower (null)",
			input: "lower(null)",
			exp:   nil,
		},
		{
			name:  "upper args count error",
			input: `upper("A", "b")`,
			err:   true,
		},
		{
			name:  "upper arg type error",
			input: `upper(1)`,
			err:   true,
		},
		{
			name:  "upper",
			input: `upper("abc")`,
			exp:   "ABC",
		},
		{
			name:  "upper (null)",
			input: "upper(null)",
			exp:   nil,
		},
		{
			name:  "custom",
			input: `reverse("abc")`,
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
