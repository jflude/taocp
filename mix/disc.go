package mix

import (
	"encoding/binary"
	"fmt"
	"io"
)

const maxDiscBlock = 40000

type Disc struct {
	rwsc readWriteSeekCloser
	name string
	here int64
	c    *Computer
}

// see Section 5.4.9
func NewDisc(rwsc readWriteSeekCloser, unit int, c *Computer) (*Disc, error) {
	return &Disc{rwsc, fmt.Sprintf("DISC%02d", unit), 0, c}, nil
}

func (d *Disc) Name() string {
	return d.name
}

func (*Disc) BlockSize() int {
	return 100
}

func (d *Disc) Read(block []Word) (int64, error) {
	duration, err := d.seekToX()
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
	return 15000 + duration, nil
}

func (d *Disc) Write(block []Word) (int64, error) {
	duration, err := d.seekToX()
	if err != nil {
		return 0, err
	}
	buf := make([]byte, 4*len(block))
	for i, j := 0, 0; i < len(block); i, j = i+1, j+4 {
		binary.LittleEndian.PutUint32(buf[j:j+4], uint32(block[i]))
	}
	_, err = d.rwsc.Write(buf)
	return 15000 + duration, err
}

func (d *Disc) Control(m int) (int64, error) {
	if m != 0 {
		return 0, ErrInvalidCommand
	}
	return d.seekToX()
}

func (d *Disc) Close() error {
	return d.rwsc.Close()
}

func (d *Disc) seekToX() (duration int64, err error) {
	x := abs64(int64(d.c.Reg[X].Int()))
	if x > maxDiscBlock {
		return 0, ErrInvalidBlock
	}
	x *= 4 * int64(d.BlockSize())
	if d.here != x {
		d.here, err = d.rwsc.Seek(x, io.SeekStart)
		duration = 60000
	}
	return
}
