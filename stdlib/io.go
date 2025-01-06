package stdlib

import (
	"bufio"
	"errors"
	"fmt"
	"os"

	"github.com/ever-eduardo/vida"
	"github.com/ever-eduardo/vida/verror"
)

func generateIO() vida.Value {
	m := &vida.Object{Value: make(map[string]vida.Value)}
	m.Value["open"] = vida.GFn(ioOpenFile())
	m.Value["exists"] = vida.GFn(ioFileExists())
	m.Value["removeFile"] = vida.GFn(ioRemoveFile())
	m.Value["fileSize"] = vida.GFn(ioFileSize())
	m.Value["R"] = vida.Integer(os.O_RDONLY)
	m.Value["W"] = vida.Integer(os.O_WRONLY)
	m.Value["RW"] = vida.Integer(os.O_RDWR)
	m.Value["A"] = vida.Integer(os.O_APPEND)
	m.Value["C"] = vida.Integer(os.O_CREATE)
	m.Value["T"] = vida.Integer(os.O_TRUNC)
	m.Value["stdout"] = &FileHandler{Handler: os.Stdout, IsOpen: true}
	m.UpdateKeys()
	return m
}

const (
	fileHandlerName     = "handler"
	argIsNotFileHandler = "argument is not a FileHandler value"
	fileAlreadyClosed   = "file is closed"
)

// IO API
func ioOpenFile() vida.GFn {
	return func(args ...vida.Value) (vida.Value, error) {
		l := len(args)
		if l == 1 {
			if fname, ok := args[0].(*vida.String); ok {
				file, err := os.Create(fname.Value)
				if err != nil {
					return vida.Error{Message: &vida.String{Value: err.Error()}}, nil
				}
				o := &vida.Object{Value: make(map[string]vida.Value)}
				o.Value[fileHandlerName] = &FileHandler{Handler: file, IsOpen: true}
				o.Value["close"] = vida.GFn(fileClose())
				o.Value["isClosed"] = vida.GFn(fileIsClosed())
				o.Value["name"] = vida.GFn(fileName())
				o.Value["writeString"] = vida.GFn(fileWriteString())
				o.Value["readLines"] = vida.GFn(fileReadLines())
				o.UpdateKeys()
				return o, nil
			}
			return vida.NilValue, nil
		}
		if len(args) > 1 {
			if path, ok := args[0].(*vida.String); ok {
				if mode, ok := args[1].(vida.Integer); ok {
					file, err := os.OpenFile(path.Value, int(mode), 0666)
					if err != nil {
						return vida.Error{Message: &vida.String{Value: err.Error()}}, nil
					}
					o := &vida.Object{Value: make(map[string]vida.Value)}
					o.Value[fileHandlerName] = &FileHandler{Handler: file, IsOpen: true}
					o.Value["close"] = vida.GFn(fileClose())
					o.Value["isClosed"] = vida.GFn(fileIsClosed())
					o.Value["name"] = vida.GFn(fileName())
					o.Value["writeString"] = vida.GFn(fileWriteString())
					o.Value["readLines"] = vida.GFn(fileReadLines())
					o.UpdateKeys()
					return o, nil
				}
			}
			return vida.NilValue, nil
		}
		return vida.NilValue, nil
	}
}

func ioFileExists() vida.GFn {
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

func ioRemoveFile() vida.GFn {
	return func(args ...vida.Value) (vida.Value, error) {
		if len(args) > 0 {
			if path, ok := args[0].(*vida.String); ok {
				err := os.Remove(path.Value)
				if err != nil {
					return vida.Error{Message: &vida.String{Value: err.Error()}}, nil
				}
				return &vida.String{Value: string(vida.Success)}, nil
			}
			return vida.NilValue, nil
		}
		return vida.NilValue, nil
	}
}

func ioFileSize() vida.GFn {
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

// Type FileHandler is a wrap over *os.File
type FileHandler struct {
	Handler *os.File
	IsOpen  bool
}

// Implementation of the interface vida.Value
func (file *FileHandler) Boolean() vida.Bool {
	return vida.Bool(file.IsOpen)
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
	if file.IsOpen {
		return fmt.Sprintf("fileHandler(%v)", file.Handler.Fd())
	}
	return "fileHandler(closed)"
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
					if file.IsOpen {
						err := file.Handler.Close()
						file.IsOpen = false
						if err != nil {
							return vida.Error{Message: &vida.String{Value: err.Error()}}, nil
						}
						return &vida.String{Value: string(vida.Success)}, nil
					} else {
						return vida.Error{Message: &vida.String{Value: fileAlreadyClosed}}, nil
					}
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
					return vida.Bool(!file.IsOpen), nil
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
					if file.IsOpen {
						scanner := bufio.NewScanner(file.Handler)
						var data []string
						for scanner.Scan() {
							data = append(data, scanner.Text())
						}
						if err := scanner.Err(); err != nil {
							return vida.Error{Message: &vida.String{Value: err.Error()}}, nil
						}
						xs := &vida.List{}
						for _, v := range data {
							xs.Value = append(xs.Value, &vida.String{Value: v})
						}
						return xs, nil
					} else {
						return vida.Error{Message: &vida.String{Value: fileAlreadyClosed}}, nil
					}
				}
				return vida.Error{Message: &vida.String{Value: argIsNotFileHandler}}, nil
			}
		}
		return vida.NilValue, nil
	}
}

func fileWriteString() vida.GFn {
	return func(args ...vida.Value) (vida.Value, error) {
		if len(args) > 1 {
			if obj, ok := args[0].(*vida.Object); ok {
				if file, ok := obj.Value[fileHandlerName].(*FileHandler); ok {
					if data, ok := args[1].(*vida.String); ok {
						if file.IsOpen {
							i, err := file.Handler.WriteString(data.Value)
							if err != nil {
								return vida.Error{Message: &vida.String{Value: err.Error()}}, nil
							}
							return vida.Integer(i), nil
						} else {
							return vida.Error{Message: &vida.String{Value: fileAlreadyClosed}}, nil
						}
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
