package module

import (
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"log"
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

type Battery struct{}

// getUeventMap reads the contents of a uevent formatted file (e.g.:
// KEY=VAL\n...) and seeks the fd pointer back to the beginning of the file so
// the fd may remain open for continued reads.
func getUeventMap(f *os.File) (map[string]string, error) {
	defer f.Seek(0, io.SeekStart)
	data, err := ioutil.ReadAll(f)
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

func (b *Battery) sendError(err error, c chan Update, position int) {
	c <- Update{
		Block: i3.Block{
			FullText: fmt.Sprintf("BAT: %s", err),
			Color:    col.Red,
		},
		Position: position,
	}
}

func (b *Battery) getBlock(capacity float64, acPluggedIn bool) i3.Block {
	var (
		fullText     string
		color        = col.Normal
		presCapacity int
	)

	roundedCapacity := math.RoundToEven(capacity)
	presCapacity = int(roundedCapacity)
	bucket := int(math.Floor(roundedCapacity / capacityBucketSize))
	capacityRune := batteryChars[bucket]

	switch true {
	case acPluggedIn:
		fullText = fmt.Sprintf("BAT: %c %c %d%%", pluggedInEmoji, capacityRune, presCapacity)
		if capacity > 80 {
			color = col.Green
		}
	case capacity < 20:
		color = col.Red
	}

	if fullText == "" {
		fullText = fmt.Sprintf("BAT: %c %d%%", capacityRune, presCapacity)
	}

	return i3.Block{
		FullText: fullText,
		Color:    color,
	}
}

func (b *Battery) Run(tx chan Update, rx chan i3.ClickEvent, position int) {
	batteries := []string{} // of the form "BAT0", "BAT1", etc

	if err := filepath.WalkDir("/sys/class/power_supply", func(path string, d fs.DirEntry, err error) error {
		base := filepath.Base(path)
		if strings.HasPrefix(base, "BAT") {
			batteries = append(batteries, base)
		}
		return nil
	}); err != nil {
		log.Println(err)
		return
	}

	acFd, err := os.Open("/sys/class/power_supply/AC/uevent")
	if err != nil {
		b.sendError(err, tx, position)
		return
	}

	batteryFDs := []*os.File{}
	for _, bat := range batteries {
		fd, err := os.Open(fmt.Sprintf("/sys/class/power_supply/%s/uevent", bat))
		if err != nil {
			b.sendError(err, tx, position)
			return
		}
		defer fd.Close()
		batteryFDs = append(batteryFDs, fd)
	}

	for {
		var (
			acPluggedIn   bool
			energyFullSum int
			energyNowSum  int
		)

		acUeventMap, err := getUeventMap(acFd)
		if err != nil {
			b.sendError(err, tx, position)
		}
		if powerSupplyOnline, ok := acUeventMap["POWER_SUPPLY_ONLINE"]; ok &&
			powerSupplyOnline == "1" {
			acPluggedIn = true
		}

	batteryLoop:
		for _, fd := range batteryFDs {
			batteryUeventMap, err := getUeventMap(fd)
			if err != nil {
				b.sendError(err, tx, position)
			}
			for k, v := range batteryUeventMap {
				switch k {
				case "POWER_SUPPLY_ENERGY_FULL":
					full, err := strconv.Atoi(v)
					if err != nil {
						continue batteryLoop
					}
					energyFullSum += full
				case "POWER_SUPPLY_ENERGY_NOW":
					now, err := strconv.Atoi(v)
					if err != nil {
						continue batteryLoop
					}
					energyNowSum += now
				case "POWER_SUPPLY_STATUS":
					if !acPluggedIn {
						acPluggedIn = (v == SysfsPowerSupplyCharging)
					}
				}
			}

		}

		capacity := float64(energyNowSum) / float64(energyFullSum) * 100

		tx <- Update{
			Block:    b.getBlock(capacity, acPluggedIn),
			Position: position,
		}
		time.Sleep(5 * time.Second)
	}
}
