package vida

type lKey struct {
	id            string
	compilerLevel int
	scopeLevel    int
}

type LocalBuilder struct {
	Map     map[lKey]byte
	History []lKey
	auxKey  lKey
}

func NewLocalBuilder() *LocalBuilder {
	return &LocalBuilder{
		Map:    make(map[lKey]byte),
		auxKey: lKey{},
	}
}

func (lb *LocalBuilder) IsLocal(id string, compilerLevel int, scopeLevel int) (byte, bool) {
	lb.auxKey.compilerLevel = compilerLevel
	lb.auxKey.id = id
	lb.auxKey.scopeLevel = scopeLevel
	b, ok := lb.Map[lb.auxKey]
	return b, ok
}

func (lb *LocalBuilder) AddLocal(id string, compilerLevel int, scopeLevel int, register byte) {
	lb.Map[lKey{id: id, compilerLevel: compilerLevel, scopeLevel: scopeLevel}] = register
}

func (lb *LocalBuilder) CleanScope(compilerLevel int, scopeLevel int) int {
	var count int
	length := len(lb.History)
	for i := length - 1; i >= 0; i-- {
		if lb.History[i].compilerLevel == compilerLevel && lb.History[i].scopeLevel == scopeLevel {
			count++
			delete(lb.Map, lb.History[i])
		} else {
			break
		}
	}
	lb.History = lb.History[:length-count]
	return count
}
