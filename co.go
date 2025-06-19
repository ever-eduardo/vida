package vida

func loadFoundationCoroutine() Value {
	m := &Object{Value: make(map[string]Value)}
	m.Value["new"] = GFn(gfnNewThread)
	m.UpdateKeys()
	return m
}

func gfnNewThread(args ...Value) (Value, error) {
	return NilValue, nil
}
