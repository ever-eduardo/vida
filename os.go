package vida

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

func loadFoundationOS() Value {
	m := &Object{Value: make(map[string]Value, 18)}
	m.Value["args"] = NativeFunction(osArgs)
	m.Value["env"] = NativeFunction(osEnviron)
	m.Value["exit"] = NativeFunction(osExit)
	m.Value["getFromEnv"] = NativeFunction(osGetEnv)
	m.Value["hostname"] = NativeFunction(osHostname)
	m.Value["pathSeparator"] = NativeFunction(osGetPathSeparator)
	m.Value["mkdir"] = NativeFunction(osMkdir)
	m.Value["mkdirAll"] = NativeFunction(osMkdirAll)
	m.Value["rm"] = NativeFunction(osRemove)
	m.Value["rmAll"] = NativeFunction(osRemoveAll)
	m.Value["name"] = NativeFunction(osName)
	m.Value["arch"] = NativeFunction(osArch)
	m.Value["run"] = NativeFunction(osRunCMD)
	m.Value["cwd"] = NativeFunction(osCwd)
	m.Value["chdir"] = NativeFunction(osChdir)
	m.Value["stdin"] = newVidaFile(os.Stdin, "[stdin]", os.O_RDONLY)
	m.Value["stdout"] = newVidaFile(os.Stdout, "[stdout]", os.O_WRONLY)
	m.Value["stderr"] = newVidaFile(os.Stderr, "[stderr]", os.O_WRONLY)
	return m
}

func osArgs(ctx *Context, args ...Value) (Value, error) {
	xs := &Array{Value: make([]Value, len(os.Args))}
	for i, v := range os.Args {
		xs.Value[i] = &String{Value: v}
	}
	return xs, nil
}

func osEnviron(ctx *Context, args ...Value) (Value, error) {
	env := os.Environ()
	xs := &Array{Value: make([]Value, len(env))}
	for i, v := range os.Environ() {
		xs.Value[i] = &String{Value: v}
	}
	return xs, nil
}

func osExit(ctx *Context, args ...Value) (Value, error) {
	code := 0
	if len(args) > 0 {
		c, ok := args[0].(Integer)
		if !ok {
			return argError("os.exit", fmt.Sprintf("expected an integer exit code argument, got %s", args[0].Type()))
		}
		code = int(c)
	}
	os.Exit(code)
	return Nil, nil
}

func osGetEnv(ctx *Context, args ...Value) (Value, error) {
	if len(args) == 0 {
		return argError("os.getFromEnv", "expected a string argument (environment variable name)")
	}
	val, ok := args[0].(*String)
	if !ok {
		return argError("os.getFromEnv", fmt.Sprintf("expected a string argument, got %s", args[0].Type()))
	}
	xs := make([]Value, 0, 2)
	if r, ok := os.LookupEnv(val.Value); ok {
		xs = append(xs, &String{Value: r})
		xs = append(xs, Bool(ok))
	} else {
		xs = append(xs, &String{Value: EmptyString})
		xs = append(xs, Bool(ok))
	}
	return &Array{Value: xs}, nil
}

func osHostname(ctx *Context, args ...Value) (Value, error) {
	if h, e := os.Hostname(); e == nil {
		return &String{Value: h}, nil
	} else {
		return &VidaError{Message: &String{Value: e.Error()}}, nil
	}
}

func osGetPathSeparator(ctx *Context, args ...Value) (Value, error) {
	return &String{Value: string(os.PathSeparator)}, nil
}

func osMkdir(ctx *Context, args ...Value) (Value, error) {
	if len(args) == 0 {
		return argError("os.mkdir", "expected a path string argument")
	}
	d, ok := args[0].(*String)
	if !ok {
		return argError("os.mkdir", fmt.Sprintf("expected a path string argument, got %s", args[0].Type()))
	}
	err := os.Mkdir(d.Value, 0660)
	if err != nil && !os.IsExist(err) {
		return &VidaError{Message: &String{Value: err.Error()}}, nil
	}
	return True, nil
}

func osMkdirAll(ctx *Context, args ...Value) (Value, error) {
	if len(args) == 0 {
		return argError("os.mkdirAll", "expected a path string argument")
	}
	d, ok := args[0].(*String)
	if !ok {
		return argError("os.mkdirAll", fmt.Sprintf("expected a path string argument, got %s", args[0].Type()))
	}
	err := os.MkdirAll(d.Value, 0660)
	if err != nil {
		return &VidaError{Message: &String{Value: err.Error()}}, nil
	}
	return True, nil
}

func osRemove(ctx *Context, args ...Value) (Value, error) {
	if len(args) == 0 {
		return argError("os.rm", "expected a path string argument")
	}
	d, ok := args[0].(*String)
	if !ok {
		return argError("os.rm", fmt.Sprintf("expected a path string argument, got %s", args[0].Type()))
	}
	err := os.Remove(d.Value)
	if err != nil {
		return &VidaError{Message: &String{Value: err.Error()}}, nil
	}
	return True, nil
}

func osRemoveAll(ctx *Context, args ...Value) (Value, error) {
	if len(args) == 0 {
		return argError("os.rmAll", "expected a path string argument")
	}
	d, ok := args[0].(*String)
	if !ok {
		return argError("os.rmAll", fmt.Sprintf("expected a path string argument, got %s", args[0].Type()))
	}
	err := os.RemoveAll(d.Value)
	if err != nil {
		return &VidaError{Message: &String{Value: err.Error()}}, nil
	}
	return True, nil
}

func osName(ctx *Context, args ...Value) (Value, error) {
	return &String{Value: runtime.GOOS}, nil
}

func osArch(ctx *Context, args ...Value) (Value, error) {
	return &String{Value: runtime.GOARCH}, nil
}

func osRunCMD(ctx *Context, args ...Value) (Value, error) {
	l := len(args)
	if l == 0 {
		return argError("os.run", "expected a command string argument")
	}
	val, ok := args[0].(*String)
	if !ok {
		return argError("os.run", fmt.Sprintf("argument 1 (command) must be a string, got %s", args[0].Type()))
	}
	arr := make([]string, 0, l-1)
	for i := 1; i < l; i++ {
		v, ok := args[i].(*String)
		if !ok {
			return argError("os.run", fmt.Sprintf("argument %d must be a string, got %s", i+1, args[i].Type()))
		}
		arr = append(arr, v.Value)
	}
	cmd := exec.Command(val.Value, arr...)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return &VidaError{Message: &String{Value: err.Error()}}, nil
	}
	return True, nil
}
