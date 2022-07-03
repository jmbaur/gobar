package main

import (
	"log"
	"os"

	"github.com/jmbaur/gobar/cmd"
)

func main() {
	log.SetPrefix("gobar: ")
	log.SetFlags(log.Lmsgprefix)
	if err := cmd.Run(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
