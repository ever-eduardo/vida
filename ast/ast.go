package ast

type Node interface {
	_node()
}

type Statement interface {
	_stmt()
}

type Expr interface {
	_expr()
}

type Ast struct {
	Statement []Statement
}

type Val struct {
	Identifier string
	Expr       Expr
}

type Loc struct {
	Identifier string
	Expr       Expr
}

type Mut struct {
	Identifier string
	Expr       Expr
}

type Reference struct {
	Identifier string
}

type Atomic struct {
	Expr Expr
}

type Boolean struct {
	Value bool
}

type Nil struct {
	Value struct{}
}

func (ast Ast) _node()       {}
func (val Val) _node()       {}
func (loc Loc) _node()       {}
func (mut Mut) _node()       {}
func (ref Reference) _node() {}
func (atomic Atomic) _node() {}
func (b Boolean) _node()     {}
func (n Nil) _node()         {}

func (val Val) _stmt() {}
func (loc Loc) _stmt() {}
func (mut Mut) _stmt() {}

func (ref Reference) _expr() {}
func (atomic Atomic) _expr() {}
func (b Boolean) _expr()     {}
func (n Nil) _expr()         {}
