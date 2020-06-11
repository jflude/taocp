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

func extractColumns(line string, from, to int, trim bool) string {
	if len(line) < from {
		return ""
	}
	s := line[from-1:]
	if len(s) > to-from+1 {
		s = s[:to-from+1]
	}
	if trim {
		s = strings.TrimSpace(s)
	}
	return s
}
