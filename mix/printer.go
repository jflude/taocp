package mix

import (
	"io"
	"os"
	"strings"
)

type Printer struct {
	c  *Computer
	wc io.WriteCloser
}

func NewPrinter(c *Computer, wc io.WriteCloser) (*Printer, error) {
	if wc == nil {
		var err error
		wc, err = os.OpenFile("printer.mix",
			os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			return nil, err
		}
	}
	return &Printer{c, wc}, nil
}

func (*Printer) Name() string {
	return "PRINTER"
}

func (*Printer) BlockSize() int {
	return 24
}

func (*Printer) Read([]Word) error {
	return ErrInvalidOperation
}

func (p *Printer) Write(block []Word) error {
	line := strings.TrimRight(ConvertToUTF8(block), " ")
	_, err := p.wc.Write([]byte(line + "\n"))
	return err
}

func (p *Printer) Control(m int) error {
	if m != 0 {
		return ErrInvalidControl
	}
	_, err := p.wc.Write([]byte("\014")) // form feed
	return err
}

func (p *Printer) BusyUntil() int64 {
	return 0
}

func (p *Printer) Close() error {
	return p.wc.Close()
}
