package gen_x64

import (
	"fmt"
	"monkey/ast"
	"monkey/compiler"
	"monkey/lexer"
	"monkey/parser"
	"os"
	"os/exec"
	"syscall"
	"testing"
)

type testCase struct {
	input    string
	expected int
}

func TestGenerator(t *testing.T) {
	tests := []testCase{
		{
			input:    `return 1`,
			expected: 1,
		},
		{
			input:    `return 1 + 1`,
			expected: 2,
		},
		{
			input:    `return 5 - 1`,
			expected: 4,
		},
		{
			input:    `return 5 * 3`,
			expected: 15,
		},
		{
			input:    `return 9 / 3`,
			expected: 3,
		},
		{
			input:    `return 10 / 3`,
			expected: 3,
		},
		{
			input:    `let a = 3; return a;`,
			expected: 3,
		},
		{
			input:    `let a = 3; let b = 1; return a;`,
			expected: 3,
		},
		{
			input:    `let a = 2; let b = 5; return b;`,
			expected: 5,
		},
	}

	for _, tt := range tests {
		// parse
		program := parse(tt.input)

		// compile(to bytecode)
		comp := compiler.New()
		err := comp.Compile(program)
		if err != nil {
			t.Fatalf("compiler error: %s", err)
		}

		// compile(x86 code generation)
		g := New(comp.Bytecode())
		err = g.Genx64()
		if err != nil {
			t.Errorf("code generation error: %s", err)
		}

		// write tmp file
		os.Remove("/tmp/monkeytmp.s")
		f, err := os.OpenFile("/tmp/monkeytmp.s", os.O_RDWR|os.O_CREATE, 0644)
		if err != nil {
			t.Errorf(err.Error())
		}
		f.Write(g.Assembly.Bytes())
		f.Close()

		// change Assembly to machine code and link.
		cmd := exec.Command("/usr/bin/gcc", "/tmp/monkeytmp.s", "-o", "/tmp/monkeytmp")
		err = cmd.Run()
		if err != nil {
			fmt.Println(g.instraction)
			fmt.Println(g.Assembly.String())
			t.Errorf("gcc error")
		}
		cmd = exec.Command("/tmp/monkeytmp")
		err = cmd.Run()
		returncode := -1

		if err != nil {
			if s, ok := err.(*exec.ExitError).Sys().(syscall.WaitStatus); ok {
				returncode = s.ExitStatus()
			} else {
				fmt.Println(g.Assembly.String())
				t.Errorf("can't get return code")
			}
		} else {
			returncode = 1
		}

		if returncode != tt.expected {
			fmt.Println(tt.input)
			fmt.Println(g.instraction)
			fmt.Println(g.Assembly.String())
			t.Errorf("return code is different got=%d, expected=%d", returncode, tt.expected)
		}

		// debug
		//fmt.Println("======================")
		//fmt.Println(tt.input)
		//fmt.Println("---")
		//fmt.Println(g.instraction)
		//fmt.Println("---")
		//fmt.Println(g.Assembly.String())
		//fmt.Println("======================")

	}

	// delete gabages
	os.Remove("/tmp/mokeytmp")
}

func parse(input string) *ast.Program {
	l := lexer.New(input)
	p := parser.New(l)
	return p.ParseProgram()
}
