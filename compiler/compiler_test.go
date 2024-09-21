package compiler_test

import (
	"testing"

	"github.com/stevecallear/mexl/ast"
	"github.com/stevecallear/mexl/ast/token"
	"github.com/stevecallear/mexl/compiler"
	"github.com/stevecallear/mexl/parser"
	"github.com/stevecallear/mexl/types"
	"github.com/stevecallear/mexl/vm"
)

type (
	testCase struct {
		name string
		node ast.Node
		exp  expectation
		err  bool
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
			name: "addition",
			node: parse("1 + 2"),
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
			name: "subtraction",
			node: parse("2 - 1"),
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
			name: "multiplication",
			node: parse("1 * 2"),
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
			name: "division",
			node: parse("2 / 1"),
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
			name: "modulus",
			node: parse("5 % 2"),
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
			name: "minus",
			node: parse("-1"),
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
			name: "addition",
			node: parse("1.1 + 2.2"),
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
			name: "subtraction",
			node: parse("2.2 - 1.1"),
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
			name: "multiplication",
			node: parse("1.1 * 2.2"),
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
			name: "division",
			node: parse("2.2 / 1.1"),
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
			name: "minus",
			node: parse("-1.1"),
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
			name: "true",
			node: parse("true"),
			exp: expectation{
				constants: []any{},
				instructions: []vm.Instructions{
					vm.Make(vm.OpTrue),
				},
			},
		},
		{
			name: "false",
			node: parse("false"),
			exp: expectation{
				constants: []any{},
				instructions: []vm.Instructions{
					vm.Make(vm.OpFalse),
				},
			},
		},
		{
			name: "int greater than",
			node: parse("1 > 2"),
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
			name: "int greater than or equal",
			node: parse("1 >= 2"),
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
			name: "int less than",
			node: parse("1 < 2"),
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
			name: "int less than or equal",
			node: parse("1 <= 2"),
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
			name: "float equal",
			node: parse("1.1 == 2.2"),
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
			name: "float not equal",
			node: parse("1.1 != 2.2"),
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
			name: "float greater than",
			node: parse("1.1 > 2.2"),
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
			name: "float greater than or equal",
			node: parse("1.1 >= 2.2"),
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
			name: "float less than",
			node: parse("1.1 < 2.2"),
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
			name: "float less than or equal",
			node: parse("1.1 <= 2.2"),
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
			name: "float equal",
			node: parse("1.1 == 2.2"),
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
			name: "float not equal",
			node: parse("1.1 != 2.2"),
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
			name: "mixed types 1",
			node: parse("2.2 > 1"),
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
			name: "mixed types 2",
			node: parse("1 < 2.2"),
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
			name: "bool equal",
			node: parse("true == false"),
			exp: expectation{
				instructions: []vm.Instructions{
					vm.Make(vm.OpTrue),
					vm.Make(vm.OpFalse),
					vm.Make(vm.OpEqual),
				},
			},
		},
		{
			name: "bool not equal",
			node: parse("true != false"),
			exp: expectation{
				instructions: []vm.Instructions{
					vm.Make(vm.OpTrue),
					vm.Make(vm.OpFalse),
					vm.Make(vm.OpNotEqual),
				},
			},
		},
		{
			name: "bang",
			node: parse("!true"),
			exp: expectation{
				instructions: []vm.Instructions{
					vm.Make(vm.OpTrue),
					vm.Make(vm.OpNot),
				},
			},
		},
		{
			name: "and",
			node: parse("true && false"),
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
			name: "or",
			node: parse("true || false"),
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
			name: "start with",
			node: parse(`"abc" sw "a"`),
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
			name: "ends with",
			node: parse(`"abc" ew "c"`),
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
			name: "in",
			node: parse(`2 in ["a", 2, 3.3]`),
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
			name: "concatenation",
			node: parse(`"mess" + "age"`),
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
			name: "constant",
			node: parse(`"message"`),
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
			name: "array",
			node: parse(`["a", 2, 3.3]`),
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
			name: "array expression",
			node: parse("[1, 3 > 2]"),
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
			name: "array index",
			node: parse("[1, 2, 3][1]"),
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
			name: "global",
			node: parse("x"),
			exp: expectation{
				identifiers: []string{"x"},
				instructions: []vm.Instructions{
					vm.Make(vm.OpGlobal, 0),
				},
			},
		},
		{
			name: "global expression",
			node: parse("x < 2"),
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
			name: "builtin",
			node: parse(`upper("test")`),
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
		{
			name: "builtin (reuse)",
			node: parse(`upper("test") eq upper("TEST")`),
			exp: expectation{
				constants:   []any{"test", "TEST"},
				identifiers: []string{"upper"},
				instructions: []vm.Instructions{
					vm.Make(vm.OpGlobal, 0),
					vm.Make(vm.OpConstant, 0),
					vm.Make(vm.OpCall, 1),
					vm.Make(vm.OpGlobal, 0),
					vm.Make(vm.OpConstant, 1),
					vm.Make(vm.OpCall, 1),
					vm.Make(vm.OpEqual),
				},
			},
		},
	}

	testCompiler(t, tests)
}

func TestNull(t *testing.T) {
	tests := []testCase{
		{
			name: "null",
			node: parse(`email eq null`),
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
			name: "or jump",
			node: parse("true or false"),
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
			name: "and jump",
			node: parse("false and true"),
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
			name: "root",
			node: parse("x"),
			exp: expectation{
				identifiers: []string{"x"},
				instructions: []vm.Instructions{
					vm.Make(vm.OpGlobal, 0),
				},
			},
		},
		{
			name: "depth",
			node: parse("x.y.z"),
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

func TestErrors(t *testing.T) {
	tests := []testCase{
		{
			name: "invalid prefix operator",
			node: &ast.PrefixExpression{
				Token: token.Token{
					Type:    token.Minus,
					Literal: "-",
				},
				Operator: "$",
				Right: &ast.IntegerLiteral{
					Token: token.Token{
						Type:    token.Ident,
						Literal: "1",
					},
					Value: 1,
				},
			},
			err: true,
		},
		{
			name: "invalid prefix operand",
			node: &ast.PrefixExpression{
				Token: token.Token{
					Type:    token.Minus,
					Literal: "-",
				},
				Operator: "-",
				Right:    new(invalidNode),
			},
			err: true,
		},
		{
			name: "invalid infix operator",
			node: &ast.InfixExpression{
				Token: token.Token{
					Type:    token.Plus,
					Literal: "+",
				},
				Left: &ast.IntegerLiteral{
					Token: token.Token{
						Type:    token.Ident,
						Literal: "1",
					},
					Value: 1,
				},
				Operator: "$",
				Right: &ast.IntegerLiteral{
					Token: token.Token{
						Type:    token.Ident,
						Literal: "1",
					},
					Value: 1,
				},
			},
			err: true,
		},
		{
			name: "invalid infix left operand",
			node: &ast.InfixExpression{
				Token: token.Token{
					Type:    token.Plus,
					Literal: "+",
				},
				Left:     new(invalidNode),
				Operator: "+",
				Right: &ast.IntegerLiteral{
					Token: token.Token{
						Type:    token.Ident,
						Literal: "1",
					},
					Value: 1,
				},
			},
			err: true,
		},
		{
			name: "invalid infix right operand",
			node: &ast.InfixExpression{
				Token: token.Token{
					Type:    token.Plus,
					Literal: "+",
				},
				Operator: "+",
				Left: &ast.IntegerLiteral{
					Token: token.Token{
						Type:    token.Ident,
						Literal: "1",
					},
					Value: 1,
				},
				Right: new(invalidNode),
			},
			err: true,
		},
		{
			name: "invalid array expression",
			node: &ast.ArrayLiteral{
				Token: token.Token{
					Type:    token.LBracket,
					Literal: "[",
				},
				Elements: []ast.Node{
					new(invalidNode),
				},
			},
			err: true,
		},
		{
			name: "invalid index expression (left)",
			node: &ast.IndexExpression{
				Token: token.Token{
					Type:    token.LBracket,
					Literal: "[",
				},
				Left: new(invalidNode),
				Index: &ast.IntegerLiteral{
					Token: token.Token{
						Type:    token.Int,
						Literal: "1",
					},
					Value: 1,
				},
			},
			err: true,
		},
		{
			name: "invalid index expression (index)",
			node: &ast.IndexExpression{
				Token: token.Token{
					Type:    token.LBracket,
					Literal: "[",
				},
				Left: &ast.ArrayLiteral{
					Token: token.Token{
						Type:    token.LBracket,
						Literal: "[",
					},
					Elements: []ast.Node{
						&ast.IntegerLiteral{
							Token: token.Token{
								Type:    token.Int,
								Literal: "1",
							},
							Value: 1,
						},
					},
				},
				Index: new(invalidNode),
			},
			err: true,
		},
		{
			name: "invalid member expression (target)",
			node: &ast.MemberExpression{
				Token: token.Token{
					Type:    token.Stop,
					Literal: ".",
				},
				Left: new(invalidNode),
				Member: &ast.Identifier{
					Token: token.Token{
						Type:    token.Ident,
						Literal: "y",
					},
					Value: "y",
				},
			},
			err: true,
		},
		{
			name: "invalid member expression (member)",
			node: &ast.MemberExpression{
				Token: token.Token{
					Type:    token.Stop,
					Literal: ".",
				},
				Left: &ast.Identifier{
					Token: token.Token{
						Type:    token.Ident,
						Literal: "x",
					},
					Value: "x",
				},
				Member: new(invalidNode),
			},
			err: true,
		},
		{
			name: "invalid call expression (target)",
			node: &ast.CallExpression{
				Token: token.Token{
					Type:    token.LParen,
					Literal: "(",
				},
				Function: new(invalidNode),
				Arguments: []ast.Node{
					&ast.IntegerLiteral{
						Token: token.Token{
							Type:    token.Int,
							Literal: "1",
						},
						Value: 1,
					},
				},
			},
			err: true,
		},
		{
			name: "invalid call expression (args)",
			node: &ast.CallExpression{
				Token: token.Token{
					Type:    token.LParen,
					Literal: "(",
				},
				Function: &ast.Identifier{
					Token: token.Token{
						Type:    token.Ident,
						Literal: "x",
					},
					Value: "x",
				},
				Arguments: []ast.Node{
					new(invalidNode),
				},
			},
			err: true,
		},
	}

	testCompiler(t, tests)
}

func testCompiler(t *testing.T, tests []testCase) {
	t.Helper()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := compiler.New()
			p, err := c.Compile(tt.node)
			if err != nil && !tt.err {
				t.Fatalf("got %v, expected nil", err)
			}
			if err == nil && tt.err {
				t.Fatal("got nil, expected error")
			}
			if err != nil {
				return
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

type invalidNode struct{}

var _ ast.Node = (*invalidNode)(nil)

func (n *invalidNode) TokenLiteral() string {
	return "invalid"
}

func (n *invalidNode) String() string {
	return "invalid"
}
