package mix

import (
	"errors"
	"fmt"
)

var ErrInvalidCharacter = errors.New("invalid character")

var mix2utf = []rune(` ABCDEFGHIΘJKLMNOPQRΦΠSTUVWXYZ0123456789.,()+-*/=$<>@;:'`)
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
	for _, v := range s {
		if f = f%5 + 1; f == 1 {
			w = append(w, NewWord(0))
		}
		c, ok := utf2mix[v]
		if !ok {
			return nil, fmt.Errorf("%w: %q",
				ErrInvalidCharacter, string(v))
		}
		w[len(w)-1].SetField(FieldSpec(f, f), NewWord(c))
	}
	return w, nil
}
