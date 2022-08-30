// Copyright (c) 2020 Justin Flude. All rights reserved.
// Use of this source code is governed by the COPYING.md file.
package mix

import (
	"io"
	"strings"
)

type Printer struct {
	wc io.WriteCloser
}

// see https://en.wikipedia.org/wiki/IBM_1403
func NewPrinter(wc io.WriteCloser) (*Printer, error) {
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
	return 100000, trimWrite(p.wc, block)
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

func trimWrite(w io.Writer, block []Word) error {
	line := strings.TrimRight(ConvertToUTF8(block), " ") + "\n"
	_, err := io.WriteString(w, line)
	return err
}
