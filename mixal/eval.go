// Copyright (c) 2020 Justin Flude. All rights reserved.
// Use of this source code is governed by the COPYING.md file.
package mixal

import "github.com/jflude/taocp/mix"

func (a *asmb) evalArg(arg int) {
	if a.exprVal == nil {
		w := mix.NewWord(arg)
		a.exprVal = &w
		return
	}
	switch a.exprOp {
	case '+':
		*a.exprVal, _ = mix.AddWord(*a.exprVal, arg)
	case '-':
		*a.exprVal, _ = mix.SubWord(*a.exprVal, arg)
	case '*':
		_, *a.exprVal = mix.MulWord(*a.exprVal, arg)
	case '/':
		var x mix.Word
		mix.ShiftBitsRight(a.exprVal, &x, 30)
		*a.exprVal, _, _ = mix.DivWord(*a.exprVal, x, arg)
	case '\\':
		*a.exprVal, _, _ = mix.DivWord(*a.exprVal, 0, arg)
	case ':':
		high, low := mix.MulWord(*a.exprVal, 8)
		high.SetField(mix.Spec(1, 5), low)
		*a.exprVal, _ = mix.AddWord(high, arg)
	default:
		panic(ErrInternal)
	}
}
