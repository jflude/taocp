package main

import (
	"bytes"
	"errors"
	"io"
	"log"
	"os"

	"github.com/jflude/gnuth/mixal"
)

func main() {
	log.SetFlags(0)
	if err := run(); err != nil {
		log.Fatalln("error:", err)
	}
}

func run() error {
	if len(os.Args) != 2 {
		return errors.New("mixal [input file]")
	}
	in, err := os.Open(os.Args[1])
	if err != nil {
		return err
	}
	var buf bytes.Buffer
	if err := mixal.Assemble(in, &buf); err != nil {
		return err
	}
	out, err := os.Create("reader.mix")
	if err != nil {
		return err
	}
	_, err = io.Copy(out, &buf)
	return err
}
