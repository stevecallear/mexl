package ast_test

import (
	"testing"

	"github.com/stevecallear/mexl/ast"
	"github.com/stevecallear/mexl/ast/token"
)

func TestNode(t *testing.T) {
	tests := []struct {
		name string
		sut  ast.Node
		lit  string
		str  string
	}{
		{
			name: "identifier",
			sut: &ast.Identifier{
				Token: token.Token{Type: token.Ident, Literal: "ident"},
				Value: "ident",
			},
			lit: "ident",
			str: "ident",
		},
		{
			name: "integer literal",
			sut: &ast.IntegerLiteral{
				Token: token.Token{Type: token.Int, Literal: "1"},
				Value: 1,
			},
			lit: "1",
			str: "1",
		},
		{
			name: "float literal",
			sut: &ast.FloatLiteral{
				Token: token.Token{Type: token.Float, Literal: "1.1"},
				Value: 1,
			},
			lit: "1.1",
			str: "1.1",
		},
		{
			name: "string literal",
			sut: &ast.StringLiteral{
				Token: token.Token{Type: token.String, Literal: "abc"},
				Value: "abc",
			},
			lit: "abc",
			str: "abc",
		},
		{
			name: "boolean",
			sut: &ast.Boolean{
				Token: token.Token{Type: token.True, Literal: "true"},
				Value: true,
			},
			lit: "true",
			str: "true",
		},
		{
			name: "null",
			sut: &ast.Null{
				Token: token.Token{Type: token.Null, Literal: "null"},
			},
			lit: "null",
			str: "null",
		},
		{
			name: "array literal",
			sut: &ast.ArrayLiteral{
				Token: token.Token{Type: token.RBracket, Literal: "]"},
				Elements: []ast.Node{
					&ast.IntegerLiteral{
						Token: token.Token{Type: token.Int, Literal: "1"},
						Value: 1,
					},
					&ast.IntegerLiteral{
						Token: token.Token{Type: token.Int, Literal: "2"},
						Value: 2,
					},
				},
			},
			lit: "]",
			str: "[1, 2]",
		},
		{
			name: "index expression",
			sut: &ast.IndexExpression{
				Token: token.Token{Type: token.RParen, Literal: "]"},
				Left: &ast.ArrayLiteral{
					Token: token.Token{Type: token.RBracket, Literal: "]"},
					Elements: []ast.Node{
						&ast.IntegerLiteral{
							Token: token.Token{Type: token.Int, Literal: "1"},
							Value: 1,
						},
						&ast.StringLiteral{
							Token: token.Token{Type: token.Int, Literal: "foo"},
							Value: "foo",
						},
					},
				},
				Index: &ast.IntegerLiteral{
					Token: token.Token{Type: token.Int, Literal: "1"},
					Value: 1,
				},
			},
			lit: "]",
			str: "([1, foo][1])",
		},
		{
			name: "member expression",
			sut: &ast.MemberExpression{
				Token: token.Token{Type: token.Stop, Literal: "."},
				Left: &ast.Identifier{
					Token: token.Token{Type: token.Ident, Literal: "x"},
					Value: "x",
				},
				Member: &ast.Identifier{
					Token: token.Token{Type: token.Ident, Literal: "y"},
					Value: "y",
				},
			},
			lit: ".",
			str: "(x.y)",
		},
		{
			name: "prefix expression",
			sut: &ast.PrefixExpression{
				Token:    token.Token{Type: token.Minus, Literal: "-"},
				Operator: "-",
				Right: &ast.IntegerLiteral{
					Token: token.Token{Type: token.Int, Literal: "5"},
					Value: 5,
				},
			},
			lit: "-",
			str: "(-5)",
		},
		{
			name: "prefix expression multi char",
			sut: &ast.PrefixExpression{
				Token:    token.Token{Type: token.Bang, Literal: "not"},
				Operator: "not",
				Right: &ast.Boolean{
					Token: token.Token{Type: token.True, Literal: "true"},
					Value: true,
				},
			},
			lit: "not",
			str: "(not true)",
		},
		{
			name: "infix expression",
			sut: &ast.InfixExpression{
				Token: token.Token{Type: token.Plus, Literal: "+"},
				Left: &ast.IntegerLiteral{
					Token: token.Token{Type: token.Int, Literal: "4"},
					Value: 4,
				},
				Operator: "+",
				Right: &ast.IntegerLiteral{
					Token: token.Token{Type: token.Int, Literal: "5"},
					Value: 5,
				},
			},
			lit: "+",
			str: "(4 + 5)",
		},
		{
			name: "call expression",
			sut: &ast.CallExpression{
				Token: token.Token{Type: token.RParen, Literal: ")"},
				Function: &ast.Identifier{
					Token: token.Token{Type: token.Ident, Literal: "fn"},
					Value: "fn",
				},
				Arguments: []ast.Node{
					&ast.IntegerLiteral{
						Token: token.Token{Type: token.Int, Literal: "4"},
						Value: 4,
					},
					&ast.IntegerLiteral{
						Token: token.Token{Type: token.Int, Literal: "5"},
						Value: 5,
					},
				},
			},
			lit: ")",
			str: "fn(4, 5)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if act, exp := tt.sut.TokenLiteral(), tt.lit; act != exp {
				t.Errorf("got %s, expected %s", act, exp)
			}
			if act, exp := tt.sut.String(), tt.str; act != exp {
				t.Errorf("got %s, expected %s", act, exp)
			}
		})
	}
}
