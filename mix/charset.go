package mix

import (
	"errors"
	"fmt"
)

var ErrInvalidChar = errors.New("mix: invalid character")

var mix2utf = []rune(` ABCDEFGHIΔJKLMNOPQRΣΠSTUVWXYZ0123456789.,()+-*/=$<>@;:'`)
var utf2mix = make(map[rune]int)

func init() {
	for i, v := range mix2utf {
		utf2mix[v] = i
	}
}

func ConvertToUTF8(w []Word) string {
	var r []rune
	for _, v := range w {
		for f := 1; f <= 5; f++ {
			r = append(r, mix2utf[v.Field(FieldSpec(f, f))])
		}
	}
	return string(r)
}

func ConvertToMIX(s string) ([]Word, error) {
	var w []Word
	f := 0
	for _, r := range s {
		if f = f%5 + 1; f == 1 {
			w = append(w, NewWord(0))
		}
		c, ok := utf2mix[r]
		if !ok {
			return nil, charError(r)
		}
		w[len(w)-1].SetField(FieldSpec(f, f), NewWord(c))
	}
	return w, nil
}

func charError(r rune) error {
	return fmt.Errorf("%w: %#U", ErrInvalidChar, r)
}

func IsChar(r rune) bool {
	_, ok := utf2mix[r]
	return ok
}

func IsDigit(r rune) bool {
	return r >= '0' && r <= '9'
}

func IsLetter(r rune) bool {
	return (r >= 'A' && r <= 'Z') || r == 'Δ' || r == 'Σ' || r == 'Π'
}

func IsSpace(r rune) bool {
	return r == ' '
}
