package vida

import (
	"fmt"
	"strings"
)

func PrintBytecode(m *Module, moduleName string) string {
	clear()
	fmt.Println("Bytecode for module", moduleName)
	var sb strings.Builder
	sb.WriteString(printHeader(m))
	var counter int
	var ip int = 8
	var s string
	for ip < len(m.MainFunction.CoreFn.Code) {
		s, ip, counter = printInstr(ip, m.MainFunction.CoreFn.Code, counter, false)
		sb.WriteString(s)
	}
	for idx, v := range m.Konstants {

		if f, ok := v.(*CoreFunction); ok {
			ip = 0
			counter = 0
			sb.WriteString(fmt.Sprintf("\n\nFunction %v/%v/%v", idx, f.Arity, f.Free))
			for ip < len(f.Code) {
				s, ip, counter = printInstr(ip, f.Code, counter, false)
				sb.WriteString(s)
			}
		}
	}
	sb.WriteString(printKonstants(m.Konstants))
	return sb.String()
}

func printHeader(m *Module) string {
	var sb strings.Builder
	sb.WriteRune(32)
	sb.WriteString(fmt.Sprintf("version %v.%v.%v", int(m.MainFunction.CoreFn.Code[4]), int(m.MainFunction.CoreFn.Code[5]), int(m.MainFunction.CoreFn.Code[6])))
	sb.WriteRune(10)
	sb.WriteRune(10)
	sb.WriteString("Main\n")
	return sb.String()
}

func printKonstants(konst []Value) string {
	var sb strings.Builder
	sb.WriteString("\n\n\nKonstants\n")
	for i, v := range konst {
		sb.WriteString(fmt.Sprintf("  %4v  [%4v]  %v\n", i+1, i, v))
	}
	return sb.String()
}

func printInstr(ip int, code []uint32, counter int, isRunningDebug bool) (string, int, int) {
	var sb strings.Builder
	var op uint32
	if !isRunningDebug {
		counter++
		op = code[ip]
		sb.WriteRune(10)
		sb.WriteString(fmt.Sprintf("  %4v  [%4v]  ", counter, ip))
		sb.WriteString(fmt.Sprintf("%7v", opcodes[op]))
	} else {
		op = code[ip]
		sb.WriteString(opcodes[op])
	}
	ip++
	switch op {
	case end:
		return sb.String(), ip, counter
	}
	return sb.String(), ip, counter
}
