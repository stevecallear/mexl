package vm_test

import (
	"testing"

	"github.com/stevecallear/mexl/compiler"
	"github.com/stevecallear/mexl/parser"
	"github.com/stevecallear/mexl/vm"
)

func TestMake(t *testing.T) {
	t.Run("should return empty for invalid opcode", func(t *testing.T) {
		act := vm.Make(vm.OpInvalid)
		if len(act) != 0 {
			t.Errorf("got %v, expected empty", act)
		}
	})
}

func TestLookup(t *testing.T) {
	t.Run("should return error for invalid opcode", func(t *testing.T) {
		_, err := vm.Lookup(byte(vm.OpInvalid))
		if err == nil {
			t.Error("got nil, expected error")
		}
	})
}

func TestInstructions_String(t *testing.T) {
	n, err := parser.New(`email ne null and email ew "@test.com"`).Parse()
	if err != nil {
		t.Fatalf("got %v, expected nil", err)
	}

	p, err := compiler.New().Compile(n)
	if err != nil {
		t.Fatalf("got %v, expected nil", err)
	}

	exp := `0000 OpGlobal 0
0002 OpNull
0003 OpNotEqual
0004 OpJumpIfFalse 14
0007 OpGlobal 0
0009 OpConstant 0
0012 OpEndsWith
0013 OpAnd
`

	act := p.Instructions.String()
	if act != exp {
		t.Errorf("got %s, expected %s", act, exp)
	}
}
