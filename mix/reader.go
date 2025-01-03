package mix

import (
	"bufio"
	"errors"
	"io"
	"unicode/utf8"
)

type CardReader struct {
	br *bufio.Reader
}

var ErrFormat = errors.New("mix: format error")

// see https://en.wikipedia.org/wiki/IBM_2540
func NewCardReader(r io.Reader) (*CardReader, error) {
	return &CardReader{bufio.NewReader(r)}, nil
}

func (*CardReader) Name() string {
	return "READER"
}

func (*CardReader) BlockSize() int {
	return 16
}

func (r *CardReader) Read(block []Word) (int64, error) {
	s, err := r.br.ReadString('\n')
	if err != nil || s == "" {
		return 0, err
	}
	if s[len(s)-1] != '\n' || utf8.RuneCountInString(s) > 81 {
		return 0, ErrFormat
	}
	s = s[:len(s)-1]
	if ch, ok := IsPunchable(s); !ok {
		return 0, charError(ch)
	}
	m, err := EncodeAsMIX(s)
	if err != nil {
		return 0, err
	}
	copy(block, m)
	for i := len(m); i < len(block); i++ {
		block[i] = 0
	}
	return 60000, nil
}

func (*CardReader) Write([]Word) (int64, error) {
	return 0, ErrInvalidCommand
}

func (r *CardReader) Control(m int) (int64, error) {
	return 0, ErrInvalidCommand
}

func (r *CardReader) Close() error {
	return nil
}
