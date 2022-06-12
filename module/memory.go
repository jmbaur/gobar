package module

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	col "github.com/jmbaur/gobar/color"
	"github.com/jmbaur/gobar/i3"
)

var digitsRe = regexp.MustCompile("[0-9]+")

type Memory struct {
	memTotal float32
}

func (m *Memory) sendError(err error, c chan Update, position int) {
	c <- Update{
		Block: i3.Block{
			FullText: fmt.Sprintf("MEM: %s", err),
			Color:    col.Red,
		},
		Position: position,
	}
}

func (m *Memory) Run(tx chan Update, rx chan i3.ClickEvent, position int) {
	f, err := os.Open("/proc/meminfo")
	if err != nil {
		m.sendError(err, tx, position)
		return
	}
	defer f.Close()

outer:
	for {
		var memAvailable float32

		data, err := io.ReadAll(f)
		if err != nil {
			m.sendError(err, tx, position)
			continue
		}
		_, err = f.Seek(0, io.SeekStart)
		if err != nil {
			m.sendError(err, tx, position)
			continue
		}

		for _, line := range strings.Split(string(data), "\n") {
			split := strings.Split(line, ":")
			if len(split) != 2 {
				continue
			}
			key := strings.TrimSpace(split[0])
			val := strings.TrimSpace(split[1])
			switch key {
			case "MemTotal":
				if m.memTotal != 0 {
					continue
				}
				memTotal, err := strconv.Atoi(digitsRe.FindString(val))
				if err != nil {
					break outer
				}
				m.memTotal = float32(memTotal)
			case "MemAvailable":
				memAvailableInt, err := strconv.Atoi(digitsRe.FindString(val))
				if err != nil {
					break outer
				}
				memAvailable = float32(memAvailableInt)
			default:
				if m.memTotal != 0 && memAvailable != 0 {
					break
				}
				continue
			}
		}

		percentUnavailable := ((m.memTotal - memAvailable) / m.memTotal) * 100

		color := col.Normal
		switch true {
		case percentUnavailable > 50:
			color = col.Yellow
		case percentUnavailable > 75:
			color = col.Red
		}

		tx <- Update{
			Block: i3.Block{
				FullText: fmt.Sprintf("MEM: %.0f%%", percentUnavailable),
				Color:    color,
			},
			Position: position,
		}
		time.Sleep(5 * time.Second)
	}
}
