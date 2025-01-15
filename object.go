package vida

const objectLibName = "object"

func loadObjectLib() Value {
	m := &Object{Value: make(map[string]Value)}
	m.Value["inject"] = injectProps()
	m.Value["extend"] = extendObject()
	m.Value["override"] = injectAndOverrideProps()
	m.Value["check"] = checkProps()
	m.UpdateKeys()
	return m
}

func injectProps() GFn {
	return func(args ...Value) (Value, error) {
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
}

func injectAndOverrideProps() GFn {
	return func(args ...Value) (Value, error) {
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
}

func checkProps() GFn {
	return func(args ...Value) (Value, error) {
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
}

func extendObject() GFn {
	return func(args ...Value) (Value, error) {
		if len(args) > 1 {
			extension := make(map[string]Value)
			for _, val := range args {
				if obj, ok := val.(*Object); ok {
					for k, v := range obj.Value {
						if _, isPresent := extension[k]; !isPresent {
							extension[k] = v
						}
					}
				}
			}
			return &Object{Value: extension}, nil
		}
		return NilValue, nil
	}
}
