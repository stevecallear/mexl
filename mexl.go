package mexl

import (
	"github.com/stevecallear/mexl/compiler"
	"github.com/stevecallear/mexl/parser"
	"github.com/stevecallear/mexl/types"
	"github.com/stevecallear/mexl/vm"
)

// Eval compiles and runs the input
func Eval(input string, env map[string]any) (any, error) {
	p, err := Compile(input)
	if err != nil {
		return nil, err
	}
	return Run(p, env)
}

// Compile compiles the input
func Compile(input string) (*vm.Program, error) {
	n, err := parser.New(input).Parse()
	if err != nil {
		return nil, err
	}

	return compiler.New().Compile(n)

}

// Run runs the compiled program
func Run(p *vm.Program, env map[string]any) (any, error) {
	m, err := types.ToMap(env)
	if err != nil {
		return nil, err
	}

	out, err := vm.New(p, m).Run()
	if err != nil {
		return nil, err
	}

	return types.ToNative(out)
}
