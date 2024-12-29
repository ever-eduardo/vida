package main

import (
	"fmt"
	"os"

	"github.com/ever-eduardo/vida"
	"github.com/ever-eduardo/vida/stdlib"
)

func main() {
	// f, err := os.Create("vida.prof")
	// handleError(err)
	// pprof.StartCPUProfile(f)
	// defer pprof.StopCPUProfile()
	clear()
	println(vida.Name(), vida.Version())
	debug := false
	ast := false
	code := false
	module := "sketchpad.vida"
	if debug {
		switch {
		case ast:
			printAST(module)
		case code:
			printMachineCode(module)
		default:
			debugPath(module)
		}
	} else {
		measuredPath(module)
	}
}

func debugPath(modulePath string) {
	i, err := vida.NewDebugger(modulePath, stdlib.LoadStdlib())
	handleError(err)
	r, err := i.Debug()
	handleError(err)
	fmt.Println(r)
}

func measuredPath(modulePath string) {
	i, err := vida.NewInterpreter(modulePath, stdlib.LoadStdlib())
	handleError(err)
	r, err := i.MeasureRunTime()
	if err != nil {
		i.PrintCallStack()
		handleError(err)
	}
	fmt.Println(r)
}

func printAST(modulePath string) {
	err := vida.PrintAST(modulePath)
	handleError(err)
}

func printMachineCode(modulePath string) {
	err := vida.PrintMachineCode(modulePath)
	handleError(err)
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
