package main

import (
	"time"

	"github.com/jmbaur/gobar/module"
)

func main() {
	battery0 := module.Battery{Name: "BAT0"}
	wifi := module.Network{Interface: "wlan0"}
	datetime := module.Datetime{Format: time.RFC1123Z}

	module.Run(battery0, wifi, datetime)
}
