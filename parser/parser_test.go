package parser_test

import (
	"fmt"
	"testing"

	"github.com/stevecallear/mexl/ast"
	"github.com/stevecallear/mexl/parser"
)

func TestIdentifierExpression(t *testing.T) {
	const input = "abc"
	s := parse(input)
	assertIdentifierExpression(t, s.Expression, "abc")
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
		s := parse(tt.input)
		assertLiteral(t, s.Expression, tt.exp)
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
		s := parse(tt.input)

		pe, ok := s.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("expected prefix expression, got %T", s.Expression)
		}

		if pe.Operator != tt.operator {
			t.Errorf("got %s, expected %s", pe.Operator, tt.operator)
		}

		assertLiteral(t, pe.Right, tt.exp)
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
		s := parse(tt.input)
		assertInfixExpression(t, s.Expression, tt.left, tt.operator, tt.right)
	}
}

func TestCallExpression(t *testing.T) {
	const input = `fn(1, "a")`
	args := []any{1, "a"}

	s := parse(input)

	c, ok := s.Expression.(*ast.CallExpression)
	if !ok {
		t.Fatalf("got %T, expected member expression", s)
	}

	assertIdentifierExpression(t, c.Function, "fn")

	if al, el := len(c.Arguments), len(args); al != el {
		t.Fatalf("got %d arguments, expected %d", al, el)
	}

	for i, a := range c.Arguments {
		assertLiteral(t, a, args[i])
	}
}

func TestIndexExpression(t *testing.T) {
	const input = "[1, 2][1]"
	arr := []any{1, 2}
	idx := 1

	s := parse(input)

	m, ok := s.Expression.(*ast.IndexExpression)
	if !ok {
		t.Fatalf("got %T, expected member expression", s)
	}

	assertArrayLiteral(t, m.Left, arr)
	assertLiteral(t, m.Index, idx)
}

func TestMemberExpression(t *testing.T) {
	const input = "x.y"
	const left = "x"
	const member = "y"

	s := parse(input)

	m, ok := s.Expression.(*ast.MemberExpression)
	if !ok {
		t.Fatalf("got %T, expected member expression", s)
	}

	assertIdentifierExpression(t, m.Left, left)
	assertIdentifierExpression(t, m.Member, member)
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

func assertIdentifierExpression(t *testing.T, e ast.Expression, value string) {
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

func assertInfixExpression(t *testing.T, e ast.Expression, left any, operator string, right any) {
	t.Helper()

	ie, ok := e.(*ast.InfixExpression)
	if !ok {
		t.Fatalf("got %T, expected infix expression", e)
		return
	}

	if !assertLiteral(t, ie.Left, left) {
		return
	}

	if ie.Operator != operator {
		t.Errorf("got %s, expected %s", ie.Operator, operator)
		return
	}

	assertLiteral(t, ie.Right, right)
}

func assertLiteral(t *testing.T, e ast.Expression, value any) bool {
	t.Helper()

	switch tv := value.(type) {
	case string:
		return assertStringLiteral(t, e, tv)
	case int:
		return assertIntegerLiteral(t, e, int64(tv))
	case int64:
		return assertIntegerLiteral(t, e, tv)
	case float64:
		return assertFloatLiteral(t, e, tv)
	case bool:
		return assertBoolean(t, e, tv)
	case []any:
		return assertArrayLiteral(t, e, tv)
	case nil:
		return assertNull(t, e)
	default:
		t.Fatalf("unexpected type: %T", value)
		return false
	}
}

func assertIntegerLiteral(t *testing.T, e ast.Expression, value int64) bool {
	t.Helper()

	l, ok := e.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("got %T, expected integer literal", e)
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

func assertFloatLiteral(t *testing.T, e ast.Expression, value float64) bool {
	t.Helper()

	l, ok := e.(*ast.FloatLiteral)
	if !ok {
		t.Errorf("got %T, expected float literal", e)
		return false
	}

	if l.Value != value {
		t.Errorf("got %v, expected %v", l.Value, value)
		return false
	}

	return true
}

func assertBoolean(t *testing.T, e ast.Expression, value bool) bool {
	t.Helper()

	l, ok := e.(*ast.Boolean)
	if !ok {
		t.Errorf("got %T, expected integer literal", e)
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

func assertStringLiteral(t *testing.T, e ast.Expression, value string) bool {
	t.Helper()

	l, ok := e.(*ast.StringLiteral)
	if !ok {
		t.Errorf("got %T, expected string literal", e)
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

func assertArrayLiteral(t *testing.T, e ast.Expression, value []any) bool {
	t.Helper()

	l, ok := e.(*ast.ArrayLiteral)
	if !ok {
		t.Errorf("got %T, expected array literal", e)
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

func assertNull(t *testing.T, e ast.Expression) bool {
	t.Helper()

	l, ok := e.(*ast.Null)
	if !ok {
		t.Errorf("got %T, expected null", e)
		return false
	}

	if l.TokenLiteral() != "null" {
		t.Errorf("got %s, expected null", l.TokenLiteral())
		return false
	}

	return true
}

func parse(input string) *ast.ExpressionStatement {
	s, err := parser.New(input).Parse()
	if err != nil {
		panic(err)
	}
	return s.(*ast.ExpressionStatement)
}
