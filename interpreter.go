package vida

import (
	"fmt"
	"time"

	"github.com/ever-eduardo/vida/ast"
)

type Interpreter struct {
	parser   *Parser
	compiler *Compiler
	vm       *VM
}

func NewInterpreter(modulePath string, stdlib map[string]func() Value) (*Interpreter, error) {
	src, err := readModule(modulePath)
	if err != nil {
		return nil, err
	}
	p := newParser(src, modulePath)
	rAst, err := p.parse()
	if err != nil {
		return nil, err
	}
	c := newCompiler(rAst, modulePath)
	m, err := c.compileModule()
	if err != nil {
		return nil, err
	}
	vm, err := newVM(m, stdlib)
	if err != nil {
		return nil, err
	}
	return &Interpreter{
		parser:   p,
		compiler: c,
		vm:       vm,
	}, nil
}

func NewDebugger(modulePath string, stdlib map[string]func() Value) (*Interpreter, error) {
	src, err := readModule(modulePath)
	if err != nil {
		return nil, err
	}
	p := newParser(src, modulePath)
	rAst, err := p.parse()
	if err != nil {
		return nil, err
	}
	fmt.Println(ast.PrintAST(rAst))
	fmt.Scanf(" ")
	c := newCompiler(rAst, modulePath)
	m, err := c.compileModule()
	if err != nil {
		return nil, err
	}
	fmt.Println(PrintBytecode(m, m.Name))
	fmt.Scanf(" ")
	vm, err := newVM(m, stdlib)
	if err != nil {
		return nil, err
	}
	return &Interpreter{
		parser:   p,
		compiler: c,
		vm:       vm,
	}, nil
}

func (i *Interpreter) Run() (Result, error) {
	return i.vm.run()
}

func (i *Interpreter) MeasureRunTime() (Result, error) {
	init := time.Now()
	r, err := i.vm.run()
	fmt.Printf("VM time = %vs\n", time.Since(init).Seconds())
	fmt.Printf("VM time = %v\n", time.Since(init))
	return r, err
}

func (i *Interpreter) Debug() (Result, error) {
	return i.vm.debug()
}
