package vm

import (
	"fmt"
	"strings"

	"github.com/stevecallear/mexl/types"
)

const stackSize = 2048

type VM struct {
	program     *Program
	environment types.Map
	stack       []types.Object
	sp          int
}

var (
	objTrue  = &types.Boolean{Value: true}
	objFalse = &types.Boolean{Value: false}
	objNull  = &types.Null{}
)

func New(p *Program, env types.Map) *VM {
	return &VM{
		program:     p,
		environment: env,
		stack:       make([]types.Object, stackSize),
		sp:          0,
	}
}

func (vm *VM) Run() (types.Object, error) {
	if err := vm.run(); err != nil {
		return nil, err
	}
	return vm.stack[vm.sp-1], nil
}

func (vm *VM) run() (err error) {
	for i := 0; i < len(vm.program.Instructions); i++ {
		op := Opcode(vm.program.Instructions[i])

		switch op {
		case OpConstant:
			cidx := readUint16(vm.program.Instructions[i+1:])
			i += 2
			vm.push(vm.program.Constants[cidx])

		case OpArray:
			alen := readUint16(vm.program.Instructions[i+1:])
			i += 2
			vm.execArray(alen)

		case OpTrue:
			vm.push(objTrue)

		case OpFalse:
			vm.push(objFalse)

		case OpNull:
			vm.push(objNull)

		case OpAdd, OpSubtract, OpMultiply, OpDivide, OpModulus:
			if err = vm.execBinaryOp(op); err != nil {
				return err
			}

		case OpEqual, OpNotEqual:
			if err = vm.execEqualityComparison(op); err != nil {
				return err
			}

		case OpLess, OpLessOrEqual, OpGreater, OpGreaterOrEqual, OpAnd, OpOr, OpStartsWith, OpEndsWith:
			if err = vm.execComparison(op); err != nil {
				return err
			}

		case OpIn:
			if err = vm.execInOp(); err != nil {
				return err
			}

		case OpNot:
			vm.execBangOp()

		case OpMinus:
			if err = vm.execMinusOp(); err != nil {
				return err
			}

		case OpIndex:
			if err = vm.execIndexExpression(); err != nil {
				return err
			}

		case OpGlobal:
			idx := readUint8(vm.program.Instructions[i+1:])
			i++
			vm.execIdentifier(vm.program.Identifiers[idx])

		case OpMember:
			idx := readUint8(vm.program.Instructions[i+1:])
			i++

			if err = vm.execMemberExpression(int(idx)); err != nil {
				return err
			}

		case OpCall:
			nargs := readUint8(vm.program.Instructions[i+1:])
			i++

			if err = vm.execCallExpression(nargs); err != nil {
				return err
			}

		case OpJumpIfFalse, OpJumpIfTrue:
			pos := readUint16(vm.program.Instructions[i+1:])
			i += 2
			i = vm.execJump(op, i, int(pos))

		default:
			return fmt.Errorf("invalid opcode: %d", op)
		}
	}

	return nil
}

func (vm *VM) execBinaryOp(op Opcode) error {
	r := vm.pop()
	l := vm.pop()

	l, r = types.Coerce(l, r)
	lt, rt := l.Type(), r.Type()

	switch {
	case lt == types.TypeInteger && rt == types.TypeInteger:
		return vm.execBinaryIntegerOp(op, l, r)

	case lt == types.TypeFloat && rt == types.TypeFloat:
		return vm.execBinaryFloatOp(op, l, r)

	case lt == types.TypeString && rt == types.TypeString:
		return vm.execBinaryStringOp(op, l, r)

	default:
		return fmt.Errorf("unsupported type for binary operation: %s, %s", lt, rt)
	}
}

func (vm *VM) execBinaryIntegerOp(op Opcode, left, right types.Object) error {
	l := left.(*types.Integer).Value
	r := right.(*types.Integer).Value

	switch op {
	case OpAdd:
		vm.push(&types.Integer{Value: l + r})

	case OpSubtract:
		vm.push(&types.Integer{Value: l - r})

	case OpMultiply:
		vm.push(&types.Integer{Value: l * r})

	case OpDivide:
		switch l % r {
		case 0:
			vm.push(&types.Integer{Value: l / r})

		default:
			vm.push(&types.Float{Value: float64(l) / float64(r)})
		}

	case OpModulus:
		vm.push(&types.Integer{Value: l % r})

	default:
		return fmt.Errorf("unknown integer operator: %d", op)
	}

	return nil
}

func (vm *VM) execBinaryFloatOp(op Opcode, left, right types.Object) error {
	l := left.(*types.Float).Value
	r := right.(*types.Float).Value

	var v float64
	switch op {
	case OpAdd:
		v = l + r

	case OpSubtract:
		v = l - r

	case OpMultiply:
		v = l * r

	case OpDivide:
		v = l / r

	default:
		return fmt.Errorf("unknown float operator: %d", op)
	}

	vm.push(&types.Float{Value: v})
	return nil
}

func (vm *VM) execBinaryStringOp(op Opcode, left, right types.Object) error {
	l := left.(*types.String).Value
	r := right.(*types.String).Value

	if op != OpAdd {
		return fmt.Errorf("unknown string operator: %d", op)
	}

	vm.push(&types.String{Value: l + r})
	return nil
}

func (vm *VM) execEqualityComparison(op Opcode) error {
	r := vm.pop()
	l := vm.pop()

	l, r = types.Coerce(l, r)

	switch op {
	case OpEqual:
		vm.push(boolToObject(l.Equal(r)))

	case OpNotEqual:
		vm.push(boolToObject(!l.Equal(r)))

	default:
		return fmt.Errorf("unknown equality comparison operator: %d (%s %s)", op, l.Type(), r.Type())
	}

	return nil
}

func (vm *VM) execComparison(op Opcode) error {
	r := vm.pop()
	l := vm.pop()

	l, r = types.Coerce(l, r)
	lt, rt := l.Type(), r.Type()

	switch {
	case lt == types.TypeInteger && rt == types.TypeInteger:
		return vm.execIntegerComparison(op, l, r)

	case lt == types.TypeFloat && rt == types.TypeFloat:
		return vm.execFloatComparison(op, l, r)

	case lt == types.TypeBoolean && rt == types.TypeBoolean:
		return vm.execBoolComparison(op, l, r)

	case lt == types.TypeString && rt == types.TypeString:
		return vm.execStringComparison(op, l, r)

	default:
		return fmt.Errorf("unknown comparison operator: %d (%s %s)", op, l.Type(), r.Type())
	}
}

func (vm *VM) execIntegerComparison(op Opcode, left, right types.Object) error {
	l := left.(*types.Integer).Value
	r := right.(*types.Integer).Value

	switch op {
	case OpLess:
		vm.push(boolToObject(l < r))

	case OpLessOrEqual:
		vm.push(boolToObject(l <= r))

	case OpGreater:
		vm.push(boolToObject(l > r))

	case OpGreaterOrEqual:
		vm.push(boolToObject(l >= r))

	default:
		return fmt.Errorf("unknown integer comparison operator: %d", op)
	}

	return nil
}

func (vm *VM) execFloatComparison(op Opcode, left, right types.Object) error {
	l := left.(*types.Float).Value
	r := right.(*types.Float).Value

	switch op {
	case OpLess:
		vm.push(boolToObject(l < r))

	case OpLessOrEqual:
		vm.push(boolToObject(l <= r))

	case OpGreater:
		vm.push(boolToObject(l > r))

	case OpGreaterOrEqual:
		vm.push(boolToObject(l >= r))

	default:
		return fmt.Errorf("unknown float comparison operator: %d", op)
	}

	return nil
}

func (vm *VM) execBoolComparison(op Opcode, left, right types.Object) error {
	l := left.(*types.Boolean).Value
	r := right.(*types.Boolean).Value

	switch op {
	case OpAnd:
		vm.push(boolToObject(l && r))

	case OpOr:
		vm.push(boolToObject(l || r))

	default:
		return fmt.Errorf("unknown boolean comparison operator: %d", op)
	}

	return nil
}

func (vm *VM) execStringComparison(op Opcode, left, right types.Object) error {
	l := left.(*types.String).Value
	r := right.(*types.String).Value

	switch op {
	case OpStartsWith:
		vm.push(boolToObject(strings.HasPrefix(l, r)))

	case OpEndsWith:
		vm.push(boolToObject(strings.HasSuffix(l, r)))

	default:
		return fmt.Errorf("unknown string comparison operator: %d", op)
	}

	return nil
}

func (vm *VM) execBangOp() {
	o := vm.pop()
	switch o {
	case objTrue:
		vm.push(objFalse)

	case objFalse:
		vm.push(objTrue)

	default:
		vm.push(objFalse)
	}
}

func (vm *VM) execMinusOp() error {
	o := vm.pop()

	switch o.Type() {
	case types.TypeNull:
		o = objNull

	case types.TypeInteger:
		v := o.(*types.Integer).Value
		o = &types.Integer{Value: -v}

	case types.TypeFloat:
		v := o.(*types.Float).Value
		o = &types.Float{Value: -v}

	default:
		return fmt.Errorf("unsupported type for negation: %s", o.Type())
	}

	vm.push(o)
	return nil
}

func (vm *VM) execInOp() error {
	r := vm.pop()
	l := vm.pop()
	lt, rt := l.Type(), r.Type()

	switch {
	case lt == types.TypeNull || rt == types.TypeNull:
		vm.push(objFalse)

	case lt == types.TypeString && rt == types.TypeString:
		vm.execStringInOp(l, r)

	case rt == types.TypeArray:
		vm.execArrayInOp(l, r)

	default:
		return fmt.Errorf("unsupported types for in operation: %s IN %s", lt, rt)
	}

	return nil
}

func (vm *VM) execStringInOp(left, right types.Object) {
	l := left.(*types.String).Value
	r := right.(*types.String).Value
	vm.push(boolToObject(strings.Contains(r, l))) // left in right
}

func (vm *VM) execArrayInOp(left, right types.Object) {
	obj := right.(types.Array)
	for _, e := range obj {
		if e.Equal(left) {
			vm.push(objTrue)
			return
		}
	}
	vm.push(objFalse)
}

func (vm *VM) execIdentifier(n string) {
	if obj, ok := vm.environment[n]; ok {
		vm.push(obj)
		return
	}

	if obj, ok := builtIns[n]; ok {
		vm.push(obj)
		return
	}

	vm.push(objNull)
}

func (vm *VM) execArray(alen uint16) {
	a := make(types.Array, alen)

	for i := len(a) - 1; i >= 0; i-- {
		a[i] = vm.pop()
	}

	vm.push(a)
}

func (vm *VM) execIndexExpression() error {
	index := vm.pop()
	left := vm.pop()

	switch {
	case left.Type() == types.TypeArray && index.Type() == types.TypeInteger:
		a := left.(types.Array)
		i := index.(*types.Integer).Value
		vm.push(a[i])
	default:
		return fmt.Errorf("index operator not supported: %s", left.Type())
	}

	return nil
}

func (vm *VM) execMemberExpression(idx int) error {
	left := vm.pop()

	switch left.Type() {
	case types.TypeMap:
		if m, ok := left.(types.Map)[vm.program.Identifiers[idx]]; ok {
			vm.push(m)
		} else {
			vm.push(objNull)
		}

	case types.TypeNull:
		vm.push(objNull)

	default:
		return fmt.Errorf("container not supported: %s", left.Type())
	}

	return nil
}

func (vm *VM) execCallExpression(nargs uint8) error {
	args := make([]types.Object, int(nargs))

	for i := len(args) - 1; i >= 0; i-- {
		args[i] = vm.pop()
	}

	fn := vm.pop()

	switch fn.Type() {
	case types.TypeFunc:
		obj, err := fn.(types.Func)(args...)
		if err != nil {
			return err
		}

		vm.push(obj)
		return nil

	default:
		return fmt.Errorf("invalid function type: %T", fn)
	}
}

func (vm *VM) execJump(op Opcode, cpos, jpos int) int {
	cond := vm.peek().(*types.Boolean)

	if (op == OpJumpIfFalse && !cond.Value) || (op == OpJumpIfTrue && cond.Value) {
		cpos = jpos - 1
	}

	return cpos
}

func (vm *VM) push(o types.Object) {
	if vm.sp >= stackSize {
		panic("stack overflow")
	}

	vm.stack[vm.sp] = o
	vm.sp++
}

func (vm *VM) pop() types.Object {
	obj := vm.peek()
	vm.sp--
	return obj
}

func (vm *VM) peek() types.Object {
	obj := vm.stack[vm.sp-1]
	return obj
}

func boolToObject(b bool) *types.Boolean {
	if b {
		return objTrue
	}
	return objFalse
}
