// Copyright (c) 2020 Justin Flude. All rights reserved.
// Use of this source code is governed by the COPYING.md file.
package mix

import "io"

type PaperTape struct {
	rwsc readWriteSeekCloser
}

// see https://en.wikipedia.org/wiki/Teletype_Model_33
func NewPaperTape(rwsc readWriteSeekCloser) (*PaperTape, error) {
	return &PaperTape{rwsc}, nil
}

func (*PaperTape) Name() string {
	return "PAPER"
}

func (*PaperTape) BlockSize() int {
	return 14
}

func (p *PaperTape) Read(block []Word) (int64, error) {
	buf := make([]byte, 5*p.BlockSize())
	if _, err := io.ReadFull(p.rwsc, buf); err != nil {
		return 0, err
	}
	m, err := ConvertToMIX(string(buf))
	if err != nil {
		return 0, err
	}
	copy(block, m)
	return 70000, nil
}

func (p *PaperTape) Write(block []Word) (int64, error) {
	_, err := io.WriteString(p.rwsc, ConvertToUTF8(block))
	return 200000, err
}

func (p *PaperTape) Control(m int) (int64, error) {
	if m != 0 {
		return 0, ErrInvalidCommand
	}
	_, err := p.rwsc.Seek(0, io.SeekStart)
	return 60000000, err
}

func (p *PaperTape) Close() error {
	return p.rwsc.Close()
}
