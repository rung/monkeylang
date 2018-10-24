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
	Assembly    *bytes.Buffer
}

func New(b *compiler.Bytecode) *Gen {
	g := &Gen{
		constants:   b.Constants,
		instraction: b.Instructions,
	}
	return g
}

func (g *Gen) Genx64() error {
	g.Assembly = &bytes.Buffer{}

	// write header
	fmt.Fprintln(g.Assembly, ".intel_syntax noprefix\n")
	fmt.Fprintln(g.Assembly, ".text")
	fmt.Fprintln(g.Assembly, ".global main")
	fmt.Fprintln(g.Assembly, "main:")

	for ip := 0; ip < len(g.instraction); ip++ {
		op := code.Opcode(g.instraction[ip])
		switch op {
		case code.OpConstant:
			constIndex := code.ReadUint16(g.instraction[ip+1:])
			ip += 2

			obj := g.constants[constIndex]
			i := obj.(*object.Integer).Value
			fmt.Fprintf(g.Assembly, "	push %d\n", i)
		case code.OpReturnValue:
			fmt.Fprintln(g.Assembly, "	pop rax")
			fmt.Fprintln(g.Assembly, "	ret")
		case code.OpAdd:
			fmt.Fprintln(g.Assembly, "	pop rbx")
			fmt.Fprintln(g.Assembly, "	pop rax")
			fmt.Fprintln(g.Assembly, "	add rax, rbx")
			fmt.Fprintln(g.Assembly, "	push rax")
		case code.OpSub:
			fmt.Fprintln(g.Assembly, "	pop rbx")
			fmt.Fprintln(g.Assembly, "	pop rax")
			fmt.Fprintln(g.Assembly, "	sub rax, rbx")
			fmt.Fprintln(g.Assembly, "	push rax")
		case code.OpMul:
			fmt.Fprintln(g.Assembly, "	pop rbx")
			fmt.Fprintln(g.Assembly, "	pop rax")
			fmt.Fprintln(g.Assembly, "	imul rbx")
			fmt.Fprintln(g.Assembly, "	push rax")
		case code.OpDiv:
			fmt.Fprintln(g.Assembly, "	pop rbx")
			fmt.Fprintln(g.Assembly, "	pop rax")
			fmt.Fprintln(g.Assembly, "	cdq")
			fmt.Fprintln(g.Assembly, "	idiv rbx")
			fmt.Fprintln(g.Assembly, "	push rax")
		}
	}

	return nil
}
