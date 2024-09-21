package parser_test

import (
	"fmt"
	"testing"

	"github.com/stevecallear/mexl/ast"
	"github.com/stevecallear/mexl/parser"
)

func TestIdentifierExpression(t *testing.T) {
	const input = "abc"
	n := parse(input)
	assertIdentifier(t, n, "abc")
}

func TestLiteral(t *testing.T) {
	tests := []struct {
		input string
		exp   any
	}{
		{"5", 5},
		{".5", 0.5},
		{"0.5", 0.5},
		{"5.5", 5.5},
		{`"abc"`, "abc"},
		{"false", false},
		{"true", true},
		{"[]", []any{}},
		{`[5, 0.5, "abc", true, [1]]`, []any{5, 0.5, "abc", true, []any{1}}},
		{"null", nil},
	}

	for _, tt := range tests {
		n := parse(tt.input)
		assertLiteral(t, n, tt.exp)
	}
}

func TestPrefixExpression(t *testing.T) {
	tests := []struct {
		input    string
		operator string
		exp      any
	}{
		{"-15", "-", 15},
		{"not true", "not", true},
		{"!true", "!", true},
	}

	for _, tt := range tests {
		n := parse(tt.input)

		e, ok := n.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("expected prefix expression, got %T", n)
		}

		if e.Operator != tt.operator {
			t.Errorf("got %s, expected %s", e.Operator, tt.operator)
		}

		assertLiteral(t, e.Right, tt.exp)
	}
}

func TestInfixExpression(t *testing.T) {
	tests := []struct {
		input    string
		left     any
		operator string
		right    any
	}{
		{"5 + 5", 5, "+", 5},
		{"5 - 5", 5, "-", 5},
		{"5 * 5", 5, "*", 5},
		{"5 / 5", 5, "/", 5},
		{"5 % 5", 5, "%", 5},
		{"5 == 5", 5, "==", 5},
		{"5 eq 5", 5, "eq", 5},
		{"5 != 5", 5, "!=", 5},
		{"5 ne 5", 5, "ne", 5},
		{"5 < 5", 5, "<", 5},
		{"5 lt 5", 5, "lt", 5},
		{"5 > 5", 5, ">", 5},
		{"5 gt 5", 5, "gt", 5},
		{"5 <= 5", 5, "<=", 5},
		{"5 le 5", 5, "le", 5},
		{"5 >= 5", 5, ">=", 5},
		{"5 ge 5", 5, "ge", 5},
		{`"abc" + "abc"`, "abc", "+", "abc"},
		{`"abc" == "abc"`, "abc", "==", "abc"},
		{`"abc" eq "abc"`, "abc", "eq", "abc"},
		{`"abc" != "abc"`, "abc", "!=", "abc"},
		{`"abc" ne "abc"`, "abc", "ne", "abc"},
		{`"abc" sw "a"`, "abc", "sw", "a"},
		{`"abc" ew "c"`, "abc", "ew", "c"},
		{"1 in [1, 2]", 1, "in", []any{1, 2}},
		{"true && true", true, "&&", true},
		{"true && false", true, "&&", false},
		{"true and true", true, "and", true},
		{"true and false", true, "and", false},
		{"true || true", true, "||", true},
		{"true || false", true, "||", false},
		{"true or true", true, "or", true},
		{"true or false", true, "or", false},
	}

	for _, tt := range tests {
		n := parse(tt.input)
		assertInfixExpression(t, n, tt.left, tt.operator, tt.right)
	}
}

func TestCallExpression(t *testing.T) {
	const input = `fn(1, "a")`
	args := []any{1, "a"}

	n := parse(input)

	e, ok := n.(*ast.CallExpression)
	if !ok {
		t.Fatalf("got %T, expected member expression", n)
	}

	assertIdentifier(t, e.Function, "fn")

	if al, el := len(e.Arguments), len(args); al != el {
		t.Fatalf("got %d arguments, expected %d", al, el)
	}

	for i, a := range e.Arguments {
		assertLiteral(t, a, args[i])
	}
}

func TestIndexExpression(t *testing.T) {
	const input = "[1, 2][1]"
	arr := []any{1, 2}
	idx := 1

	n := parse(input)

	e, ok := n.(*ast.IndexExpression)
	if !ok {
		t.Fatalf("got %T, expected index expression", n)
	}

	assertArrayLiteral(t, e.Left, arr)
	assertLiteral(t, e.Index, idx)
}

func TestMemberExpression(t *testing.T) {
	const input = "x.y"
	const left = "x"
	const member = "y"

	n := parse(input)

	e, ok := n.(*ast.MemberExpression)
	if !ok {
		t.Fatalf("got %T, expected member expression", n)
	}

	assertIdentifier(t, e.Left, left)
	assertIdentifier(t, e.Member, member)
}

func TestOperatorPrecedence(t *testing.T) {
	tests := []struct {
		input string
		exp   string
	}{
		{"3 + 4 * 5 == 3 * 1 + 4 * 5", "((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))"},
		{"5 lt 4 ne 3 gt 4", "((5 lt 4) ne (3 gt 4))"},
		{"x le 3", "(x le 3)"},
		{"x ge 3", "(x ge 3)"},
		{"true", "true"},
		{"false", "false"},
		{"3 gt 5 eq false", "((3 gt 5) eq false)"},
		{"3 lt 5 eq true", "((3 lt 5) eq true)"},
		{"true and true", "(true and true)"},
		{"true and (true or false)", "(true and (true or false))"},
		{"not true and false", "((not true) and false)"},
		{`fn("abc") eq "ABC"`, "(fn(abc) eq ABC)"},
		{`y eq "a" and x lt 3`, "((y eq a) and (x lt 3))"},
		{`fn(y) eq "a" and x lt 3`, "((fn(y) eq a) and (x lt 3))"},
		{`fn("abc") sw "A" and x lt 3`, "((fn(abc) sw A) and (x lt 3))"},
		{"2 in [1, 2, 3]", "(2 in [1, 2, 3])"},
		{"[1, 2, 3][0]", "([1, 2, 3][0])"},
		{"x + y + z", "((x + y) + z)"},
		{"x.y.z", "((x.y).z)"},
		{"1 in x.y", "(1 in (x.y))"},
	}

	for _, tt := range tests {
		s, err := parser.New(tt.input).Parse()
		if err != nil {
			t.Fatalf("got %v, expected nil", err)
		}

		if act := s.String(); act != tt.exp {
			t.Errorf("got %s, expected %s", act, tt.exp)
		}
	}
}

func TestErrors(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "invalid expression",
			input: "!",
		},
		{
			name:  "multiple expressions",
			input: "x y",
		},
		{
			name:  "invalid prefix",
			input: "$x",
		},
		{
			name:  "invalid infix",
			input: "1 $ 2",
		},
		{
			name:  "invalid group",
			input: "(1 + 2",
		},
		{
			name:  "invalid index",
			input: "x[1",
		},
		{
			name:  "invalid call",
			input: "x(1",
		},
		{
			name:  "invalid float",
			input: "1.2.3",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := parser.New(tt.input).Parse()
			if err == nil {
				t.Fatal("got nil, expected error")
			}
		})
	}
}

func assertIdentifier(t *testing.T, e ast.Node, value string) {
	t.Helper()

	i, ok := e.(*ast.Identifier)
	if !ok {
		t.Fatalf("expected identifier, got %T", e)
		return
	}

	if i.Value != value {
		t.Errorf("got %s, expected %s", i.Value, value)
		return
	}

	if i.TokenLiteral() != value {
		t.Errorf("got %s, expected %s", i.TokenLiteral(), value)
		return
	}
}

func assertInfixExpression(t *testing.T, n ast.Node, left any, operator string, right any) {
	t.Helper()

	e, ok := n.(*ast.InfixExpression)
	if !ok {
		t.Fatalf("got %T, expected infix expression", n)
		return
	}

	if !assertLiteral(t, e.Left, left) {
		return
	}

	if e.Operator != operator {
		t.Errorf("got %s, expected %s", e.Operator, operator)
		return
	}

	assertLiteral(t, e.Right, right)
}

func assertLiteral(t *testing.T, n ast.Node, value any) bool {
	t.Helper()

	switch tv := value.(type) {
	case string:
		return assertStringLiteral(t, n, tv)
	case int:
		return assertIntegerLiteral(t, n, int64(tv))
	case int64:
		return assertIntegerLiteral(t, n, tv)
	case float64:
		return assertFloatLiteral(t, n, tv)
	case bool:
		return assertBoolean(t, n, tv)
	case []any:
		return assertArrayLiteral(t, n, tv)
	case nil:
		return assertNull(t, n)
	default:
		t.Fatalf("unexpected type: %T", value)
		return false
	}
}

func assertIntegerLiteral(t *testing.T, n ast.Node, value int64) bool {
	t.Helper()

	l, ok := n.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("got %T, expected integer literal", n)
		return false
	}

	if l.Value != value {
		t.Errorf("got %d, expected %d", l.Value, value)
		return false
	}

	if l.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("got %s, expected %d", l.TokenLiteral(), value)
		return false
	}

	return true
}

func assertFloatLiteral(t *testing.T, n ast.Node, value float64) bool {
	t.Helper()

	l, ok := n.(*ast.FloatLiteral)
	if !ok {
		t.Errorf("got %T, expected float literal", n)
		return false
	}

	if l.Value != value {
		t.Errorf("got %v, expected %v", l.Value, value)
		return false
	}

	return true
}

func assertBoolean(t *testing.T, n ast.Node, value bool) bool {
	t.Helper()

	l, ok := n.(*ast.Boolean)
	if !ok {
		t.Errorf("got %T, expected integer literal", n)
		return false
	}

	if l.Value != value {
		t.Errorf("got %v, expected %v", l.Value, value)
		return false
	}

	if l.TokenLiteral() != fmt.Sprintf("%v", value) {
		t.Errorf("got %s, expected %v", l.TokenLiteral(), value)
		return false
	}

	return true
}

func assertStringLiteral(t *testing.T, n ast.Node, value string) bool {
	t.Helper()

	l, ok := n.(*ast.StringLiteral)
	if !ok {
		t.Errorf("got %T, expected string literal", n)
		return false
	}

	if l.Value != value {
		t.Errorf("got %s, expected %s", l.Value, value)
		return false
	}

	if l.TokenLiteral() != value {
		t.Errorf("got %s, expected %s", l.TokenLiteral(), value)
		return false
	}

	return true
}

func assertArrayLiteral(t *testing.T, n ast.Node, value []any) bool {
	t.Helper()

	l, ok := n.(*ast.ArrayLiteral)
	if !ok {
		t.Errorf("got %T, expected array literal", n)
		return false
	}

	if al, el := len(l.Elements), len(value); al != el {
		t.Errorf("got %d elements, expected %d", al, el)
	}

	for i, ee := range l.Elements {
		if !assertLiteral(t, ee, value[i]) {
			return false
		}
	}

	return true
}

func assertNull(t *testing.T, n ast.Node) bool {
	t.Helper()

	nn, ok := n.(*ast.Null)
	if !ok {
		t.Errorf("got %T, expected null", n)
		return false
	}

	if nn.TokenLiteral() != "null" {
		t.Errorf("got %s, expected null", nn.TokenLiteral())
		return false
	}

	return true
}

func parse(input string) ast.Node {
	n, err := parser.New(input).Parse()
	if err != nil {
		panic(err)
	}
	return n
}
