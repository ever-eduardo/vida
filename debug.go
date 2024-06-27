package vida

import (
	"fmt"
)

func (vm VM) Inspect(ip int) {
	clear()
	fmt.Println("Running", vm.Module.Name)
	fmt.Println("Store")
	for k, v := range vm.Module.Store {
		fmt.Println(k, " : ", v)
	}
	fmt.Println("Konst")
	for i, v := range vm.Module.Konstants {
		fmt.Println(i, " : ", v)
	}
	fmt.Println("Code")
	for i, v := range vm.Module.Code {
		if i == ip {
			fmt.Printf("[%v : %v], ", i, v)
		} else {
			fmt.Printf("%v, ", v)
		}
	}
	fmt.Println()
	fmt.Println("Stack")
	for _, v := range vm.Stack {
		fmt.Printf("%v, ", v)
	}
	fmt.Println()
	fmt.Printf("Press 'Enter' to continue => ")
	fmt.Scanf(" ")
}

func (vm VM) Debug() (Result, error) {
	frame := &vm.Frames[vm.fp]
	frame.code = vm.Module.Code
	return Success, nil
}

func clear() {
	fmt.Printf("\u001B[H")
	fmt.Printf("\u001B[2J")
}
