package main

import (
	"fmt"
	"io/ioutil"
	"monkey/ast"
	"monkey/compiler"
	"monkey/gen_x64"
	"monkey/lexer"
	"monkey/parser"
	"os"
)

func main() {

	fp, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}
	input, err := ioutil.ReadAll(fp)
	if err != nil {
		panic(err)
	}
	// parse
	program := parse(string(input))

	// compile(to bytecode)
	comp := compiler.New()
	err = comp.Compile(program)
	if err != nil {
		panic("compiler error")
	}

	// compile(x86 code generation)
	g := gen_x64.New(comp.Bytecode())
	err = g.Genx64()
	if err != nil {
		panic("code generation error")
	}

	// write to standard output
	fmt.Println(g.Assembly.String())
}

func parse(input string) *ast.Program {
	l := lexer.New(input)
	p := parser.New(l)
	return p.ParseProgram()
}
