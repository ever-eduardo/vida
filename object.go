package vida

func loadObjectLib() Value {
	m := &Object{Value: make(map[string]Value)}
	m.Value["inject"] = GFn(injectProps)
	m.Value["extract"] = GFn(extractProps)
	m.Value["override"] = GFn(injectAndOverrideProps)
	m.Value["check"] = GFn(checkProps)
	m.Value["delProp"] = GFn(deleteProperty)
	m.UpdateKeys()
	return m
}

func injectProps(args ...Value) (Value, error) {
	if len(args) > 1 {
		if o, ok := args[0].(*Object); ok {
			modified := false
			for _, v := range args[1:] {
				if m, ok := v.(*Object); ok && m != o {
					for k, x := range m.Value {
						if _, isPresent := o.Value[k]; !isPresent {
							modified = true
							o.Value[k] = x
						}
					}
				}
			}
			if modified {
				o.UpdateKeys()
			}
			return o, nil
		}
	}
	return NilValue, nil
}

func injectAndOverrideProps(args ...Value) (Value, error) {
	if len(args) > 1 {
		if o, ok := args[0].(*Object); ok {
			modified := false
			for _, v := range args[1:] {
				if m, ok := v.(*Object); ok && m != o {
					for k, x := range m.Value {
						modified = true
						o.Value[k] = x
					}
				}
			}
			if modified {
				o.UpdateKeys()
			}
			return o, nil
		}
	}
	return NilValue, nil
}

func checkProps(args ...Value) (Value, error) {
	if len(args) > 1 {
		if o, ok := args[0].(*Object); ok {
			for _, v := range args[1:] {
				if m, ok := v.(*Object); ok && m != o {
					for k := range m.Value {
						if _, isPresent := o.Value[k]; !isPresent {
							return Bool(false), nil
						}
					}
				}
			}
			return Bool(true), nil
		}
	}
	return NilValue, nil
}

func extractProps(args ...Value) (Value, error) {
	if len(args) > 1 {
		if o, ok := args[0].(*Object); ok {
			for _, v := range args[1:] {
				if m, ok := v.(*Object); ok && m != o {
					for k := range m.Value {
						delete(o.Value, k)
					}
				}
			}
			o.UpdateKeys()
			return o, nil
		}
	}
	return NilValue, nil
}

func deleteProperty(args ...Value) (Value, error) {
	if len(args) >= 2 {
		if o, ok := args[0].(*Object); ok {
			delete(o.Value, args[1].String())
			o.UpdateKeys()
		}
	}
	return NilValue, nil
}
