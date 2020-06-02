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
	if !IsPunchable(s) {
		return 0, ErrInvalidCharacter
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

func IsPunchable(s string) bool {
	// see Ex. 26, Section 1.3.1 for characters which cannot be punched
	return !strings.ContainsAny(s, "ΦΠ$<>@;:'")
}

func OverPunch(digit rune) rune {
	return []rune("ΘJKLMNOPQR")[digit-'0']
}
