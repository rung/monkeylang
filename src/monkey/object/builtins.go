package object

import "fmt"

var Builtins = []struct {
	Name     string
	Builtin  *Builtin
	Assembly string
}{
	{
		Name: "len",
		Builtin: &Builtin{
			Fn: func(args ...Object) Object {
				if len(args) != 1 {
					return newError("wrong number of arguments. got=%d, want=1",
						len(args))
				}

				switch arg := args[0].(type) {
				case *String:
					return &Integer{Value: int64(len(arg.Value))}
				case *Array:
					return &Integer{Value: int64(len(arg.Elements))}
				default:
					return newError("argument to `len` not supported, got %s",
						args[0].Type())
				}
			},
		},
	},
	{
		Name: "puts",
		Builtin: &Builtin{
			Fn: func(args ...Object) Object {
				for _, arg := range args {
					fmt.Println(arg.Inspect())
				}
				return nil
			},
		},
		Assembly: `.global puts
puts:
	#header
	push rbp
	mov rbp, rsp

	#save register
	push rdi
	push rsi
	push rdx
	
	#strlen start
	xor rcx, rcx
	sub rcx, 1

	xor rax, rax
	lea rdx, [rbp+16]
	mov bl, [rdx]
	cmp bl, 0
	je .L1
.L0:
	add rdx, 1
	mov bl, [rdx]
	cmp bl, 0
	loopne .L0
.L1:
	not rcx	
	mov rdx, rcx
	# strlen end

	# write(1, "string", strlen) // printf
	mov rax, 1
	mov rdi, 1
	mov rsi, rbp
	add rsi, 16
	syscall
	push rax

	#footar
	mov rsp, rbp
	pop rbp
	ret
`,
	},
	{
		Name: "first",
		Builtin: &Builtin{
			Fn: func(args ...Object) Object {
				if len(args) != 1 {
					return newError("wrong number of arguments. got=%d, want=1",
						len(args))
				}
				if args[0].Type() != ARRAY_OBJ {
					return newError("argument to `first` must be ARRAY, got %s",
						args[0].Type())
				}

				arr := args[0].(*Array)
				if len(arr.Elements) > 0 {
					return arr.Elements[0]
				}
				return nil
			},
		},
	},
	{
		Name: "last",
		Builtin: &Builtin{
			Fn: func(args ...Object) Object {
				if len(args) != 1 {
					return newError("wrong number of arguments. got=%d, want=1",
						len(args))
				}
				if args[0].Type() != ARRAY_OBJ {
					return newError("argument to `last` must be ARRAY, got %s",
						args[0].Type())
				}

				arr := args[0].(*Array)
				length := len(arr.Elements)
				if len(arr.Elements) > 0 {
					return arr.Elements[length-1]
				}
				return nil
			},
		},
	},
	{
		Name: "rest",
		Builtin: &Builtin{
			Fn: func(args ...Object) Object {
				if len(args) != 1 {
					return newError("wrong number of arguments. got=%d, want=1",
						len(args))
				}
				if args[0].Type() != ARRAY_OBJ {
					return newError("argument to `rest` must be ARRAY, got %s",
						args[0].Type())
				}

				arr := args[0].(*Array)
				length := len(arr.Elements)
				if length > 0 {
					newElements := make([]Object, length-1, length-1)
					copy(newElements, arr.Elements[1:])
					return &Array{Elements: newElements}
				}

				return nil
			},
		},
	}, {
		Name: "push",
		Builtin: &Builtin{
			Fn: func(args ...Object) Object {
				if len(args) != 2 {
					return newError("wrong number of arguments. got=%d, want=2",
						len(args))
				}
				if args[0].Type() != ARRAY_OBJ {
					return newError("argument to `push` must be ARRAY, got %s",
						args[0].Type())
				}

				arr := args[0].(*Array)
				length := len(arr.Elements)

				newElements := make([]Object, length+1, length+1)
				copy(newElements, arr.Elements)
				newElements[length] = args[1]

				return &Array{Elements: newElements}

			},
		},
	},
}

func newError(format string, a ...interface{}) *Error {
	return &Error{Message: fmt.Sprintf(format, a...)}
}

func GetBuiltinByName(name string) *Builtin {
	for _, def := range Builtins {
		if def.Name == name {
			return def.Builtin
		}
	}
	return nil
}
