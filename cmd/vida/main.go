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
	debug := true
	module := "sketchpad.vida"
	if debug {
		debugPath(module)
	} else {
		normalPath(module)
	}
}

func debugPath(modulePath string) {
	i, err := vida.NewDebugger(modulePath, stdlib.LoadStdlib())
	handleError(err)
	r, err := i.Debug()
	handleError(err)
	fmt.Println(r)
}

func normalPath(modulePath string) {
	i, err := vida.NewInterpreter(modulePath, stdlib.LoadStdlib())
	handleError(err)
	r, err := i.Measure()
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
