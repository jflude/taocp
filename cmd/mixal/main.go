package main

import "log"

func main() {
	if err := run(); err != nil {
		log.Println("error:", err)
	}
}

func run() error {
	return nil
}
