package ast

import (
	"strings"
)

func PrintAST(node Node) string {
	var sb strings.Builder
	printAST(node, &sb, 0)
	return sb.String()
}

func buildIndent(sb *strings.Builder, level int) {
	for range level {
		sb.WriteRune(32)
	}
}

func printAST(node Node, sb *strings.Builder, level int) {
	switch n := node.(type) {
	case *Ast:
		sb.WriteRune(10)
		sb.WriteRune(10)
		sb.WriteString("AST")
		for i := range len(n.Statement) {
			printAST(n.Statement[i], sb, level+1)
		}
		sb.WriteRune(10)
		sb.WriteRune(10)
	case *Loc:
		sb.WriteRune(10)
		buildIndent(sb, level+1)
		sb.WriteString("Loc")
		sb.WriteRune(10)
		buildIndent(sb, level+2)
		sb.WriteString(n.Identifier)
		sb.WriteRune(10)
		printAST(n.Expr, sb, level+1)
	case *Set:
		sb.WriteRune(10)
		buildIndent(sb, level+1)
		sb.WriteString("Set")
		buildIndent(sb, level+2)
		printAST(n.LHS, sb, level+1)
		sb.WriteRune(10)
		printAST(n.Expr, sb, level+1)
	case *Reference:
		buildIndent(sb, level+1)
		sb.WriteString("Ref")
		sb.WriteRune(10)
		buildIndent(sb, level+2)
		sb.WriteString(n.Value)
	case *Identifier:
		sb.WriteRune(10)
		buildIndent(sb, level+1)
		sb.WriteString("Id")
		sb.WriteRune(10)
		buildIndent(sb, level+2)
		sb.WriteString(n.Value)
	case *Boolean:
		buildIndent(sb, level+1)
		sb.WriteString("Bool")
		sb.WriteRune(10)
		buildIndent(sb, level+2)
		if n.Value {
			sb.WriteString("true")
		} else {
			sb.WriteString("false")
		}
	case *Nil:
		buildIndent(sb, level+1)
		sb.WriteString("nil")
	case *PrefixExpr:
		buildIndent(sb, level+1)
		sb.WriteString("Prefix")
		sb.WriteRune(10)
		buildIndent(sb, level+2)
		sb.WriteString(n.Op.String())
		sb.WriteRune(10)
		printAST(n.Expr, sb, level+1)
	case *BinaryExpr:
		buildIndent(sb, level+1)
		sb.WriteString("Binary")
		sb.WriteRune(10)
		buildIndent(sb, level+2)
		sb.WriteString(n.Op.String())
		sb.WriteRune(10)
		printAST(n.Lhs, sb, level+2)
		sb.WriteRune(10)
		printAST(n.Rhs, sb, level+2)
	}
}
