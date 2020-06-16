package mix

import (
	"io"
	"os"
)

type CardReader struct {
	rc io.ReadCloser
}

// see https://en.wikipedia.org/wiki/IBM_2540
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
	if r, ok := IsPunchable(s); !ok {
		return 0, charError(r)
	}
	m, err := ConvertToMIX(s)
	if err != nil {
		return 0, err
	}
	copy(block, m)
	return 60000, nil
}

func (*CardReader) Write([]Word) (int64, error) {
	return 0, ErrInvalidCommand
}

func (r *CardReader) Control(m int) (int64, error) {
	return 0, ErrInvalidCommand
}

func (r *CardReader) Close() error {
	return r.rc.Close()
}
