package mix

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

const maxDiscBlock = 50000

type Disc struct {
	f    *os.File
	name string
	here int64
	c    *Computer
}

// see https://www.ibm.com/ibm/history/exhibits/storage/storage_2314.html
func NewDisc(file string, unit int, c *Computer) (*Disc, error) {
	f, err := os.OpenFile(file, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}
	return &Disc{f, fmt.Sprintf("DISC%02d", unit), 0, c}, nil
}

func (d *Disc) Name() string {
	return d.name
}

func (*Disc) BlockSize() int {
	return 100
}

func (d *Disc) Read(block []Word) (int64, error) {
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
	return 3000 + dur, nil // TODO: check timing
}

func (d *Disc) Write(block []Word) (int64, error) {
	dur, err := d.seekToX()
	if err != nil {
		return 0, err
	}
	buf := make([]byte, 4*len(block))
	for i, j := 0, 0; i < len(block); i, j = i+1, j+4 {
		binary.LittleEndian.PutUint32(buf[j:j+4], uint32(block[i]))
	}
	_, err = d.f.Write(buf)
	return 3000 + dur, err // TODO: check timing
}

func (d *Disc) Control(m int) (int64, error) {
	if m != 0 {
		return 0, ErrInvalidCommand
	}
	return d.seekToX()
}

func (d *Disc) Close() error {
	return d.f.Close()
}

func (d *Disc) seekToX() (dur int64, err error) {
	x := abs64(int64(d.c.Reg[X].Int()))
	if x > maxDiscBlock {
		return 0, ErrInvalidBlock
	}
	x *= 4 * int64(d.BlockSize())
	if d.here != x {
		d.here, err = d.f.Seek(x, io.SeekStart)
		dur = 20000 // TODO: check timing
	}
	return
}
