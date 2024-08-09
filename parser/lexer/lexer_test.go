package lexer_test

import (
	"reflect"
	"testing"

	"github.com/stevecallear/mexl/ast/token"
	"github.com/stevecallear/mexl/parser/lexer"
)

func TestLexer_NextToken(t *testing.T) {
	tests := []struct {
		name  string
		input string
		exp   []token.Token
	}{
		{
			name:  "operators",
			input: "+ - * / % not and or eq ne lt gt le ge sw ew in ! && || == != < > <= >= !",
			exp: []token.Token{
				{Type: token.Plus, Literal: "+"},
				{Type: token.Minus, Literal: "-"},
				{Type: token.Asterisk, Literal: "*"},
				{Type: token.Slash, Literal: "/"},
				{Type: token.Percent, Literal: "%"},
				{Type: token.Bang, Literal: "not"},
				{Type: token.And, Literal: "and"},
				{Type: token.Or, Literal: "or"},
				{Type: token.Equal, Literal: "eq"},
				{Type: token.NotEqual, Literal: "ne"},
				{Type: token.LessThan, Literal: "lt"},
				{Type: token.GreaterThan, Literal: "gt"},
				{Type: token.LessThanOrEqual, Literal: "le"},
				{Type: token.GreaterThanOrEqual, Literal: "ge"},
				{Type: token.StartsWith, Literal: "sw"},
				{Type: token.EndsWith, Literal: "ew"},
				{Type: token.In, Literal: "in"},
				{Type: token.Bang, Literal: "!"},
				{Type: token.And, Literal: "&&"},
				{Type: token.Or, Literal: "||"},
				{Type: token.Equal, Literal: "=="},
				{Type: token.NotEqual, Literal: "!="},
				{Type: token.LessThan, Literal: "<"},
				{Type: token.GreaterThan, Literal: ">"},
				{Type: token.LessThanOrEqual, Literal: "<="},
				{Type: token.GreaterThanOrEqual, Literal: ">="},
				{Type: token.Bang, Literal: "!"}, // test EOF peek
			},
		},
		{
			name:  "delimiters",
			input: ".,()[]",
			exp: []token.Token{
				{Type: token.Stop, Literal: "."},
				{Type: token.Comma, Literal: ","},
				{Type: token.LParen, Literal: "("},
				{Type: token.RParen, Literal: ")"},
				{Type: token.LBracket, Literal: "["},
				{Type: token.RBracket, Literal: "]"},
			},
		},
		{
			name:  "whitespace",
			input: "eq eq\teq\neq\req",
			exp: []token.Token{
				{Type: token.Equal, Literal: "eq"},
				{Type: token.Equal, Literal: "eq"},
				{Type: token.Equal, Literal: "eq"},
				{Type: token.Equal, Literal: "eq"},
				{Type: token.Equal, Literal: "eq"},
			},
		},
		{
			name:  "integers",
			input: "1 123 -123 (123)",
			exp: []token.Token{
				{Type: token.Int, Literal: "1"},
				{Type: token.Int, Literal: "123"},
				{Type: token.Minus, Literal: "-"},
				{Type: token.Int, Literal: "123"},
				{Type: token.LParen, Literal: "("},
				{Type: token.Int, Literal: "123"},
				{Type: token.RParen, Literal: ")"},
			},
		},
		{
			name:  "floats",
			input: ".1 0.1 1.1 -0.1 (1.1)",
			exp: []token.Token{
				{Type: token.Float, Literal: ".1"},
				{Type: token.Float, Literal: "0.1"},
				{Type: token.Float, Literal: "1.1"},
				{Type: token.Minus, Literal: "-"},
				{Type: token.Float, Literal: "0.1"},
				{Type: token.LParen, Literal: "("},
				{Type: token.Float, Literal: "1.1"},
				{Type: token.RParen, Literal: ")"},
			},
		},
		{
			name:  "idents",
			input: "abc ABC x.y",
			exp: []token.Token{
				{Type: token.Ident, Literal: "abc"},
				{Type: token.Ident, Literal: "ABC"},
				{Type: token.Ident, Literal: "x"},
				{Type: token.Stop, Literal: "."},
				{Type: token.Ident, Literal: "y"},
			},
		},
		{
			name:  "booleans",
			input: "true false",
			exp: []token.Token{
				{Type: token.True, Literal: "true"},
				{Type: token.False, Literal: "false"},
			},
		},
		{
			name:  "null",
			input: "null",
			exp: []token.Token{
				{Type: token.Null, Literal: "null"},
			},
		},
		{
			name:  "strings",
			input: "\"abc 123\" \"abc",
			exp: []token.Token{
				{Type: token.String, Literal: "abc 123"},
				{Type: token.Illegal, Literal: "abc"},
			},
		},
		{
			name:  "calls",
			input: `fn("abc")`,
			exp: []token.Token{
				{Type: token.Ident, Literal: "fn"},
				{Type: token.LParen, Literal: "("},
				{Type: token.String, Literal: "abc"},
				{Type: token.RParen, Literal: ")"},
			},
		},
		{
			name:  "illegal tokens",
			input: "$ =$ !$ <$ >$ &$ |$",
			exp: []token.Token{
				{Type: token.Illegal, Literal: "$"},
				{Type: token.Illegal, Literal: "=$"},
				{Type: token.Bang, Literal: "!"},
				{Type: token.Illegal, Literal: "$"},
				{Type: token.LessThan, Literal: "<"},
				{Type: token.Illegal, Literal: "$"},
				{Type: token.GreaterThan, Literal: ">"},
				{Type: token.Illegal, Literal: "$"},
				{Type: token.Illegal, Literal: "&$"},
				{Type: token.Illegal, Literal: "|$"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sut := lexer.New(tt.input)
			var act []token.Token

			for tok := sut.NextToken(); tok.Type != token.EOF; tok = sut.NextToken() {
				act = append(act, tok)
			}

			if !reflect.DeepEqual(act, tt.exp) {
				t.Errorf("got %v, expected %v", act, tt.exp)
			}
		})
	}
}
