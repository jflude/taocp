package mix

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

type Drum struct {
	f    *os.File
	name string
	here int64
	c    *Computer
}

func NewDrum(file string, unit int, c *Computer) (*Drum, error) {
	f, err := os.OpenFile(file, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}
	return &Drum{f, fmt.Sprintf("DRUM%02d", unit), 0, c}, nil
}

func (d *Drum) Name() string {
	return d.name
}

func (*Drum) BlockSize() int {
	return 100
}

func (d *Drum) Read(block []Word) (int64, error) {
	dur, err := d.seekToX()
	if err != nil {
		return 0, err
	}
	buf := make([]byte, 4*len(block))
	if _, err := io.ReadFull(d.f, buf); err != nil {
		return 0, err
	}
	for i, j := 0, 0; i < len(block); i, j = i+1, j+4 {
		block[i] = Word(binary.LittleEndian.Uint32(buf[j : j+4]))
	}
	return 3000 + dur, nil
}

func (d *Drum) Write(block []Word) (int64, error) {
	dur, err := d.seekToX()
	if err != nil {
		return 0, err
	}
	buf := make([]byte, 4*len(block))
	for i, j := 0, 0; i < len(block); i, j = i+1, j+4 {
		binary.LittleEndian.PutUint32(buf[j:j+4], uint32(block[i]))
	}
	_, err = d.f.Write(buf)
	return 3000 + dur, err
}

func (d *Drum) Control(m int) (int64, error) {
	if m != 0 {
		return 0, ErrInvalidCommand
	}
	return d.seekToX()
}

func (d *Drum) Close() error {
	return d.f.Close()
}

func (d *Drum) seekToX() (dur int64, err error) {
	x := int64(d.c.Reg[X].Field(045).Int() * 4 * d.BlockSize())
	if d.here != x {
		d.here, err = d.f.Seek(x, io.SeekStart)
		dur = 20000
	}
	return
}
