package mixal

import "github.com/jflude/gnuth/mix"

func (a *asmb) newSegment(orig int) {
	a.obj.orig = append(a.obj.orig, orig)
	a.obj.seg = append(a.obj.seg, nil)
}

func (a *asmb) emit(w mix.Word) {
	if a.obj.orig == nil {
		a.newSegment(0)
	}
	i := len(a.obj.seg) - 1
	a.obj.seg[i] = append(a.obj.seg[i], w)
	if a.self++; abs(a.self) >= mix.MemorySize {
		panic(mix.ErrInvalidAddress)
	}
}
