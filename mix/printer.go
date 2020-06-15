package mix

import (
	"io"
	"os"
	"strings"
)

type Printer struct {
	wc io.WriteCloser
}

func NewPrinter(file string) (*Printer, error) {
	wc, err := os.OpenFile(file, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}
	return &Printer{wc}, nil
}

func (*Printer) Name() string {
	return "PRINTER"
}

func (*Printer) BlockSize() int {
	return 24
}

func (*Printer) Read([]Word) (int64, error) {
	return 0, ErrInvalidCommand
}

func (p *Printer) Write(block []Word) (int64, error) {
	line := strings.TrimRight(ConvertToUTF8(block), " ")
	_, err := p.wc.Write([]byte(line + "\n"))
	return 400000, err
}

func (p *Printer) Control(m int) (int64, error) {
	if m != 0 {
		return 0, ErrInvalidCommand
	}
	_, err := p.wc.Write([]byte("\014")) // form feed
	return 1000000, err
}

func (p *Printer) Close() error {
	return p.wc.Close()
}
