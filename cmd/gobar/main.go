package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/jmbaur/gobar/config"
	"github.com/jmbaur/gobar/module"
)

func must(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func main() {
	exe, err := os.Executable()
	must(err)

	log.SetPrefix(fmt.Sprintf("%s: ", filepath.Base(exe)))
	log.SetFlags(log.Lmsgprefix)
	log.Println("running")

	configFile := flag.String("config", "", "Path to gobar.yaml config file")
	flag.Parse()

	cfg, err := config.GetConfig(*configFile)
	must(err)

	must(module.Run(cfg))
}
