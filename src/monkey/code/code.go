package code

import (
	"encoding/binary"
	"fmt"
)

type Instructions []byte

type Opcode byte

const (
	OpConstant Opcode = iota
)

type Definition struct {
	Name          string
	OperandWidths []int
}

var definitions = map[Opcode]*Definition{
	OpConstant: {"OpConstant", []int{2}},
}

func Lookup(op byte) (*Definition, error) {
	def, ok := definitions[Opcode(op)]
	if !ok {
		return nil, fmt.Errorf("opcode %d undefined", op)
	}

	return def, nil
}

func Make(op Opcode, operands ...int) []byte {
	def, ok := definitions[op]
	if !ok {
		return []byte{}
	}

	instructionLen := 1
	for _, w := range def.OperandWidths {
		instructionLen += w
	}

	// instruction {Opcode(1byte), argument1(1byte), argument2(1byte)}
	instruction := make([]byte, instructionLen)
	instruction[0] = byte(op)

	offset := 1
	// instructionに順番に値を突っ込んでいく(big endianで値は渡される)
	for i, o := range operands {
		width := def.OperandWidths[i]
		switch width {
		case 2:
			// PutUint16 put uint16 to []byte
			binary.BigEndian.PutUint16(instruction[offset:], uint16(o))
		}
		offset += width
	}

	return instruction
}
