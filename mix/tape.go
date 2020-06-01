package mix

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"strings"
)

type Tape struct {
	c    *Computer
	f    *os.File
	name string
}

func NewTape(c *Computer, f *os.File, unit int) (*Tape, error) {
	n := fmt.Sprintf("TAPE%02d", unit)
	if f == nil {
		file := strings.ToLower(n) + ".mix"
		var err error
		f, err = os.OpenFile(file, os.O_CREATE|os.O_RDWR, 0644)
		if err != nil {
			return nil, err
		}
	}
	return &Tape{c, f, n}, nil
}

func (t *Tape) Name() string {
	return t.name
}

func (*Tape) BlockSize() int {
	return 100
}

func (t *Tape) Read(block []Word) error {
	buf := make([]byte, 4*len(block))
	if _, err := io.ReadFull(t.f, buf); err != nil {
		return err
	}
	for i, j := 0, 0; i < len(block); i, j = i+1, j+4 {
		block[i] = Word(binary.LittleEndian.Uint32(buf[j : j+4]))
	}
	return nil
}

func (t *Tape) Write(block []Word) error {
	buf := make([]byte, 4*len(block))
	for i, j := 0, 0; i < len(block); i, j = i+1, j+4 {
		binary.LittleEndian.PutUint32(buf[j:j+4], uint32(block[i]))
	}
	_, err := t.f.Write(buf)
	return err
}

func (t *Tape) Control(m int) error {
	var p, wh int
	switch {
	case m < 0:
		p = -4 * t.BlockSize()
		wh = io.SeekCurrent
	case m == 0:
		p = 0
		wh = io.SeekStart
	case m > 0:
		p = 4 * t.BlockSize()
		wh = io.SeekCurrent
	}
	_, err := t.f.Seek(int64(p), wh)
	return err
}

func (t *Tape) BusyUntil() int64 {
	return 0
}

func (t *Tape) Close() error {
	return t.f.Close()
}
