package mix

import "os"

type Teletype struct{}

// see https://en.wikipedia.org/wiki/Teletype_Model_33
func NewTeletype(file string) (*Teletype, error) {
	return &Teletype{}, nil
}

func (*Teletype) Name() string {
	return "TTY"
}

func (*Teletype) BlockSize() int {
	return 14
}

func (t *Teletype) Read(block []Word) (int64, error) {
	buf := make([]byte, 5*t.BlockSize())
	if n, err := os.Stdin.Read(buf); n == 0 {
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
		block[i] = 0
	}
	return 7000000, nil
}

func (t *Teletype) Write(block []Word) (int64, error) {
	return 7000000, trimWrite(os.Stdout, block)
}

func (t *Teletype) Control(m int) (int64, error) {
	return 0, ErrInvalidCommand
}

func (t *Teletype) Close() error {
	return nil
}
