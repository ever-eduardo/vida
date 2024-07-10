package ast

import (
	"fmt"
	"strings"
)

func PrintAST(node Node) string {
	var sb strings.Builder
	printAST(node, &sb, zeroLevel)
	return sb.String()
}

func buildIndent(sb *strings.Builder, level int) {
	for range level {
		sb.WriteRune(space)
	}
}

func printAST(node Node, sb *strings.Builder, level int) {
	switch n := node.(type) {
	case *Ast:
		sb.WriteRune(nl)
		sb.WriteRune(nl)
		sb.WriteString("AST")
		for i := range len(n.Statement) {
			printAST(n.Statement[i], sb, level+oneLevel)
		}
		sb.WriteRune(nl)
		sb.WriteRune(nl)
	case *Loc:
		sb.WriteRune(nl)
		buildIndent(sb, level+oneLevel)
		sb.WriteString("Loc")
		sb.WriteRune(nl)
		buildIndent(sb, level+twoLevels)
		sb.WriteString(n.Identifier)
		sb.WriteRune(nl)
		printAST(n.Expr, sb, level+oneLevel)
	case *Set:
		sb.WriteRune(nl)
		buildIndent(sb, level+oneLevel)
		sb.WriteString("Set")
		buildIndent(sb, level+twoLevels)
		printAST(n.LHS, sb, level+oneLevel)
		sb.WriteRune(nl)
		printAST(n.Expr, sb, level+oneLevel)
	case *Reference:
		buildIndent(sb, level+oneLevel)
		sb.WriteString("Ref")
		sb.WriteRune(nl)
		buildIndent(sb, level+twoLevels)
		sb.WriteString(n.Value)
	case *Identifier:
		sb.WriteRune(nl)
		buildIndent(sb, level+oneLevel)
		sb.WriteString("Id")
		sb.WriteRune(nl)
		buildIndent(sb, level+twoLevels)
		sb.WriteString(n.Value)
	case *Boolean:
		buildIndent(sb, level+oneLevel)
		sb.WriteString("Bool")
		sb.WriteRune(nl)
		buildIndent(sb, level+twoLevels)
		if n.Value {
			sb.WriteString("true")
		} else {
			sb.WriteString("false")
		}
	case *Nil:
		buildIndent(sb, level+oneLevel)
		sb.WriteString("Nil")
		sb.WriteRune(nl)
		buildIndent(sb, level+twoLevels)
		sb.WriteString("nil")
	case *Integer:
		buildIndent(sb, level+oneLevel)
		sb.WriteString("Integer")
		sb.WriteRune(nl)
		buildIndent(sb, level+twoLevels)
		sb.WriteString(fmt.Sprint(n.Value))
	case *Float:
		buildIndent(sb, level+oneLevel)
		sb.WriteString("Float")
		sb.WriteRune(nl)
		buildIndent(sb, level+twoLevels)
		sb.WriteString(fmt.Sprint(n.Value))
	case *String:
		buildIndent(sb, level+oneLevel)
		sb.WriteString("String")
		sb.WriteRune(nl)
		buildIndent(sb, level+twoLevels)
		sb.WriteString(n.Value)
	case *List:
		buildIndent(sb, level+oneLevel)
		sb.WriteString("List")
		sb.WriteRune(nl)
		if len(n.ExprList) == 0 {
			buildIndent(sb, level+twoLevels)
			sb.WriteString("[]")
		} else {
			for _, v := range n.ExprList {
				printAST(v, sb, level+oneLevel)
				sb.WriteRune(nl)
			}
		}
	case *PrefixExpr:
		buildIndent(sb, level+oneLevel)
		sb.WriteString("Prefix")
		sb.WriteRune(nl)
		buildIndent(sb, level+twoLevels)
		sb.WriteString(n.Op.String())
		sb.WriteRune(nl)
		printAST(n.Expr, sb, level+oneLevel)
	case *BinaryExpr:
		buildIndent(sb, level+oneLevel)
		sb.WriteString("Binary")
		sb.WriteRune(nl)
		buildIndent(sb, level+twoLevels)
		sb.WriteString(n.Op.String())
		sb.WriteRune(nl)
		printAST(n.Lhs, sb, level+twoLevels)
		sb.WriteRune(nl)
		printAST(n.Rhs, sb, level+twoLevels)
	case *IndexGet:
		buildIndent(sb, level+oneLevel)
		sb.WriteString("IndexGet")
		sb.WriteRune(nl)
		printAST(n.Indexable, sb, level+oneLevel)
		sb.WriteRune(nl)
		printAST(n.Index, sb, level+oneLevel)
	case *Slice:
		buildIndent(sb, level+oneLevel)
		sb.WriteString("Slice")
		sb.WriteRune(nl)
		printAST(n.Value, sb, level+oneLevel)
		sb.WriteRune(nl)
		printAST(n.First, sb, level+oneLevel)
		sb.WriteRune(nl)
		printAST(n.Last, sb, level+oneLevel)
	case *Block:
		sb.WriteRune(nl)
		buildIndent(sb, level+oneLevel)
		sb.WriteString("Block")
		for i := range len(n.Statement) {
			printAST(n.Statement[i], sb, level+twoLevels)
		}
	default:
		buildIndent(sb, level+oneLevel)
		sb.WriteString("Node")
	}
}

const space rune = 32
const nl rune = 10
const oneLevel = 2
const twoLevels = 4
const zeroLevel = 0
