package vida

import (
	"encoding/binary"
	"fmt"
	"strings"

	"github.com/ever-eduardo/vida/token"
)

func PrintBytecode(m *Module, moduleName string) string {
	clear()
	fmt.Println("Bytecode for module", moduleName)
	var sb strings.Builder
	sb.WriteString(printHeader(m))
	var counter int
	var ip int = 8
	var s string
	for ip < len(m.Function.Code) {
		s, ip, counter = printInstr(ip, m.Function.Code, counter, false)
		sb.WriteString(s)
	}
	for idx, v := range m.Konstants {

		if f, ok := v.(*Function); ok {
			ip = 0
			counter = 0
			sb.WriteString(fmt.Sprintf("\n\nFunction %v", idx))
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
	sb.Write(m.Function.Code[:4])
	sb.WriteRune(32)
	sb.WriteString(fmt.Sprintf("version %v.%v.%v", int(m.Function.Code[4]), int(m.Function.Code[5]), int(m.Function.Code[6])))
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

func printInstr(ip int, code []byte, counter int, isRunningDebug bool) (string, int, int) {
	var sb strings.Builder
	var op byte
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
	case setG:
		sb.WriteRune(32)
		sb.WriteString(fmt.Sprintf("%4v", int(code[ip])))
		ip++
		sb.WriteRune(32)
		sb.WriteString(fmt.Sprintf("%4v", binary.NativeEndian.Uint16(code[ip:])))
		ip += 2
		sb.WriteRune(32)
		sb.WriteString(fmt.Sprintf("%4v", binary.NativeEndian.Uint16(code[ip:])))
		ip += 2
	case setL:
		sb.WriteRune(32)
		sb.WriteString(fmt.Sprintf("%4v", int(code[ip])))
		ip++
		sb.WriteRune(32)
		sb.WriteString(fmt.Sprintf("%4v", binary.NativeEndian.Uint16(code[ip:])))
		ip += 2
		sb.WriteRune(32)
		sb.WriteString(fmt.Sprintf("%4v", int(code[ip])))
		ip++
	case setF:
		sb.WriteRune(32)
		sb.WriteString(fmt.Sprintf("%4v", int(code[ip])))
		ip++
		sb.WriteRune(32)
		sb.WriteString(fmt.Sprintf("%4v", binary.NativeEndian.Uint16(code[ip:])))
		ip += 2
		sb.WriteRune(32)
		sb.WriteString(fmt.Sprintf("%4v", binary.NativeEndian.Uint16(code[ip:])))
		ip += 2
	case move:
		sb.WriteRune(32)
		sb.WriteString(fmt.Sprintf("%4v", int(code[ip])))
		ip++
		sb.WriteRune(32)
		sb.WriteString(fmt.Sprintf("%4v", int(code[ip])))
		ip++
	case prefix:
		sb.WriteRune(32)
		sb.WriteString(fmt.Sprintf("%4v", int(code[ip])))
		ip++
		sb.WriteRune(32)
		sb.WriteString(fmt.Sprintf("%4v", int(code[ip])))
		ip++
		sb.WriteRune(32)
		sb.WriteString(fmt.Sprintf("%4v", binary.NativeEndian.Uint16(code[ip:])))
		ip += 2
		sb.WriteRune(32)
		sb.WriteString(fmt.Sprintf("%4v", int(code[ip])))
		ip++
	case equals:
		sb.WriteRune(32)
		sb.WriteString(fmt.Sprintf("%4v", int(code[ip])))
		ip++
		sb.WriteRune(32)
		sb.WriteString(fmt.Sprintf("%4v", int(code[ip])))
		ip++
		sb.WriteRune(32)
		sb.WriteString(fmt.Sprintf("%4v", int(code[ip])))
		ip++
		sb.WriteRune(32)
		sb.WriteString(fmt.Sprintf("%4v", binary.NativeEndian.Uint16(code[ip:])))
		ip += 2
		sb.WriteRune(32)
		sb.WriteString(fmt.Sprintf("%4v", binary.NativeEndian.Uint16(code[ip:])))
		ip += 2
		sb.WriteRune(32)
		sb.WriteString(fmt.Sprintf("%4v", int(code[ip])))
		ip++
	case binop:
		sb.WriteRune(32)
		sb.WriteString(fmt.Sprintf("%4v", token.Tokens[int(code[ip])]))
		ip++
		sb.WriteRune(32)
		sb.WriteString(fmt.Sprintf("%4v", int(code[ip])))
		ip++
		sb.WriteRune(32)
		sb.WriteString(fmt.Sprintf("%4v", int(code[ip])))
		ip++
		sb.WriteRune(32)
		sb.WriteString(fmt.Sprintf("%4v", binary.NativeEndian.Uint16(code[ip:])))
		ip += 2
		sb.WriteRune(32)
		sb.WriteString(fmt.Sprintf("%4v", binary.NativeEndian.Uint16(code[ip:])))
		ip += 2
		sb.WriteRune(32)
		sb.WriteString(fmt.Sprintf("%4v", int(code[ip])))
		ip++
	case list:
		sb.WriteRune(32)
		sb.WriteString(fmt.Sprintf("%4v", int(code[ip])))
		ip++
		sb.WriteRune(32)
		sb.WriteString(fmt.Sprintf("%4v", int(code[ip])))
		ip++
		sb.WriteRune(32)
		sb.WriteString(fmt.Sprintf("%4v", int(code[ip])))
		ip++
	case doc:
		sb.WriteRune(32)
		sb.WriteString(fmt.Sprintf("%4v", int(code[ip])))
		ip++
		sb.WriteRune(32)
		sb.WriteString(fmt.Sprintf("%4v", int(code[ip])))
		ip++
		sb.WriteRune(32)
		sb.WriteString(fmt.Sprintf("%4v", int(code[ip])))
		ip++
	case iGet:
		sb.WriteRune(32)
		sb.WriteString(fmt.Sprintf("%4v", int(code[ip])))
		ip++
		sb.WriteRune(32)
		sb.WriteString(fmt.Sprintf("%4v", int(code[ip])))
		ip++
		sb.WriteRune(32)
		sb.WriteString(fmt.Sprintf("%4v", binary.NativeEndian.Uint16(code[ip:])))
		ip += 2
		sb.WriteRune(32)
		sb.WriteString(fmt.Sprintf("%4v", binary.NativeEndian.Uint16(code[ip:])))
		ip += 2
		sb.WriteRune(32)
		sb.WriteString(fmt.Sprintf("%4v", int(code[ip])))
		ip++
	case iSet:
		sb.WriteRune(32)
		sb.WriteString(fmt.Sprintf("%4v", int(code[ip])))
		ip++
		sb.WriteRune(32)
		sb.WriteString(fmt.Sprintf("%4v", int(code[ip])))
		ip++
		sb.WriteRune(32)
		sb.WriteString(fmt.Sprintf("%4v", binary.NativeEndian.Uint16(code[ip:])))
		ip += 2
		sb.WriteRune(32)
		sb.WriteString(fmt.Sprintf("%4v", binary.NativeEndian.Uint16(code[ip:])))
		ip += 2
		sb.WriteRune(32)
		sb.WriteString(fmt.Sprintf("%4v", int(code[ip])))
		ip++
		sb.WriteRune(32)
		sb.WriteString(fmt.Sprintf("%4v", int(code[ip])))
		ip++
	case slice:
		sb.WriteRune(32)
		sb.WriteString(fmt.Sprintf("%4v", int(code[ip])))
		ip++
		sb.WriteRune(32)
		sb.WriteString(fmt.Sprintf("%4v", int(code[ip])))
		ip++
		sb.WriteRune(32)
		sb.WriteString(fmt.Sprintf("%4v", int(code[ip])))
		ip++
		sb.WriteRune(32)
		sb.WriteString(fmt.Sprintf("%4v", int(code[ip])))
		ip++
		sb.WriteRune(32)
		sb.WriteString(fmt.Sprintf("%4v", binary.NativeEndian.Uint16(code[ip:])))
		ip += 2
		sb.WriteRune(32)
		sb.WriteString(fmt.Sprintf("%4v", binary.NativeEndian.Uint16(code[ip:])))
		ip += 2
		sb.WriteRune(32)
		sb.WriteString(fmt.Sprintf("%4v", binary.NativeEndian.Uint16(code[ip:])))
		ip += 2
		sb.WriteRune(32)
		sb.WriteString(fmt.Sprintf("%4v", int(code[ip])))
		ip++
	case forSet:
		sb.WriteRune(32)
		sb.WriteString(fmt.Sprintf("%4v", int(code[ip])))
		ip++
		sb.WriteRune(32)
		sb.WriteString(fmt.Sprintf("%4v", binary.NativeEndian.Uint16(code[ip:])))
		ip += 2
	case iForSet:
		sb.WriteRune(32)
		sb.WriteString(fmt.Sprintf("%4v", int(code[ip])))
		ip++
		sb.WriteRune(32)
		sb.WriteString(fmt.Sprintf("%4v", int(code[ip])))
		ip++
		sb.WriteRune(32)
		sb.WriteString(fmt.Sprintf("%4v", binary.NativeEndian.Uint16(code[ip:])))
		ip += 2
		sb.WriteRune(32)
		sb.WriteString(fmt.Sprintf("%4v", binary.NativeEndian.Uint16(code[ip:])))
		ip += 2
	case forLoop:
		sb.WriteRune(32)
		sb.WriteString(fmt.Sprintf("%4v", int(code[ip])))
		ip++
		sb.WriteRune(32)
		sb.WriteString(fmt.Sprintf("%4v", binary.NativeEndian.Uint16(code[ip:])))
		ip += 2
	case iForLoop:
		sb.WriteRune(32)
		sb.WriteString(fmt.Sprintf("%4v", int(code[ip])))
		ip++
		sb.WriteRune(32)
		sb.WriteString(fmt.Sprintf("%4v", binary.NativeEndian.Uint16(code[ip:])))
		ip += 2
	case jump:
		sb.WriteRune(32)
		sb.WriteString(fmt.Sprintf("%4v", binary.NativeEndian.Uint16(code[ip:])))
		ip += 2
	case testF:
		sb.WriteRune(32)
		sb.WriteString(fmt.Sprintf("%4v", int(code[ip])))
		ip++
		sb.WriteRune(32)
		sb.WriteString(fmt.Sprintf("%4v", binary.NativeEndian.Uint16(code[ip:])))
		ip += 2
		sb.WriteRune(32)
		sb.WriteString(fmt.Sprintf("%4v", binary.NativeEndian.Uint16(code[ip:])))
		ip += 2
	case ret:
		sb.WriteRune(32)
		sb.WriteString(fmt.Sprintf("%4v", int(code[ip])))
		ip++
		sb.WriteRune(32)
		sb.WriteString(fmt.Sprintf("%4v", binary.NativeEndian.Uint16(code[ip:])))
		ip += 2
	case fun:
		sb.WriteRune(32)
		sb.WriteString(fmt.Sprintf("%4v", binary.NativeEndian.Uint16(code[ip:])))
		ip += 2
		sb.WriteRune(32)
		sb.WriteString(fmt.Sprintf("%4v", int(code[ip])))
		ip++
	case call:
		sb.WriteRune(32)
		sb.WriteString(fmt.Sprintf("%4v", int(code[ip])))
		ip++
		sb.WriteRune(32)
		sb.WriteString(fmt.Sprintf("%4v", int(code[ip])))
		ip++
	}
	return sb.String(), ip, counter
}
