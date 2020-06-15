package mix

import (
	"io"
	"os"
)

type Teletype struct {
	rwc io.ReadWriteCloser
}

// see https://en.wikipedia.org/wiki/Teletype_Model_33
func NewTeletype(file string) (*Teletype, error) {
	var rwc io.ReadWriteCloser
	if file != "" {
		var err error
		rwc, err = os.OpenFile(file, os.O_CREATE|os.O_RDWR, 0644)
		if err != nil {
			return nil, err
		}
	}
	return &Teletype{rwc}, nil
}

func (*Teletype) Name() string {
	return "TTY"
}

func (*Teletype) BlockSize() int {
	return 14
}

func (t *Teletype) Read(block []Word) (int64, error) {
	var r io.Reader
	if t.rwc != nil {
		r = t.rwc
	} else {
		r = os.Stdin
	}
	buf := make([]byte, 5*t.BlockSize())
	if n, err := r.Read(buf); n == 0 {
		return 0, err
	}
	if buf[len(buf)-1] == '\010' {
		buf = buf[:len(buf)-1]
	}
	m, err := ConvertToMIX(string(buf))
	if err != nil {
		return 0, err
	}
	copy(block, m)
	for i := len(m); i < len(block); i++ {
		block[i] = NewWord(0)
	}
	return 7000000, nil
}

func (t *Teletype) Write(block []Word) (int64, error) {
	var w io.Writer
	if t.rwc != nil {
		w = t.rwc
	} else {
		w = os.Stdout
	}
	return 7000000, trimWrite(w, block)
}

func (t *Teletype) Control(m int) (int64, error) {
	if m != 0 {
		return 0, ErrInvalidCommand
	}
	var rwc io.ReadWriteCloser
	if t.rwc != nil {
		rwc = t.rwc
	} else {
		rwc = os.Stdin
	}
	if s, ok := rwc.(io.Seeker); ok {
		_, err := s.Seek(0, io.SeekStart)
		return 60000000, err
	}
	return 0, ErrInvalidCommand
}

func (t *Teletype) Close() error {
	if t.rwc == nil {
		return nil
	}
	return t.rwc.Close()
}
