package main

import (
	"fmt"
	"os"

	"github.com/alkemist-17/vida"
	"github.com/alkemist-17/vida/lib"
)

func main() {
	clear()
	fmt.Println(vida.Name(), vida.Version())
	basePath := "./"
	scripts, err := os.ReadDir(basePath)
	handleError(err, basePath)
	count := 0
	for _, v := range scripts {
		if !v.IsDir() && v.Name() != "main.go" && v.Name() != "tests.exe" {
			count++
			fmt.Printf("Running file '%v'\n", v.Name())
			executeScript(v.Name())
			fmt.Printf("\n\n\n")
		}
	}
	fmt.Printf("All %v tests were ok!\n\n\n", count)
}

func executeScript(path string) {
	i, err := vida.NewInterpreter(path, lib.Loadlibs())
	handleError(err, path)
	r, err := i.MeasureRunTime()
	handleError(err, path)
	fmt.Println(r)
}

func handleError(err error, path string) {
	if err != nil {
		fmt.Println(err)
		fmt.Println(path)
		os.Exit(0)
	}
}

func clear() {
	fmt.Printf("\u001B[H")
	fmt.Printf("\u001B[2J")
}
