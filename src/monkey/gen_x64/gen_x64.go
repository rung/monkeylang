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
	constants []object.Object
	//Global    *bytes.Buffer

	labelcnt int

	frame  []*Frame
	fcnt   int
	fIndex map[int]int
}

type Frame struct {
	instraction code.Instructions
	Assembly    *bytes.Buffer
	symbolnum   int
	paramNum    int
}

func (g *Gen) currentFrame() *Frame {
	return g.frame[g.fcnt]
}

func (g *Gen) pushFrame(obj *object.CompiledFunction, constIndex, paramNum int) {
	f := &Frame{
		instraction: obj.Instructions,
		symbolnum:   obj.NumLocals,
		paramNum:    paramNum,
	}
	g.fcnt++
	g.fIndex[constIndex] = g.fcnt
	g.frame = append(g.frame, f)
}

func (g *Gen) popFrame() *Frame {
	g.fcnt--
	return g.frame[g.fcnt]
}

func New(b *compiler.Bytecode) *Gen {
	f := &Frame{
		instraction: b.Instructions,
		symbolnum:   b.SymbolNum,
		paramNum:    0,
	}

	g := &Gen{
		constants: b.Constants,
		frame:     []*Frame{f},
		fcnt:      0,
		fIndex:    make(map[int]int),
	}
	return g
}

func (g *Gen) Assembly() *bytes.Buffer {
	b := &bytes.Buffer{}

	// write header
	fmt.Fprintln(b, ".intel_syntax noprefix\n")
	fmt.Fprintln(b, ".text")

	// write function
	for i := 0; i <= g.fcnt; i++ {
		fmt.Fprintln(b, g.frame[i].Assembly.String())
	}

	return b
}

func (g *Gen) Genx64() error {
	cf := g.currentFrame()
	cf.Assembly = &bytes.Buffer{}
	currentFCnt := g.fcnt

	if currentFCnt == 0 {
		fmt.Fprintln(cf.Assembly, ".global main")
		fmt.Fprintln(cf.Assembly, "main:")
	} else {
		fmt.Fprintf(cf.Assembly, ".global function%d\n", currentFCnt)
		fmt.Fprintf(cf.Assembly, "function%d:\n", currentFCnt)
	}

	fmt.Fprintln(cf.Assembly, "	push rbp")
	fmt.Fprintln(cf.Assembly, "	mov rbp, rsp")

	// treat parameter
	//  x64 and Monkey VM ABI are different.
	//  Monkey: Parameter is used as local binding.
	//   [Mokey Stack Layout]
	//     Local binding1 (as local1)
	//     ---------------
	//     Argument1 (as local0)
	//     —------------- <- base pointer
	//     Previous function's stack
	//
	//   x64: Parameter is below of Return Pointer
	//     Local binding1
	//     -----------------
	//     Previous Base pointer
	//     -----------------
	//     Return Poiner
	//     -----------------
	//     Argument1

	//   so, out x64 compiler's leyout is below
	//     Local binding1
	//     --------------------
	//     Argument1 <- copy from below's Argument1
	//	   -----------------
	//	   Previous Base pointer
	//     -----------------
	//     Return Poiner
	//     -----------------
	//	   Argument1
	//
	// in x64, we move Argument1 to below of Local bindings. bacause bytecode requires it.
	paramNum := cf.paramNum
	// move Argument to above of Return pointer.
	n := 16 + (8 * (paramNum - 1))
	for i := 0; i < paramNum; i++ {
		n -= (8 * i)
		fmt.Fprintf(cf.Assembly, "	push [rbp+%d]\n", n)
	}

	// 変数分を先に引いておく
	//  symbolnum contains paramNum, so have to sub cf.paramNum
	fmt.Fprintf(cf.Assembly, "	sub rsp, %d\n", (cf.symbolnum-cf.paramNum)*8)

	for ip := 0; ip < len(cf.instraction); ip++ {
		op := code.Opcode(cf.instraction[ip])

		l, ok := reserveLabel[ip]
		if ok {
			fmt.Fprintf(cf.Assembly, ".LABEL%d:\n", l)
		}

		switch op {
		case code.OpConstant:
			constIndex := code.ReadUint16(cf.instraction[ip+1:])
			ip += 2

			obj := g.constants[constIndex]
			i := obj.(*object.Integer).Value
			fmt.Fprintf(cf.Assembly, "	push %d\n", i)

		case code.OpReturnValue:
			fmt.Fprintln(cf.Assembly, "	pop rax")
			fmt.Fprintln(cf.Assembly, "	mov rsp, rbp")
			fmt.Fprintln(cf.Assembly, "	pop rbp")
			fmt.Fprintln(cf.Assembly, "	ret")

		case code.OpAdd:
			fmt.Fprintln(cf.Assembly, "	pop rbx")
			fmt.Fprintln(cf.Assembly, "	pop rax")
			fmt.Fprintln(cf.Assembly, "	add rax, rbx")
			fmt.Fprintln(cf.Assembly, "	push rax")
		case code.OpSub:
			fmt.Fprintln(cf.Assembly, "	pop rbx")
			fmt.Fprintln(cf.Assembly, "	pop rax")
			fmt.Fprintln(cf.Assembly, "	sub rax, rbx")
			fmt.Fprintln(cf.Assembly, "	push rax")
		case code.OpMul:
			fmt.Fprintln(cf.Assembly, "	pop rbx")
			fmt.Fprintln(cf.Assembly, "	pop rax")
			fmt.Fprintln(cf.Assembly, "	imul rbx")
			fmt.Fprintln(cf.Assembly, "	push rax")
		case code.OpDiv:
			fmt.Fprintln(cf.Assembly, "	pop rbx")
			fmt.Fprintln(cf.Assembly, "	pop rax")
			fmt.Fprintln(cf.Assembly, "	cdq")
			fmt.Fprintln(cf.Assembly, "	idiv rbx")
			fmt.Fprintln(cf.Assembly, "	push rax")
		case code.OpMinus:
			fmt.Fprintln(cf.Assembly, "	mov rax, 0")
			fmt.Fprintln(cf.Assembly, "	pop rbx")
			fmt.Fprintln(cf.Assembly, "	sub rax, rbx")
			fmt.Fprintln(cf.Assembly, "	push rax")
		case code.OpEqual:
			// Trueだったら0, それ以外は0以外をpush
			fmt.Fprintln(cf.Assembly, "	pop rbx")
			fmt.Fprintln(cf.Assembly, "	pop rax")
			// cmp命令でZFを立てるのではなく、sub演算の結果をstackに積む
			fmt.Fprintln(cf.Assembly, "	sub rax, rbx")
			fmt.Fprintln(cf.Assembly, "	push rax")
		case code.OpNotEqual:
			// Trueだったら0以外、それ以外は0をpush
			fmt.Fprintln(cf.Assembly, "	pop rax")
			fmt.Fprintln(cf.Assembly, "	pop rbx")
			// cmp命令でZFを立てるのではなく、sub演算の結果をstackに積む
			fmt.Fprintln(cf.Assembly, "	cmp rax, rbx")
			// rax, rbxが一致しなかったら0をpush
			fmt.Fprintf(cf.Assembly, "	jne .LABEL%d\n", g.labelcnt)
			fmt.Fprintln(cf.Assembly, "	push 1")
			fmt.Fprintf(cf.Assembly, "	jmp .LABEL%d\n", g.labelcnt+1)
			fmt.Fprintf(cf.Assembly, ".LABEL%d:\n", g.labelcnt)
			fmt.Fprintln(cf.Assembly, "	push 0")
			fmt.Fprintf(cf.Assembly, ".LABEL%d:\n", g.labelcnt+1)
			g.labelcnt += 2

		case code.OpGreaterThan:
			fmt.Fprintln(cf.Assembly, "	pop rbx")
			fmt.Fprintln(cf.Assembly, "	pop rax")
			fmt.Fprintln(cf.Assembly, "	cmp rax, rbx")
			fmt.Fprintf(cf.Assembly, "	jle .LABEL%d\n", g.labelcnt)
			fmt.Fprintln(cf.Assembly, "	push 0")
			fmt.Fprintf(cf.Assembly, "	jmp .LABEL%d\n", g.labelcnt+1)
			fmt.Fprintf(cf.Assembly, ".LABEL%d:\n", g.labelcnt)
			fmt.Fprintln(cf.Assembly, "	push 1")
			fmt.Fprintf(cf.Assembly, ".LABEL%d:\n", g.labelcnt+1)
			g.labelcnt += 2

		case code.OpSetGlobal:
			globalIndex := code.ReadUint16(cf.instraction[ip+1:])
			ip += 2
			fmt.Fprintln(cf.Assembly, "	pop rax")
			fmt.Fprintf(cf.Assembly, "	mov [rbp-%d] ,rax\n", (globalIndex+1)*8)

		case code.OpSetLocal:
			globalIndex := code.ReadUint8(cf.instraction[ip+1:])
			ip += 1
			fmt.Fprintln(cf.Assembly, "	pop rax")
			fmt.Fprintf(cf.Assembly, "	mov [rbp-%d] ,rax\n", (globalIndex+1)*8)

		case code.OpGetGlobal:
			globalIndex := code.ReadUint16(cf.instraction[ip+1:])
			ip += 2

			fmt.Fprintf(cf.Assembly, "	mov rax, [rbp-%d]\n", (globalIndex+1)*8)
			fmt.Fprintln(cf.Assembly, "	push rax")

		case code.OpGetLocal:
			globalIndex := code.ReadUint8(cf.instraction[ip+1:])
			ip += 1

			fmt.Fprintf(cf.Assembly, "	mov rax, [rbp-%d]\n", (globalIndex+1)*8)
			fmt.Fprintln(cf.Assembly, "	push rax")

		case code.OpNull:
			fmt.Fprintln(cf.Assembly, "	push 0")

		case code.OpJump:
			fmt.Fprintf(cf.Assembly, "	jmp .LABEL%d\n", g.labelcnt)
			bytecodeNo := int(code.ReadUint16(cf.instraction[ip+1:]))

			pushLabel(g.labelcnt, bytecodeNo)

			g.labelcnt++
			ip += 2

		case code.OpJumpNotTruthy:
			fmt.Fprintln(cf.Assembly, "	pop rax")
			fmt.Fprintln(cf.Assembly, "	cmp rax, 0")
			fmt.Fprintf(cf.Assembly, "	jne .LABEL%d\n", g.labelcnt)
			bytecodeNo := int(code.ReadUint16(cf.instraction[ip+1:]))

			pushLabel(g.labelcnt, bytecodeNo)

			g.labelcnt++
			ip += 2

		case code.OpPop:
			fmt.Fprintln(cf.Assembly, "	pop rax")

		case code.OpClosure:
			constIndex := code.ReadUint16(cf.instraction[ip+1:])
			numFree := code.ReadUint8(cf.instraction[ip+3:])
			ip += 3

			err := g.pushClosure(int(constIndex), int(numFree))
			if err != nil {
				return nil
			}
			fmt.Fprintf(cf.Assembly, "	lea rax, function%d[rip]\n", g.fcnt)
			fmt.Fprintln(cf.Assembly, "	push rax")

		case code.OpCall:
			paramNum := code.ReadUint8(cf.instraction[ip+1:])
			ip += 1

			fmt.Fprintf(cf.Assembly, "	mov rax, [rsp+%d]\n", paramNum*8)
			fmt.Fprintln(cf.Assembly, "	call rax")
			fmt.Fprintf(cf.Assembly, "	add rsp, %d\n", 8+paramNum*8) // pop paramNum * 8
			fmt.Fprintln(cf.Assembly, "	push rax")

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

func (g *Gen) pushClosure(constIndex int, numFree int) error {
	constant := g.constants[constIndex]
	function, ok := constant.(*object.CompiledFunction)
	if !ok {
		return fmt.Errorf("not a function: %+v", constant)
	}
	paramNum := function.NumParameters

	err := g.writeFunction(function, constIndex, paramNum)
	if err != nil {
		fmt.Errorf("writing function error: %+v", err)
	}

	return nil
}

func (g *Gen) writeFunction(f *object.CompiledFunction, constIndex, paramNum int) error {
	g.pushFrame(f, constIndex, paramNum)
	err := g.Genx64()
	return err
}
