package vida

import (
	"fmt"
	"strings"

	"github.com/ever-eduardo/vida/token"
)

func PrintBytecode(m *Module, moduleName string) string {
	clear()
	fmt.Println("Machine Code for", moduleName)
	var sb strings.Builder
	sb.WriteString(printHeader(m))
	var s string
	for i := 1; i < len(m.MainFunction.CoreFn.Code); i++ {
		s = printInstr(m.MainFunction.CoreFn.Code[i], uint64(i), false)
		sb.WriteString(s)
	}
	for idx, v := range m.Konstants {
		if f, ok := v.(*CoreFunction); ok {
			sb.WriteString(fmt.Sprintf("\n\nFunction %v/%v/%v", idx, f.Arity, f.Free))
			var s string
			for i := 1; i < len(f.Code); i++ {
				s = printInstr(f.Code[i], uint64(i), false)
				sb.WriteString(s)
			}
		}
	}
	sb.WriteString(printKonstants(m.Konstants))
	return sb.String()
}

func printHeader(m *Module) string {
	var sb strings.Builder
	var major, minor, patch uint64
	major = m.MainFunction.CoreFn.Code[0] >> 24 & 255
	minor = m.MainFunction.CoreFn.Code[0] >> 16 & 255
	patch = m.MainFunction.CoreFn.Code[0] >> 8 & 255
	sb.WriteString(fmt.Sprintf("Version %v.%v.%v", major, minor, patch))
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

func printInstr(instr, ip uint64, isRunningDebug bool) string {
	var sb strings.Builder
	var op, A, B, P uint64
	op = instr >> shift56
	A = instr >> shift16 & clean16
	B = instr & clean16
	P = instr >> shift32 & clean24
	if !isRunningDebug {
		sb.WriteRune(10)
		sb.WriteString(fmt.Sprintf("  [%3v]  ", ip))
		sb.WriteString(fmt.Sprintf("%7v", opcodes[op]))
	} else {
		sb.WriteString(opcodes[op])
	}
	switch op {
	case end:
		return sb.String()
	case storeG:
		sb.WriteString(fmt.Sprintf(" %3v %3v %3v", P, A, B))
	case loadG, loadF, loadK, move, storeF, list, object, iGet:
		sb.WriteString(fmt.Sprintf(" %3v %3v", A, B))
	case prefix, equals:
		sb.WriteString(fmt.Sprintf(" %3v %3v %3v", token.Token(P).String(), A, B))
	case binopG, binopK, binop:
		sb.WriteString(fmt.Sprintf(" %3v %3v %3v %3v", token.Token(P>>shift16).String(), P&clean16, A, B))
	}
	return sb.String()
}
