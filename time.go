package vida

import (
	"time"
)

func loadFoundationTime() Value {
	m := &Object{Value: make(map[string]Value)}
	m.Value["timestamp"] = GFn(timestampNano)
	m.Value["epochNano"] = GFn(timestampNano)
	m.Value["epochMilli"] = GFn(timestampMilli)
	m.Value["epochMicro"] = GFn(timestampMicro)
	m.Value["epochSec"] = GFn(timestamp)
	m.Value["now"] = GFn(timeNow)
	m.Value["sleep"] = GFn(sleep)
	m.Value["millisecond"] = Integer(time.Millisecond)
	m.Value["nanosecond"] = Integer(time.Nanosecond)
	m.Value["microsecond"] = Integer(time.Microsecond)
	m.Value["hour"] = Integer(time.Hour)
	m.Value["minute"] = Integer(time.Minute)
	m.Value["second"] = Integer(time.Second)
	m.Value["RFC3339"] = &String{Value: time.RFC3339}
	m.Value["RFC3339Nano"] = &String{Value: time.RFC3339Nano}
	m.Value["RFC1123"] = &String{Value: time.RFC1123}
	m.Value["RFC1123Z"] = &String{Value: time.RFC1123Z}
	m.Value["RFC822"] = &String{Value: time.RFC822}
	m.Value["RFC822Z"] = &String{Value: time.RFC822Z}
	m.Value["RFC850"] = &String{Value: time.RFC850}
	m.Value["Epoch"] = &String{Value: time.UnixDate}
	m.Value["Local"] = &String{Value: time.Local.String()}
	m.Value["UTC"] = &String{Value: time.UTC.String()}
	m.UpdateKeys()
	return m
}

func sleep(args ...Value) (Value, error) {
	if len(args) > 0 {
		val, ok := args[0].(Integer)
		if ok {
			time.Sleep(time.Duration(val))
		}
	}
	return NilValue, nil
}

func timestampNano(args ...Value) (Value, error) {
	return Integer(time.Now().UnixNano()), nil
}

func timestampMilli(args ...Value) (Value, error) {
	return Integer(time.Now().UnixMilli()), nil
}

func timestampMicro(args ...Value) (Value, error) {
	return Integer(time.Now().UnixMicro()), nil
}

func timestamp(args ...Value) (Value, error) {
	return Integer(time.Now().Unix()), nil
}

func timeNow(args ...Value) (Value, error) {
	switch len(args) {
	case 0:
		return &String{Value: time.Now().Format(time.RFC3339)}, nil
	case 1:
		if f, ok := args[0].(*String); ok && f.Value == time.Local.String() {
			return &String{Value: time.Now().Local().Format(time.RFC3339)}, nil
		} else if ok && f.Value == time.UTC.String() {
			return &String{Value: time.Now().UTC().Format(time.RFC3339)}, nil
		} else {
			r := time.Now().Format(f.Value)
			if len(r) > 0 {
				return &String{Value: r}, nil
			}
		}
	}
	return NilValue, nil
}
