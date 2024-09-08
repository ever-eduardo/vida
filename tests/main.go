package main

import (
	"fmt"
	"os"
	"time"

	"github.com/ever-eduardo/vida"
	"github.com/ever-eduardo/vida/stdlib"
)

func main() {
	clear()
	fmt.Println(vida.Name(), vida.Version())
	basePath := "./"
	modules, err := os.ReadDir(basePath)
	handleError(err)
	count := 0
	for _, v := range modules {
		if !v.IsDir() && v.Name() != "main.go" && v.Name() != "tests.exe" {
			count++
			fmt.Printf("Running file '%v'\n", v.Name())
			executeModule(v.Name())
			fmt.Printf("\n\n\n")
		}
	}
	fmt.Printf("All %v tests were ok!\n\n\n", count)
}

func executeModule(modulePath string) {
	var src []byte
	var err error
	init := time.Now()
	src, err = vida.ReadModule(modulePath)
	handleError(err)
	p := vida.NewParser(src, modulePath)
	rAst, err := p.Parse()
	handleError(err)
	c := vida.NewCompiler(rAst, modulePath)
	m, e := c.CompileModule()
	handleError(e)
	fmt.Printf("Compiler time = %vs\n", time.Since(init).Seconds())
	fmt.Printf("Compiler time = %v\n", time.Since(init))
	init = time.Now()
	vm, err := vida.NewVM(m, stdlib.LoadStdlib())
	handleError(err)
	res, err := vm.Run()
	fmt.Printf("VM time = %vs\n", time.Since(init).Seconds())
	fmt.Printf("VM time = %v\n", time.Since(init))
	handleError(err)
	fmt.Println(res)
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
