package vida

import (
	"time"
)

func loadFoundationTime() Value {
	m := &Object{Value: make(map[string]Value)}
	m.Value["sleep"] = GFn(sleep)
	m.Value["millisecond"] = Integer(time.Millisecond)
	m.Value["nanosecond"] = Integer(time.Nanosecond)
	m.Value["microsecond"] = Integer(time.Microsecond)
	m.Value["hour"] = Integer(time.Hour)
	m.Value["minute"] = Integer(time.Minute)
	m.Value["second"] = Integer(time.Second)
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
