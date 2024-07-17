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
		case equals:
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
		case doc:
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
		case iSet:
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
		case forSet:
			sb.WriteRune(32)
			sb.WriteString(fmt.Sprintf("%3v", int(code[ip])))
			ip++
			sb.WriteRune(32)
			sb.WriteString(fmt.Sprintf("%3v", binary.NativeEndian.Uint16(code[ip:])))
			ip += 2
		case iForSet:
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
		case forLoop, iForLoop:
			sb.WriteRune(32)
			sb.WriteString(fmt.Sprintf("%3v", int(code[ip])))
			ip++
			sb.WriteRune(32)
			sb.WriteString(fmt.Sprintf("%3v", binary.NativeEndian.Uint16(code[ip:])))
			ip += 2
		case jump:
			sb.WriteRune(32)
			sb.WriteString(fmt.Sprintf("%3v", binary.NativeEndian.Uint16(code[ip:])))
			ip += 2
		case testF:
			sb.WriteRune(32)
			sb.WriteString(fmt.Sprintf("%3v", int(code[ip])))
			ip++
			sb.WriteRune(32)
			sb.WriteString(fmt.Sprintf("%3v", binary.NativeEndian.Uint16(code[ip:])))
			ip += 2
			sb.WriteRune(32)
			sb.WriteString(fmt.Sprintf("%3v", binary.NativeEndian.Uint16(code[ip:])))
			ip += 2
		case ret:
			sb.WriteRune(32)
			sb.WriteString(fmt.Sprintf("%3v", int(code[ip])))
			ip++
			sb.WriteRune(32)
			sb.WriteString(fmt.Sprintf("%3v", binary.NativeEndian.Uint16(code[ip:])))
			ip += 2
			sb.WriteRune(32)
			sb.WriteString(fmt.Sprintf("%3v", int(code[ip])))
			ip++
		case fun:
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

func printInstr(ip int, code []byte) string {
	var sb strings.Builder
	op := code[ip]
	sb.WriteString(opcodes[op])
	ip++
	switch op {
	case end:
		return sb.String()
	case setG:
		sb.WriteRune(32)
		sb.WriteString(fmt.Sprintf("%3v", int(code[ip])))
		ip++
		sb.WriteRune(32)
		sb.WriteString(fmt.Sprintf("%3v", binary.NativeEndian.Uint16(code[ip:])))
		ip += 2
		sb.WriteRune(32)
		sb.WriteString(fmt.Sprintf("%3v", binary.NativeEndian.Uint16(code[ip:])))
	case setL:
		sb.WriteRune(32)
		sb.WriteString(fmt.Sprintf("%3v", int(code[ip])))
		ip++
		sb.WriteRune(32)
		sb.WriteString(fmt.Sprintf("%3v", binary.NativeEndian.Uint16(code[ip:])))
		ip += 2
		sb.WriteRune(32)
		sb.WriteString(fmt.Sprintf("%3v", int(code[ip])))
	case move:
		sb.WriteRune(32)
		sb.WriteString(fmt.Sprintf("%3v", int(code[ip])))
		ip++
		sb.WriteRune(32)
		sb.WriteString(fmt.Sprintf("%3v", int(code[ip])))
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
	case equals:
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
	case binop:
		sb.WriteRune(32)
		sb.WriteString(fmt.Sprintf("%3v", token.Tokens[int(code[ip])]))
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
	case list:
		sb.WriteRune(32)
		sb.WriteString(fmt.Sprintf("%3v", int(code[ip])))
		ip++
		sb.WriteRune(32)
		sb.WriteString(fmt.Sprintf("%3v", int(code[ip])))
		ip++
		sb.WriteRune(32)
		sb.WriteString(fmt.Sprintf("%3v", int(code[ip])))
	case doc:
		sb.WriteRune(32)
		sb.WriteString(fmt.Sprintf("%3v", int(code[ip])))
		ip++
		sb.WriteRune(32)
		sb.WriteString(fmt.Sprintf("%3v", int(code[ip])))
		ip++
		sb.WriteRune(32)
		sb.WriteString(fmt.Sprintf("%3v", int(code[ip])))
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
	case iSet:
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
		sb.WriteRune(32)
		sb.WriteString(fmt.Sprintf("%3v", int(code[ip])))
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
	case forSet:
		sb.WriteRune(32)
		sb.WriteString(fmt.Sprintf("%3v", int(code[ip])))
		ip++
		sb.WriteRune(32)
		sb.WriteString(fmt.Sprintf("%3v", binary.NativeEndian.Uint16(code[ip:])))
	case iForSet:
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
	case forLoop, iForLoop:
		sb.WriteRune(32)
		sb.WriteString(fmt.Sprintf("%3v", int(code[ip])))
		ip++
		sb.WriteRune(32)
		sb.WriteString(fmt.Sprintf("%3v", binary.NativeEndian.Uint16(code[ip:])))
	case jump:
		sb.WriteRune(32)
		sb.WriteString(fmt.Sprintf("%3v", binary.NativeEndian.Uint16(code[ip:])))
	case testF:
		sb.WriteRune(32)
		sb.WriteString(fmt.Sprintf("%3v", int(code[ip])))
		ip++
		sb.WriteRune(32)
		sb.WriteString(fmt.Sprintf("%3v", binary.NativeEndian.Uint16(code[ip:])))
		ip += 2
		sb.WriteRune(32)
		sb.WriteString(fmt.Sprintf("%3v", binary.NativeEndian.Uint16(code[ip:])))
	case ret:
		sb.WriteRune(32)
		sb.WriteString(fmt.Sprintf("%3v", int(code[ip])))
		ip++
		sb.WriteRune(32)
		sb.WriteString(fmt.Sprintf("%3v", binary.NativeEndian.Uint16(code[ip:])))
		ip += 2
		sb.WriteRune(32)
		sb.WriteString(fmt.Sprintf("%3v", int(code[ip])))
	case fun:
		sb.WriteRune(32)
		sb.WriteString(fmt.Sprintf("%3v", binary.NativeEndian.Uint16(code[ip:])))
		ip += 2
		sb.WriteRune(32)
		sb.WriteString(fmt.Sprintf("%3v", binary.NativeEndian.Uint16(code[ip:])))
		ip += 2
		sb.WriteRune(32)
		sb.WriteString(fmt.Sprintf("%3v", int(code[ip])))
	}
	return sb.String()
}
