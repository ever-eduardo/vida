package vida

import (
	"fmt"
	"time"

	"github.com/alkemist-17/vida/ast"
	"github.com/alkemist-17/vida/lexer"
	"github.com/alkemist-17/vida/token"
)

type Interpreter struct {
	parser   *parser
	compiler *compiler
	vm       *vM
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
	c := newMainCompiler(rAst, modulePath)
	m, err := c.compileModule()
	if err != nil {
		return nil, err
	}
	vm, err := newVM(m, stdlib, c.linesMap)
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
	fmt.Print("\n\nPress 'Enter' to continue => ")
	fmt.Scanf(" ")
	c := newMainCompiler(rAst, modulePath)
	m, err := c.compileModule()
	if err != nil {
		return nil, err
	}
	fmt.Println(PrintBytecode(m, m.MainFunction.CoreFn.ModuleName))
	fmt.Print("\n\nPress 'Enter' to continue => ")
	fmt.Scanf(" ")
	vm, err := newVM(m, stdlib, c.linesMap)
	if err != nil {
		return nil, err
	}
	return &Interpreter{
		parser:   p,
		compiler: c,
		vm:       vm,
	}, nil
}

func PrintAST(modulePath string) error {
	src, err := readModule(modulePath)
	if err != nil {
		return err
	}
	p := newParser(src, modulePath)
	rAst, err := p.parse()
	if err != nil {
		return err
	}
	fmt.Println(ast.PrintAST(rAst))
	return nil
}

func PrintTokens(modulePath string) error {
	src, err := readModule(modulePath)
	if err != nil {
		return err
	}
	l := lexer.New(src, modulePath)
	hadError := false
	fmt.Printf("%4v %-15v %-2v\n\n", "line", "token", "repr")
	for {
		line, tok, lit := l.Next()
		if l.LexicalError.Message != "" {
			hadError = true
			break
		}
		fmt.Printf("%4v %-15v %-2v\n", line, tok, lit)
		if tok == token.EOF {
			fmt.Println()
			break
		}
	}
	if hadError {
		return l.LexicalError
	}
	return nil
}

func PrintMachineCode(modulePath string) error {
	src, err := readModule(modulePath)
	if err != nil {
		return err
	}
	p := newParser(src, modulePath)
	rAst, err := p.parse()
	if err != nil {
		return err
	}
	c := newMainCompiler(rAst, modulePath)
	m, err := c.compileModule()
	if err != nil {
		return err
	}
	fmt.Println(PrintBytecode(m, m.MainFunction.CoreFn.ModuleName))
	return nil
}

func (i *Interpreter) Run() (Result, error) {
	return i.vm.run()
}

func (i *Interpreter) MeasureRunTime() (Result, error) {
	init := time.Now()
	r, err := i.vm.run()
	end := time.Since(init)
	fmt.Printf("\n\nThe interpreter has finished.\n\n")
	fmt.Printf("Time = %vs\n", end.Seconds())
	fmt.Printf("Time = %v\n\n", end)
	return r, err
}

func (i *Interpreter) Debug() (Result, error) {
	return i.vm.debug()
}

func (i *Interpreter) PrintCallStack() {
	i.vm.printCallStack()
}
