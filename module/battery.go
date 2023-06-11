package module

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"time"

	col "github.com/jmbaur/gobar/color"
	"github.com/jmbaur/gobar/i3"
)

// Battery is a module that prints the capacity of batteries. Only works on
// Linux.
type Battery struct {
	batteries []batteryInfo
}

type batteryInfo struct {
	fd       *os.File
	capacity int
	name     string
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

		if bat.capacity < 10 {
			color = c.Red()
		} else if bat.capacity < 20 {
			color = c.Yellow()
		}

		text := fmt.Sprintf("%s: %d%%", bat.name, bat.capacity)

		blocks = append(blocks, i3.Block{
			Name:      "battery",
			Instance:  bat.name,
			FullText:  text,
			Color:     color,
			ShortText: text,
			MinWidth:  len(text),
			Urgent:    bat.capacity < 5,
		})
	}
	tx <- blocks
}

// Run implements Module.
func (b *Battery) Run(tx chan []i3.Block, rx chan i3.ClickEvent, c col.Color) {
	if err := filepath.WalkDir("/sys/class/power_supply", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		base := filepath.Base(path)
		if base == "power_supply" {
			return nil
		}

		typeFile, fErr := os.Open(filepath.Join(path, "type"))
		if fErr != nil {
			return err
		}
		defer typeFile.Close()

		typeContents, readErr := io.ReadAll(typeFile)
		if readErr != nil {
			return err
		}

		// don't include a power_supply that is not classified as a battery
		if string(bytes.TrimSpace(typeContents)) != "Battery" {
			return nil
		}

		// don't include a battery that doesn't have the capacity file
		if _, err := os.Stat(filepath.Join(path, "capacity")); err != nil {
			return nil
		}

		b.batteries = append(b.batteries, batteryInfo{name: base})

		return nil
	}); err != nil {
		b.print(tx, err, c)
		return
	}

	for i, bat := range b.batteries {
		fd, err := os.Open(fmt.Sprintf("/sys/class/power_supply/%s/capacity", bat.name))
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
		// no click support for battery
		case <-rx:
		case <-ready:
			{
				for i, bat := range b.batteries {
					bat.fd.Seek(0, io.SeekStart)
					data, err := io.ReadAll(bat.fd)
					if err != nil {
						b.print(tx, err, c)
					}
					capacity, err := strconv.Atoi(string(bytes.TrimSpace(data)))
					if err != nil {
						continue
					}
					b.batteries[i].capacity = capacity
				}

				b.print(tx, nil, c)

				go func() {
					time.Sleep(5 * time.Second)
					ready <- struct{}{}
				}()
			}
		}
	}
}
