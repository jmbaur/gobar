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
	SysfsPowerSupplyCharging = "Charging"
	SysfsPowerSupplyFull     = "Full"
	pluggedInEmoji           = '\U0001F50C'
	batteryChars             = []rune{
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

func (b *Battery) print(c chan i3.Block, idx int, err error) {
	if err != nil {
		if idx < 0 {
			c <- i3.Block{
				Name:     "battery",
				Instance: "battery",
				FullText: fmt.Sprintf("BAT: %s", err),
				Color:    col.Red,
			}
		} else {
			name := b.batteries[idx].name
			c <- i3.Block{
				Name:      "battery",
				Instance:  name,
				FullText:  fmt.Sprintf("%s: %s", name, err),
				ShortText: fmt.Sprintf("%s: %s", name, err),
				Color:     col.Red,
			}
		}
		return
	}

	bat := b.batteries[idx]

	color := col.Normal

	bucket := int(math.Floor(float64(bat.capacity) / capacityBucketSize))
	capacityRune := batteryChars[bucket]

	if bat.capacity < 10 {
		color = col.Red
	} else if bat.capacity < 20 {
		color = col.Yellow
	}

	fullText := fmt.Sprintf("%s: %c %d%% (%s)", bat.name, capacityRune, bat.capacity, bat.status)
	shortText := fmt.Sprintf("%s: %d%%", bat.name, bat.capacity)
	if !bat.verbose {
		fullText = shortText
	}

	c <- i3.Block{
		Name:      "battery",
		Instance:  bat.name,
		FullText:  fullText,
		Color:     color,
		ShortText: shortText,
		MinWidth:  len(shortText),
		Urgent:    bat.capacity < 5,
	}
}

func (b *Battery) Run(tx chan i3.Block, rx chan i3.ClickEvent) {
	if err := filepath.WalkDir("/sys/class/power_supply", func(path string, d fs.DirEntry, err error) error {
		base := filepath.Base(path)
		if strings.HasPrefix(base, "BAT") {
			b.batteries = append(b.batteries, batteryInfo{name: base})
		}
		return nil
	}); err != nil {
		b.print(tx, -1, err)
		return
	}

	for i, bat := range b.batteries {
		fd, err := os.Open(fmt.Sprintf("/sys/class/power_supply/%s/uevent", bat.name))
		if err != nil {
			b.print(tx, i, err)
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
			case i3.LeftClick, i3.RightClick:
				for i, bat := range b.batteries {
					if click.Instance == bat.name {
						b.batteries[i].verbose = !bat.verbose
						b.print(tx, i, nil)
					}
				}
			}
		case <-ready:
			for i, bat := range b.batteries {
				batteryUeventMap, err := getUeventMap(bat.fd)
				if err != nil {
					b.print(tx, i, err)
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

				b.print(tx, i, nil)
			}
			go func() {
				time.Sleep(5 * time.Second)
				ready <- struct{}{}
			}()
		}
	}
}
