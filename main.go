package main

import (
	"log"

	"github.com/jmbaur/gobar/cmd"
)

func main() {
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}
