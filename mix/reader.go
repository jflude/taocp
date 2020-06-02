package mix

import (
	"io"
	"os"
)

type CardReader struct {
	rc io.ReadCloser
}

func NewCardReader(file string) (*CardReader, error) {
	rc, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	return &CardReader{rc}, nil
}

func (*CardReader) Name() string {
	return "READER"
}

func (*CardReader) BlockSize() int {
	return 16
}

func (r *CardReader) Read(block []Word) (int64, error) {
	buf := make([]byte, 5*r.BlockSize())
	if _, err := io.ReadFull(r.rc, buf); err != nil {
		return 0, err
	}
	s := string(buf)
	if !isPunchable(s) {
		return 0, ErrInvalidCharacter
	}
	m, err := ConvertToMIX(s)
	if err != nil {
		return 0, err
	}
	copy(block, m)
	return 400000, nil
}

func (*CardReader) Write([]Word) (int64, error) {
	return 0, ErrInvalidOperation
}

func (r *CardReader) Control(m int) (int64, error) {
	return 0, ErrInvalidControl
}

func (r *CardReader) Close() error {
	return r.rc.Close()
}
