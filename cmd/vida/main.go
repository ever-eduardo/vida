package main

import (
	"fmt"
	"os"
	"time"

	"github.com/ever-eduardo/vida"
	"github.com/ever-eduardo/vida/ast"
)

func main() {
	clear()
	println(vida.Name(), vida.Version())
	debug := true
	module := "sketchpad.vida"
	if debug {
		debugPath(module)
	} else {
		normalPath(module)
	}
}

func debugPath(module string) {
	var src []byte
	var err error
	src, err = vida.ReadModule(module)
	handleError(err)
	p := vida.NewParser(src, module)
	r, err := p.Parse()
	handleError(err)
	fmt.Println(ast.PrintAST(r))
	fmt.Scanf(" ")
	c := vida.NewCompiler(r, module)
	m := c.CompileModule()
	fmt.Println(vida.PrintBytecode(m, m.Name))
	fmt.Scanf(" ")
	vm, err := vida.NewVM(m)
	handleError(err)
	res, err := vm.Debug()
	handleError(err)
	fmt.Println(res)
}

func normalPath(module string) {
	var src []byte
	var err error
	init := time.Now()
	src, err = vida.ReadModule(module)
	handleError(err)
	p := vida.NewParser(src, module)
	r, err := p.Parse()
	handleError(err)
	c := vida.NewCompiler(r, module)
	m := c.CompileModule()
	fmt.Printf("Compiler time = %vs\n", time.Since(init).Seconds())
	fmt.Printf("Compiler time = %v\n", time.Since(init))
	init = time.Now()
	vm, err := vida.NewVM(m)
	handleError(err)
	res, err := vm.Run()
	fmt.Printf("VM time = %vs\n", time.Since(init).Seconds())
	fmt.Printf("VM time = %v\n", time.Since(init))
	handleError(err)
	fmt.Println(res)
	fmt.Printf("Result %v\n", vm.Module.Store["r"])
}

func handleError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
}

func clear() {
	fmt.Printf("\u001B[H")
	fmt.Printf("\u001B[2J")
}
