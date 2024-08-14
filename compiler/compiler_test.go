package compiler_test

import (
	"testing"

	"github.com/stevecallear/mexl/ast"
	"github.com/stevecallear/mexl/compiler"
	"github.com/stevecallear/mexl/parser"
	"github.com/stevecallear/mexl/types"
	"github.com/stevecallear/mexl/vm"
)

type (
	testCase struct {
		name  string
		input string
		exp   expectation
	}

	expectation struct {
		instructions []vm.Instructions
		constants    []any
		identifiers  []string
	}
)

func TestIntegerArithmetic(t *testing.T) {
	tests := []testCase{
		{
			name:  "addition",
			input: "1 + 2",
			exp: expectation{
				constants: []any{1, 2},
				instructions: []vm.Instructions{
					vm.Make(vm.OpConstant, 0),
					vm.Make(vm.OpConstant, 1),
					vm.Make(vm.OpAdd),
				},
			},
		},
		{
			name:  "subtraction",
			input: "2 - 1",
			exp: expectation{
				constants: []any{2, 1},
				instructions: []vm.Instructions{
					vm.Make(vm.OpConstant, 0),
					vm.Make(vm.OpConstant, 1),
					vm.Make(vm.OpSubtract),
				},
			},
		},
		{
			name:  "multiplication",
			input: "1 * 2",
			exp: expectation{
				constants: []any{1, 2},
				instructions: []vm.Instructions{
					vm.Make(vm.OpConstant, 0),
					vm.Make(vm.OpConstant, 1),
					vm.Make(vm.OpMultiply),
				},
			},
		},
		{
			name:  "division",
			input: "2 / 1",
			exp: expectation{
				constants: []any{2, 1},
				instructions: []vm.Instructions{
					vm.Make(vm.OpConstant, 0),
					vm.Make(vm.OpConstant, 1),
					vm.Make(vm.OpDivide),
				},
			},
		},
		{
			name:  "modulus",
			input: "5 % 2",
			exp: expectation{
				constants: []any{5, 2},
				instructions: []vm.Instructions{
					vm.Make(vm.OpConstant, 0),
					vm.Make(vm.OpConstant, 1),
					vm.Make(vm.OpModulus),
				},
			},
		},
		{
			name:  "minus",
			input: "-1",
			exp: expectation{
				constants: []any{1},
				instructions: []vm.Instructions{
					vm.Make(vm.OpConstant, 0),
					vm.Make(vm.OpMinus),
				},
			},
		},
	}

	testCompiler(t, tests)
}

func TestFloatArithmetic(t *testing.T) {
	tests := []testCase{
		{
			name:  "addition",
			input: "1.1 + 2.2",
			exp: expectation{
				constants: []any{1.1, 2.2},
				instructions: []vm.Instructions{
					vm.Make(vm.OpConstant, 0),
					vm.Make(vm.OpConstant, 1),
					vm.Make(vm.OpAdd),
				},
			},
		},
		{
			name:  "subtraction",
			input: "2.2 - 1.1",
			exp: expectation{
				constants: []any{2.2, 1.1},
				instructions: []vm.Instructions{
					vm.Make(vm.OpConstant, 0),
					vm.Make(vm.OpConstant, 1),
					vm.Make(vm.OpSubtract),
				},
			},
		},
		{
			name:  "multiplication",
			input: "1.1 * 2.2",
			exp: expectation{
				constants: []any{1.1, 2.2},
				instructions: []vm.Instructions{
					vm.Make(vm.OpConstant, 0),
					vm.Make(vm.OpConstant, 1),
					vm.Make(vm.OpMultiply),
				},
			},
		},
		{
			name:  "division",
			input: "2.2 / 1.1",
			exp: expectation{
				constants: []any{2.2, 1.1},
				instructions: []vm.Instructions{
					vm.Make(vm.OpConstant, 0),
					vm.Make(vm.OpConstant, 1),
					vm.Make(vm.OpDivide),
				},
			},
		},
		{
			name:  "minus",
			input: "-1.1",
			exp: expectation{
				constants: []any{1.1},
				instructions: []vm.Instructions{
					vm.Make(vm.OpConstant, 0),
					vm.Make(vm.OpMinus),
				},
			},
		},
	}

	testCompiler(t, tests)
}

func TestBooleanExpressions(t *testing.T) {
	tests := []testCase{
		{
			name:  "true",
			input: "true",
			exp: expectation{
				constants: []any{},
				instructions: []vm.Instructions{
					vm.Make(vm.OpTrue),
				},
			},
		},
		{
			name:  "false",
			input: "false",
			exp: expectation{
				constants: []any{},
				instructions: []vm.Instructions{
					vm.Make(vm.OpFalse),
				},
			},
		},
		{
			name:  "int greater than",
			input: "1 > 2",
			exp: expectation{
				constants: []any{1, 2},
				instructions: []vm.Instructions{
					vm.Make(vm.OpConstant, 0),
					vm.Make(vm.OpConstant, 1),
					vm.Make(vm.OpGreater),
				},
			},
		},
		{
			name:  "int greater than or equal",
			input: "1 >= 2",
			exp: expectation{
				constants: []any{1, 2},
				instructions: []vm.Instructions{
					vm.Make(vm.OpConstant, 0),
					vm.Make(vm.OpConstant, 1),
					vm.Make(vm.OpGreaterOrEqual),
				},
			},
		},
		{
			name:  "int less than",
			input: "1 < 2",
			exp: expectation{
				constants: []any{1, 2},
				instructions: []vm.Instructions{
					vm.Make(vm.OpConstant, 0),
					vm.Make(vm.OpConstant, 1),
					vm.Make(vm.OpLess),
				},
			},
		},
		{
			name:  "int less than or equal",
			input: "1 <= 2",
			exp: expectation{
				constants: []any{1, 2},
				instructions: []vm.Instructions{
					vm.Make(vm.OpConstant, 0),
					vm.Make(vm.OpConstant, 1),
					vm.Make(vm.OpLessOrEqual),
				},
			},
		},
		{
			name:  "float equal",
			input: "1.1 == 2.2",
			exp: expectation{
				constants: []any{1.1, 2.2},
				instructions: []vm.Instructions{
					vm.Make(vm.OpConstant, 0),
					vm.Make(vm.OpConstant, 1),
					vm.Make(vm.OpEqual),
				},
			},
		},
		{
			name:  "float not equal",
			input: "1.1 != 2.2",
			exp: expectation{
				constants: []any{1.1, 2.2},
				instructions: []vm.Instructions{
					vm.Make(vm.OpConstant, 0),
					vm.Make(vm.OpConstant, 1),
					vm.Make(vm.OpNotEqual),
				},
			},
		},
		{
			name:  "float greater than",
			input: "1.1 > 2.2",
			exp: expectation{
				constants: []any{1.1, 2.2},
				instructions: []vm.Instructions{
					vm.Make(vm.OpConstant, 0),
					vm.Make(vm.OpConstant, 1),
					vm.Make(vm.OpGreater),
				},
			},
		},
		{
			name:  "float greater than or equal",
			input: "1.1 >= 2.2",
			exp: expectation{
				constants: []any{1.1, 2.2},
				instructions: []vm.Instructions{
					vm.Make(vm.OpConstant, 0),
					vm.Make(vm.OpConstant, 1),
					vm.Make(vm.OpGreaterOrEqual),
				},
			},
		},
		{
			name:  "float less than",
			input: "1.1 < 2.2",
			exp: expectation{
				constants: []any{1.1, 2.2},
				instructions: []vm.Instructions{
					vm.Make(vm.OpConstant, 0),
					vm.Make(vm.OpConstant, 1),
					vm.Make(vm.OpLess),
				},
			},
		},
		{
			name:  "float less than or equal",
			input: "1.1 <= 2.2",
			exp: expectation{
				constants: []any{1.1, 2.2},
				instructions: []vm.Instructions{
					vm.Make(vm.OpConstant, 0),
					vm.Make(vm.OpConstant, 1),
					vm.Make(vm.OpLessOrEqual),
				},
			},
		},
		{
			name:  "float equal",
			input: "1.1 == 2.2",
			exp: expectation{
				constants: []any{1.1, 2.2},
				instructions: []vm.Instructions{
					vm.Make(vm.OpConstant, 0),
					vm.Make(vm.OpConstant, 1),
					vm.Make(vm.OpEqual),
				},
			},
		},
		{
			name:  "float not equal",
			input: "1.1 != 2.2",
			exp: expectation{
				constants: []any{1.1, 2.2},
				instructions: []vm.Instructions{
					vm.Make(vm.OpConstant, 0),
					vm.Make(vm.OpConstant, 1),
					vm.Make(vm.OpNotEqual),
				},
			},
		},
		{
			name:  "mixed types 1",
			input: "2.2 > 1",
			exp: expectation{
				constants: []any{2.2, 1},
				instructions: []vm.Instructions{
					vm.Make(vm.OpConstant, 0),
					vm.Make(vm.OpConstant, 1),
					vm.Make(vm.OpGreater),
				},
			},
		},
		{
			name:  "mixed types 2",
			input: "1 < 2.2",
			exp: expectation{
				constants: []any{1, 2.2},
				instructions: []vm.Instructions{
					vm.Make(vm.OpConstant, 0),
					vm.Make(vm.OpConstant, 1),
					vm.Make(vm.OpLess),
				},
			},
		},
		{
			name:  "bool equal",
			input: "true == false",
			exp: expectation{
				instructions: []vm.Instructions{
					vm.Make(vm.OpTrue),
					vm.Make(vm.OpFalse),
					vm.Make(vm.OpEqual),
				},
			},
		},
		{
			name:  "bool not equal",
			input: "true != false",
			exp: expectation{
				instructions: []vm.Instructions{
					vm.Make(vm.OpTrue),
					vm.Make(vm.OpFalse),
					vm.Make(vm.OpNotEqual),
				},
			},
		},
		{
			name:  "bang",
			input: "!true",
			exp: expectation{
				instructions: []vm.Instructions{
					vm.Make(vm.OpTrue),
					vm.Make(vm.OpNot),
				},
			},
		},
		{
			name:  "and",
			input: "true && false",
			exp: expectation{
				instructions: []vm.Instructions{
					vm.Make(vm.OpTrue),
					vm.Make(vm.OpJumpIfFalse, 6),
					vm.Make(vm.OpFalse),
					vm.Make(vm.OpAnd),
				},
			},
		},
		{
			name:  "or",
			input: "true || false",
			exp: expectation{
				instructions: []vm.Instructions{
					vm.Make(vm.OpTrue),
					vm.Make(vm.OpJumpIfTrue, 6),
					vm.Make(vm.OpFalse),
					vm.Make(vm.OpOr),
				},
			},
		},
		{
			name:  "start with",
			input: `"abc" sw "a"`,
			exp: expectation{
				constants: []any{"abc", "a"},
				instructions: []vm.Instructions{
					vm.Make(vm.OpConstant, 0),
					vm.Make(vm.OpConstant, 1),
					vm.Make(vm.OpStartsWith),
				},
			},
		},
		{
			name:  "ends with",
			input: `"abc" ew "c"`,
			exp: expectation{
				constants: []any{"abc", "c"},
				instructions: []vm.Instructions{
					vm.Make(vm.OpConstant, 0),
					vm.Make(vm.OpConstant, 1),
					vm.Make(vm.OpEndsWith),
				},
			},
		},
		{
			name:  "in",
			input: `2 in ["a", 2, 3.3]`,
			exp: expectation{
				constants: []any{2, "a", 2, 3.3},
				instructions: []vm.Instructions{
					vm.Make(vm.OpConstant, 0),
					vm.Make(vm.OpConstant, 1),
					vm.Make(vm.OpConstant, 2),
					vm.Make(vm.OpConstant, 3),
					vm.Make(vm.OpArray, 3),
					vm.Make(vm.OpIn),
				},
			},
		},
	}

	testCompiler(t, tests)
}

func TestStringExpressions(t *testing.T) {
	tests := []testCase{
		{
			name:  "concatenation",
			input: `"mess" + "age"`,
			exp: expectation{
				constants: []any{"mess", "age"},
				instructions: []vm.Instructions{
					vm.Make(vm.OpConstant, 0),
					vm.Make(vm.OpConstant, 1),
					vm.Make(vm.OpAdd),
				},
			},
		},
		{
			name:  "constant",
			input: "\"message\"",
			exp: expectation{
				constants: []any{"message"},
				instructions: []vm.Instructions{
					vm.Make(vm.OpConstant, 0),
				},
			},
		},
	}

	testCompiler(t, tests)
}

func TestArrays(t *testing.T) {
	tests := []testCase{
		{
			name:  "array",
			input: `["a", 2, 3.3]`,
			exp: expectation{
				constants: []any{"a", 2, 3.3},
				instructions: []vm.Instructions{
					vm.Make(vm.OpConstant, 0),
					vm.Make(vm.OpConstant, 1),
					vm.Make(vm.OpConstant, 2),
					vm.Make(vm.OpArray, 3),
				},
			},
		},
		{
			name:  "array expression",
			input: `[1, 3 > 2]`,
			exp: expectation{
				constants: []any{1, 3, 2},
				instructions: []vm.Instructions{
					vm.Make(vm.OpConstant, 0),
					vm.Make(vm.OpConstant, 1),
					vm.Make(vm.OpConstant, 2),
					vm.Make(vm.OpGreater),
					vm.Make(vm.OpArray, 2),
				},
			},
		},
		{
			name:  "array index",
			input: "[1, 2, 3][1]",
			exp: expectation{
				constants: []any{1, 2, 3, 1},
				instructions: []vm.Instructions{
					vm.Make(vm.OpConstant, 0),
					vm.Make(vm.OpConstant, 1),
					vm.Make(vm.OpConstant, 2),
					vm.Make(vm.OpArray, 3),
					vm.Make(vm.OpConstant, 3),
					vm.Make(vm.OpIndex),
				},
			},
		},
	}

	testCompiler(t, tests)
}

func TestGlobals(t *testing.T) {
	tests := []testCase{
		{
			name:  "global",
			input: "x",
			exp: expectation{
				identifiers: []string{"x"},
				instructions: []vm.Instructions{
					vm.Make(vm.OpGlobal, 0),
				},
			},
		},
		{
			name:  "global expression",
			input: "x < 2",
			exp: expectation{
				constants:   []any{2},
				identifiers: []string{"x"},
				instructions: []vm.Instructions{
					vm.Make(vm.OpGlobal, 0),
					vm.Make(vm.OpConstant, 0),
					vm.Make(vm.OpLess),
				},
			},
		},
	}

	testCompiler(t, tests)
}

func TestCallExpressions(t *testing.T) {
	tests := []testCase{
		{
			name:  "builtin",
			input: `upper("test")`,
			exp: expectation{
				constants:   []any{"test"},
				identifiers: []string{"upper"},
				instructions: []vm.Instructions{
					vm.Make(vm.OpGlobal, 0),
					vm.Make(vm.OpConstant, 0),
					vm.Make(vm.OpCall, 1),
				},
			},
		},
	}

	testCompiler(t, tests)
}

func TestNull(t *testing.T) {
	tests := []testCase{
		{
			name:  "null",
			input: `email eq null`,
			exp: expectation{
				identifiers: []string{"email"},
				instructions: []vm.Instructions{
					vm.Make(vm.OpGlobal, 0),
					vm.Make(vm.OpNull),
					vm.Make(vm.OpEqual),
				},
			},
		},
	}

	testCompiler(t, tests)
}

func TestShortCircuiting(t *testing.T) {
	tests := []testCase{
		{
			name:  "or jump",
			input: "true or false",
			exp: expectation{
				instructions: []vm.Instructions{
					vm.Make(vm.OpTrue),
					vm.Make(vm.OpJumpIfTrue, 6),
					vm.Make(vm.OpFalse),
					vm.Make(vm.OpOr),
				},
			},
		},
		{
			name:  "and jump",
			input: "false and true",
			exp: expectation{
				instructions: []vm.Instructions{
					vm.Make(vm.OpFalse),
					vm.Make(vm.OpJumpIfFalse, 6),
					vm.Make(vm.OpTrue),
					vm.Make(vm.OpAnd),
				},
			},
		},
	}

	testCompiler(t, tests)
}

func TestMapAccess(t *testing.T) {
	tests := []testCase{
		{
			name:  "root",
			input: "x",
			exp: expectation{
				identifiers: []string{"x"},
				instructions: []vm.Instructions{
					vm.Make(vm.OpGlobal, 0),
				},
			},
		},
		{
			name:  "depth",
			input: "x.y.z",
			exp: expectation{
				identifiers: []string{"x", "y", "z"},
				instructions: []vm.Instructions{
					vm.Make(vm.OpGlobal, 0),
					vm.Make(vm.OpMember, 1),
					vm.Make(vm.OpMember, 2),
				},
			},
		},
	}

	testCompiler(t, tests)
}

func testCompiler(t *testing.T, tests []testCase) {
	t.Helper()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			statement := parse(tt.input)

			c := compiler.New()
			p, err := c.Compile(statement)
			if err != nil {
				t.Fatalf("got %v, expected nil", err)
			}

			assertInstructions(t, p.Instructions, tt.exp.instructions)
			assertConstants(t, p.Constants, tt.exp.constants)
			assertSymbols(t, p.Identifiers, tt.exp.identifiers)
		})
	}
}

func assertInstructions(t *testing.T, act vm.Instructions, exp []vm.Instructions) {
	t.Helper()

	fexp := flattenInstructions(exp)

	if al, el := len(act), len(fexp); al != el {
		t.Errorf("got %d instructions, expected %d\n%s\n%s", al, el, act.String(), fexp.String())
		return
	}

	for i, in := range fexp {
		if act[i] != in {
			t.Errorf("got %q at %d, expected %q\n%s\n%s", act[i], i, in, act, fexp)
		}
	}
}

func assertConstants(t *testing.T, act []types.Object, exp []any) {
	t.Helper()

	if al, el := len(act), len(exp); al != el {
		t.Errorf("got %d constants, expected %d", al, el)
		return
	}

	for i, constant := range exp {
		assertObject(t, act[i], constant)
	}
}

func assertSymbols(t *testing.T, act, exp []string) {
	t.Helper()

	if al, el := len(act), len(exp); al != el {
		t.Errorf("got %d symbols, expected %d", al, el)
		return
	}

	for i, symbol := range exp {
		if act, exp := act[i], symbol; act != exp {
			t.Errorf("got symbol %s, expected %s", act, exp)
		}
	}
}

func assertObject(t *testing.T, act types.Object, exp any) {
	t.Helper()

	switch exp := exp.(type) {
	case int:
		assertIntegerObject(t, act, int64(exp))
	case float64:
		assertFloatObject(t, act, exp)
	case bool:
		assertBooleanObject(t, act, exp)
	case string:
		assertStringObject(t, act, exp)
	default:
		t.Errorf("invalid assertion type: %T", exp)
	}
}

func assertIntegerObject(t *testing.T, act types.Object, exp int64) {
	t.Helper()

	obj, ok := act.(*types.Integer)
	if !ok {
		t.Errorf("got %T, expected types.Integer", act)
		return
	}

	if obj.Value != exp {
		t.Errorf("got %d, expected %d", obj.Value, exp)
	}
}

func assertFloatObject(t *testing.T, act types.Object, exp float64) {
	t.Helper()

	obj, ok := act.(*types.Float)
	if !ok {
		t.Errorf("got %T, expected types.Float", act)
		return
	}

	if obj.Value != exp {
		t.Errorf("got %v, expected %v", obj.Value, exp)
	}
}

func assertBooleanObject(t *testing.T, act types.Object, exp bool) {
	t.Helper()

	obj, ok := act.(*types.Boolean)
	if !ok {
		t.Errorf("got %T, expected types.Boolean", act)
		return
	}

	if obj.Value != exp {
		t.Errorf("got %v, expected %v", obj.Value, exp)
	}
}

func assertStringObject(t *testing.T, act types.Object, exp string) {
	obj, ok := act.(*types.String)
	if !ok {
		t.Errorf("got %T, expected types.String", act)
		return
	}

	if obj.Value != exp {
		t.Errorf("got %s, expected %s", obj.Value, exp)
	}
}

func flattenInstructions(is []vm.Instructions) vm.Instructions {
	out := vm.Instructions{}
	for _, i := range is {
		out = append(out, i...)
	}
	return out
}

func parse(input string) ast.Node {
	n, err := parser.New(input).Parse()
	if err != nil {
		panic(err)
	}
	return n
}
