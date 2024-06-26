package main

import (
	"fmt"
	"time"

	"github.com/ever-eduardo/vida"
)

func main() {
	println(vida.Name(), vida.Version())
	debug := false
	moduleName := "test.vida"
	src, e := vida.ReadFile(moduleName)
	if e != nil {
		fmt.Println(e)
		return
	}
	c := vida.NewCompiler(src, moduleName)
	if m, err := c.Compile(); err == nil {
		vm := vida.NewVM(m)
		if debug {
			r := vm.Debug()
			fmt.Println(r)
		} else {
			i := time.Now()
			r := vm.Run()
			fmt.Printf("Time %v\n", time.Since(i))
			fmt.Printf("Time %vs\n", time.Since(i).Seconds())
			fmt.Println(r)
			fmt.Println("Store")
			for k, v := range vm.Module.Store {
				fmt.Println(k, " : ", v)
			}
			fmt.Println("Konst")
			for i, v := range vm.Module.Konstants {
				fmt.Println(i, " : ", v)
			}
		}
	} else {
		fmt.Println(err)
	}
}
