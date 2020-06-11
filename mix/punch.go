package mix

import (
	"io"
	"os"
	"strings"
)

type CardPunch struct {
	wc io.WriteCloser
}

func NewCardPunch(file string) (*CardPunch, error) {
	wc, err := os.OpenFile(file, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}
	return &CardPunch{wc}, nil
}

func (*CardPunch) Name() string {
	return "PUNCH"
}

func (*CardPunch) BlockSize() int {
	return 16
}

func (*CardPunch) Read([]Word) (int64, error) {
	return 0, ErrInvalidOperation
}

func (p *CardPunch) Write(block []Word) (int64, error) {
	s := ConvertToUTF8(block)
	if r, ok := IsPunchable(s); !ok {
		return 0, charError(r)
	}
	_, err := p.wc.Write([]byte(s))
	return 600000, err
}

func (p *CardPunch) Control(m int) (int64, error) {
	return 0, ErrInvalidControl
}

func (p *CardPunch) Close() error {
	return p.wc.Close()
}

func IsPunchable(s string) (rune, bool) {
	// see Ex. 26, Section 1.3.1 for characters which cannot be punched
	if i := strings.IndexAny(s, "ΦΠ$<>@;:'"); i != -1 {
		return rune(s[i]), false
	}
	return 0, true
}

func OverPunch(digit rune) rune {
	return []rune("ΘJKLMNOPQR")[digit-'0']
}
