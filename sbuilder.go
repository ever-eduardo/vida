package vida

var dummy struct{}

type lKey struct {
	id            string
	compilerLevel int
	scopeLevel    int
	register      byte
}

type symbolBuilder struct {
	History   []lKey
	GlobalSet map[string]struct{}
}

func newSymbolBuilder() *symbolBuilder {
	return &symbolBuilder{
		GlobalSet: make(map[string]struct{}),
	}
}

func (sb *symbolBuilder) isLocal(id string) (byte, bool) {
	for i := len(sb.History) - 1; i >= 0; i-- {
		if sb.History[i].id == id {
			return sb.History[i].register, true
		}
	}
	return 0, false
}

func (sb *symbolBuilder) isGlobal(id string) (isGlobal bool) {
	_, isGlobal = sb.GlobalSet[id]
	return
}

func (sb *symbolBuilder) addLocal(id string, compilerLevel int, scopeLevel int, register byte) {
	sb.History = append(sb.History,
		lKey{
			id:            id,
			compilerLevel: compilerLevel,
			scopeLevel:    scopeLevel,
			register:      register,
		})
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
		} else {
			break
		}
	}
	if count > 0 {
		sb.History = sb.History[:length-count]
	}
	return count
}
