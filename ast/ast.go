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

	Statement interface {
		Node
		isStatement()
	}

	Expression interface {
		Node
		isExpression()
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
		Elements []Expression
	}

	IndexExpression struct {
		Token token.Token
		Left  Expression
		Index Expression
	}

	MemberExpression struct {
		Token  token.Token
		Left   Expression
		Member Expression
	}

	PrefixExpression struct {
		Token    token.Token
		Operator string
		Right    Expression
	}

	InfixExpression struct {
		Token    token.Token
		Left     Expression
		Operator string
		Right    Expression
	}

	CallExpression struct {
		Token     token.Token
		Function  Expression
		Arguments []Expression
	}

	ExpressionStatement struct {
		Token      token.Token
		Expression Expression
	}
)

var (
	_ Expression = (*Identifier)(nil)
	_ Expression = (*IntegerLiteral)(nil)
	_ Expression = (*FloatLiteral)(nil)
	_ Expression = (*StringLiteral)(nil)
	_ Expression = (*Boolean)(nil)
	_ Expression = (*Null)(nil)
	_ Expression = (*ArrayLiteral)(nil)
	_ Expression = (*IndexExpression)(nil)
	_ Expression = (*MemberExpression)(nil)
	_ Expression = (*PrefixExpression)(nil)
	_ Expression = (*InfixExpression)(nil)
	_ Expression = (*CallExpression)(nil)
	_ Statement  = (*ExpressionStatement)(nil)
)

func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}

func (i *Identifier) String() string {
	return i.Value
}

func (i *Identifier) isExpression() {}

func (i *IntegerLiteral) TokenLiteral() string {
	return i.Token.Literal
}

func (i *IntegerLiteral) String() string {
	return i.Token.Literal
}

func (i *IntegerLiteral) isExpression() {}

func (f *FloatLiteral) TokenLiteral() string {
	return f.Token.Literal
}

func (f *FloatLiteral) String() string {
	return f.Token.Literal
}

func (f *FloatLiteral) isExpression() {}

func (i *StringLiteral) TokenLiteral() string {
	return i.Token.Literal
}

func (i *StringLiteral) String() string {
	return i.Token.Literal
}

func (i *StringLiteral) isExpression() {}

func (b *Boolean) TokenLiteral() string {
	return b.Token.Literal
}

func (b *Boolean) String() string {
	return b.Token.Literal
}

func (b *Boolean) isExpression() {}

func (b *Null) TokenLiteral() string {
	return b.Token.Literal
}

func (b *Null) String() string {
	return b.Token.Literal
}

func (b *Null) isExpression() {}

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

func (a *ArrayLiteral) isExpression() {}

func (i *IndexExpression) TokenLiteral() string {
	return i.Token.Literal
}

func (i *IndexExpression) String() string {
	return "(" + i.Left.String() + "[" + i.Index.String() + "])"
}

func (i *IndexExpression) isExpression() {}

func (m *MemberExpression) TokenLiteral() string {
	return m.Token.Literal
}

func (m *MemberExpression) String() string {
	return "(" + m.Left.String() + "." + m.Member.String() + ")"
}

func (m *MemberExpression) isExpression() {}

func (e *PrefixExpression) TokenLiteral() string {
	return e.Token.Literal
}

func (e *PrefixExpression) String() string {
	if len(e.Operator) == 1 {
		return "(" + e.Operator + e.Right.String() + ")"
	}
	return "(" + e.Operator + " " + e.Right.String() + ")"
}

func (e *PrefixExpression) isExpression() {}

func (e *InfixExpression) TokenLiteral() string {
	return e.Token.Literal
}

func (e *InfixExpression) String() string {
	return "(" + e.Left.String() + " " + e.Operator + " " + e.Right.String() + ")"
}

func (e *InfixExpression) isExpression() {}

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

func (e *CallExpression) isExpression() {}

func (s *ExpressionStatement) TokenLiteral() string {
	return s.Token.Literal
}

func (s *ExpressionStatement) String() string {
	return s.Expression.String()
}

func (s *ExpressionStatement) isStatement() {}
