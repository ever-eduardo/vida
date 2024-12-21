package main

import (
	"fmt"
	"os"

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
	i, err := vida.NewInterpreter(modulePath, stdlib.LoadStdlib())
	handleError(err)
	r, err := i.MeasureRunTime()
	handleError(err)
	fmt.Println(r)
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
