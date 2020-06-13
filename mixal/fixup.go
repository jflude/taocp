package mixal

import "github.com/jflude/gnuth/mix"

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
		panic(ErrInternalError)
	}
	for _, r := range refs {
		w := a.obj.findWord(r)
		if w == nil {
			panic(ErrInternalError)
		}
		w.SetField(mix.FieldSpec(0, 2), val)
	}
	delete(a.fixups, sym)
}
