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

// Memory provides information on RAM and swap usage for the system. Only works
// on Linux.
type Memory struct {
	percentMemUnavailable  float32
	percentSwapUnavailable float32
	currentLabel           string
}

func (m *Memory) print(tx chan []i3.Block, err error, c col.Color) {
	if err != nil {
		tx <- []i3.Block{{
			Name:     "memory",
			Instance: "memory",
			FullText: fmt.Sprintf("MEM: %s", err),
			Color:    c.Red(),
			Urgent:   true,
		}}
	} else {
		urgent := false
		var percent float32
		if m.currentLabel == "SWAP" {
			percent = m.percentSwapUnavailable
		} else {
			percent = m.percentMemUnavailable
		}
		color := c.Normal()
		switch true {
		case percent > 50:
			color = c.Yellow()
		case percent > 75:
			color = c.Red()
			urgent = true
		}
		tx <- []i3.Block{{
			Name:     "memory",
			Instance: "memory",
			FullText: fmt.Sprintf("%s: %d%%", m.currentLabel, int(percent)),
			Color:    color,
			Urgent:   urgent,
		}}
	}
}

// Run implements Module.
func (m *Memory) Run(tx chan []i3.Block, rx chan i3.ClickEvent, c col.Color) {
	f, err := os.Open("/proc/meminfo")
	if err != nil {
		m.print(tx, err, c)
		return
	}

	ready := make(chan struct{}, 1)
	defer func() {
		f.Close()
		close(ready)
	}()

	go func() {
		ready <- struct{}{}
	}()

	m.currentLabel = "MEM"

outer:
	for {
		select {
		case click := <-rx:
			switch true {
			case click.Button == i3.LeftClick || click.Button == i3.RightClick:
				if m.currentLabel == "SWAP" {
					m.currentLabel = "MEM"
				} else {
					m.currentLabel = "SWAP"
				}
				m.print(tx, nil, c)
			}
		case <-ready:
			var memTotal, memAvailable, swapTotal, swapFree float32

			data, err := io.ReadAll(f)
			if err != nil {
				m.print(tx, err, c)
				continue
			}
			_, err = f.Seek(0, io.SeekStart)
			if err != nil {
				m.print(tx, err, c)
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
					memTotalInt, err := strconv.Atoi(digitsRe.FindString(val))
					if err != nil {
						break outer
					}
					memTotal = float32(memTotalInt)
				case "MemAvailable":
					memAvailableInt, err := strconv.Atoi(digitsRe.FindString(val))
					if err != nil {
						break outer
					}
					memAvailable = float32(memAvailableInt)
				case "SwapTotal":
					swapTotalInt, err := strconv.Atoi(digitsRe.FindString(val))
					if err != nil {
						break outer
					}
					swapTotal = float32(swapTotalInt)
				case "SwapFree":
					swapFreeInt, err := strconv.Atoi(digitsRe.FindString(val))
					if err != nil {
						break outer
					}
					swapFree = float32(swapFreeInt)
				default:
					if memTotal != 0 && memAvailable != 0 {
						break
					}
					continue
				}
			}

			m.percentMemUnavailable = ((memTotal - memAvailable) / memTotal) * 100
			m.percentSwapUnavailable = ((swapTotal - swapFree) / swapTotal) * 100

			m.print(tx, nil, c)

			go func() {
				time.Sleep(5 * time.Second)
				ready <- struct{}{}
			}()
		}
	}
}
