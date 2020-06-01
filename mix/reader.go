package mix

import (
	"io"
	"os"
)

type CardReader struct {
	c  *Computer
	rc io.ReadCloser
}

func NewCardReader(c *Computer, rc io.ReadCloser) (*CardReader, error) {
	if rc == nil {
		var err error
		if rc, err = os.Open("reader.mix"); err != nil {
			return nil, err
		}
	}
	return &CardReader{c, rc}, nil
}

func (*CardReader) Name() string {
	return "READER"
}

func (*CardReader) BlockSize() int {
	return 16
}

func (r *CardReader) Read(block []Word) error {
	buf := make([]byte, 5*r.BlockSize())
	if _, err := r.rc.Read(buf); err != nil {
		return err
	}
	s := string(buf)
	if !isPunchable(s) {
		return ErrInvalidCharacter
	}
	m, err := ConvertToMIX(s)
	if err != nil {
		return err
	}
	copy(block, m)
	return nil
}

func (*CardReader) Write([]Word) error {
	return ErrInvalidOperation
}

func (r *CardReader) Control(m int) error {
	return ErrInvalidControl
}

func (r *CardReader) BusyUntil() int64 {
	return 0
}

func (c *CardReader) Close() error {
	return c.rc.Close()
}
