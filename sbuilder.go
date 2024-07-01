package vida

var dummy struct{}

type lKey struct {
	id            string
	compilerLevel int
	scopeLevel    int
}

type symbolBuilder struct {
	GlobalSet  map[string]struct{}
	LocalMap   map[lKey]byte
	History    []lKey
	privateKey lKey
}

func newSymbolBuilder() *symbolBuilder {
	return &symbolBuilder{
		GlobalSet:  make(map[string]struct{}),
		LocalMap:   make(map[lKey]byte),
		privateKey: lKey{},
	}
}

func (sb *symbolBuilder) isLocal(id string, compilerLevel int, scopeLevel int) (reg byte, isLocal bool) {
	sb.privateKey.compilerLevel = compilerLevel
	sb.privateKey.id = id
	sb.privateKey.scopeLevel = scopeLevel
	reg, isLocal = sb.LocalMap[sb.privateKey]
	return
}

func (sb *symbolBuilder) isGlobal(id string) (isGlobal bool) {
	_, isGlobal = sb.GlobalSet[id]
	return
}

func (sb *symbolBuilder) addLocal(id string, compilerLevel int, scopeLevel int, register byte) {
	sb.LocalMap[lKey{id: id, compilerLevel: compilerLevel, scopeLevel: scopeLevel}] = register
}

func (sb *symbolBuilder) addGlobal(id string) {
	sb.GlobalSet[id] = dummy
}

func (sb *symbolBuilder) clearLocals(compilerLevel int, scopeLevel int) int {
	var count int
	length := len(sb.History)
	for i := length - 1; i >= 0; i-- {
		if sb.History[i].compilerLevel == compilerLevel && sb.History[i].scopeLevel == scopeLevel {
			count++
			delete(sb.LocalMap, sb.History[i])
			continue
		}
		break
	}
	sb.History = sb.History[:length-count]
	return count
}
