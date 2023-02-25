// Copyright (c) 2020 Justin Flude. All rights reserved.
// Use of this source code is governed by the COPYING.md file.
package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

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
	c := mix.NewComputer()
	if env, ok := os.LookupEnv("MIX_ENABLE_INTERRUPTS"); ok {
		env = strings.ToLower(env)
		c.Interrupts = (env[0] == '1' || env[0] == 'y' || env == "true")
	}
	var unit [mix.DeviceCount]string
	binding := *mix.DefaultBinding
	for i, v := range binding {
		if v != nil {
			unit[i] = v.(string)
		}
	}
	var op bool
	var trace string
	flag.IntVar(&c.BootFrom, "boot", c.BootFrom, "unit to boot from")
	flag.BoolVar(&c.Interrupts, "int", c.Interrupts, "enable interrupts")
	flag.BoolVar(&op, "op", op, "involve the operator")
	flag.StringVar(&trace, "trace", trace, "output trace to file")
	flag.IntVar(&c.Trigger, "trigger", c.Trigger,
		"trace addresses above this")
	flag.StringVar(&unit[0], "t0", unit[0], "")
	flag.StringVar(&unit[1], "t1", unit[1], "")
	flag.StringVar(&unit[2], "t2", unit[2], "")
	flag.StringVar(&unit[3], "t3", unit[3], "")
	flag.StringVar(&unit[4], "t4", unit[4], "")
	flag.StringVar(&unit[5], "t5", unit[5], "")
	flag.StringVar(&unit[6], "t6", unit[6], "")
	flag.StringVar(&unit[7], "t7", unit[7], "")
	flag.StringVar(&unit[8], "d8", unit[8], "")
	flag.StringVar(&unit[9], "d9", unit[9], "")
	flag.StringVar(&unit[10], "d10", unit[10], "")
	flag.StringVar(&unit[11], "d11", unit[11], "")
	flag.StringVar(&unit[12], "d12", unit[12], "")
	flag.StringVar(&unit[13], "d13", unit[13], "")
	flag.StringVar(&unit[14], "d14", unit[14], "")
	flag.StringVar(&unit[15], "d15", unit[15], "")
	flag.StringVar(&unit[16], "rdr", unit[16], "")
	flag.StringVar(&unit[17], "pun", unit[17], "")
	flag.StringVar(&unit[18], "prt", unit[18], "")
	flag.StringVar(&unit[19], "tty", unit[19], "")
	flag.StringVar(&unit[20], "pap", unit[20], "")
	flag.Parse()
	for i, v := range unit {
		if v != "" {
			binding[i] = v
		} else {
			binding[i] = nil
		}
	}
	if err = c.Bind(&binding); err != nil {
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
	if trace != "" {
		c.Tracer, err = os.OpenFile(trace,
			os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			return
		}
		defer func() {
			if err2 := c.Tracer.Close(); err2 != nil {
				if err == nil {
					err = err2
				} else {
					log.Println("error:", err2)
				}
			}
		}()
	}
	if !op {
		err = reportGo(c)
		return
	}
	for {
		var y bool
		if y, err = yesOrNo("Go (Y/n)? "); err != nil || !y {
			break
		}
		if err = reportGo(c); !errors.Is(err, mix.ErrHalted) {
			break
		}
	}
	return
}

func yesOrNo(prompt string) (bool, error) {
	fmt.Print(prompt)
	var s string
	_, err := fmt.Scanf("%s", &s)
	if err != nil && !strings.Contains(err.Error(), "unexpected newline") {
		return false, err
	}
	return s == "" || s[0] == 'y' || s[0] == 'Y', nil
}

func reportGo(c *mix.Computer) error {
	err := c.GoButton()
	if c.Elapsed > 0 {
		err = fmt.Errorf("%w (elapsed: %du, idle: %du)",
			err, c.Elapsed, c.Idle)
	}
	return err
}
