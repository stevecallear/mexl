package vm_test

import (
	"testing"

	"github.com/stevecallear/mexl/compiler"
	"github.com/stevecallear/mexl/parser"
)

func TestInstructions_String(t *testing.T) {
	n, err := parser.New(`email ne null and email ew "@test.com"`).Parse()
	if err != nil {
		t.Fatalf("got %v, expected nil", err)
	}

	p, err := compiler.New().Compile(n)
	if err != nil {
		t.Fatalf("got %v, expected nil", err)
	}

	exp := `0000 OpFetch 0
0002 OpNull
0003 OpNotEqual
0004 OpJumpIfFalse 14
0007 OpFetch 0
0009 OpConstant 0
0012 OpEndsWith
0013 OpAnd
`

	act := p.Instructions.String()
	if act != exp {
		t.Errorf("got %s, expected %s", act, exp)
	}
}
