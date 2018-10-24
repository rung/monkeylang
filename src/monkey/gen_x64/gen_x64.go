package gen_x64

import (
	"bytes"
	"fmt"
	"monkey/code"
	"monkey/compiler"
	"monkey/object"
)

var reserveLabel = make(map[int]int)

type Gen struct {
	constants   []object.Object
	instraction code.Instructions
	Assembly    *bytes.Buffer
	Global      *bytes.Buffer

	global [30]int
	sp     int

	labelcnt int
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

		l, ok := reserveLabel[ip]
		if ok {
			fmt.Fprintf(g.Assembly, ".LABEL%d:\n", l)
		}

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
			g.sp--
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
		case code.OpEqual:
			// Trueだったら0, それ以外は0以外をpush
			fmt.Fprintln(g.Assembly, "	pop rbx")
			fmt.Fprintln(g.Assembly, "	pop rax")
			// cmp命令でZFを立てるのではなく、sub演算の結果をstackに積む
			fmt.Fprintln(g.Assembly, "	sub rax, rbx")
			fmt.Fprintln(g.Assembly, "	push rax")
			g.sp--
		case code.OpNotEqual:
			// Trueだったら0以外、それ以外は0をpush
			fmt.Fprintln(g.Assembly, "	pop rax")
			fmt.Fprintln(g.Assembly, "	pop rbx")
			// cmp命令でZFを立てるのではなく、sub演算の結果をstackに積む
			fmt.Fprintln(g.Assembly, "	cmp rax, rbx")
			// rax, rbxが一致しなかったら0をpush
			fmt.Fprintf(g.Assembly, "	jne .LABEL%d\n", g.labelcnt)
			fmt.Fprintln(g.Assembly, "	push 1")
			fmt.Fprintf(g.Assembly, "	jmp .LABEL%d\n", g.labelcnt+1)
			fmt.Fprintf(g.Assembly, ".LABEL%d:\n", g.labelcnt)
			fmt.Fprintln(g.Assembly, "	push 0")
			fmt.Fprintf(g.Assembly, ".LABEL%d:\n", g.labelcnt+1)
			g.labelcnt += 2
			g.sp--

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

		case code.OpNull:
			fmt.Fprintln(g.Assembly, "	push 0")

		case code.OpJump:
			fmt.Fprintf(g.Assembly, "	jmp .LABEL%d\n", g.labelcnt)
			bytecodeNo := int(code.ReadUint16(g.instraction[ip+1:]))

			pushLabel(g.labelcnt, bytecodeNo)

			g.labelcnt++
			ip += 2

		case code.OpJumpNotTruthy:
			fmt.Fprintln(g.Assembly, "	pop rax")
			fmt.Fprintln(g.Assembly, "	cmp rax, 0")
			fmt.Fprintf(g.Assembly, "	jne .LABEL%d\n", g.labelcnt)
			bytecodeNo := int(code.ReadUint16(g.instraction[ip+1:]))

			pushLabel(g.labelcnt, bytecodeNo)

			g.labelcnt++
			ip += 2

		case code.OpPop:
			fmt.Fprintln(g.Assembly, "	pop rax")

		default:
			return fmt.Errorf("non-supported opcode")
		}
	}

	return nil
}

// 指定したbytecodeのline noの箇所に、ラベルを吐く
func pushLabel(labelcnt, b_line int) {
	reserveLabel[b_line] = labelcnt
	return
}
