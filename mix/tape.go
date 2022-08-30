// Copyright (c) 2020 Justin Flude. All rights reserved.
// Use of this source code is governed by the COPYING.md file.
package mix

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

const maxTapeBlock = 46000

type Tape struct {
	rwsc readWriteSeekCloser
	name string
	here int64
}

var ErrInvalidBlock = errors.New("mix: invalid block")

// see Section 5.4.6
func NewTape(rwsc readWriteSeekCloser, unit int) (*Tape, error) {
	return &Tape{rwsc, fmt.Sprintf("TAPE%02d", unit), 0}, nil
}

func (t *Tape) Name() string {
	return t.name
}

func (*Tape) BlockSize() int {
	return 100
}

func (t *Tape) Read(block []Word) (int64, error) {
	if t.isPastEnd(4 * int64(t.BlockSize())) {
		return 0, ErrInvalidBlock
	}
	buf := make([]byte, 4*len(block))
	if _, err := io.ReadFull(t.rwsc, buf); err != nil {
		return 0, err
	}
	for i, j := 0, 0; i < len(block); i, j = i+1, j+4 {
		block[i] = Word(binary.LittleEndian.Uint32(buf[j : j+4]))
	}
	t.here += int64(4 * t.BlockSize())
	return 8000, nil
}

func (t *Tape) Write(block []Word) (int64, error) {
	if t.isPastEnd(4 * int64(t.BlockSize())) {
		return 0, ErrInvalidBlock
	}
	buf := make([]byte, 4*len(block))
	for i, j := 0, 0; i < len(block); i, j = i+1, j+4 {
		binary.LittleEndian.PutUint32(buf[j:j+4], uint32(block[i]))
	}
	_, err := t.rwsc.Write(buf)
	if err == nil {
		t.here += int64(4 * t.BlockSize())
	}
	return 8000, err
}

func (t *Tape) Control(m int) (int64, error) {
	var off, delay int64
	var wh int
	if m == 0 {
		off = 0
		wh = io.SeekStart
		delay = 5000 * t.here / int64(4*t.BlockSize())
	} else {
		off = int64(4 * t.BlockSize() * m)
		if t.isPastEnd(off) {
			return 0, ErrInvalidBlock
		}
		if t.here+off < 0 {
			off = -t.here
		}
		wh = io.SeekCurrent
		delay = 5000 * abs64(int64(m))
	}
	var err error
	t.here, err = t.rwsc.Seek(off, wh)
	return delay, err
}

func (t *Tape) Close() error {
	return t.rwsc.Close()
}

func (t *Tape) isPastEnd(offset int64) bool {
	return t.here+offset > 4*int64(t.BlockSize())*maxTapeBlock
}
