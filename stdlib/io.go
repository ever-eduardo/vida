package stdlib

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/ever-eduardo/vida"
	"github.com/ever-eduardo/vida/verror"
)

func generateIO() vida.Value {
	m := &vida.Object{Value: make(map[string]vida.Value)}
	m.Value["open"] = vida.GFn(openFile())
	m.Value["create"] = vida.GFn(createFile())
	m.Value["exists"] = vida.GFn(exists())
	m.Value["remove"] = vida.GFn(remove())
	m.Value["size"] = vida.GFn(fsize())
	m.Value["fprint"] = vida.GFn(fprint())
	m.Value["fprintf"] = vida.GFn(fprintf())
	m.Value["isFile"] = vida.GFn(isFile())
	m.Value["createTemp"] = vida.GFn(tempfile())
	m.Value["tempDir"] = &vida.String{Value: os.TempDir()}
	m.Value["ok"] = success
	m.Value["R"] = vida.Integer(os.O_RDONLY)
	m.Value["W"] = vida.Integer(os.O_WRONLY)
	m.Value["RW"] = vida.Integer(os.O_RDWR)
	m.Value["A"] = vida.Integer(os.O_APPEND)
	m.Value["C"] = vida.Integer(os.O_CREATE)
	m.Value["T"] = vida.Integer(os.O_TRUNC)
	m.Value["stdin"] = &FileHandler{Handler: os.Stdin}
	m.Value["stdout"] = &FileHandler{Handler: os.Stdout}
	m.Value["stderr"] = &FileHandler{Handler: os.Stderr}
	m.UpdateKeys()
	return m
}

const (
	fileHandlerName     = "handler"
	argIsNotFileHandler = "argument is not a FileHandler value"
	fileAlreadyClosed   = "file is already closed"
	noStringFormat      = "no string format given"
	noOrClosedFH        = argIsNotFileHandler + " or " + fileAlreadyClosed
	expectedBytes       = "expected a value of type bytes"
)

var success = &vida.String{Value: string(vida.Success)}

func generateFileHandlerObject(file *os.File) vida.Value {
	o := &vida.Object{Value: make(map[string]vida.Value)}
	o.Value[fileHandlerName] = &FileHandler{Handler: file}
	o.Value["close"] = vida.GFn(fileClose())
	o.Value["isClosed"] = vida.GFn(fileIsClosed())
	o.Value["name"] = vida.GFn(fileName())
	o.Value["write"] = vida.GFn(fileWrite())
	o.Value["lines"] = vida.GFn(fileReadLines())
	o.Value["read"] = vida.GFn(fileRead())
	o.UpdateKeys()
	return o
}

// File API
func openFile() vida.GFn {
	return func(args ...vida.Value) (vida.Value, error) {
		l := len(args)
		if l == 1 {
			if fname, ok := args[0].(*vida.String); ok {
				file, err := os.OpenFile(fname.Value, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
				if err != nil {
					file.Close()
					return vida.Error{Message: &vida.String{Value: err.Error()}}, nil
				}
				return generateFileHandlerObject(file), nil
			}
			return vida.NilValue, nil
		}
		if len(args) > 1 {
			if path, ok := args[0].(*vida.String); ok {
				if mode, ok := args[1].(vida.Integer); ok {
					file, err := os.OpenFile(path.Value, int(mode), 0666)
					if err != nil {
						file.Close()
						return vida.Error{Message: &vida.String{Value: err.Error()}}, nil
					}
					return generateFileHandlerObject(file), nil
				}
			}
			return vida.NilValue, nil
		}
		return vida.NilValue, nil
	}
}

func createFile() vida.GFn {
	return func(args ...vida.Value) (vida.Value, error) {
		if len(args) > 0 {
			if fname, ok := args[0].(*vida.String); ok {
				file, err := os.Create(fname.Value)
				if err != nil {
					file.Close()
					return vida.Error{Message: &vida.String{Value: err.Error()}}, nil
				}
				return generateFileHandlerObject(file), nil
			}
			return vida.NilValue, nil
		}
		return vida.NilValue, nil
	}
}

func exists() vida.GFn {
	return func(args ...vida.Value) (vida.Value, error) {
		if len(args) > 0 {
			if path, ok := args[0].(*vida.String); ok {
				_, err := os.Stat(path.Value)
				if errors.Is(err, os.ErrNotExist) {
					return vida.Bool(false), nil
				}
				return vida.Bool(true), nil
			}
			return vida.NilValue, nil
		}
		return vida.NilValue, nil
	}
}

func remove() vida.GFn {
	return func(args ...vida.Value) (vida.Value, error) {
		if len(args) > 0 {
			if path, ok := args[0].(*vida.String); ok {
				err := os.Remove(path.Value)
				if err != nil {
					return vida.Error{Message: &vida.String{Value: err.Error()}}, nil
				}
				return success, nil
			}
			return vida.NilValue, nil
		}
		return vida.NilValue, nil
	}
}

func fsize() vida.GFn {
	return func(args ...vida.Value) (vida.Value, error) {
		if len(args) > 0 {
			if path, ok := args[0].(*vida.String); ok {
				fileInfo, err := os.Stat(path.Value)
				if errors.Is(err, os.ErrNotExist) {
					return vida.Error{Message: &vida.String{Value: err.Error()}}, nil
				}
				return vida.Integer(fileInfo.Size()), nil
			}
			return vida.NilValue, nil
		}
		return vida.NilValue, nil
	}
}

func fprint() vida.GFn {
	return func(args ...vida.Value) (vida.Value, error) {
		if len(args) > 1 {
			switch handler := args[0].(type) {
			case *vida.Object:
				if fileHandler, ok := handler.Value[fileHandlerName].(*FileHandler); ok && !fileHandler.IsClosed {
					var s []any
					for _, v := range args[1:] {
						s = append(s, v)
					}
					n, err := fmt.Fprint(fileHandler.Handler, s...)
					if err != nil {
						fileHandler.IsClosed = true
						fileHandler.Handler.Close()
						return vida.Error{Message: &vida.String{Value: err.Error()}}, nil
					}
					return vida.Integer(n), nil
				}
				return vida.Error{Message: &vida.String{Value: noOrClosedFH}}, nil
			case *FileHandler:
				if handler.IsClosed {
					return vida.Error{Message: &vida.String{Value: fileAlreadyClosed}}, nil
				}
				var s []any
				for _, v := range args[1:] {
					s = append(s, v)
				}
				n, err := fmt.Fprint(handler.Handler, s...)
				if err != nil {
					handler.IsClosed = true
					handler.Handler.Close()
					return vida.Error{Message: &vida.String{Value: err.Error()}}, nil
				}
				return vida.Integer(n), nil
			}
		}
		return vida.NilValue, nil
	}
}

func fprintf() vida.GFn {
	return func(args ...vida.Value) (vida.Value, error) {
		if len(args) > 2 {
			switch handler := args[0].(type) {
			case *vida.Object:
				if fileHandler, ok := handler.Value[fileHandlerName].(*FileHandler); ok && !fileHandler.IsClosed {
					if formatstr, ok := args[1].(*vida.String); ok {
						var s []any
						for _, v := range args[2:] {
							s = append(s, v)
						}
						n, err := fmt.Fprintf(fileHandler.Handler, formatstr.Value, s...)
						if err != nil {
							fileHandler.IsClosed = true
							fileHandler.Handler.Close()
							return vida.Error{Message: &vida.String{Value: err.Error()}}, nil
						}
						return vida.Integer(n), nil
					}
					return vida.Error{Message: &vida.String{Value: noStringFormat}}, nil
				}
				return vida.Error{Message: &vida.String{Value: noOrClosedFH}}, nil
			case *FileHandler:
				if formatstr, ok := args[1].(*vida.String); ok {
					if handler.IsClosed {
						return vida.Error{Message: &vida.String{Value: fileAlreadyClosed}}, nil
					}
					var s []any
					for _, v := range args[2:] {
						s = append(s, v)
					}
					n, err := fmt.Fprintf(handler.Handler, formatstr.Value, s...)
					if err != nil {
						handler.IsClosed = true
						handler.Handler.Close()
						return vida.Error{Message: &vida.String{Value: err.Error()}}, nil
					}
					return vida.Integer(n), nil
				}
				return vida.Error{Message: &vida.String{Value: noStringFormat}}, nil
			}
		}
		return vida.NilValue, nil
	}
}

func isFile() vida.GFn {
	return func(args ...vida.Value) (vida.Value, error) {
		if len(args) > 0 {
			if _, ok := args[0].(*FileHandler); ok {
				return vida.Bool(ok), nil
			}
			return vida.Bool(false), nil
		}
		return vida.NilValue, nil
	}
}

func tempfile() vida.GFn {
	return func(args ...vida.Value) (vida.Value, error) {
		if len(args) > 1 {
			if dir, ok := args[0].(*vida.String); ok {
				if pattern, ok := args[1].(*vida.String); ok {
					f, err := os.CreateTemp(dir.Value, pattern.Value)
					if err != nil {
						f.Close()
						return vida.Error{Message: &vida.String{Value: err.Error()}}, nil
					}
					return generateFileHandlerObject(f), nil
				}
			}
		}
		return vida.NilValue, nil
	}
}

// Type FileHandler is a wrap over *os.File
type FileHandler struct {
	vida.ReferenceSemanticsImpl
	Handler  *os.File
	IsClosed bool
}

// Implementation of the interface vida.Value
func (file *FileHandler) Boolean() vida.Bool {
	return vida.Bool(!file.IsClosed)
}

func (file *FileHandler) Prefix(uint64) (vida.Value, error) {
	return vida.NilValue, verror.ErrPrefixOpNotDefined
}

func (file *FileHandler) Binop(uint64, vida.Value) (vida.Value, error) {
	return vida.NilValue, verror.ErrBinaryOpNotDefined
}

func (file *FileHandler) IGet(vida.Value) (vida.Value, error) {
	return vida.NilValue, verror.ErrValueNotIndexable
}

func (file *FileHandler) ISet(vida.Value, vida.Value) error {
	return verror.ErrValueIsConstant
}

func (file *FileHandler) Equals(other vida.Value) vida.Bool {
	if v, ok := other.(*FileHandler); ok {
		return v.Handler.Fd() == file.Handler.Fd()
	}
	return vida.Bool(false)
}

func (file *FileHandler) IsIterable() vida.Bool {
	return false
}

func (file *FileHandler) Iterator() vida.Value {
	return vida.NilValue
}

func (file *FileHandler) IsCallable() vida.Bool {
	return false
}

func (file *FileHandler) String() string {
	if file.IsClosed {
		return "fileHandler(closed)"
	}
	return fmt.Sprintf("fileHandler(%v)", file.Handler.Fd())
}

func (file *FileHandler) Type() string {
	return "fileHandler"
}

func (file *FileHandler) Clone() vida.Value {
	return file
}

// FileHandler Methods
func fileClose() vida.GFn {
	return func(args ...vida.Value) (vida.Value, error) {
		if len(args) > 0 {
			if obj, ok := args[0].(*vida.Object); ok {
				if file, ok := obj.Value[fileHandlerName].(*FileHandler); ok {
					if file.Handler.Fd() == os.Stdout.Fd() ||
						file.Handler.Fd() == os.Stdin.Fd() ||
						file.Handler.Fd() == os.Stderr.Fd() {
						return vida.Error{Message: &vida.String{Value: "cannot close file open system files"}}, nil
					}
					if file.IsClosed {
						return vida.Error{Message: &vida.String{Value: fileAlreadyClosed}}, nil
					}
					err := file.Handler.Close()
					file.IsClosed = true
					if err != nil {
						return vida.Error{Message: &vida.String{Value: err.Error()}}, nil
					}
					return success, nil
				}
				return vida.Error{Message: &vida.String{Value: argIsNotFileHandler}}, nil
			}
		}
		return vida.NilValue, nil
	}
}

func fileIsClosed() vida.GFn {
	return func(args ...vida.Value) (vida.Value, error) {
		if len(args) > 0 {
			if obj, ok := args[0].(*vida.Object); ok {
				if file, ok := obj.Value[fileHandlerName].(*FileHandler); ok {
					return vida.Bool(file.IsClosed), nil
				}
				return vida.Error{Message: &vida.String{Value: argIsNotFileHandler}}, nil
			}
		}
		return vida.NilValue, nil
	}
}

func fileName() vida.GFn {
	return func(args ...vida.Value) (vida.Value, error) {
		if len(args) > 0 {
			if obj, ok := args[0].(*vida.Object); ok {
				if file, ok := obj.Value[fileHandlerName].(*FileHandler); ok {
					return &vida.String{Value: file.Handler.Name()}, nil
				}
				return vida.Error{Message: &vida.String{Value: argIsNotFileHandler}}, nil
			}
		}
		return vida.NilValue, nil
	}
}

func fileReadLines() vida.GFn {
	return func(args ...vida.Value) (vida.Value, error) {
		if len(args) > 0 {
			if obj, ok := args[0].(*vida.Object); ok {
				if file, ok := obj.Value[fileHandlerName].(*FileHandler); ok {
					if file.IsClosed {
						return vida.Error{Message: &vida.String{Value: fileAlreadyClosed}}, nil
					}
					scanner := bufio.NewScanner(file.Handler)
					var data []string
					for scanner.Scan() {
						data = append(data, scanner.Text())
					}
					if err := scanner.Err(); err != nil {
						file.IsClosed = true
						file.Handler.Close()
						return vida.Error{Message: &vida.String{Value: err.Error()}}, nil
					}
					xs := &vida.List{}
					for _, v := range data {
						xs.Value = append(xs.Value, &vida.String{Value: v})
					}
					return xs, nil
				}
				return vida.Error{Message: &vida.String{Value: argIsNotFileHandler}}, nil
			}
		}
		return vida.NilValue, nil
	}
}

func fileRead() vida.GFn {
	return func(args ...vida.Value) (vida.Value, error) {
		if len(args) > 1 {
			if obj, ok := args[0].(*vida.Object); ok {
				if file, ok := obj.Value[fileHandlerName].(*FileHandler); ok {
					if file.IsClosed {
						return vida.Error{Message: &vida.String{Value: fileAlreadyClosed}}, nil
					}
					if b, ok := args[1].(*vida.Bytes); ok {
						n, err := file.Handler.Read(b.Value)
						if err != nil && !errors.Is(err, io.EOF) {
							file.Handler.Close()
							file.IsClosed = true
							return vida.Error{Message: &vida.String{Value: err.Error()}}, nil
						}
						return vida.Integer(n), nil
					}
					return vida.Error{Message: &vida.String{Value: expectedBytes}}, nil
				}
				return vida.Error{Message: &vida.String{Value: argIsNotFileHandler}}, nil
			}
		}
		return vida.NilValue, nil
	}
}

func fileWrite() vida.GFn {
	return func(args ...vida.Value) (vida.Value, error) {
		if len(args) > 1 {
			if obj, ok := args[0].(*vida.Object); ok {
				if file, ok := obj.Value[fileHandlerName].(*FileHandler); ok {
					if file.IsClosed {
						return vida.Error{Message: &vida.String{Value: fileAlreadyClosed}}, nil
					}
					if data, ok := args[1].(*vida.String); ok {
						i, err := file.Handler.WriteString(data.Value)
						if err != nil {
							file.IsClosed = true
							file.Handler.Close()
							return vida.Error{Message: &vida.String{Value: err.Error()}}, nil
						}
						return vida.Integer(i), nil
					} else if data, ok := args[1].(*vida.Bytes); ok {
						i, err := file.Handler.Write(data.Value)
						if err != nil {
							file.IsClosed = true
							file.Handler.Close()
							return vida.Error{Message: &vida.String{Value: err.Error()}}, nil
						}
						return vida.Integer(i), nil
					} else {
						return vida.Error{Message: &vida.String{Value: "expected data of type string"}}, nil
					}
				}
				return vida.Error{Message: &vida.String{Value: argIsNotFileHandler}}, nil
			}
		}
		return vida.NilValue, nil
	}
}
