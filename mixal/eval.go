package mixal

import "github.com/jflude/gnuth/mix"

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
		x := a.exprVal.ShiftBytesRight(5)
		*a.exprVal, _, _ = mix.DivWord(*a.exprVal, x, arg)
	case '\\':
		*a.exprVal, _, _ = mix.DivWord(*a.exprVal, 0, arg)
	case ':':
		hi, lo := mix.MulWord(*a.exprVal, 8)
		hi.SetField(mix.FieldSpec(1, 5), lo)
		*a.exprVal, _ = mix.AddWord(hi, arg)
	default:
		panic(ErrInternal)
	}
}
