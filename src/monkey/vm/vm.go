package vm

import (
	"fmt"
	"monkey/code"
	"monkey/compiler"
	"monkey/object"
)

const StackSize = 2048

type VM struct {
	constants    []object.Object
	instructions code.Instructions

	stack []object.Object
	// stack pointer
	//   Always points to the next value. Top of stack is stack[sp-1]
	//   A new element would be stored at stack[sp],
	sp int
}

func New(bytecode *compiler.Bytecode) *VM {
	return &VM{
		instructions: bytecode.Instructions,
		constants:    bytecode.Constants,

		stack: make([]object.Object, StackSize),
		sp:    0,
	}
}

func (vm *VM) StackTop() object.Object {
	if vm.sp == 0 {
		return nil
	}

	return vm.stack[vm.sp-1]
}

func (vm *VM) Run() error {
	// fetch cycle
	for ip := 0; ip < len(vm.instructions); ip++ {
		// byte to opcode (decode)
		op := code.Opcode(vm.instructions[ip])

		switch op {
		case code.OpConstant:
			// Constant(定数)のインデックスがOpConstsntのOprandには入っている
			constIndex := code.ReadUint16(vm.instructions[ip+1:])
			ip += 2

			// 取り出したインデックスのObjectをStackにpushする
			err := vm.push(vm.constants[constIndex])
			if err != nil {
				return err
			}
		case code.OpAdd:
			right := vm.pop()
			left := vm.pop()
			leftValue := left.(*object.Integer).Value
			rightValue := right.(*object.Integer).Value

			result := leftValue + rightValue
			vm.push(&object.Integer{Value: result})
		case code.OpPop:
			vm.pop()
		}
	}
	return nil
}

func (vm *VM) push(o object.Object) error {
	if vm.sp >= StackSize {
		return fmt.Errorf("stack overflow")
	}

	vm.stack[vm.sp] = o
	vm.sp++

	return nil
}

// Memo: spが0を下回ったときのエラーハンドリングも入れたほうがよい
func (vm *VM) pop() object.Object {
	o := vm.stack[vm.sp-1]
	vm.sp--
	return o
}

// PopしたあとにPop前の一番上のstackを取りだす (for test)
func (vm *VM) LastPoppedStackElem() object.Object {
	return vm.stack[vm.sp]
}
