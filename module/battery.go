package module

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"os"
	"strconv"
	"strings"
	"time"

	col "github.com/jmbaur/gobar/color"
	"github.com/jmbaur/gobar/i3"
)

var (
	pluggedInEmoji     = '\U0001F50C'
	batteryChars       = []rune{'\u2581', '\u2582', '\u2583', '\u2584', '\u2585', '\u2586', '\u2587', '\u2588'}
	capacityBucketSize = float64(101) / float64(len(batteryChars))
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
		capacity int
		status   string
		color    = col.Normal
		fullText string
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

		bucket := int(math.Floor(float64(capacity) / capacityBucketSize))
		capacityRune := batteryChars[bucket]

		switch true {
		case (status == "Charging" || status == "Full") && capacity > 80:
			fullText = fmt.Sprintf("BAT%d: %c %c %d%%", b.Index, pluggedInEmoji, capacityRune, capacity)
			color = col.Green
		case status != "Charging" && capacity < 20:
			color = col.Red
			fallthrough
		default:
			fullText = fmt.Sprintf("BAT%d: %c %d%%", b.Index, capacityRune, capacity)
		}

		tx <- Update{
			Block: i3.Block{
				FullText: fullText,
				Color:    color,
			},
			Position: position,
		}
		time.Sleep(5 * time.Second)
	}
}
