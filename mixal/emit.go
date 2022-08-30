// Copyright (c) 2020 Justin Flude. All rights reserved.
// Use of this source code is governed by the COPYING.md file.
package mixal

import "github.com/jflude/taocp/mix"

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

func abs(v int) int {
	if v < 0 {
		return -v
	}
	return v
}
