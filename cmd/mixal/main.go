package main

import (
	"bytes"
	"errors"
	"flag"
	"io"
	"log"
	"os"
	"strings"

	"github.com/jflude/taocp/mixal"
)

func main() {
	log.SetFlags(0)
	if err := run(); err != nil {
		log.Fatalln("error:", err)
	}
}

func run() (err error) {
	var interrupts bool
	if env, ok := os.LookupEnv("MIX_ENABLE_INTERRUPTS"); ok {
		env = strings.ToLower(env)
		interrupts = (env[0] == '1' || env[0] == 'y' || env == "true")
	}
	flag.BoolVar(&interrupts, "int", interrupts, "enable interrupts")
	flag.Parse()
	if flag.NArg() != 1 {
		err = errors.New("mixal [-int] INPUT-FILE")
		return
	}
	var in io.ReadCloser
	if flag.Arg(0) == "-" {
		in = os.Stdin
	} else {
		if in, err = os.Open(flag.Arg(0)); err != nil {
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
	}
	var buf bytes.Buffer
	if err = mixal.Assemble(in, &buf, interrupts); err != nil {
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
