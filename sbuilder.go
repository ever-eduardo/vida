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
	LocalMap  map[string]int
	Count     int
}

func newSymbolBuilder() *symbolBuilder {
	return &symbolBuilder{
		GlobalSet: make(map[string]struct{}),
		LocalMap:  make(map[string]int),
		Count:     -1,
	}
}

func (sb *symbolBuilder) isLocal(id string) (byte, bool) {
	if i, isPresent := sb.LocalMap[id]; isPresent {
		return sb.History[i].register, true
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
	sb.Count++
	sb.LocalMap[id] = sb.Count
}

func (sb *symbolBuilder) addGlobal(id string) {
	sb.GlobalSet[id] = dummy
}

func (sb *symbolBuilder) clearLocals(compilerLevel int, scopeLevel int) int {
	var count int
	length := len(sb.History)
	for i := length - 1; i >= 0; i-- {
		lkey := sb.History[i]
		if lkey.compilerLevel == compilerLevel && lkey.scopeLevel == scopeLevel {
			count++
			if n := sb.LocalMap[lkey.id]; n-1 < 0 {
				delete(sb.LocalMap, lkey.id)
			} else {
				sb.LocalMap[lkey.id]--
			}
			continue
		}
		break
	}
	if count > 0 {
		sb.History = sb.History[:length-count]
		sb.Count -= count
	}
	return count
}
