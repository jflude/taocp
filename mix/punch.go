package mix

import (
	"io"
	"strings"
)

type CardPunch struct {
	wc io.WriteCloser
}

// see https://en.wikipedia.org/wiki/IBM_2540
func NewCardPunch(wc io.WriteCloser) (*CardPunch, error) {
	return &CardPunch{wc}, nil
}

func (*CardPunch) Name() string {
	return "PUNCH"
}

func (*CardPunch) BlockSize() int {
	return 16
}

func (*CardPunch) Read([]Word) (int64, error) {
	return 0, ErrInvalidCommand
}

func (p *CardPunch) Write(block []Word) (int64, error) {
	s := strings.TrimRight(EncodeAsUTF8(block), " ")
	if ch, ok := IsPunchable(s); !ok {
		return 0, charError(ch)
	}
	_, err := io.WriteString(p.wc, s+"\n")
	return 200000, err
}

func (p *CardPunch) Control(m int) (int64, error) {
	return 0, ErrInvalidCommand
}

func (p *CardPunch) Close() error {
	return p.wc.Close()
}

func IsPunchable(s string) (rune, bool) {
	// see Ex. 26, Section 1.3.1 for characters which cannot be punched
	// (see also https://homepage.divms.uiowa.edu/~jones/cards/codes.html)
	// I am assuming a colon can be punched as it is required by MIXAL.
	if i := strings.IndexAny(s, "ΦΠ$<>@;'"); i != -1 {
		return rune(s[i]), false
	}
	return 0, true
}

func OverPunch(digit rune) rune {
	// see Ex. 26, Section 1.3.1 for digits which have been overpunched
	return []rune("ΔJKLMNOPQR")[digit-'0']
}
