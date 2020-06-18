package mixal

import (
	"fmt"
	"io"
	//"log"

	"github.com/jflude/gnuth/mix"
)

const wordsPerCard = 7

type object struct {
	start int
	orig  []int
	seg   [][]mix.Word
}

var loader = " O O6 Y O6    I   B= D O4 Z IQ Z I3 Z EN    E   EU 0BB= H IU   EJ  CA. ACB=   EU\n 1A-H V A=  CEU 0AEH 1AEN    E  CLU  ABG H IH A A= J B. A  9                    \n"
var transfer = "TRANS0%04d                                                                      \n"

func (o *object) findWord(address int) *mix.Word {
	for i, orig := range o.orig {
		if address >= orig && address < orig+len(o.seg[i]) {
			return &o.seg[i][address-orig]
		}
	}
	return nil
}

func (o *object) writeCards(w io.Writer) error {
	if _, err := io.WriteString(w, loader); err != nil {
		return err
	}
	for i := range o.orig {
		if err := o.writeSeg(w, i); err != nil {
			return err
		}
	}
	_, err := io.WriteString(w, fmt.Sprintf(transfer, o.start))
	return err
}

func (o *object) writeSeg(w io.Writer, n int) error {
	orig := o.orig[n]
	for i := 0; i < len(o.seg[n]); i += wordsPerCard {
		c := len(o.seg[n]) - i
		if c > wordsPerCard {
			c = wordsPerCard
		}
		s := fmt.Sprintf("ABCDE%1d%04d", c, orig)
		if _, err := io.WriteString(w, s); err != nil {
			return err
		}
		for j := 0; j < c; j++ {
			v := o.seg[n][i+j]
			//log.Printf("NewWord(%#o), // %5d:", v.Int(), orig+j)
			s = fmt.Sprintf("%010d", abs(v.Int()))
			if v.Sign() == -1 {
				d := mix.OverPunch(rune(s[len(s)-1]))
				s = s[:len(s)-1] + string(d)
			}
			if _, err := io.WriteString(w, s); err != nil {
				return err
			}
		}
		orig += c
		for j := c; j < wordsPerCard; j++ {
			_, err := io.WriteString(w, "          ")
			if err != nil {
				return err
			}
		}
		if _, err := io.WriteString(w, "\n"); err != nil {
			return err
		}
	}
	return nil
}
