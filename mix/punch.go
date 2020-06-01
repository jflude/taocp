package mix

import (
	"io"
	"os"
	"strings"
)

type CardPunch struct {
	c  *Computer
	wc io.WriteCloser
}

func NewCardPunch(c *Computer, wc io.WriteCloser) (*CardPunch, error) {
	if wc == nil {
		var err error
		wc, err = os.OpenFile("punch.mix",
			os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			return nil, err
		}
	}
	return &CardPunch{c, wc}, nil
}

func (*CardPunch) Name() string {
	return "PUNCH"
}

func (*CardPunch) BlockSize() int {
	return 16
}

func (*CardPunch) Read([]Word) error {
	return ErrInvalidOperation
}

func (p *CardPunch) Write(block []Word) error {
	s := ConvertToUTF8(block)
	if !isPunchable(s) {
		return ErrInvalidCharacter
	}
	_, err := p.wc.Write([]byte(s))
	return err
}

func (p *CardPunch) Control(m int) error {
	return ErrInvalidControl
}

func (p *CardPunch) BusyUntil() int64 {
	return 0
}

func (p *CardPunch) Close() error {
	return p.wc.Close()
}

func isPunchable(s string) bool {
	return !strings.ContainsAny(s, "ΘΦΠ")
}
