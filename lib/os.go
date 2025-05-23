package lib

import (
	"os"
	"os/exec"
	"runtime"

	"github.com/alkemist-17/vida"
)

func generateOS() vida.Value {
	m := &vida.Object{Value: make(map[string]vida.Value)}
	m.Value["args"] = vida.GFn(args())
	m.Value["env"] = vida.GFn(environ())
	m.Value["exit"] = vida.GFn(exit())
	m.Value["getFromEnv"] = vida.GFn(getEnv())
	m.Value["pwd"] = vida.GFn(getWD())
	m.Value["hostname"] = vida.GFn(hostname())
	m.Value["pathSeparator"] = &vida.String{Value: string(os.PathSeparator)}
	m.Value["mkdir"] = vida.GFn(mkdir())
	m.Value["mkdirAll"] = vida.GFn(mkdirAll())
	m.Value["rm"] = vida.GFn(rm())
	m.Value["rmAll"] = vida.GFn(rmAll())
	m.Value["name"] = vida.GFn(osName())
	m.Value["arch"] = vida.GFn(osArch())
	m.Value["run"] = vida.GFn(runCMD())
	m.Value["stdin"] = &FileHandler{Handler: os.Stdin}
	m.Value["stdout"] = &FileHandler{Handler: os.Stdout}
	m.Value["stderr"] = &FileHandler{Handler: os.Stderr}
	m.UpdateKeys()
	return m
}

func args() vida.GFn {
	return func(args ...vida.Value) (vida.Value, error) {
		xs := &vida.List{}
		for _, v := range os.Args {
			xs.Value = append(xs.Value, &vida.String{Value: v})
		}
		return xs, nil
	}
}

func environ() vida.GFn {
	return func(args ...vida.Value) (vida.Value, error) {
		xs := &vida.List{}
		for _, v := range os.Environ() {
			xs.Value = append(xs.Value, &vida.String{Value: v})
		}
		return xs, nil
	}
}

func exit() vida.GFn {
	return func(args ...vida.Value) (vida.Value, error) {
		os.Exit(0)
		return vida.NilValue, nil
	}
}

func getEnv() vida.GFn {
	return func(args ...vida.Value) (vida.Value, error) {
		if len(args) > 0 {
			if val, ok := args[0].(*vida.String); ok {
				xs := make([]vida.Value, 0, 2)
				if r, ok := os.LookupEnv(val.Value); ok {
					xs = append(xs, &vida.String{Value: r})
					xs = append(xs, vida.Bool(ok))
				} else {
					xs = append(xs, &vida.String{Value: ""})
					xs = append(xs, vida.Bool(ok))
				}
				return &vida.List{Value: xs}, nil
			}
		}
		return vida.NilValue, nil
	}
}

func getWD() vida.GFn {
	return func(args ...vida.Value) (vida.Value, error) {
		if d, e := os.Getwd(); e == nil {
			return &vida.String{Value: d}, nil
		} else {
			return vida.Error{Message: &vida.String{Value: e.Error()}}, nil
		}
	}
}

func hostname() vida.GFn {
	return func(args ...vida.Value) (vida.Value, error) {
		if h, e := os.Hostname(); e == nil {
			return &vida.String{Value: h}, nil
		} else {
			return vida.Error{Message: &vida.String{Value: e.Error()}}, nil
		}
	}
}

func mkdir() vida.GFn {
	return func(args ...vida.Value) (vida.Value, error) {
		if len(args) > 0 {
			if d, ok := args[0].(*vida.String); ok {
				err := os.Mkdir(d.Value, 0660)
				if err != nil && !os.IsExist(err) {
					return vida.Error{Message: &vida.String{Value: err.Error()}}, nil
				}
				return Success, nil
			}
		}
		return vida.NilValue, nil
	}
}

func mkdirAll() vida.GFn {
	return func(args ...vida.Value) (vida.Value, error) {
		if len(args) > 0 {
			if d, ok := args[0].(*vida.String); ok {
				err := os.MkdirAll(d.Value, 0660)
				if err != nil {
					return vida.Error{Message: &vida.String{Value: err.Error()}}, nil
				}
				return Success, nil
			}
		}
		return vida.NilValue, nil
	}
}

func rm() vida.GFn {
	return func(args ...vida.Value) (vida.Value, error) {
		if len(args) > 0 {
			if d, ok := args[0].(*vida.String); ok {
				err := os.Remove(d.Value)
				if err != nil {
					return vida.Error{Message: &vida.String{Value: err.Error()}}, nil
				}
				return Success, nil
			}
		}
		return vida.NilValue, nil
	}
}

func rmAll() vida.GFn {
	return func(args ...vida.Value) (vida.Value, error) {
		if len(args) > 0 {
			if d, ok := args[0].(*vida.String); ok {
				err := os.RemoveAll(d.Value)
				if err != nil {
					return vida.Error{Message: &vida.String{Value: err.Error()}}, nil
				}
				return Success, nil
			}
		}
		return vida.NilValue, nil
	}
}

func osName() vida.GFn {
	return func(args ...vida.Value) (vida.Value, error) {
		return &vida.String{Value: runtime.GOOS}, nil
	}
}

func osArch() vida.GFn {
	return func(args ...vida.Value) (vida.Value, error) {
		return &vida.String{Value: runtime.GOARCH}, nil
	}
}

func runCMD() vida.GFn {
	return func(args ...vida.Value) (vida.Value, error) {
		l := len(args)
		if l > 0 {
			if val, ok := args[0].(*vida.String); ok {
				var arr []string
				for i := 1; i < l; i++ {
					if v, ok := args[i].(*vida.String); ok {
						arr = append(arr, v.Value)
					}
				}
				cmd := exec.Command(val.Value, arr...)
				cmd.Stdout = os.Stdout
				cmd.Stdin = os.Stdin
				cmd.Stderr = os.Stderr
				err := cmd.Run()
				if err != nil {
					return &vida.Error{Message: &vida.String{Value: err.Error()}}, nil
				}
				return Success, nil
			}
		}
		return vida.NilValue, nil
	}
}
