package ast

import "github.com/ever-eduardo/vida/token"

type Node interface {
	_node()
}

type Ast struct {
	Statement []Node
}

type Loc struct {
	Identifier string
	Expr       Node
}

type Set struct {
	LHS  Node
	Expr Node
}

type Reference struct {
	Value string
}

type Identifier struct {
	Value string
}

type Boolean struct {
	Value bool
}

type Integer struct {
	Value int64
}

type Float struct {
	Value float64
}

type Nil struct {
	Value struct{}
}

type PrefixExpr struct {
	Expr Node
	Op   token.Token
}

type BinaryExpr struct {
	Lhs Node
	Rhs Node
	Op  token.Token
}

type Block struct {
	Statement []Node
}

func (ast *Ast) _node()       {}
func (loc *Loc) _node()       {}
func (mut *Set) _node()       {}
func (ref *Reference) _node() {}
func (id *Identifier) _node() {}
func (b *Boolean) _node()     {}
func (n *Nil) _node()         {}
func (n *PrefixExpr) _node()  {}
func (n *BinaryExpr) _node()  {}
func (n *Block) _node()       {}
func (n *Integer) _node()     {}
func (n *Float) _node()       {}
