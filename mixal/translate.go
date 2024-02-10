package mixal

import (
	"bufio"
	"io"
	"log"

	"github.com/jflude/taocp/mix"
)

type parseFunc func(*asmb, string, string, string)

func (a *asmb) translate(r io.Reader, parse parseFunc) error {
	a.symbols = make(map[string]mix.Word)
	a.fixups = make(map[string][]int)
	var err error
	sc := bufio.NewScanner(r)
	for sc.Scan() {
		if err2 := a.processCard(sc.Text(), parse); err2 != nil {
			if err != nil {
				log.Println("error:", err)
			}
			err = err2
		}
	}
	return err
}
