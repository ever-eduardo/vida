package main

import (
	"fmt"

	"github.com/ever-eduardo/vida"
)

func main() {
	println(vida.Name(), vida.Version())
	module := "hello.vida"
	src, err := vida.ReadFile(module)
	if err != nil {
		fmt.Println(err)
		return
	}
	p := vida.NewParser(src, module)
	ast, err := p.Parse()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(ast)
}
