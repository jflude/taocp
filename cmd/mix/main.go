// Copyright (c) 2020 Justin Flude. All rights reserved.
// Use of this source code is governed by the COPYING.md file.
package main

import (
	"errors"
	"flag"
	"fmt"
	"log"

	"github.com/jflude/taocp/mix"
)

func main() {
	log.SetFlags(0)
	if err := run(); err != nil {
		if errors.Is(err, mix.ErrHalted) {
			log.Println(err)
		} else {
			log.Fatalln("error:", err)
		}
	}
}

func run() (err error) {
	var opt [mix.DeviceCount]string
	b := *mix.DefaultBinding
	for i, v := range b {
		if v != nil {
			opt[i] = v.(string)
		} else {
			opt[i] = "stdin/out"
		}
	}
	c := mix.NewComputer()
	flag.BoolVar(&c.Interrupts, "int", false, "enable interrupts")
	flag.StringVar(&opt[0], "t0", opt[0], "")
	flag.StringVar(&opt[1], "t1", opt[1], "")
	flag.StringVar(&opt[2], "t2", opt[2], "")
	flag.StringVar(&opt[3], "t3", opt[3], "")
	flag.StringVar(&opt[4], "t4", opt[4], "")
	flag.StringVar(&opt[5], "t5", opt[5], "")
	flag.StringVar(&opt[6], "t6", opt[6], "")
	flag.StringVar(&opt[7], "t7", opt[7], "")
	flag.StringVar(&opt[8], "d8", opt[8], "")
	flag.StringVar(&opt[9], "d9", opt[9], "")
	flag.StringVar(&opt[10], "d10", opt[10], "")
	flag.StringVar(&opt[11], "d11", opt[11], "")
	flag.StringVar(&opt[12], "d12", opt[12], "")
	flag.StringVar(&opt[13], "d13", opt[13], "")
	flag.StringVar(&opt[14], "d14", opt[14], "")
	flag.StringVar(&opt[15], "d15", opt[15], "")
	flag.StringVar(&opt[16], "rdr", opt[16], "")
	flag.StringVar(&opt[17], "pun", opt[17], "")
	flag.StringVar(&opt[18], "prt", opt[18], "")
	flag.StringVar(&opt[19], "tty", opt[19], "")
	flag.StringVar(&opt[20], "pap", opt[20], "")
	flag.Parse()
	for i, v := range opt {
		if v != "" {
			b[i] = v
		} else {
			b[i] = nil
		}
	}
	if err = c.Bind(&b); err != nil {
		return
	}
	defer func() {
		if err2 := c.Shutdown(); err2 != nil {
			if err == nil {
				err = err2
			} else {
				log.Println("error:", err2)
			}
		}
	}()
	if err = c.GoButton(16); c.Elapsed > 0 {
		err = fmt.Errorf("%w (elapsed: %du, idle: %du)",
			err, c.Elapsed, c.Idle)
	}
	return
}
