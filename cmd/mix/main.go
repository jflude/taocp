package main

import "log"

func main() {
	if err := run(); err != nil {
		log.Print("error:", err)
	}
}

func run() error {
	return nil
}
