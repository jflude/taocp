package mixal

import (
	"fmt"
	"io"

	"github.com/jflude/gnuth/mix"
)

const wordsPerCard = 7

type object struct {
	start int
	orig  []int
	seg   [][]mix.Word
}

var loader = ` O O6 Y O6    I   B= D O4 Z IQ Z I3 Z EN    E   EU 0BB= H IU   EJ  CA. ACB=   EU 1A-H V A=  CEU 0AEH 1AEN    E  CLU  ABG H IH A A= J B. A  9                    `
var transfer = `TRANS0%04d                                                                      `

func (obj *object) writeCards(w io.Writer) error {
	if _, err := io.WriteString(w, loader); err != nil {
		return err
	}
	for i := 0; i < len(obj.orig); i++ {
		orig := obj.orig[i]
		for j := 0; j < len(obj.seg[i]); j += wordsPerCard {
			n := len(obj.seg[i]) - j
			if n > wordsPerCard {
				n = wordsPerCard
			}
			s := fmt.Sprintf("ABCDE%1d%04d", n, orig)
			if _, err := io.WriteString(w, s); err != nil {
				return err
			}
			for _, v := range obj.seg[i][j : j+n] {
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
			for j := n; j < wordsPerCard; j++ {
				_, err := io.WriteString(w, "          ")
				if err != nil {
					return err
				}
			}
		}
	}
	_, err := io.WriteString(w, fmt.Sprintf(transfer, obj.start))
	return err
}
