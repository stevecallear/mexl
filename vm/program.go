package vm

import "github.com/stevecallear/mexl/types"

type (
	Instructions []byte

	Program struct {
		Instructions Instructions
		Constants    []types.Object
		Identifiers  []string
	}
)
