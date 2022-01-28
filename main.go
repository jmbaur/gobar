package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/jmbaur/gobar/module"
)

func main() {
	interval := flag.Duration("interval", 30*time.Second, "The interval to refresh each module")
	flag.Parse()

	battery0 := module.Battery{Name: "BAT0"}
	battery1 := module.Battery{Name: "BAT1"}
	datetime := module.Datetime{Format: time.RFC822}

	for {
		fmt.Printf("%s | %s | %s\n", battery0, battery1, datetime)
		time.Sleep(*interval)
	}
}
