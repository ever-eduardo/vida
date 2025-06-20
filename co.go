package vida

func loadFoundationCoroutine() Value {
	m := &Object{Value: make(map[string]Value)}
	m.Value["new"] = GFn(gfnNewThread)
	m.Value["stack"] = getStackSizes()
	m.UpdateKeys()
	return m
}

func gfnNewThread(args ...Value) (Value, error) {
	l := len(args)
	if l == 1 {
		if fn, ok := args[0].(*Function); ok {
			return newThread(fn, ((*clbu)[mainThIndex].(*Thread)).Script, femtoStack), nil
		}
	} else if l > 1 {
		if fn, ok := args[0].(*Function); ok {
			if s, ok := args[1].(Integer); ok && femtoStack <= s && s <= fullStack {
				return newThread(fn, ((*clbu)[mainThIndex].(*Thread)).Script, int(s)), nil
			}
		}
	}
	return NilValue, nil
}

func getStackSizes() *Object {
	m := &Object{Value: make(map[string]Value)}
	m.Value["of1024"] = Integer(fullStack)
	m.Value["of512"] = Integer(halfStack)
	m.Value["of256"] = Integer(quarterStack)
	m.Value["of128"] = Integer(microStack)
	m.Value["of64"] = Integer(milliStack)
	m.Value["of32"] = Integer(nanoStack)
	m.Value["of16"] = Integer(picoStack)
	m.Value["of8"] = Integer(femtoStack)
	m.UpdateKeys()
	return m
}
