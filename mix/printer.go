package mix

import (
	"os"
	"strings"
)

type Printer struct {
	file *os.File
}

func NewPrinter() *Printer {
	return new(Printer)
}

func (*Printer) Name() string {
	return "Printer"
}

func (*Printer) BlockSize() int {
	return 24
}

func (*Printer) Read([]Word) error {
	return ErrInvalidOperation
}

func (p *Printer) Write(w []Word) error {
	if err := p.checkOpen(); err != nil {
		return err
	}
	line := strings.TrimRight(ConvertToUTF8(w), " ")
	_, err := p.file.WriteString(line + "\n")
	return err
}

func (p *Printer) Control(op int) error {
	if op != 0 {
		return ErrInvalidControl
	}
	if err := p.checkOpen(); err != nil {
		return err
	}
	_, err := p.file.WriteString("\014") // form feed
	return err
}

func (p *Printer) BusyUntil(now int64) int64 {
	return 0
}

func (p *Printer) checkOpen() error {
	if p.file == nil {
		var err error
		if p.file, err = os.Create("printer.out"); err != nil {
			return err
		}
	}
	return nil
}
