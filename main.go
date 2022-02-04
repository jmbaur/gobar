package main

import (
	"flag"
	"io"
	"log"
	"os"
	"time"

	"github.com/jmbaur/gobar/module"
)

func main() {
	debug := flag.Bool("debug", false, "Enable debug mode")
	// configFilePath := flag.String("config", "", "Path to config file")
	flag.Parse()

	var logger *log.Logger
	if *debug {
		logger = log.New(os.Stderr, "", log.LstdFlags)
	} else {
		logger = log.New(io.Discard, "", log.LstdFlags)
	}

	battery0 := module.Battery{Name: "BAT0", Log: logger}
	battery1 := module.Battery{Name: "BAT1", Log: logger}
	wifi := module.Network{Interface: "wlp3s0", Log: logger}
	ethernet := module.Network{Interface: "enp0s31f6", Log: logger}
	datetime := module.Datetime{Format: time.RFC1123Z, Log: logger}

	module.Run(logger, "|", battery0, battery1, wifi, ethernet, datetime)
}
