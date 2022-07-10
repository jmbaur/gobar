package module

import (
	"log"
	"time"

	col "github.com/jmbaur/gobar/color"
	"github.com/jmbaur/gobar/i3"
	"golang.org/x/exp/slices"
)

// Datetime is a module for printing the date and time.
type Datetime struct {
	// For example: Local, UTC, Europe/Zurich, etc.
	Timezones []string `mapstructure:"timezones"`
	// Whether to show all timezones at once. If false, the timezones can be
	// toggled with a middle click.
	ShowAllTimezones bool `mapstructure:"show_all_timezones"`

	currentLocation *time.Location
	locations       []*time.Location
	shortFormat     string
	longFormat      string
	verbose         bool
}

func (d *Datetime) print(tx chan []i3.Block, t time.Time) {
	longFormat := d.longFormat
	if !d.verbose {
		longFormat = d.shortFormat
	}

	blocks := []i3.Block{}

	if d.ShowAllTimezones {
		for _, loc := range d.locations {
			blocks = append(blocks, i3.Block{
				Name:      "datetime",
				Instance:  loc.String(),
				FullText:  t.In(loc).Format(longFormat),
				Color:     col.Normal,
				ShortText: t.In(loc).Format(d.shortFormat),
				MinWidth:  len(d.shortFormat),
			})
		}
	} else {
		blocks = []i3.Block{{
			Name:      "datetime",
			Instance:  "datetime",
			FullText:  t.In(d.currentLocation).Format(longFormat),
			Color:     col.Normal,
			ShortText: t.In(d.currentLocation).Format(d.shortFormat),
			MinWidth:  len(d.shortFormat),
		}}
	}

	tx <- blocks
}

func (d *Datetime) Run(tx chan []i3.Block, rx chan i3.ClickEvent) {
	d.verbose = true
	d.shortFormat = "15:04:05 MST"
	d.longFormat = time.RFC1123

	if len(d.Timezones) == 0 {
		d.Timezones = append(d.Timezones, "Local")
		d.shortFormat = "15:04:05"
		d.longFormat = "Mon, 02 Jan 2006 15:04:05"
	}

	for _, tz := range d.Timezones {
		loc, err := time.LoadLocation(tz)
		if err != nil {
			log.Printf("error parsing timezone: %v", err)
		}
		d.locations = append(d.locations, loc)
	}

	d.currentLocation = d.locations[0]

	ready := make(chan struct{}, 1)
	defer close(ready)

	go func() {
		ready <- struct{}{}
	}()

	now := time.Now()

	for {
		select {
		case click := <-rx:
			direction := 0
			switch click.Button {
			case i3.MiddleClick:
				// TODO(jared): don't make this a global switch for all blocks
				// in the module
				d.verbose = !d.verbose
			case i3.LeftClick:
				direction = 1
			case i3.RightClick:
				direction = -1
			}

			if d.ShowAllTimezones {
				continue
			}

			idx := slices.Index(d.locations, d.currentLocation)
			if idx < 0 {
				continue
			}

			idx += direction
			if idx >= len(d.locations) {
				idx = 0
			} else if idx < 0 {
				idx = len(d.locations) - 0
			}
			d.currentLocation = d.locations[idx]

			d.print(tx, now)
		case <-ready:
			now = time.Now()
			d.print(tx, now)
			go func() {
				time.Sleep(1 * time.Second)
				ready <- struct{}{}
			}()
		}
	}
}
