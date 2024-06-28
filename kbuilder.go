package vida

type KonstBuilder struct {
	stringMap  map[string]int
	booleanMap map[bool]int
	index      int
	nilIndex   int
	Konstants  []Value
}

func newKonstBuilder() *KonstBuilder {
	return &KonstBuilder{
		stringMap:  make(map[string]int),
		booleanMap: make(map[bool]int),
		nilIndex:   -1,
	}
}

func (kb *KonstBuilder) StringIndex(value string) int {
	idx, isPresent := kb.stringMap[value]
	if isPresent {
		return idx
	}
	i := kb.index
	kb.Konstants = append(kb.Konstants, value)
	kb.stringMap[value] = i
	kb.index++
	return i
}

func (kb *KonstBuilder) BooleanIndex(value bool) int {
	idx, isPresent := kb.booleanMap[value]
	if isPresent {
		return idx
	}
	i := kb.index
	kb.Konstants = append(kb.Konstants, value)
	kb.booleanMap[value] = i
	kb.index++
	return i
}

func (kb *KonstBuilder) NilIndex() int {
	if kb.nilIndex != -1 {
		return kb.nilIndex
	}
	kb.nilIndex = kb.index
	kb.Konstants = append(kb.Konstants, globalNil)
	kb.index++
	return kb.nilIndex
}
