package vida

import (
	"bufio"
	"fmt"
	"os"

	"github.com/ever-eduardo/vida/verror"
)

var NilValue = Nil{}
var LibLoaders map[string]func() Value

var coreLibNames = []string{
	"print",
	"len",
	"append",
	"mkls",
	"load",
	"type",
	"assert",
	"fmt",
	"input",
	"clone",
	"del",
	"error",
	"exception",
	"err",
}

func loadCoreLib() []Value {
	clib := make([]Value, 0, 32)
	clib = append(clib,
		GFn(gfnPrint),
		GFn(gfnLen),
		GFn(gfnAppend),
		GFn(gfnMakeList),
		GFn(gfnLoadLib),
		GFn(gfnType),
		GFn(gfnAssert),
		GFn(gfnFormat),
		GFn(gfnReadLine),
		GFn(gfnClone),
		GFn(gfnDel),
		GFn(gfnError),
		GFn(gfnExcept),
		GFn(gfnIsError))
	return clib
}

func gfnPrint(args ...Value) (Value, error) {
	var s []any
	for _, v := range args {
		s = append(s, v.String())
	}
	fmt.Fprintln(os.Stdout, s...)
	return NilValue, nil
}

func gfnLen(args ...Value) (Value, error) {
	if len(args) > 0 {
		switch v := args[0].(type) {
		case *List:
			return Integer(len(v.Value)), nil
		case *Object:
			return Integer(len(v.Value)), nil
		case String:
			if v.Runes == nil {
				v.Runes = []rune(v.Value)
			}
			return Integer(len(v.Runes)), nil
		}
	}
	return NilValue, nil
}

func gfnType(args ...Value) (Value, error) {
	if len(args) > 0 {
		return String{Value: args[0].Type()}, nil
	}
	return NilValue, nil
}

func gfnFormat(args ...Value) (Value, error) {
	if len(args) > 1 {
		switch v := args[0].(type) {
		case String:
			s, e := formatValue(v.Value, args[1:]...)
			return String{Value: s}, e
		}
	}
	return NilValue, nil
}

func gfnAssert(args ...Value) (Value, error) {
	if len(args) == 1 {
		if args[0].Boolean() {
			return NilValue, nil
		} else {
			return NilValue, verror.ErrAssertionFailure
		}
	} else if len(args) > 1 {
		if args[0].Boolean() {
			return NilValue, nil
		} else {
			return NilValue, fmt.Errorf("%v: %v", verror.ErrAssertionFailure, args[1].String())
		}
	}
	return NilValue, nil
}

func gfnAppend(args ...Value) (Value, error) {
	if len(args) >= 2 {
		switch v := args[0].(type) {
		case *List:
			v.Value = append(v.Value, args[1:]...)
			return v, nil
		}
	}
	return NilValue, nil
}

func gfnMakeList(args ...Value) (Value, error) {
	largs := len(args)
	if largs > 0 {
		switch v := args[0].(type) {
		case Integer:
			var init Value = NilValue
			if largs > 1 {
				init = args[1]
			}
			if v >= 0 && v <= 0x7FFF_FFFF {
				arr := make([]Value, v)
				for i := range v {
					arr[i] = init
				}
				return &List{Value: arr}, nil
			}
		}
	}
	return NilValue, nil
}

func gfnReadLine(args ...Value) (Value, error) {
	if len(args) > 0 {
		fmt.Print(args[0])
	} else {
		fmt.Print("Input >> ")
	}
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		return String{Value: scanner.Text()}, nil
	}
	if err := scanner.Err(); err != nil {
		return NilValue, err
	}
	return NilValue, nil
}

func gfnClone(args ...Value) (Value, error) {
	if len(args) > 0 {
		return args[0].Clone(), nil
	}
	return NilValue, nil
}

func gfnDel(args ...Value) (Value, error) {
	if len(args) >= 2 {
		if o, ok := args[0].(*Object); ok {
			delete(o.Value, args[1].String())
			o.UpdateKeys()
		}
	}
	return NilValue, nil
}

func gfnLoadLib(args ...Value) (Value, error) {
	if len(args) > 0 {
		if v, ok := args[0].(String); ok {
			if l, isPresent := LibLoaders[v.Value]; isPresent {
				return l(), nil
			}
		}
	}
	return NilValue, nil
}

func gfnExcept(args ...Value) (Value, error) {
	if len(args) > 0 {
		return NilValue, fmt.Errorf("%v: %v", verror.Exception, args[0].String())
	}
	return NilValue, fmt.Errorf("%v", verror.Exception)
}

func gfnError(args ...Value) (Value, error) {
	if len(args) > 0 {
		return Error{Message: args[0]}, nil
	}
	return Error{Message: NilValue}, nil
}
func gfnIsError(args ...Value) (Value, error) {
	if len(args) > 0 {
		_, ok := args[0].(Error)
		return Bool(ok), nil
	}
	return Bool(false), nil
}
