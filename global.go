package vida

import (
	"bufio"
	"fmt"
	"os"

	"github.com/ever-eduardo/vida/verror"
)

var NilValue = Nil{}
var stdlibLoader map[string]func() Value

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
	"isError",
}

func loadCoreLib(store *[]Value) {
	*store = append(*store,
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
	argsLength := len(args)
	if argsLength == 1 {
		if args[0].Boolean() {
			return NilValue, nil
		}
		err := fmt.Errorf("%s", fmt.Sprintf("\n\n  [%v]\n\n", verror.AssertionErrType))
		return NilValue, err
	}
	if argsLength > 1 {
		if args[0].Boolean() {
			return NilValue, nil
		}
		err := fmt.Errorf("%s", fmt.Sprintf("\n\n  [%v]\n   Message : %v\n\n", verror.AssertionErrType, args[1].String()))
		return NilValue, err

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
	return &List{}, nil
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
			if l, isPresent := stdlibLoader[v.Value]; isPresent {
				return l(), nil
			}
		}
	}
	return NilValue, nil
}

func gfnExcept(args ...Value) (Value, error) {
	if len(args) > 0 {
		err := fmt.Errorf("%s", fmt.Sprintf("\n\n  [%v]\n   Message : %v\n\n", verror.ExceptionErrType, args[0].String()))
		return NilValue, err
	}
	err := fmt.Errorf("%s", fmt.Sprintf("\n\n  [%v]\n\n", verror.ExceptionErrType))
	return NilValue, err
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

var coreLibDescription = []string{
	`
	Print one or more values separated by a comma.
	Always return nil.
	Examples: print(value), print(a, b, c) -> nil
	`,
	`
	Return an integer representing the length of lists, 
	objects or strings. In case of a string value, 
	the function returns the number of unicode codepoints.
	Example: len(value) -> int
	`,
	`
	Append one of more values separated by comma 
	at the end of a list.
	Return the list passed as first argument.
	Examples: let xs be a list, then 
	append(xs, value), append(xs, a, b, c) -> xs
	`,
	`
	Create a list. 
	Receive 0, 1 or 2 arguments. 
	Whith zero argumeents, return an empty list. 
	With 1 argument n, with n of type intenger,
	return a list of n elements initialized to nil.
	With 2 argumeents (n, m), with n of type integer,
	and m of any type, return a list of n elements all 
	initialized to the m value.
	Examples: 
		mkls() -> [],
		mkls(10) -> [nil, ..., nil],
		mkls(n, v) -> [v, v, ... , v]
	`,
	`
	Load a specific library from the stdlib.
	Receive an argument n of type string, and return an object
	containing the constants and functionality. If thee library
	does not exist, return nil.
	Example: load("math"), load("random")
	`,
	`
	Return the type of a value as string.
	Example: type(123) -> "int"
	`,
	`
	Make an assertion about an expression, and optionally print a
	message in case of assertion failure.
	If the assertion fails, then it will always produce 
	a runtime error. Otherwise, it just return a nil value.
	Example: assert(false), assert(true)
	`,
	`
	Return a string with the given format.
	The most common verb formats are: %v, %T, %f, %d, %b, %x
	Example: fmt("This is the number %v", 15)
	`,
	`
	Show a prompt and wait for an input from the user.
	If no prompt is given, it shows a default one.
	Return a string representing the user input.
	Example: input("Write something here") -> string
	`,
	`
	Make a copy of value semantics values or a deep copy 
	of a reference semantics values.
	Example: clone(someValue)
	`,
	`
	Delete a key from an object.
	Example: let xs be an objet containing a key val, then
	del(xs, "val") deletes the key val.
	`,
	`
	Create an error value. An error value may be used to signal
	some behavior considered an error. The boolean value of an
	error value is always false. When an argument is give, it will
	be the printable message for the client of the functionality
	with the unexpected behavior.
	Example: 
		ret error(message)
	        let result = f()
		if not result {handle the error} or
		if result {handle the returned value}
	`,
	`
	Create an exception to signal some exceptional or unexpected
	behavior. It will always generate a runtime error. 
	When an argumentis given, it is shown in the error message.
	Example: exception(message)
	`,
	`
	Help to explicitly check for an error value.
	Example: if isError(value) {handle the error here}
	`,
}

func PrintCoreLibInformation() {
	fmt.Printf("CoreLib:\n\n")
	for i := 0; i < len(coreLibNames); i++ {
		fmt.Printf("  %v %v\n\n", coreLibNames[i], coreLibDescription[i])
	}
}
