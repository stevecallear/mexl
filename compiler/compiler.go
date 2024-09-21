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
		c.emit(vm.OpGlobal, c.addIdentifier(node.Value))

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

	default:
		return fmt.Errorf("invalid ast node: %T", n)
	}

	return nil
}

func (c *Compiler) compileExpressions(exprs []ast.Node) error {
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
	var op, jop vm.Opcode

	switch n.Operator {
	case "+":
		op = vm.OpAdd

	case "-":
		op = vm.OpSubtract

	case "*":
		op = vm.OpMultiply

	case "/":
		op = vm.OpDivide

	case "%":
		op = vm.OpModulus

	case "eq", "==":
		op = vm.OpEqual

	case "ne", "!=":
		op = vm.OpNotEqual

	case "lt", "<":
		op = vm.OpLess

	case "le", "<=":
		op = vm.OpLessOrEqual

	case "gt", ">":
		op = vm.OpGreater

	case "ge", ">=":
		op = vm.OpGreaterOrEqual

	case "and", "&&":
		op = vm.OpAnd
		jop = vm.OpJumpIfFalse

	case "or", "||":
		op = vm.OpOr
		jop = vm.OpJumpIfTrue

	case "sw":
		op = vm.OpStartsWith

	case "ew":
		op = vm.OpEndsWith

	case "in":
		op = vm.OpIn

	default:
		return fmt.Errorf("unknown infix operator: %s", n.Operator)
	}

	return func() error {
		if err = c.compile(n.Left); err != nil {
			return err
		}

		if jop > 0 {
			jpos := c.emit(jop, jumpPlaceholder)
			defer c.patchJump(jpos)
		}

		if err = c.compile(n.Right); err != nil {
			return err
		}

		c.emit(op)
		return nil
	}()
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
