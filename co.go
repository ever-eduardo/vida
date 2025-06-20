package vida

func loadFoundationCoroutine() Value {
	m := &Object{Value: make(map[string]Value)}
	m.Value["new"] = GFn(gfnNewThread)
	m.Value["stack"] = getStackSizes()
	m.Value["state"] = GFn(gfnGetThreadState)
	m.Value["ready"] = Integer(Ready)
	m.Value["running"] = Integer(Running)
	m.Value["suspended"] = Integer(Suspended)
	m.Value["waiting"] = Integer(Waiting)
	m.Value["closed"] = Integer(Closed)
	m.UpdateKeys()
	return m
}

func gfnNewThread(args ...Value) (Value, error) {
	l := len(args)
	if l == 1 {
		if fn, ok := args[0].(*Function); ok {
			return newThread(fn, ((*clbu)[mainThIndex].(*Thread)).Script, primeStack), nil
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

func gfnGetThreadState(args ...Value) (Value, error) {
	if len(args) > 0 {
		if th, ok := args[0].(*Thread); ok {
			return Integer(th.State), nil
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
	m.Value["of23"] = Integer(primeStack)
	m.Value["of16"] = Integer(picoStack)
	m.Value["of8"] = Integer(femtoStack)
	m.UpdateKeys()
	return m
}
