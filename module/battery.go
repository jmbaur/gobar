package module

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jmbaur/gobar/color"
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
	return fmt.Sprintf("%s", bytes.Trim(data, "\n")), nil
}

func (b Battery) sendError(err error, c chan Update, position int) {
	c <- Update{
		Block: i3.Block{
			FullText: fmt.Sprintf("BAT%d: %s", b.Index, err),
			Color:    color.Red,
		},
		Position: position,
	}
}

func (b Battery) Run(c chan Update, position int) {
	log.Println("battery", position)
	fd, err := os.Open(fmt.Sprintf("/sys/class/power_supply/BAT%d/uevent", b.Index))
	if err != nil {
		b.sendError(err, c, position)
		return
	}
	defer fd.Close()

	col := color.Normal
	var capacity int

	for {
		data, err := getFileContents(fd)
		if err != nil {
			b.sendError(err, c, position)
		}

		for _, line := range strings.Split(data, "\n") {
			split := strings.Split(line, "=")
			if len(split) != 2 {
				continue
			}
			key := split[0]
			val := split[1]
			switch key {
			case "POWER_SUPPLY_CAPACITY":
				capacity, err = strconv.Atoi(val)
				if err != nil {
					continue
				}
				if capacity > 80 {
					col = color.Green
				} else if capacity < 20 {
					col = color.Red
				}
			}
		}

		c <- Update{
			Block: i3.Block{
				FullText: fmt.Sprintf("BAT%d: %d%%", b.Index, capacity),
				Color:    col,
			},
			Position: position,
		}
		time.Sleep(5 * time.Second)
	}
}
