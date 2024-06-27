package main

import (
	"fmt"

	"github.com/ever-eduardo/vida"
	"github.com/ever-eduardo/vida/ast"
)

func main() {
	println(vida.Name(), vida.Version())
	module := "sketchpad.vida"
	src, err := vida.ReadFile(module)
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
	fmt.Println(ast.PrintAST(r))
}
