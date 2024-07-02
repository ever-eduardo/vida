package vida

import (
	"encoding/binary"
	"fmt"
	"strings"
)

func PrintBytecode(m *Module, moduleName string) string {
	clear()
	fmt.Println("Bytecode for module", moduleName)
	return processBytecode(m.Code, m.Konstants)
}

func processBytecode(code []byte, konst []Value) string {
	var sb strings.Builder
	sb.Write(code[:4])
	sb.WriteRune(32)
	sb.WriteString(fmt.Sprintf("version %v.%v.%v", int(code[4]), int(code[5]), int(code[6])))
	sb.WriteRune(10)
	sb.WriteRune(10)
	ip := 8
	var op byte
	for {
		op = code[ip]
		sb.WriteRune(10)
		sb.WriteString(fmt.Sprintf(" [%3v]  ", ip))
		sb.WriteString(fmt.Sprintf("%7v", opcodes[op]))
		ip++
		switch op {
		case setG:
			sb.WriteRune(32)
			sb.WriteString(fmt.Sprintf("%v", int(code[ip])))
			ip++
			sb.WriteRune(32)
			sb.WriteString(fmt.Sprintf("%v", binary.NativeEndian.Uint16(code[ip:])))
			ip += 2
			sb.WriteRune(32)
			sb.WriteString(fmt.Sprintf("%v", binary.NativeEndian.Uint16(code[ip:])))
			ip += 2
		case setL:
			sb.WriteRune(32)
			sb.WriteString(fmt.Sprintf("%v", int(code[ip])))
			ip++
			sb.WriteRune(32)
			sb.WriteString(fmt.Sprintf("%v", binary.NativeEndian.Uint16(code[ip:])))
			ip += 2
			sb.WriteRune(32)
			sb.WriteString(fmt.Sprintf("%v", int(code[ip])))
			ip++
		case move:
			sb.WriteRune(32)
			sb.WriteString(fmt.Sprintf("%v", int(code[ip])))
			ip++
			sb.WriteRune(32)
			sb.WriteString(fmt.Sprintf("%v", int(code[ip])))
			ip++
		case not:
			sb.WriteRune(32)
			sb.WriteString(fmt.Sprintf("%v", int(code[ip])))
			ip++
			sb.WriteRune(32)
			sb.WriteString(fmt.Sprintf("%v", binary.NativeEndian.Uint16(code[ip:])))
			ip += 2
			sb.WriteRune(32)
			sb.WriteString(fmt.Sprintf("%v", int(code[ip])))
			ip++
		case end:
			sb.WriteRune(10)
			sb.WriteRune(10)
			sb.WriteRune(10)
			printKonstants(konst, &sb)
			return sb.String()
		}
	}
}

func printKonstants(konst []Value, sb *strings.Builder) {
	sb.WriteString("Konstants\n")
	for i, v := range konst {
		sb.WriteString(fmt.Sprintf(" [%3v]  %v\n", i, v))
	}
}
