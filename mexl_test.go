package mexl_test

import (
	"fmt"
	"log"
	"reflect"
	"testing"

	"github.com/stevecallear/mexl"
	"github.com/stevecallear/mexl/types"
)

func ExampleEval() {
	const input = `lower(user.email) ew "@email.com" or "beta" in user.roles`

	env := map[string]any{
		"user": map[string]any{
			"email": "Test@Email.com",
			"roles": []any{"admin", "beta"},
		},
	}

	out, err := mexl.Eval(input, env)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(out)
	// Output: true
}

func ExampleRun() {
	const input = `lower(user.email) ew "@email.com" or "beta" in user.roles`

	program, err := mexl.Compile(input)
	if err != nil {
		log.Fatal(err)
	}

	env := map[string]any{
		"user": map[string]any{
			"email": "Test@Email.com",
			"roles": []any{"admin", "beta"},
		},
	}

	out, err := mexl.Run(program, env)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(out)
	// Output: true
}

func TestEval(t *testing.T) {
	tests := []struct {
		name  string
		input string
		env   map[string]any
		exp   any
		err   bool
	}{
		{
			name:  "invalid input",
			input: "£",
			err:   true,
		},
		{
			name:  "invalid env",
			input: "x",
			env:   map[string]any{"x": struct{}{}},
			err:   true,
		},
		{
			name:  "vm error",
			input: `1 lt "a"`,
			err:   true,
		},
		{
			name:  "bool",
			input: "1 lt 2",
			exp:   true,
		},
		{
			name:  "int",
			input: "1 + 2",
			exp:   int64(3),
		},
		{
			name:  "float",
			input: "2.2 / 1.1",
			exp:   float64(2),
		},
		{
			name:  "string",
			input: `lower("ABC") + "def"`,
			exp:   "abcdef",
		},
		{
			name:  "env",
			input: `lower(user.email) ew "@email.com"`,
			env: map[string]any{
				"user": map[string]any{
					"email": "Test@Email.com",
				},
			},
			exp: true,
		},
		{
			name:  "func",
			input: `reverse("abc")`,
			env: map[string]any{
				"reverse": func(args ...types.Object) (types.Object, error) {
					runes := []rune(args[0].(*types.String).Value)
					for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
						runes[i], runes[j] = runes[j], runes[i]
					}
					return &types.String{Value: string(runes)}, nil
				},
			},
			exp: "cba",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			act, err := mexl.Eval(tt.input, tt.env)
			if err != nil && !tt.err {
				t.Errorf("got %v, expected nil", err)
			}
			if err == nil && tt.err {
				t.Error("got nil, expected error")
			}
			if !reflect.DeepEqual(act, tt.exp) {
				t.Errorf("got %v, expected %v", act, tt.exp)
			}
		})
	}
}
