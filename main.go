package main

import (
	"flag"
	"time"

	"github.com/jmbaur/gobar/module"
)

func main() {
	flag.Parse()

	battery0 := module.Battery{Name: "BAT0"}
	battery1 := module.Battery{Name: "BAT1"}
	datetime := module.Datetime{Format: time.RFC1123Z}

	module.Run("|", battery0, battery1, datetime)
}
