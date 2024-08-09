package vm

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type (
	Opcode byte

	Definition struct {
		Name          string
		OperandWidths []int
	}
)

const (
	OpInvalid Opcode = iota
	OpConstant
	OpArray
	OpAdd
	OpSubtract
	OpMultiply
	OpDivide
	OpModulus
	OpPop
	OpTrue
	OpFalse
	OpNull
	OpMinus
	OpNot
	OpAnd
	OpOr
	OpEqual
	OpNotEqual
	OpLess
	OpLessOrEqual
	OpGreater
	OpGreaterOrEqual
	OpStartsWith
	OpEndsWith
	OpIn
	OpIndex
	OpFetch
	OpMember
	OpCall
	OpJumpIfTrue
	OpJumpIfFalse
)

var definitions = map[Opcode]*Definition{
	OpConstant:       {"OpConstant", []int{2}},
	OpArray:          {"OpArray", []int{2}},
	OpAdd:            {"OpAdd", []int{}},
	OpSubtract:       {"OpSubtract", []int{}},
	OpMultiply:       {"OpMultiply", []int{}},
	OpDivide:         {"OpDivide", []int{}},
	OpModulus:        {"OpModulus", []int{}},
	OpPop:            {"OpPop", []int{}},
	OpTrue:           {"OpTrue", []int{}},
	OpFalse:          {"OpFalse", []int{}},
	OpNull:           {"OpNull", []int{}},
	OpMinus:          {"OpMinus", []int{}},
	OpNot:            {"OpNot", []int{}},
	OpAnd:            {"OpAnd", []int{}},
	OpOr:             {"OpOr", []int{}},
	OpEqual:          {"OpEqual", []int{}},
	OpNotEqual:       {"OpNotEqual", []int{}},
	OpLess:           {"OpLess", []int{}},
	OpLessOrEqual:    {"OpLessOrEqual", []int{}},
	OpGreater:        {"OpGreater", []int{}},
	OpGreaterOrEqual: {"OpGreaterOrEqual", []int{}},
	OpStartsWith:     {"OpStartsWith", []int{}},
	OpEndsWith:       {"OpEndsWith", []int{}},
	OpIn:             {"OpIn", []int{}},
	OpIndex:          {"OpIndex", []int{}},
	OpFetch:          {"OpFetch", []int{1}},
	OpMember:         {"OpMember", []int{1}},
	OpCall:           {"OpCall", []int{1}},
	OpJumpIfTrue:     {"OpJumpIfTrue", []int{2}},
	OpJumpIfFalse:    {"OpJumpIfFalse", []int{2}},
}

func Make(op Opcode, operands ...int) []byte {
	d, ok := definitions[op]
	if !ok {
		return []byte{}
	}

	instructionLen := 1
	for _, w := range d.OperandWidths {
		instructionLen += w
	}

	instruction := make([]byte, instructionLen)
	instruction[0] = byte(op)

	offset := 1
	for i, o := range operands {
		width := d.OperandWidths[i]
		switch width {
		case 1:
			instruction[offset] = byte(o)
		case 2:
			binary.BigEndian.PutUint16(instruction[offset:], uint16(o))
		}
		offset += width
	}

	return instruction
}

func Lookup(op byte) (*Definition, error) {
	d, ok := definitions[Opcode(op)]
	if !ok {
		return nil, fmt.Errorf("opcode %d undefined", op)
	}
	return d, nil
}

func readOperands(def *Definition, ins Instructions) ([]int, int) {
	operands := make([]int, len(def.OperandWidths))
	offset := 0

	for i, width := range def.OperandWidths {
		switch width {
		case 1:
			operands[i] = int(readUint8(ins[offset:]))
		case 2:
			operands[i] = int(readUint16(ins[offset:]))
		}
		offset += width
	}

	return operands, offset
}

func readUint8(ins Instructions) uint8 {
	return uint8(ins[0])
}

func readUint16(ins Instructions) uint16 {
	return binary.BigEndian.Uint16(ins)
}

func (ins Instructions) String() string {
	var out bytes.Buffer

	i := 0
	for i < len(ins) {
		def, err := Lookup(ins[i])
		if err != nil {
			fmt.Fprintf(&out, "ERROR: %s\n", err.Error())
			continue
		}

		operands, n := readOperands(def, ins[i+1:])
		fmt.Fprintf(&out, "%04d %s\n", i, formatInstruction(def, operands))

		i += 1 + n
	}

	return out.String()
}

func formatInstruction(def *Definition, operands []int) string {
	count := len(def.OperandWidths)

	if len(operands) != count {
		return fmt.Sprintf("ERROR: operand count mismatch: %d, %d", len(operands), count)
	}

	switch count {
	case 0:
		return def.Name
	case 1:
		return fmt.Sprintf("%s %d", def.Name, operands[0])
	default:
		return fmt.Sprintf("ERROR: unhandled operand count: %d", count)
	}
}
