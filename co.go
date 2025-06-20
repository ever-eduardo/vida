package vida

func loadFoundationCoroutine() Value {
	m := &Object{Value: make(map[string]Value)}
	m.Value["new"] = GFn(gfnNewThread)
	m.Value["stack"] = getStackSizes()
	m.UpdateKeys()
	return m
}

func gfnNewThread(args ...Value) (Value, error) {
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
