package main

import (
	"fmt"
	"time"

	"github.com/ever-eduardo/vida"
)

func main() {
	println(vida.Name(), vida.Version())
	debug := true
	moduleName := "test.vida"
	src, e := vida.ReadFile(moduleName)
	if e != nil {
		fmt.Println(e)
		return
	}
	c := vida.NewCompiler(src, moduleName)
	if m, err := c.Compile(); err == nil {
		vm, abierr := vida.NewVM(m)
		if abierr != nil {
			fmt.Println(abierr)
			return
		}
		if debug {
			if r, vmerr := vm.Debug(); vmerr == nil {
				fmt.Println(r)
			} else {
				fmt.Println(r, vmerr)
			}
		} else {
			i := time.Now()
			if r, vmerr := vm.Run(); vmerr == nil {
				s := time.Since(i)
				fmt.Println(r)
				fmt.Printf("Time %v\n", s)
				fmt.Printf("Time %vs\n", s.Seconds())
				fmt.Println("Store")
				for k, v := range vm.Module.Store {
					fmt.Println(k, " : ", v)
				}
				fmt.Println("Konst")
				for i, v := range vm.Module.Konstants {
					fmt.Println(i, " : ", v)
				}
			} else {
				fmt.Println(r, vmerr)
			}
		}
	} else {
		fmt.Println(err)
	}
}
