package vm_test

import (
	"strconv"
	"testing"

	"github.com/stevecallear/mexl/compiler"
	"github.com/stevecallear/mexl/parser"
	"github.com/stevecallear/mexl/types"
	"github.com/stevecallear/mexl/vm"
)

type testCase struct {
	name string
	prog *vm.Program
	env  types.Map
	exp  any
	err  bool
}

func newTestCase(input string, exp any) testCase {
	return testCase{
		prog: compile(input),
		exp:  exp,
	}
}

func TestIntegerArithmetic(t *testing.T) {
	tests := []testCase{
		newTestCase("1", 1),
		newTestCase("2", 2),
		newTestCase("1 + 2", 3),
		newTestCase("1 - 2", -1),
		newTestCase("1 * 2", 2),
		newTestCase("4 / 2", 2),
		newTestCase("3 / 2", 1.5),
		newTestCase("4 % 2", 0),
		newTestCase("5 % 2", 1),
		newTestCase("50 / 2 * 2 + 10 - 5", 55),
		newTestCase("5 + 5 + 5 + 5 - 10", 10),
		newTestCase("2 * 2 * 2 * 2 * 2", 32),
		newTestCase("5 * 2 + 10", 20),
		newTestCase("5 + 2 * 10", 25),
		newTestCase("5 * (2 + 10)", 60),
		newTestCase("-5", -5),
		newTestCase("-10", -10),
		newTestCase("-50 + 100 + -50", 0),
		newTestCase("(5 + 10 * 2 + 15 / 3) * 2 + -10", 50),
	}

	testVM(t, tests)
}

func TestFloatArithmetic(t *testing.T) {
	tests := []testCase{
		newTestCase("1.1", 1.1),
		newTestCase("2.2", 2.2),
		newTestCase("1.0 + 2.2", 3.2),
		newTestCase("1.1 - 2.2", -1.1),
		newTestCase("1.0 * 2.2", 2.2),
		newTestCase("3.0 / 1.5", 2.0),
		newTestCase("-1.1", -1.1),
	}

	testVM(t, tests)
}

func TestArithmeticConversion(t *testing.T) {
	tests := []testCase{
		newTestCase("1.1 + 2", 3.1),
		newTestCase("3 - 1.5", 1.5),
		newTestCase("1.1 * 2", 2.2),
		newTestCase("3 / 1.5", 2.0),
	}

	testVM(t, tests)
}

func TestBooleanExpressions(t *testing.T) {
	tests := []testCase{
		newTestCase("true", true),
		newTestCase("false", false),
		newTestCase("1 < 2", true),
		newTestCase("1 > 2", false),
		newTestCase("1 < 1", false),
		newTestCase("1 > 1", false),
		newTestCase("1 <= 2", true),
		newTestCase("1 >= 2", false),
		newTestCase("1 <= 1", true),
		newTestCase("1 >= 1", true),
		newTestCase("1 == 1", true),
		newTestCase("1 != 1", false),
		newTestCase("1 == 2", false),
		newTestCase("1 != 2", true),
		newTestCase("1.1 < 2.2", true),
		newTestCase("1.1 > 2.2", false),
		newTestCase("1.1 <= 2.2", true),
		newTestCase("1.1 >= 2.2", false),
		newTestCase("1.1 <= 1.1", true),
		newTestCase("1.1 >= 1.1", true),
		newTestCase("1.1 == 1.1", true),
		newTestCase("1.1 != 1.1", false),
		newTestCase("1.1 == 2.2", false),
		newTestCase("1.1 != 2.2", true),
		newTestCase("1 == 1.0", true),
		newTestCase("2.2 == 2", false),
		newTestCase("1 <= 2.2", true),
		newTestCase("1.1 >= 2", false),
		newTestCase("1 <= 1.0", true),
		newTestCase("1.0 >= 1", true),
		newTestCase("true == true", true),
		newTestCase("false == false", true),
		newTestCase("true == false", false),
		newTestCase("true != false", true),
		newTestCase("false != true", true),
		newTestCase("(1 < 2) == true", true),
		newTestCase("(1 < 2) == false", false),
		newTestCase("(1 > 2) == true", false),
		newTestCase("(1 > 2) == false", true),
		newTestCase("!true", false),
		newTestCase("!false", true),
		newTestCase("!5", false),
		newTestCase("!!true", true),
		newTestCase("!!false", false),
		newTestCase("!!5", true),
		newTestCase("!5.5", false),
		newTestCase("!!5.5", true),
		newTestCase("true && false", false),
		newTestCase("true && true", true),
		newTestCase("true || false", true),
		newTestCase("1 > 2 || 2 < 1", false),
		newTestCase(`"abc" eq "a"`, false),
		newTestCase(`"abc" eq "abc"`, true),
		newTestCase(`"abc" ne "a"`, true),
		newTestCase(`"abc" ne "abc"`, false),
		newTestCase(`"abc" sw "a"`, true),
		newTestCase(`"abc" sw "c"`, false),
		newTestCase(`"abc" ew "a"`, false),
		newTestCase(`"abc" ew "c"`, true),
		newTestCase(`1 in ["a", 2, 3.3]`, false),
		newTestCase(`2 in ["a", 2, 3.3]`, true),
		newTestCase(`"a" in "bc"`, false),
		newTestCase(`"b" in "abc"`, true),
		newTestCase(`email eq null`, true),
		newTestCase(`email ne null`, false),
		newTestCase(`1 eq null`, false),
		newTestCase(`1 ne null`, true),
		newTestCase(`false and x eq 1`, false),
		newTestCase(`true or x eq 1`, true),
	}

	testVM(t, tests)
}

func TestStringExpressions(t *testing.T) {
	tests := []testCase{
		{
			name: "constant",
			prog: compile(`"abc"`),
			exp:  "abc",
		},
		{
			name: "concatenation",
			prog: compile(`"abc" + "def"`),
			exp:  "abcdef",
		},
	}

	testVM(t, tests)
}

func TestArrays(t *testing.T) {
	tests := []testCase{
		{
			name: "constant",
			prog: compile(`[1, "b", 3.3]`),
			exp:  []any{1, "b", 3.3},
		},
		{
			name: "index operator",
			prog: compile(`[1, 2, 3][1]`),
			exp:  2,
		},
	}

	testVM(t, tests)
}

func TestMaps(t *testing.T) {
	tests := []testCase{
		{
			name: "root",
			prog: compile("x"),
			env: types.Map{
				"x": &types.Integer{Value: 1},
			},
			exp: 1,
		},
		{
			name: "member",
			prog: compile("x.y.z"),
			env: types.Map{
				"x": types.Map{
					"y": types.Map{
						"z": &types.Integer{Value: 1},
					},
				},
			},
			exp: 1,
		},
		{
			name: "invalid root",
			prog: compile("invalid.x"),
			exp:  nil,
		},
		{
			name: "invalid member",
			prog: compile("x.invalid"),
			env: types.Map{
				"x": types.Map{
					"y": &types.Integer{Value: 1},
				},
			},
			exp: nil,
		},
	}

	testVM(t, tests)
}

func TestGlobals(t *testing.T) {
	tests := []testCase{
		{
			name: "constant",
			prog: compile("x"),
			env: types.Map{
				"x": &types.Boolean{Value: true},
			},
			exp: true,
		},
		{
			name: "boolean expression",
			prog: compile("x < 4"),
			env: types.Map{
				"x": &types.Integer{Value: 5},
			},
			exp: false,
		},
		{
			name: "arithmetic operation",
			prog: compile("x + y"),
			env: types.Map{
				"x": &types.Integer{Value: 5},
				"y": &types.Integer{Value: 3},
			},
			exp: 8,
		},
	}

	testVM(t, tests)
}

func TestNull(t *testing.T) {
	tests := []testCase{
		newTestCase("null + 1", 1),
		newTestCase("null - 1", -1),
		newTestCase("0.5 * null", 0.0),
		newTestCase("null * 1.2", 0.0),
		newTestCase("null lt 1 ", true),
		newTestCase("null gt 0.5", false),
		newTestCase(`null ew "abc"`, false),
		newTestCase("-null", nil),
		newTestCase("null in [1]", false),
		newTestCase("1 in null", false),
	}

	testVM(t, tests)
}

func TestErrors(t *testing.T) {
	tests := []testCase{
		{
			name: "invalid binary op types",
			prog: compile(`"a" + 1`),
			err:  true,
		},
		{
			name: "invalid float binary op type",
			prog: compile("1.1 % 2.2"),
			err:  true,
		},
		{
			name: "invalid string binary op type",
			prog: compile(`"a" - "b"`),
			err:  true,
		},
		{
			name: "invalid comparison op type",
			prog: compile("[1] le [2]"),
			err:  true,
		},
		{
			name: "invalid int comparison op",
			prog: compile("1 sw 2"),
			err:  true,
		},
		{
			name: "invalid float comparison op",
			prog: compile("1.1 sw 2.2"),
			err:  true,
		},
		{
			name: "invalid bool comparison op",
			prog: compile("true ge false"),
			err:  true,
		},
		{
			name: "invalid string comparison op",
			prog: compile(`"a" ge "b"`),
			err:  true,
		},
		{
			name: "invalid negation type",
			prog: compile(`-"a"`),
			err:  true,
		},
		{
			name: "invalid in op type",
			prog: compile("1 in 2"),
			err:  true,
		},
		{
			name: "invalid index op type",
			prog: compile(`"a"[0]`),
			err:  true,
		},
		{
			name: "invalid member op type",
			prog: compile(`"a".x`),
			err:  true,
		},
		{
			name: "invalid call op type",
			prog: compile(`"a"(1)`),
			err:  true,
		},
	}

	testVM(t, tests)
}

func testVM(t *testing.T, tests []testCase) {
	t.Helper()

	for i, tt := range tests {
		n := tt.name
		if n == "" {
			n = strconv.Itoa(i)
		}

		t.Run(n, func(t *testing.T) {
			out, err := vm.New(tt.prog, tt.env).Run()
			if err != nil && !tt.err {
				t.Fatalf("got %v, expected nil", err)
			}
			if err == nil && tt.err {
				t.Fatal("got nil, expected error")
			}
			if err == nil {
				assertObject(t, out, tt.exp)
			}
		})
	}
}

func assertObject(t *testing.T, act types.Object, exp any) {
	switch exp := exp.(type) {
	case nil:
		assertNullObject(t, act)
	case int:
		assertIntegerObject(t, act, int64(exp))
	case float64:
		assertFloatObject(t, act, exp)
	case bool:
		assertBooleanObject(t, act, exp)
	case string:
		assertStringObject(t, act, exp)
	case []any:
		assertArrayObject(t, act, exp)
	default:
		t.Errorf("got %s, expected %T", act.Type(), exp)
	}
}

func assertNullObject(t *testing.T, act any) {
	_, ok := act.(*types.Null)
	if !ok {
		t.Errorf("got %T, expected types.Null", act)
	}
}

func assertIntegerObject(t *testing.T, act any, exp int64) {
	obj, ok := act.(*types.Integer)
	if !ok {
		t.Errorf("got %T, expected types.Integer", act)
		return
	}

	if obj.Value != exp {
		t.Errorf("got %d, expected %d", obj.Value, exp)
	}
}

func assertFloatObject(t *testing.T, act any, exp float64) {
	obj, ok := act.(*types.Float)
	if !ok {
		t.Errorf("got %T, expected types.Float", act)
		return
	}

	if obj.Value != exp {
		t.Errorf("got %v, expected %v", obj.Value, exp)
	}
}

func assertBooleanObject(t *testing.T, act any, exp bool) {
	obj, ok := act.(*types.Boolean)
	if !ok {
		t.Errorf("got %T, expected types.Boolean", act)
		return
	}

	if obj.Value != exp {
		t.Errorf("got %v, expected %v", obj.Value, exp)
	}
}

func assertStringObject(t *testing.T, act any, exp string) {
	obj, ok := act.(*types.String)
	if !ok {
		t.Errorf("got %T, expected types.String", act)
		return
	}

	if obj.Value != exp {
		t.Errorf("got %v, expected %v", obj.Value, exp)
	}
}

func assertArrayObject(t *testing.T, act any, exp []any) {
	obj, ok := act.(types.Array)
	if !ok {
		t.Errorf("got %T, expected types.Array", act)
		return
	}

	for i, exp := range exp {
		assertObject(t, obj[i], exp)
	}
}

func compile(input string) *vm.Program {
	n, err := parser.New(input).Parse()
	if err != nil {
		panic(err)
	}

	p, err := compiler.New().Compile(n)
	if err != nil {
		panic(err)
	}

	return p
}
