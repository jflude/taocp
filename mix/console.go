// Copyright (c) 2020 Justin Flude. All rights reserved.
// Use of this source code is governed by the COPYING.md file.
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
