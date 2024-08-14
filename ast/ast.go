package ast

import (
	"strings"

	"github.com/stevecallear/mexl/ast/token"
)

type (
	Node interface {
		TokenLiteral() string
		String() string
	}

	Identifier struct {
		Token token.Token
		Value string
	}

	IntegerLiteral struct {
		Token token.Token
		Value int64
	}

	FloatLiteral struct {
		Token token.Token
		Value float64
	}

	StringLiteral struct {
		Token token.Token
		Value string
	}

	Boolean struct {
		Token token.Token
		Value bool
	}

	Null struct {
		Token token.Token
	}

	ArrayLiteral struct {
		Token    token.Token
		Elements []Node
	}

	IndexExpression struct {
		Token token.Token
		Left  Node
		Index Node
	}

	MemberExpression struct {
		Token  token.Token
		Left   Node
		Member Node
	}

	PrefixExpression struct {
		Token    token.Token
		Operator string
		Right    Node
	}

	InfixExpression struct {
		Token    token.Token
		Left     Node
		Operator string
		Right    Node
	}

	CallExpression struct {
		Token     token.Token
		Function  Node
		Arguments []Node
	}
)

func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}

func (i *Identifier) String() string {
	return i.Value
}

func (i *IntegerLiteral) TokenLiteral() string {
	return i.Token.Literal
}

func (i *IntegerLiteral) String() string {
	return i.Token.Literal
}

func (f *FloatLiteral) TokenLiteral() string {
	return f.Token.Literal
}

func (f *FloatLiteral) String() string {
	return f.Token.Literal
}

func (i *StringLiteral) TokenLiteral() string {
	return i.Token.Literal
}

func (i *StringLiteral) String() string {
	return i.Token.Literal
}

func (b *Boolean) TokenLiteral() string {
	return b.Token.Literal
}

func (b *Boolean) String() string {
	return b.Token.Literal
}

func (b *Null) TokenLiteral() string {
	return b.Token.Literal
}

func (b *Null) String() string {
	return b.Token.Literal
}

func (a *ArrayLiteral) TokenLiteral() string {
	return a.Token.Literal
}

func (a *ArrayLiteral) String() string {
	es := make([]string, len(a.Elements))
	for i, e := range a.Elements {
		es[i] = e.String()
	}
	return "[" + strings.Join(es, ", ") + "]"
}

func (i *IndexExpression) TokenLiteral() string {
	return i.Token.Literal
}

func (i *IndexExpression) String() string {
	return "(" + i.Left.String() + "[" + i.Index.String() + "])"
}

func (m *MemberExpression) TokenLiteral() string {
	return m.Token.Literal
}

func (m *MemberExpression) String() string {
	return "(" + m.Left.String() + "." + m.Member.String() + ")"
}

func (e *PrefixExpression) TokenLiteral() string {
	return e.Token.Literal
}

func (e *PrefixExpression) String() string {
	if len(e.Operator) == 1 {
		return "(" + e.Operator + e.Right.String() + ")"
	}
	return "(" + e.Operator + " " + e.Right.String() + ")"
}

func (e *InfixExpression) TokenLiteral() string {
	return e.Token.Literal
}

func (e *InfixExpression) String() string {
	return "(" + e.Left.String() + " " + e.Operator + " " + e.Right.String() + ")"
}

func (e *CallExpression) TokenLiteral() string {
	return e.Token.Literal
}

func (e *CallExpression) String() string {
	as := make([]string, len(e.Arguments))
	for i, a := range e.Arguments {
		as[i] = a.String()
	}
	return e.Function.String() + "(" + strings.Join(as, ", ") + ")"
}
