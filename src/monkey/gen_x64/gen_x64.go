package gen_x64

import (
	"bytes"
	"fmt"
	"monkey/code"
	"monkey/compiler"
	"monkey/object"
)

type Gen struct {
	constants   []object.Object
	instraction code.Instructions
	assembly    *bytes.Buffer
}

func New(b *compiler.Bytecode) *Gen {
	g := &Gen{
		constants:   b.Constants,
		instraction: b.Instructions,
	}
	return g
}

func (g *Gen) Genx64() error {
	g.assembly = &bytes.Buffer{}

	// write header
	fmt.Fprintln(g.assembly, ".intel_syntax noprefix\n")
	fmt.Fprintln(g.assembly, ".text")
	fmt.Fprintln(g.assembly, ".global main")
	fmt.Fprintln(g.assembly, "main:")

	for ip := 0; ip < len(g.instraction); ip++ {
		op := code.Opcode(g.instraction[ip])
		switch op {
		case code.OpConstant:
			constIndex := code.ReadUint16(g.instraction[ip+1:])
			ip += 2

			obj := g.constants[constIndex]
			i := obj.(*object.Integer).Value
			fmt.Fprintf(g.assembly, "	push %d\n", i)
		case code.OpReturnValue:
			fmt.Fprintln(g.assembly, "	pop rax")
			fmt.Fprintln(g.assembly, "	ret")
		case code.OpAdd:
			fmt.Fprintln(g.assembly, "	pop rbx")
			fmt.Fprintln(g.assembly, "	pop rax")
			fmt.Fprintln(g.assembly, "	add rax, rbx")
			fmt.Fprintln(g.assembly, "	push rax")
		}
	}

	return nil
}
