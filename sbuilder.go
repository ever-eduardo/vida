package vida

var dummy struct{}

type lKey struct {
	id    string
	level int
	scope int
	reg   byte
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

func (sb *symbolBuilder) isLocal(id string) (byte, bool, lKey) {
	var k lKey
	for i := len(sb.History) - 1; i >= 0; i-- {
		if sb.History[i].id == id {
			k = sb.History[i]
			return sb.History[i].reg, true, k
		}
	}
	return 0, false, k
}

func (sb *symbolBuilder) isGlobal(id string) (isGlobal bool) {
	_, isGlobal = sb.GlobalSet[id]
	return
}

func (sb *symbolBuilder) addLocal(id string, level int, scope int, reg byte) {
	sb.History = append(sb.History,
		lKey{
			id:    id,
			level: level,
			scope: scope,
			reg:   reg,
		})
}

func (sb *symbolBuilder) addGlobal(id string) {
	sb.GlobalSet[id] = dummy
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
