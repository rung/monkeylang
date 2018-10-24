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
	Global      *bytes.Buffer

	global [30]int
	sp     int
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

	fmt.Fprintln(g.Assembly, "	mov rbp, rsp")

	for ip := 0; ip < len(g.instraction); ip++ {
		op := code.Opcode(g.instraction[ip])
		switch op {
		case code.OpConstant:
			constIndex := code.ReadUint16(g.instraction[ip+1:])
			ip += 2

			obj := g.constants[constIndex]
			i := obj.(*object.Integer).Value
			fmt.Fprintf(g.Assembly, "	push %d\n", i)
			g.sp++
		case code.OpReturnValue:
			fmt.Fprintln(g.Assembly, "	pop rax")
			fmt.Fprintln(g.Assembly, "	mov rsp, rbp")
			fmt.Fprintln(g.Assembly, "	ret")
			g.sp = 0
		case code.OpAdd:
			fmt.Fprintln(g.Assembly, "	pop rbx")
			fmt.Fprintln(g.Assembly, "	pop rax")
			fmt.Fprintln(g.Assembly, "	add rax, rbx")
			fmt.Fprintln(g.Assembly, "	push rax")
			g.sp--
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
			g.sp--
		case code.OpDiv:
			fmt.Fprintln(g.Assembly, "	pop rbx")
			fmt.Fprintln(g.Assembly, "	pop rax")
			fmt.Fprintln(g.Assembly, "	cdq")
			fmt.Fprintln(g.Assembly, "	idiv rbx")
			fmt.Fprintln(g.Assembly, "	push rax")
			g.sp--
		case code.OpMinus:
			fmt.Fprintln(g.Assembly, "	mov rax, 0")
			fmt.Fprintln(g.Assembly, "	pop rbx")
			fmt.Fprintln(g.Assembly, "	sub rax, rbx")
			fmt.Fprintln(g.Assembly, "	push rax")

		case code.OpSetGlobal:
			globalIndex := code.ReadUint16(g.instraction[ip+1:])
			ip += 2
			g.global[globalIndex] = g.sp

		case code.OpGetGlobal:
			globalIndex := code.ReadUint16(g.instraction[ip+1:])
			ip += 2
			sp := g.global[globalIndex]

			fmt.Fprintf(g.Assembly, "	mov rax, [rbp-%d]\n", sp*8)
			fmt.Fprintln(g.Assembly, "	push rax")
		}
	}

	return nil
}
