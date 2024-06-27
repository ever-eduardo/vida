package ast

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

type Nil struct {
	Value struct{}
}

func (ast *Ast) _node()       {}
func (loc *Loc) _node()       {}
func (mut *Set) _node()       {}
func (ref *Reference) _node() {}
func (id *Identifier) _node() {}
func (b *Boolean) _node()     {}
func (n *Nil) _node()         {}
