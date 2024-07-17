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

type ReferenceStmt struct {
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

type String struct {
	Value string
}

type Nil struct {
	Value struct{}
}

type List struct {
	ExprList []Node
}

type Property struct {
	Value string
}

type Pair struct {
	Key   Node
	Value Node
}

type Document struct {
	Pairs []*Pair
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

type IGet struct {
	Indexable Node
	Index     Node
}

type IGetStmt struct {
	Index Node
}

type Slice struct {
	Value Node
	First Node
	Last  Node
	Mode  int
}

type Select struct {
	Selectable Node
	Selector   Node
}

type SelectStmt struct {
	Selector Node
}

type ISet struct {
	Index Node
	Expr  Node
}

type For struct {
	Init  Node
	End   Node
	Step  Node
	Id    string
	Block Node
}

type ForState struct {
	Value string
}

type IFor struct {
	Key   string
	Value string
	Expr  Node
	Block Node
}

type Branch struct {
	Elifs []Node
	If    Node
	Else  Node
}

type If struct {
	Condition Node
	Block     Node
}

type Else struct {
	Block Node
}

type While struct {
	Condition Node
	Block     Node
}

type Break struct{}

type Continue struct{}

type Fun struct {
	Args     []string
	HasArrow bool
	Body     Node
	Expr     Node
}

type Ret struct {
	Expr Node
}

func (ast *Ast) _node()           {}
func (loc *Loc) _node()           {}
func (mut *Set) _node()           {}
func (ref *Reference) _node()     {}
func (ref *ReferenceStmt) _node() {}
func (id *Identifier) _node()     {}
func (b *Boolean) _node()         {}
func (n *Nil) _node()             {}
func (n *PrefixExpr) _node()      {}
func (n *BinaryExpr) _node()      {}
func (n *Block) _node()           {}
func (n *Integer) _node()         {}
func (n *Float) _node()           {}
func (n *String) _node()          {}
func (n *List) _node()            {}
func (n *IGet) _node()            {}
func (n *IGetStmt) _node()        {}
func (n *Slice) _node()           {}
func (n *Document) _node()        {}
func (n *Pair) _node()            {}
func (n *Property) _node()        {}
func (n *Select) _node()          {}
func (n *SelectStmt) _node()      {}
func (n *ISet) _node()            {}
func (n *For) _node()             {}
func (n *IFor) _node()            {}
func (n *ForState) _node()        {}
func (n *Branch) _node()          {}
func (n *If) _node()              {}
func (n *Else) _node()            {}
func (n *While) _node()           {}
func (n *Break) _node()           {}
func (n *Continue) _node()        {}
func (n *Fun) _node()             {}
func (n *Ret) _node()             {}
