// Copyright (c) 2020 Justin Flude. All rights reserved.
// Use of this source code is governed by the COPYING.md file.
package main

import (
	"bytes"
	"errors"
	"io"
	"log"
	"os"

	"github.com/jflude/taocp/mixal"
)

func main() {
	log.SetFlags(0)
	if err := run(); err != nil {
		log.Fatalln("error:", err)
	}
}

func run() (err error) {
	if len(os.Args) != 2 {
		err = errors.New("mixal [input file]")
		return
	}
	in, err := os.Open(os.Args[1])
	if err != nil {
		return
	}
	defer func() {
		if err2 := in.Close(); err2 != nil {
			if err == nil {
				err = err2
			} else {
				log.Println("error:", err2)
			}
		}
	}()
	var buf bytes.Buffer
	if err = mixal.Assemble(in, &buf); err != nil {
		return
	}
	out, err := os.Create("reader.mix")
	if err != nil {
		return
	}
	defer func() {
		if err2 := out.Close(); err2 != nil {
			if err == nil {
				err = err2
			} else {
				log.Println("error:", err2)
			}
		}
	}()
	_, err = io.Copy(out, &buf)
	return
}
