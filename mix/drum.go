// Copyright (c) 2020 Justin Flude. All rights reserved.
// Use of this source code is governed by the COPYING.md file.
package mix

import (
	"encoding/binary"
	"fmt"
	"io"
)

const maxDrumBlock = 4000

type Drum struct {
	rwsc readWriteSeekCloser
	name string
	here int64
	c    *Computer
}

// see Section 5.4.9
func NewDrum(rwsc readWriteSeekCloser, unit int, c *Computer) (*Drum, error) {
	return &Drum{rwsc, fmt.Sprintf("DRUM%02d", unit), 0, c}, nil
}

func (d *Drum) Name() string {
	return d.name
}

func (*Drum) BlockSize() int {
	return 100
}

func (d *Drum) Read(block []Word) (int64, error) {
	delay, err := d.seekToX()
	if err != nil {
		return 0, err
	}
	buf := make([]byte, 4*len(block))
	if _, err := io.ReadFull(d.rwsc, buf); err != nil {
		return 0, err
	}
	for i, j := 0, 0; i < len(block); i, j = i+1, j+4 {
		block[i] = Word(binary.LittleEndian.Uint32(buf[j : j+4]))
	}
	return 10000 + delay, nil
}

func (d *Drum) Write(block []Word) (int64, error) {
	delay, err := d.seekToX()
	if err != nil {
		return 0, err
	}
	buf := make([]byte, 4*len(block))
	for i, j := 0, 0; i < len(block); i, j = i+1, j+4 {
		binary.LittleEndian.PutUint32(buf[j:j+4], uint32(block[i]))
	}
	_, err = d.rwsc.Write(buf)
	return 10000 + delay, err
}

func (d *Drum) Control(m int) (int64, error) {
	if m != 0 {
		return 0, ErrInvalidCommand
	}
	return d.seekToX()
}

func (d *Drum) Close() error {
	return d.rwsc.Close()
}

func (d *Drum) seekToX() (delay int64, err error) {
	x := abs64(int64(d.c.Reg[X].Int()))
	if x > maxDrumBlock {
		return 0, ErrInvalidBlock
	}
	x *= 4 * int64(d.BlockSize())
	if d.here != x {
		d.here, err = d.rwsc.Seek(x, io.SeekStart)
	}
	return 0, nil
}
