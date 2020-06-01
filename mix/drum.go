package mix

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"strings"
)

type Drum struct {
	c    *Computer
	f    *os.File
	name string
	pos  int64
}

func NewDrum(c *Computer, f *os.File, unit int) (*Drum, error) {
	n := fmt.Sprintf("DRUM%02d", unit)
	if f == nil {
		file := strings.ToLower(n) + ".mix"
		var err error
		f, err = os.OpenFile(file, os.O_CREATE|os.O_RDWR, 0644)
		if err != nil {
			return nil, err
		}
	}
	return &Drum{c, f, n, 0}, nil
}

func (d *Drum) Name() string {
	return d.name
}

func (*Drum) BlockSize() int {
	return 100
}

func (d *Drum) Read(block []Word) error {
	if err := d.seekToX(); err != nil {
		return err
	}
	buf := make([]byte, 4*len(block))
	if _, err := io.ReadFull(d.f, buf); err != nil {
		return err
	}
	for i, j := 0, 0; i < len(block); i, j = i+1, j+4 {
		block[i] = Word(binary.LittleEndian.Uint32(buf[j : j+4]))
	}
	return nil
}

func (d *Drum) Write(block []Word) error {
	if err := d.seekToX(); err != nil {
		return err
	}
	buf := make([]byte, 4*len(block))
	for i, j := 0, 0; i < len(block); i, j = i+1, j+4 {
		binary.LittleEndian.PutUint32(buf[j:j+4], uint32(block[i]))
	}
	_, err := d.f.Write(buf)
	return err
}

func (d *Drum) Control(m int) error {
	if m != 0 {
		return ErrInvalidControl
	}
	return d.seekToX()
}

func (d *Drum) BusyUntil() int64 {
	return 0
}

func (d *Drum) Close() error {
	return d.f.Close()
}

func (d *Drum) seekToX() (err error) {
	x := int64(d.c.Reg[X].Field(045).Int() * 4 * d.BlockSize())
	if d.pos != x {
		d.pos, err = d.f.Seek(x, io.SeekStart)
	}
	return
}
