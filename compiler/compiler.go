package compiler

import (
	"fmt"
	"slices"

	"github.com/stevecallear/mexl/ast"
	"github.com/stevecallear/mexl/types"
	"github.com/stevecallear/mexl/vm"
)

type Compiler struct {
	instructions vm.Instructions
	constants    []types.Object
	identifiers  []string
}

const jumpPlaceholder = 9999

func New() *Compiler {
	return &Compiler{
		instructions: vm.Instructions{},
		constants:    []types.Object{},
		identifiers:  []string{},
	}
}

func (c *Compiler) Compile(n ast.Node) (*vm.Program, error) {
	if err := c.compile(n); err != nil {
		return nil, err
	}

	return &vm.Program{
		Instructions: c.instructions,
		Constants:    c.constants,
		Identifiers:  c.identifiers,
	}, nil
}

func (c *Compiler) compile(n ast.Node) (err error) {
	switch node := n.(type) {
	case *ast.ExpressionStatement:
		if err = c.compile(node.Expression); err != nil {
			return err
		}
		c.emit(vm.OpPop)

	case *ast.InfixExpression:
		if err = c.compileInfixExpression(node); err != nil {
			return err
		}

	case *ast.PrefixExpression:
		if err = c.compilePrefixExpression(node); err != nil {
			return err
		}

	case *ast.Boolean:
		if node.Value {
			c.emit(vm.OpTrue)
		} else {
			c.emit(vm.OpFalse)
		}

	case *ast.Null:
		c.emit(vm.OpNull)
	case *ast.IntegerLiteral:
		obj := &types.Integer{Value: node.Value}
		c.emit(vm.OpConstant, c.addConstant(obj))

	case *ast.FloatLiteral:
		obj := &types.Float{Value: node.Value}
		c.emit(vm.OpConstant, c.addConstant(obj))

	case *ast.StringLiteral:
		obj := &types.String{Value: node.Value}
		c.emit(vm.OpConstant, c.addConstant(obj))

	case *ast.ArrayLiteral:
		if err = c.compileExpressions(node.Elements); err != nil {
			return err
		}
		c.emit(vm.OpArray, len(node.Elements))

	case *ast.Identifier:
		c.emit(vm.OpFetch, c.addIdentifier(node.Value))

	case *ast.IndexExpression:
		if err = c.compile(node.Left); err != nil {
			return err
		}
		if err = c.compile(node.Index); err != nil {
			return err
		}
		c.emit(vm.OpIndex)

	case *ast.MemberExpression:
		if err = c.compileMemberExpression(node); err != nil {
			return err
		}

	case *ast.CallExpression:
		if err = c.compile(node.Function); err != nil {
			return err
		}
		if err = c.compileExpressions(node.Arguments); err != nil {
			return err
		}
		c.emit(vm.OpCall, len(node.Arguments))

	}

	return nil
}

func (c *Compiler) compileExpressions(exprs []ast.Expression) error {
	for _, expr := range exprs {
		if err := c.compile(expr); err != nil {
			return err
		}
	}
	return nil
}

func (c *Compiler) compilePrefixExpression(n *ast.PrefixExpression) (err error) {
	if err = c.compile(n.Right); err != nil {
		return err
	}

	switch n.Operator {
	case "not", "!":
		c.emit(vm.OpNot)

	case "-":
		c.emit(vm.OpMinus)

	default:
		return fmt.Errorf("unknown prefix operator: %s", n.Operator)
	}

	return nil
}

func (c *Compiler) compileInfixExpression(n *ast.InfixExpression) (err error) {
	compile := func() error {
		if err = c.compile(n.Left); err != nil {
			return err
		}
		return c.compile(n.Right)
	}

	switch n.Operator {
	case "+":
		err = compile()
		c.emit(vm.OpAdd)

	case "-":
		err = compile()
		c.emit(vm.OpSubtract)

	case "*":
		err = compile()
		c.emit(vm.OpMultiply)

	case "/":
		err = compile()
		c.emit(vm.OpDivide)

	case "%":
		err = compile()
		c.emit(vm.OpModulus)

	case "eq", "==":
		err = compile()
		c.emit(vm.OpEqual)

	case "ne", "!=":
		err = compile()
		c.emit(vm.OpNotEqual)

	case "lt", "<":
		err = compile()
		c.emit(vm.OpLess)

	case "le", "<=":
		err = compile()
		c.emit(vm.OpLessOrEqual)

	case "gt", ">":
		err = compile()
		c.emit(vm.OpGreater)

	case "ge", ">=":
		err = compile()
		c.emit(vm.OpGreaterOrEqual)

	case "and", "&&":
		if err = c.compile(n.Left); err != nil {
			return err
		}
		jpos := c.emit(vm.OpJumpIfFalse, jumpPlaceholder)
		if err = c.compile(n.Right); err != nil {
			return err
		}
		c.emit(vm.OpAnd)
		c.patchJump(jpos)

	case "or", "||":
		if err = c.compile(n.Left); err != nil {
			return err
		}
		jpos := c.emit(vm.OpJumpIfTrue, jumpPlaceholder)
		if err = c.compile(n.Right); err != nil {
			return err
		}
		c.emit(vm.OpOr)
		c.patchJump(jpos)

	case "sw":
		err = compile()
		c.emit(vm.OpStartsWith)

	case "ew":
		err = compile()
		c.emit(vm.OpEndsWith)

	case "in":
		err = compile()
		c.emit(vm.OpIn)

	default:
		err = fmt.Errorf("unknown infix operator: %s", n.Operator)
	}

	return err
}

func (c *Compiler) compileMemberExpression(n *ast.MemberExpression) error {
	if err := c.compile(n.Left); err != nil {
		return err
	}

	ident, ok := n.Member.(*ast.Identifier)
	if !ok {
		return fmt.Errorf("invalid member type: %T", n.Member)
	}

	c.emit(vm.OpMember, c.addIdentifier(ident.Value))
	return nil
}

func (c *Compiler) patchJump(pos int) {
	op := vm.Opcode(c.instructions[pos])
	ins := vm.Make(op, len(c.instructions))

	for i := range ins {
		c.instructions[pos+i] = ins[i]
	}

}

func (c *Compiler) emit(op vm.Opcode, operands ...int) int {
	pos := len(c.instructions)

	ins := vm.Make(op, operands...)
	c.instructions = append(c.instructions, ins...)

	return pos
}

func (c *Compiler) addConstant(obj types.Object) int {
	c.constants = append(c.constants, obj)
	return len(c.constants) - 1
}

func (c *Compiler) addIdentifier(s string) int {
	if i := slices.Index(c.identifiers, s); i >= 0 {
		return i
	}
	c.identifiers = append(c.identifiers, s)
	return len(c.identifiers) - 1
}
