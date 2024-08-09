package token

type (
	Type uint8

	Token struct {
		Type    Type
		Literal string
	}
)

const (
	Illegal Type = iota
	EOF

	Ident
	Int
	Float
	String
	Plus

	True
	False
	Null

	Minus
	Asterisk
	Slash
	Percent

	And
	Or
	Bang
	Equal
	LessThan
	GreaterThan
	NotEqual
	LessThanOrEqual
	GreaterThanOrEqual
	StartsWith
	EndsWith
	In

	LParen
	RParen
	LBracket
	RBracket

	Stop
	Comma
)
