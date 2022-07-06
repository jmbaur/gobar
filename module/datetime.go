package module

import (
	"log"
	"time"

	col "github.com/jmbaur/gobar/color"
	"github.com/jmbaur/gobar/i3"
)

// Datetime is a module for printing the date and time.
type Datetime struct {
	// Timezones is a list of zones that will be printed (e.g. Europe/Zurich)
	Timezones   []string
	shortFormat string
	longFormat  string
}

func (d *Datetime) print(tx chan i3.Block, t time.Time) {
	tx <- i3.Block{
		Name:      "datetime",
		Instance:  "datetime",
		FullText:  t.Format(d.longFormat),
		Color:     col.Normal,
		ShortText: t.Format(d.shortFormat),
		MinWidth:  len(d.shortFormat),
	}
}

func (d *Datetime) Run(tx chan i3.Block, rx chan i3.ClickEvent) {
	d.shortFormat = "15:04:05 MST"
	d.longFormat = time.RFC1123

	if len(d.Timezones) == 0 {
		d.Timezones = append(d.Timezones, "Local")
		d.shortFormat = "15:04:05"
		d.longFormat = "Mon, 02 Jan 2006 15:04:05"
	}

	tzs := []*time.Location{}
	for _, tz := range d.Timezones {
		loc, err := time.LoadLocation(tz)
		if err != nil {
			log.Printf("error parsing timezone: %v", err)
		}
		tzs = append(tzs, loc)
	}

	ready := make(chan struct{}, 1)

	// Make sure the first time through the loop, the content is printed
	// immediately.
	go func() {
		ready <- struct{}{}
	}()

	tzIndex := 0

	now := time.Now()

	for {
		select {
		case click := <-rx:
			switch click.Button {
			case i3.LeftClick, i3.RightClick:
				if tzIndex == len(tzs)-1 {
					tzIndex = 0
				} else {
					tzIndex++
				}
			}
			d.print(tx, now.In(tzs[tzIndex]))
		case <-ready:
			now = time.Now()
			d.print(tx, now.In(tzs[tzIndex]))
			go func() {
				time.Sleep(1 * time.Second)
				ready <- struct{}{}
			}()
		}
	}
}
