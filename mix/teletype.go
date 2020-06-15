package mix

import (
	"io"
	"os"
	"strings"
)

type Teletype struct {
	rwc io.ReadWriteCloser
}

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
	if _, err := r.Read(buf); err != nil {
		return 0, err
	}
	// TODO: pad with spaces after a short line terminated with a CR
	m, err := ConvertToMIX(string(buf))
	if err != nil {
		return 0, err
	}
	copy(block, m)
	return 7000000, nil
}

func (t *Teletype) Write(block []Word) (int64, error) {
	var w io.Writer
	if t.rwc != nil {
		w = t.rwc
	} else {
		w = os.Stdout
	}
	line := strings.TrimRight(ConvertToUTF8(block), " ")
	_, err := w.Write([]byte(line + "\n"))
	return 7000000, err
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
