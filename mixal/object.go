package mixal

import (
	"fmt"
	"io"

	"github.com/jflude/gnuth/mix"
)

type chunk struct {
	orig int
	data []mix.Word
}

type object []*chunk

var loader = ` O O6 Y O6    I   B= D O4 Z IQ Z I3 Z EN    E   EU 0BB= H IU   EJ  CA. ACB=   EU 1A-H V A=  CEU 0AEH 1AEN    E  CLU  ABG H IH A A= J B. A  9                    `
var transfer = `TRANS0%04d                                                                      `

func (obj object) outputCards(w io.Writer, start int) error {
	if _, err := io.WriteString(w, loader); err != nil {
		return err
	}
	for _, chunk := range obj {
		orig := chunk.orig
		for i := 0; i < len(chunk.data); i += 7 {
			n := len(chunk.data) - i
			if n > 7 {
				n = 7
			}
			s := fmt.Sprintf("ABCDE%1d%04d", n, orig)
			if _, err := io.WriteString(w, s); err != nil {
				return err
			}
			for _, v := range chunk.data[i : i+n] {
				s = fmt.Sprintf("%010d", v.Int())
				if v.Sign() == -1 {
					d := mix.OverPunch(rune(s[len(s)-1]))
					s = s[:len(s)-1] + string(d)
				}
				if _, err := io.WriteString(w, s); err != nil {
					return err
				}
			}
			orig += n
			for i := n; i < 7; i++ {
				_, err := io.WriteString(w, "          ")
				if err != nil {
					return err
				}
			}
		}
	}
	_, err := io.WriteString(w, fmt.Sprintf(transfer, start))
	return err
}
