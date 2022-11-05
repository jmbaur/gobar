package module

import (
	"fmt"
	"io"
	"io/fs"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	col "github.com/jmbaur/gobar/color"
	"github.com/jmbaur/gobar/i3"
)

var (
	batteryChars = []rune{
		'\u2581',
		'\u2582',
		'\u2583',
		'\u2584',
		'\u2585',
		'\u2586',
		'\u2587',
		'\u2588',
	}
	capacityBucketSize = float64(100) / float64(len(batteryChars)-1)
)

// Battery is a module that prints the capacity and charging status of
// batteries. Only works on Linux.
type Battery struct {
	batteries []batteryInfo
}

type batteryInfo struct {
	fd           *os.File
	verbose      bool
	capacity     int
	name, status string
}

// getUeventMap reads the contents of a uevent formatted file (e.g.:
// KEY=VAL\n...) and seeks the fd pointer back to the beginning of the file so
// the fd may remain open for continued reads.
func getUeventMap(f *os.File) (map[string]string, error) {
	defer f.Seek(0, io.SeekStart)
	data, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	ueventMap := make(map[string]string)
	for _, line := range strings.Split(string(data), "\n") {
		split := strings.Split(line, "=")
		if len(split) != 2 {
			continue
		}
		key := split[0]
		val := split[1]
		ueventMap[key] = val
	}
	return ueventMap, nil
}

func (b *Battery) print(tx chan []i3.Block, err error, c col.Color) {
	if err != nil {
		tx <- []i3.Block{{
			Name:     "battery",
			Instance: "battery",
			FullText: fmt.Sprintf("BAT: %s", err),
			Color:    c.Red(),
		}}
		return
	}

	blocks := []i3.Block{}
	for _, bat := range b.batteries {
		color := c.Normal()

		bucket := int(math.Floor(float64(bat.capacity) / capacityBucketSize))
		capacityRune := batteryChars[bucket]

		if bat.capacity < 10 {
			color = c.Red()
		} else if bat.capacity < 20 {
			color = c.Yellow()
		}

		fullText := fmt.Sprintf("%s: %c %d%% (%s)", bat.name, capacityRune, bat.capacity, bat.status)
		shortText := fmt.Sprintf("%s: %d%%", bat.name, bat.capacity)
		if !bat.verbose {
			fullText = shortText
		}

		blocks = append(blocks, i3.Block{
			Name:      "battery",
			Instance:  bat.name,
			FullText:  fullText,
			Color:     color,
			ShortText: shortText,
			MinWidth:  len(shortText),
			Urgent:    bat.capacity < 5,
		})
	}
	tx <- blocks
}

// Run implements Module.
func (b *Battery) Run(tx chan []i3.Block, rx chan i3.ClickEvent, c col.Color) {
	if err := filepath.WalkDir("/sys/class/power_supply", func(path string, d fs.DirEntry, err error) error {
		base := filepath.Base(path)
		if strings.HasPrefix(base, "BAT") {
			b.batteries = append(b.batteries, batteryInfo{name: base})
		}
		return nil
	}); err != nil {
		b.print(tx, err, c)
		return
	}

	for i, bat := range b.batteries {
		fd, err := os.Open(fmt.Sprintf("/sys/class/power_supply/%s/uevent", bat.name))
		if err != nil {
			b.print(tx, err, c)
			return
		}
		b.batteries[i].fd = fd
	}

	ready := make(chan struct{}, 1)

	defer func() {
		close(ready)
		for _, bat := range b.batteries {
			bat.fd.Close()
		}
	}()

	go func() {
		ready <- struct{}{}
	}()

	for {
		select {
		case click := <-rx:
			switch click.Button {
			case i3.MiddleClick:
				for i, bat := range b.batteries {
					if click.Instance == bat.name {
						b.batteries[i].verbose = !bat.verbose
						b.print(tx, nil, c)
					}
				}
			}
		case <-ready:
			for i, bat := range b.batteries {
				batteryUeventMap, err := getUeventMap(bat.fd)
				if err != nil {
					b.print(tx, err, c)
				}
				for k, v := range batteryUeventMap {
					switch k {
					case "POWER_SUPPLY_CAPACITY":
						c, err := strconv.Atoi(v)
						if err != nil {
							continue
						}
						b.batteries[i].capacity = c
					case "POWER_SUPPLY_STATUS":
						b.batteries[i].status = v
					}
				}

			}

			b.print(tx, nil, c)

			go func() {
				time.Sleep(5 * time.Second)
				ready <- struct{}{}
			}()
		}
	}
}
