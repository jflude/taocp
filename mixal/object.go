// Copyright (c) 2020 Justin Flude. All rights reserved.
// Use of this source code is governed by the COPYING.md file.
package mixal

import (
	"fmt"
	"io"

	//"log"

	"github.com/jflude/taocp/mix"
)

const wordsPerCard = 7

type object struct {
	start int
	orig  []int
	seg   [][]mix.Word
}

// see the answer to Ex. 26, Section 1.3.1
var regularLoader = " O O6 Z O6    I C O4 0 EH A  F F CF 0  E   EU 0 IH G BB   EJ  CA. Z EU   EH E BA\n" +
	"   EU 2A-H S BB  C U 1AEH 2AEN V  E  CLU  ABG Z EH E BB J B. A  9               \n"
var interruptLoader = " O O6 Z O6    I C O4 0 EH A  F F CF 0  E   EU 0 IH G BB   EJ  CA. Z EU   EH E BA\n" +
	"   EU 2A-H S BB  C U 1AEH 2AEN V  E  CLU  ABG Z EH E BB J B. A  9               \n" // TODO replace with a version that works
var transfer = "TRANS0%04d                                                                      \n"

func (o *object) findWord(address int) *mix.Word {
	for i, orig := range o.orig {
		if address >= orig && address < orig+len(o.seg[i]) {
			return &o.seg[i][address-orig]
		}
	}
	return nil
}

func (o *object) writeCards(w io.Writer, interrupts bool) error {
	var loader = regularLoader
	if interrupts {
		loader = interruptLoader
	}
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
		s := fmt.Sprintf("ABCDE%1d%04d", c, abs(orig))
		if orig < 0 {
			s = overPunchNegative(s)
		}
		if _, err := io.WriteString(w, s); err != nil {
			return err
		}
		for j := 0; j < c; j++ {
			v := o.seg[n][i+j]
			//log.Printf("NewWord(%#o), // %5d:", v.Int(), orig+j)
			s = fmt.Sprintf("%010d", abs(v.Int()))
			if v.Sign() == -1 {
				s = overPunchNegative(s)
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

func overPunchNegative(n string) string {
	return n[:len(n)-1] + string(mix.OverPunch(rune(n[len(n)-1])))
}
