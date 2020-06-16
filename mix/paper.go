package mix

import (
	"io"
	"os"
)

type PaperTape struct {
	rwc io.ReadWriteCloser
}

// see https://en.wikipedia.org/wiki/Teletype_Model_33
func NewPaperTape(file string) (*PaperTape, error) {
	rwc, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	return &PaperTape{rwc}, nil
}

func (*PaperTape) Name() string {
	return "PAPER"
}

func (*PaperTape) BlockSize() int {
	return 14
}

func (p *PaperTape) Read(block []Word) (int64, error) {
	buf := make([]byte, 5*p.BlockSize())
	if _, err := io.ReadFull(p.rwc, buf); err != nil {
		return 0, err
	}
	m, err := ConvertToMIX(string(buf))
	if err != nil {
		return 0, err
	}
	copy(block, m)
	return 70000, nil
}

func (p *PaperTape) Write(block []Word) (int64, error) {
	_, err := io.WriteString(p.rwc, ConvertToUTF8(block))
	return 200000, err
}

func (p *PaperTape) Control(m int) (int64, error) {
	if m != 0 {
		return 0, ErrInvalidCommand
	}
	if s, ok := p.rwc.(io.Seeker); ok {
		_, err := s.Seek(0, io.SeekStart)
		return 60000000, err
	}
	return 0, ErrInvalidCommand
}

func (p *PaperTape) Close() error {
	return p.rwc.Close()
}
