package mix

import "os"

type console struct{}

func (c console) Read(p []byte) (int, error) {
	return os.Stdin.Read(p)
}

func (c console) Write(p []byte) (int, error) {
	return os.Stdout.Write(p)
}

func (c console) Close() error {
	return nil
}
