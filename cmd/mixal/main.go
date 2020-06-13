package main

import (
	"errors"
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
	out, err := os.Create("reader.mix")
	if err != nil {
		return err
	}
	return mixal.Assemble(in, out)
}
