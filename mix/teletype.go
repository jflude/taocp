package mix

import (
	"io"
	"os"
	"strings"
)

type Teletype struct {
	c   *Computer
	rwc io.ReadWriteCloser
}

func NewTeletype(c *Computer, rwc io.ReadWriteCloser) (*Teletype, error) {
	return &Teletype{c, rwc}, nil
}

func (*Teletype) Name() string {
	return "TTY"
}

func (*Teletype) BlockSize() int {
	return 14
}

func (t *Teletype) Read(block []Word) error {
	var r io.Reader
	if t.rwc != nil {
		r = t.rwc
	} else {
		r = os.Stdin
	}
	buf := make([]byte, 5*t.BlockSize())
	if _, err := r.Read(buf); err != nil {
		return err
	}
	m, err := ConvertToMIX(string(buf))
	if err != nil {
		return err
	}
	copy(block, m)
	return nil
}

func (t *Teletype) Write(block []Word) error {
	var w io.Writer
	if t.rwc != nil {
		w = t.rwc
	} else {
		w = os.Stdout
	}
	line := strings.TrimRight(ConvertToUTF8(block), " ")
	_, err := w.Write([]byte(line + "\n"))
	return err
}

func (t *Teletype) Control(m int) error {
	if m != 0 {
		return ErrInvalidControl
	}
	var rwc io.ReadWriteCloser
	if t.rwc != nil {
		rwc = t.rwc
	} else {
		rwc = os.Stdin
	}
	if s, ok := rwc.(io.Seeker); ok {
		_, err := s.Seek(0, io.SeekStart)
		return err
	}
	return ErrInvalidOperation
}

func (t *Teletype) BusyUntil() int64 {
	return 0
}

func (t *Teletype) Close() error {
	if t.rwc == nil {
		return nil
	}
	return t.rwc.Close()
}
