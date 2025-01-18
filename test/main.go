package main

import (
	"fmt"
	"os"

	"github.com/ever-eduardo/vida"
	"github.com/ever-eduardo/vida/lib"
)

func main() {
	clear()
	fmt.Println(vida.Name(), vida.Version())
	basePath := "./"
	modules, err := os.ReadDir(basePath)
	handleError(err, basePath)
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
	i, err := vida.NewInterpreter(modulePath, lib.Loadlibs())
	handleError(err, modulePath)
	r, err := i.MeasureRunTime()
	handleError(err, modulePath)
	fmt.Println(r)
}

func handleError(err error, modulePath string) {
	if err != nil {
		fmt.Println(err)
		fmt.Println(modulePath)
		os.Exit(0)
	}
}

func clear() {
	fmt.Printf("\u001B[H")
	fmt.Printf("\u001B[2J")
}
