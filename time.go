package vida

import (
	"time"

	"github.com/alkemist-17/vida/token"
	"github.com/alkemist-17/vida/verror"
)

type Time time.Time

func (t Time) Boolean() Bool {
	return Bool(true)
}

func (t Time) Prefix(op uint64) (Value, error) {
	switch op {
	case uint64(token.NOT):
		return Bool(false), nil
	default:
		return NilValue, verror.ErrPrefixOpNotDefined
	}
}

func (t Time) Binop(op uint64, rhs Value) (Value, error) {
	switch op {
	case uint64(token.AND):
		return rhs, nil
	case uint64(token.OR):
		return t, nil
	case uint64(token.IN):
		return IsMemberOf(t, rhs)
	default:
		return NilValue, verror.ErrBinaryOpNotDefined
	}
}

func (t Time) IGet(index Value) (Value, error) {
	return NilValue, verror.ErrValueNotIndexable
}

func (t Time) ISet(index, val Value) error {
	return verror.ErrValueNotIndexable
}

func (t Time) Equals(other Value) Bool {
	if o, ok := other.(Time); ok {
		return Bool(time.Time(t).Equal(time.Time(o)))
	}
	return false
}

func (t Time) IsIterable() Bool {
	return false
}

func (t Time) IsCallable() Bool {
	return false
}

func (t Time) Call(args ...Value) (Value, error) {
	return NilValue, verror.ErrNotImplemented
}

func (t Time) Iterator() Value {
	return NilValue
}

func (t Time) String() string {
	return time.Time(t).String()
}

func (t Time) Type() string {
	return "time"
}

func (t Time) Clone() Value {
	return t
}

func loadFoundationTime() Value {
	m := &Object{Value: make(map[string]Value)}
	m.Value["unixNano"] = GFn(timestampNano)
	m.Value["unixMilli"] = GFn(timestampMilli)
	m.Value["unixMicro"] = GFn(timestampMicro)
	m.Value["unixSec"] = GFn(timestamp)
	m.Value["now"] = GFn(timeNow)
	m.Value["date"] = GFn(timeDate)
	m.Value["format"] = GFn(timeFormat)
	m.Value["getYear"] = GFn(timeGetYear)
	m.Value["getMonth"] = GFn(timeGetMonth)
	m.Value["getDay"] = GFn(timeGetDay)
	m.Value["getHours"] = GFn(timeGetHours)
	m.Value["getMinutes"] = GFn(timeGetMinutes)
	m.Value["getSeconds"] = GFn(timeGetSeconds)
	m.Value["getNanoseconds"] = GFn(timeGetNanoseconds)
	m.Value["getLocation"] = GFn(timeGetLocation)
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
	m.Value["Unix"] = &String{Value: time.UnixDate}
	m.Value["ANSIC"] = &String{Value: time.ANSIC}
	m.Value["DateTime"] = &String{Value: time.DateTime}
	m.Value["DateOnly"] = &String{Value: time.DateOnly}
	m.Value["TimeOnly"] = &String{Value: time.TimeOnly}
	m.Value["Kitchen"] = &String{Value: time.Kitchen}
	m.Value["Layout"] = &String{Value: time.Layout}
	m.Value["Stamp"] = &String{Value: time.Stamp}
	m.Value["StampMicro"] = &String{Value: time.StampMicro}
	m.Value["StampMilli"] = &String{Value: time.StampMilli}
	m.Value["StampNano"] = &String{Value: time.StampNano}
	m.Value["RubyDate"] = &String{Value: time.RubyDate}
	m.Value["nowIn"] = GFn(timeIn)
	m.Value["dateIn"] = GFn(dateIn)
	m.Value["toUnixNano"] = GFn(timeToUnixNano)
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
		return Time(time.Now()), nil
	case 1:
		if f, ok := args[0].(*String); ok && f.Value == time.Local.String() {
			return Time(time.Now().Local()), nil
		} else if ok && f.Value == time.UTC.String() {
			return Time(time.Now().UTC()), nil
		} else {
			r := time.Now().Format(f.Value)
			if len(r) > 0 {
				return &String{Value: r}, nil
			}
		}
	case 2:
		if f, ok := args[0].(*String); ok {
			if l, ok := args[1].(*String); ok {
				switch l.Value {
				case time.Local.String():
					return &String{Value: time.Now().Local().Format(f.Value)}, nil
				case time.UTC.String():
					return &String{Value: time.Now().UTC().Format(f.Value)}, nil
				}
			}
		}
	}
	return NilValue, nil
}

func timeDate(args ...Value) (Value, error) {
	switch len(args) {
	case 0:
		return Time(time.Now()), nil
	case 8:
		y, ok_0 := args[0].(Integer)
		m, ok_1 := args[1].(Integer)
		d, ok_2 := args[2].(Integer)
		h, ok_3 := args[3].(Integer)
		min, ok_4 := args[4].(Integer)
		sec, ok_5 := args[5].(Integer)
		nsec, ok_6 := args[6].(Integer)
		loc, ok_7 := args[7].(*String)
		if ok_0 && ok_1 && ok_2 && ok_3 && ok_4 && ok_5 && ok_6 && ok_7 {
			if loc.Value == time.Local.String() {
				return Time(time.Date(int(y), time.Month(m), int(d), int(h), int(min), int(sec), int(nsec), time.Local)), nil
			} else if loc.Value == time.UTC.String() {
				return Time(time.Date(int(y), time.Month(m), int(d), int(h), int(min), int(sec), int(nsec), time.UTC)), nil
			}
		}
	}
	return NilValue, nil
}

func timeFormat(args ...Value) (Value, error) {
	if len(args) > 1 {
		if t, ok := args[0].(Time); ok {
			if f, ok := args[1].(*String); ok {
				return &String{Value: time.Time(t).Format(f.Value)}, nil
			}
		}
	}
	return NilValue, nil
}

func timeGetYear(args ...Value) (Value, error) {
	if len(args) > 0 {
		if t, ok := args[0].(Time); ok {
			return Integer(time.Time(t).Year()), nil
		}
	}
	return NilValue, nil

}

func timeGetMonth(args ...Value) (Value, error) {
	if len(args) > 0 {
		if t, ok := args[0].(Time); ok {
			return Integer(time.Time(t).Month()), nil
		}
	}
	return NilValue, nil

}

func timeGetDay(args ...Value) (Value, error) {
	if len(args) > 0 {
		if t, ok := args[0].(Time); ok {
			return Integer(time.Time(t).Day()), nil
		}
	}
	return NilValue, nil

}

func timeGetHours(args ...Value) (Value, error) {
	if len(args) > 0 {
		if t, ok := args[0].(Time); ok {
			return Integer(time.Time(t).Hour()), nil
		}
	}
	return NilValue, nil

}

func timeGetMinutes(args ...Value) (Value, error) {
	if len(args) > 0 {
		if t, ok := args[0].(Time); ok {
			return Integer(time.Time(t).Minute()), nil
		}
	}
	return NilValue, nil

}

func timeGetSeconds(args ...Value) (Value, error) {
	if len(args) > 0 {
		if t, ok := args[0].(Time); ok {
			return Integer(time.Time(t).Second()), nil
		}
	}
	return NilValue, nil

}

func timeGetNanoseconds(args ...Value) (Value, error) {
	if len(args) > 0 {
		if t, ok := args[0].(Time); ok {
			return Integer(time.Time(t).Nanosecond()), nil
		}
	}
	return NilValue, nil

}

func timeGetLocation(args ...Value) (Value, error) {
	if len(args) > 0 {
		if t, ok := args[0].(Time); ok {
			return &String{Value: time.Time(t).Location().String()}, nil
		}
	}
	return NilValue, nil
}

func timeIn(args ...Value) (Value, error) {
	if len(args) > 0 {
		if zone, ok := args[0].(*String); ok {
			location, e := time.LoadLocation(zone.Value)
			if e != nil {
				return NilValue, nil
			}
			return Time(time.Now().In(location)), nil
		}
	}
	return Time(time.Now().UTC()), nil
}

func dateIn(args ...Value) (Value, error) {
	if len(args) > 1 {
		if t, ok := args[0].(Time); ok {
			if zone, ok := args[1].(*String); ok {
				location, e := time.LoadLocation(zone.Value)
				if e != nil {
					return NilValue, nil
				}
				return Time(time.Time(t).In(location)), nil
			}
		}
	}
	return NilValue, nil
}

func timeToUnixNano(args ...Value) (Value, error) {
	if len(args) > 0 {
		if t, ok := args[0].(Time); ok {
			return Integer(time.Time(t).UnixNano()), nil
		}
	}
	return NilValue, nil
}
