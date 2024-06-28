package main

import (
	"fmt"

	"github.com/ever-eduardo/vida"
	"github.com/ever-eduardo/vida/ast"
)

func main() {
	println(vida.Name(), vida.Version())
	debug := true
	module := "sketchpad.vida"
	testModule := "../../tests/setSK.vida"
	var src []byte
	var err error
	if debug {
		src, err = vida.ReadFile(module)
	} else {
		src, err = vida.ReadFile(testModule)
	}
	if err != nil {
		fmt.Println(err)
		return
	}
	p := vida.NewParser(src, module)
	r, err := p.Parse()
	if err != nil {
		fmt.Println(err)
		return
	}
	if debug {
		fmt.Println(ast.PrintAST(r))
		fmt.Scanf(" ")
	}
	c := vida.NewCompiler(r, module)
	m := c.CompileModule()
	vm, vmerr := vida.NewVM(m)
	if vmerr != nil {
		fmt.Println(vmerr)
		return
	}
	var res vida.Result
	var rterr error
	if debug {
		res, rterr = vm.Debug()
	} else {
		res, rterr = vm.Run()
	}
	if rterr != nil {
		fmt.Println(rterr)
		return
	}
	fmt.Println(res)
}
