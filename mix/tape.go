package mix

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"strings"
)

type Tape struct {
	f    *os.File
	name string
}

func NewTape(f *os.File, unit int) (*Tape, error) {
	n := fmt.Sprintf("TAPE%02d", unit)
	if f == nil {
		file := strings.ToLower(n) + ".mix"
		var err error
		f, err = os.OpenFile(file, os.O_CREATE|os.O_RDWR, 0644)
		if err != nil {
			return nil, err
		}
	}
	return &Tape{f, n}, nil
}

func (t *Tape) Name() string {
	return t.name
}

func (*Tape) BlockSize() int {
	return 100
}

func (t *Tape) Read(block []Word) (int64, error) {
	buf := make([]byte, 4*len(block))
	if _, err := io.ReadFull(t.f, buf); err != nil {
		return 0, err
	}
	for i, j := 0, 0; i < len(block); i, j = i+1, j+4 {
		block[i] = Word(binary.LittleEndian.Uint32(buf[j : j+4]))
	}
	return 6000, nil
}

func (t *Tape) Write(block []Word) (int64, error) {
	buf := make([]byte, 4*len(block))
	for i, j := 0, 0; i < len(block); i, j = i+1, j+4 {
		binary.LittleEndian.PutUint32(buf[j:j+4], uint32(block[i]))
	}
	_, err := t.f.Write(buf)
	return 6000, err
}

func (t *Tape) Control(m int) (int64, error) {
	var pos, wh int
	var dur int64
	if m == 0 {
		pos = 0
		wh = io.SeekStart
		dur = 60000000

	} else {
		pos = 4 * t.BlockSize() * m
		wh = io.SeekCurrent
		dur = 30000
	}
	_, err := t.f.Seek(int64(pos), wh)
	return dur, err
}

func (t *Tape) Close() error {
	return t.f.Close()
}
