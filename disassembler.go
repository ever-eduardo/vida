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
	var counter int
	sb.Write(code[:4])
	sb.WriteRune(32)
	sb.WriteString(fmt.Sprintf("version %v.%v.%v", int(code[4]), int(code[5]), int(code[6])))
	sb.WriteRune(10)
	sb.WriteRune(10)
	ip := 8
	var op byte
	for {
		counter++
		op = code[ip]
		sb.WriteRune(10)
		sb.WriteString(fmt.Sprintf("  %3v  [%3v]  ", counter, ip))
		sb.WriteString(fmt.Sprintf("%7v", opcodes[op]))
		ip++
		switch op {
		case setG:
			sb.WriteRune(32)
			sb.WriteString(fmt.Sprintf("%3v", int(code[ip])))
			ip++
			sb.WriteRune(32)
			sb.WriteString(fmt.Sprintf("%3v", binary.NativeEndian.Uint16(code[ip:])))
			ip += 2
			sb.WriteRune(32)
			sb.WriteString(fmt.Sprintf("%3v", binary.NativeEndian.Uint16(code[ip:])))
			ip += 2
		case setL:
			sb.WriteRune(32)
			sb.WriteString(fmt.Sprintf("%3v", int(code[ip])))
			ip++
			sb.WriteRune(32)
			sb.WriteString(fmt.Sprintf("%3v", binary.NativeEndian.Uint16(code[ip:])))
			ip += 2
			sb.WriteRune(32)
			sb.WriteString(fmt.Sprintf("%3v", int(code[ip])))
			ip++
		case move:
			sb.WriteRune(32)
			sb.WriteString(fmt.Sprintf("%3v", int(code[ip])))
			ip++
			sb.WriteRune(32)
			sb.WriteString(fmt.Sprintf("%3v", int(code[ip])))
			ip++
		case prefix:
			sb.WriteRune(32)
			sb.WriteString(fmt.Sprintf("%3v", int(code[ip])))
			ip++
			sb.WriteRune(32)
			sb.WriteString(fmt.Sprintf("%3v", int(code[ip])))
			ip++
			sb.WriteRune(32)
			sb.WriteString(fmt.Sprintf("%3v", binary.NativeEndian.Uint16(code[ip:])))
			ip += 2
			sb.WriteRune(32)
			sb.WriteString(fmt.Sprintf("%3v", int(code[ip])))
			ip++
		case binop:
			sb.WriteRune(32)
			sb.WriteString(fmt.Sprintf("%3v", int(code[ip])))
			ip++
			sb.WriteRune(32)
			sb.WriteString(fmt.Sprintf("%3v", int(code[ip])))
			ip++
			sb.WriteRune(32)
			sb.WriteString(fmt.Sprintf("%3v", int(code[ip])))
			ip++
			sb.WriteRune(32)
			sb.WriteString(fmt.Sprintf("%3v", binary.NativeEndian.Uint16(code[ip:])))
			ip += 2
			sb.WriteRune(32)
			sb.WriteString(fmt.Sprintf("%3v", binary.NativeEndian.Uint16(code[ip:])))
			ip += 2
			sb.WriteRune(32)
			sb.WriteString(fmt.Sprintf("%3v", int(code[ip])))
			ip++
		case list:
			sb.WriteRune(32)
			sb.WriteString(fmt.Sprintf("%3v", int(code[ip])))
			ip++
			sb.WriteRune(32)
			sb.WriteString(fmt.Sprintf("%3v", int(code[ip])))
			ip++
			sb.WriteRune(32)
			sb.WriteString(fmt.Sprintf("%3v", int(code[ip])))
			ip++
		case record:
			sb.WriteRune(32)
			sb.WriteString(fmt.Sprintf("%3v", int(code[ip])))
			ip++
			sb.WriteRune(32)
			sb.WriteString(fmt.Sprintf("%3v", int(code[ip])))
			ip++
			sb.WriteRune(32)
			sb.WriteString(fmt.Sprintf("%3v", int(code[ip])))
			ip++
		case iGet:
			sb.WriteRune(32)
			sb.WriteString(fmt.Sprintf("%3v", int(code[ip])))
			ip++
			sb.WriteRune(32)
			sb.WriteString(fmt.Sprintf("%3v", int(code[ip])))
			ip++
			sb.WriteRune(32)
			sb.WriteString(fmt.Sprintf("%3v", binary.NativeEndian.Uint16(code[ip:])))
			ip += 2
			sb.WriteRune(32)
			sb.WriteString(fmt.Sprintf("%3v", binary.NativeEndian.Uint16(code[ip:])))
			ip += 2
			sb.WriteRune(32)
			sb.WriteString(fmt.Sprintf("%3v", int(code[ip])))
			ip++
		case slice:
			sb.WriteRune(32)
			sb.WriteString(fmt.Sprintf("%3v", int(code[ip])))
			ip++
			sb.WriteRune(32)
			sb.WriteString(fmt.Sprintf("%3v", int(code[ip])))
			ip++
			sb.WriteRune(32)
			sb.WriteString(fmt.Sprintf("%3v", int(code[ip])))
			ip++
			sb.WriteRune(32)
			sb.WriteString(fmt.Sprintf("%3v", int(code[ip])))
			ip++
			sb.WriteRune(32)
			sb.WriteString(fmt.Sprintf("%3v", binary.NativeEndian.Uint16(code[ip:])))
			ip += 2
			sb.WriteRune(32)
			sb.WriteString(fmt.Sprintf("%3v", binary.NativeEndian.Uint16(code[ip:])))
			ip += 2
			sb.WriteRune(32)
			sb.WriteString(fmt.Sprintf("%3v", binary.NativeEndian.Uint16(code[ip:])))
			ip += 2
			sb.WriteRune(32)
			sb.WriteString(fmt.Sprintf("%3v", int(code[ip])))
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
		sb.WriteString(fmt.Sprintf("  %3v  [%3v]  %v\n", i+1, i, v))
	}
}
