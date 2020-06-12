package mixal

import "strings"

type token struct {
	kind int
	val  interface{}
}

const (
	symbol = iota + 128
	number
	asterisk
	binary
	operator
)

func (a *asmb) addToken(kind int, val interface{}) {
	a.tokens = append(a.tokens, token{kind, val})
}

func (a *asmb) lastKind() int {
	return a.tokens[len(a.tokens)-1].kind
}

func (a *asmb) lastString() string {
	return a.tokens[len(a.tokens)-1].val.(string)
}

func (a *asmb) lastQuantity() int {
	switch t := a.tokens[len(a.tokens)-1]; t.kind {
	case symbol:
		if n, ok := a.symbols[t.val.(string)]; ok {
			return n
		}
		panic(ErrInternalError)
	case asterisk:
		fallthrough
	case number:
		return t.val.(int)
	default:
		panic(ErrInternalError)
	}
}

func (a *asmb) extractColumns(from, to int, trim bool) string {
	if len(a.input) < from {
		return ""
	}
	s := a.input[from-1:]
	if len(s) > to-from+1 {
		s = s[:to-from+1]
	}
	if trim {
		s = strings.TrimSpace(s)
	}
	return s
}
