package vida

type lKey struct {
	id    string
	level int
	scope int
	reg   int
}

type symbolBuilder struct {
	History   []lKey
	GlobalSet map[string]int
	index     int
}

func newSymbolBuilder() *symbolBuilder {
	sb := &symbolBuilder{
		GlobalSet: make(map[string]int),
	}
	for _, v := range coreLibNames {
		sb.addGlobal(v)
	}
	return sb
}

func (sb *symbolBuilder) isLocal(id string) (int, bool, lKey) {
	var k lKey
	for i := len(sb.History) - 1; i >= 0; i-- {
		if sb.History[i].id == id {
			k = sb.History[i]
			return sb.History[i].reg, true, k
		}
	}
	return 0, false, k
}

func (sb *symbolBuilder) isGlobal(id string) (idx int, isGlobal bool) {
	idx, isGlobal = sb.GlobalSet[id]
	return
}

func (sb *symbolBuilder) addLocal(id string, level int, scope int, reg int) {
	sb.History = append(sb.History,
		lKey{
			id:    id,
			level: level,
			scope: scope,
			reg:   reg,
		})
}

func (sb *symbolBuilder) addGlobal(id string) (int, bool) {
	idx, isPresent := sb.GlobalSet[id]
	if isPresent {
		return idx, isPresent
	}
	i := sb.index
	sb.GlobalSet[id] = i
	sb.index++
	return i, isPresent
}

func (sb *symbolBuilder) clearLocals(compilerLevel int, scopeLevel int) int {
	var count int
	length := len(sb.History)
	for i := length - 1; i >= 0; i-- {
		if sb.History[i].level == compilerLevel && sb.History[i].scope == scopeLevel {
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
