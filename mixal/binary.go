// Copyright (c) 2020 Justin Flude. All rights reserved.
// Use of this source code is governed by the COPYING.md file.
package mixal

func (a *asmb) matchBinaryOp() bool {
	if len(a.input) > 1 {
		if a.input[:2] == "//" {
			a.addToken('\\', nil)
			a.input = a.input[2:]
			return true
		}
	}
	if len(a.input) > 0 {
		switch a.input[0] {
		case '+', '-', '*', '/', ':':
			return a.matchChar(a.input[0])
		}
	}
	return false
}

func (a *asmb) matchChar(ch byte) bool {
	if len(a.input) == 0 || a.input[0] != ch {
		return false
	}
	a.addToken(int(ch), nil)
	a.input = a.input[1:]
	return true
}
