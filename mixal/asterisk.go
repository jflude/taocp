package mixal

func (a *asmb) matchAsterisk() bool {
	if len(a.input) > 0 {
		if a.input[0] == '*' {
			a.addToken(asterisk, a.self)
			a.input = a.input[1:]
			return true
		}
	}
	return false
}
