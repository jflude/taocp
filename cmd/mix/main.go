package main

import (
	"flag"
	"log"

	"github.com/jflude/gnuth/mix"
)

func main() {
	if err := run(); err != nil {
		log.Println("error:", err)
	}
}

func run() (err error) {
	b := *mix.DefaultBinding
	flag.StringVar(&b[0], "t00", b[0], "")
	flag.StringVar(&b[1], "t01", b[1], "")
	flag.StringVar(&b[2], "t02", b[2], "")
	flag.StringVar(&b[3], "t03", b[3], "")
	flag.StringVar(&b[4], "t04", b[4], "")
	flag.StringVar(&b[5], "t05", b[5], "")
	flag.StringVar(&b[6], "t06", b[6], "")
	flag.StringVar(&b[7], "t07", b[7], "")
	flag.StringVar(&b[8], "d08", b[8], "")
	flag.StringVar(&b[9], "d09", b[9], "")
	flag.StringVar(&b[10], "d10", b[10], "")
	flag.StringVar(&b[11], "d11", b[11], "")
	flag.StringVar(&b[12], "d12", b[12], "")
	flag.StringVar(&b[13], "d13", b[13], "")
	flag.StringVar(&b[14], "d14", b[14], "")
	flag.StringVar(&b[15], "d15", b[15], "")
	flag.StringVar(&b[16], "rdr", b[16], "")
	flag.StringVar(&b[17], "pun", b[17], "")
	flag.StringVar(&b[18], "prt", b[18], "")
	flag.StringVar(&b[19], "tty", b[19], " (default stdin/stdout)")
	flag.Parse()
	c := mix.NewComputer(&b)
	defer func() {
		if err2 := c.Shutdown(); err2 != nil {
			if err == nil {
				err = err2
			} else {
				log.Println("error:", err)
			}
		}
	}()
	return c.GoButton(16)
}
