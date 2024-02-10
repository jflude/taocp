package mixal

import "github.com/jflude/taocp/mix"

func (a *asmb) addFixUp(sym string) {
	a.fixups[sym] = append(a.fixups[sym], a.self)
}

func (a *asmb) patchFixUps(sym string) {
	refs := a.fixups[sym]
	if refs == nil {
		return
	}
	val, ok := a.symbols[sym]
	if !ok {
		panic(ErrInternal)
	}
	for _, r := range refs {
		w := a.obj.findWord(r)
		if w == nil {
			panic(ErrInternal)
		}
		w.SetField(mix.Spec(0, 2), val)
	}
	delete(a.fixups, sym)
}
