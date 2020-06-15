package mix

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
)

const maxTapeBlock = 300000

type Tape struct {
	f    *os.File
	name string
	here int64
}

var ErrInvalidBlock = errors.New("mix: invalid block")

// see https://www.ibm.com/ibm/history/exhibits/storage/storage_2420.html
func NewTape(file string, unit int) (*Tape, error) {
	f, err := os.OpenFile(file, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}
	return &Tape{f, fmt.Sprintf("TAPE%02d", unit), 0}, nil
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
	if _, err := io.ReadFull(t.f, buf); err != nil {
		return 0, err
	}
	for i, j := 0, 0; i < len(block); i, j = i+1, j+4 {
		block[i] = Word(binary.LittleEndian.Uint32(buf[j : j+4]))
	}
	t.here += int64(4 * t.BlockSize())
	return 6000, nil // TODO: check timing
}

func (t *Tape) Write(block []Word) (int64, error) {
	if t.isPastEnd(4 * int64(t.BlockSize())) {
		return 0, ErrInvalidBlock
	}
	buf := make([]byte, 4*len(block))
	for i, j := 0, 0; i < len(block); i, j = i+1, j+4 {
		binary.LittleEndian.PutUint32(buf[j:j+4], uint32(block[i]))
	}
	_, err := t.f.Write(buf)
	if err == nil {
		t.here += int64(4 * t.BlockSize())
	}
	return 6000, err // TODO: check timing
}

func (t *Tape) Control(m int) (int64, error) {
	var off, dur int64
	var wh int
	if m == 0 {
		off = 0
		wh = io.SeekStart
		dur = 60000000 // TODO: check timing
	} else {
		off = int64(4 * t.BlockSize() * m)
		if t.isPastEnd(off) {
			return 0, ErrInvalidBlock
		}
		if t.here+off < 0 {
			off = -t.here
		}
		wh = io.SeekCurrent
		dur = 30000 // TODO: check timing
	}
	var err error
	t.here, err = t.f.Seek(off, wh)
	return dur, err
}

func (t *Tape) Close() error {
	return t.f.Close()
}

func (t *Tape) isPastEnd(offset int64) bool {
	return t.here+offset > 4*int64(t.BlockSize())*maxTapeBlock
}
