package module

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"

	col "github.com/jmbaur/gobar/color"
	"github.com/jmbaur/gobar/i3"
)

type Battery struct {
	Index int
}

func getFileContents(f *os.File) (string, error) {
	defer f.Seek(0, io.SeekStart)
	data, err := ioutil.ReadAll(f)
	if err != nil {
		return "", err
	}
	return string(bytes.Trim(data, "\n")), nil
}

func (b *Battery) sendError(err error, c chan Update, position int) {
	c <- Update{
		Block: i3.Block{
			FullText: fmt.Sprintf("BAT%d: %s", b.Index, err),
			Color:    col.Red,
		},
		Position: position,
	}
}

func (b *Battery) Run(tx chan Update, rx chan i3.ClickEvent, position int) {
	fd, err := os.Open(fmt.Sprintf("/sys/class/power_supply/BAT%d/uevent", b.Index))
	if err != nil {
		b.sendError(err, tx, position)
		return
	}
	defer fd.Close()

	var (
		capacity   int
		status     string
		statusRune rune
		color      = col.Normal
	)

	for {
		data, err := getFileContents(fd)
		if err != nil {
			b.sendError(err, tx, position)
		}

		for _, line := range strings.Split(data, "\n") {
			split := strings.Split(line, "=")
			if len(split) != 2 {
				continue
			}
			key := split[0]
			val := split[1]
			switch key {
			case "POWER_SUPPLY_STATUS":
				status = val
			case "POWER_SUPPLY_CAPACITY":
				maybeCapacity, err := strconv.Atoi(val)
				if err != nil {
					continue
				}
				capacity = maybeCapacity
			}
		}

		switch true {
		case capacity > 80:
			color = col.Green
		case capacity < 20:
			color = col.Red
		default:
			color = col.Normal
		}

		switch status {
		case "Charging":
			statusRune = '\u2191'
		case "Discharging":
			statusRune = '\u2193'
		case "Not charging":
			statusRune = '\u26cb'
		case "Full":
			statusRune = '\u25cf'
		case "Unknown":
			fallthrough
		default:
			statusRune = '\u003f'
		}

		tx <- Update{
			Block: i3.Block{
				FullText: fmt.Sprintf("BAT%d: %c %d%%", b.Index, statusRune, capacity),
				Color:    color,
			},
			Position: position,
		}
		time.Sleep(5 * time.Second)
	}
}
