package module

import (
	"log"
	"time"

	col "github.com/jmbaur/gobar/color"
	"github.com/jmbaur/gobar/i3"
	"golang.org/x/exp/slices"
)

type locationInfo struct {
	loc *time.Location
}

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

func (d *Datetime) print(tx chan []i3.Block, t time.Time, c col.Color) {
	blocks := []i3.Block{}

	if d.ShowAllTimezones {
		for _, loc := range d.locations {
			longFormat := d.longFormat
			if !d.verbose {
				longFormat = d.shortFormat
			}
			blocks = append(blocks, i3.Block{
				Name:      "datetime",
				Instance:  loc.String(),
				FullText:  t.In(loc).Format(longFormat),
				Color:     c.Normal(),
				ShortText: t.In(loc).Format(d.shortFormat),
				MinWidth:  len(d.shortFormat),
			})
		}
	} else {
		longFormat := d.longFormat
		if !d.verbose {
			longFormat = d.shortFormat
		}
		blocks = []i3.Block{{
			Name:      "datetime",
			Instance:  d.currentLocation.String(),
			FullText:  t.In(d.currentLocation).Format(longFormat),
			Color:     c.Normal(),
			ShortText: t.In(d.currentLocation).Format(d.shortFormat),
			MinWidth:  len(d.shortFormat),
		}}
	}

	tx <- blocks
}

// Run implements Module.
func (d *Datetime) Run(tx chan []i3.Block, rx chan i3.ClickEvent, c col.Color) {
	d.shortFormat = "15:04:05 MST"
	d.longFormat = time.RFC1123

	if len(d.Timezones) == 0 {
		d.Timezones = append(d.Timezones, "Local")
		d.shortFormat = "15:04:05"
		d.longFormat = "Mon, 02 Jan 2006 15:04:05"
	}

	// Avoid adding duplicate timezones to our list of timezones to use while
	// running. For example, if the configuration has "Local" and "UTC" set,
	// but the local timezone _is_ in UTC, then we should only have one
	// timezone in our running list of timezones.
	{
		tzMap := map[int]struct{}{}
		for _, tz := range d.Timezones {
			loc, err := time.LoadLocation(tz)
			if err != nil {
				log.Printf("error parsing timezone: %v", err)
				continue
			}
			_, offset := time.Now().In(loc).Zone()
			if _, ok := tzMap[offset]; ok {
				continue
			}
			tzMap[offset] = struct{}{}
			d.locations = append(d.locations, loc)
		}
	}

	// If all configured timezones fail to parse, ensure that at least the
	// default timezone works.
	if len(d.locations) == 0 {
		d.locations = []*time.Location{{}}
	}

	// Start at the first configured timezone.
	d.currentLocation = d.locations[0]

	ready := make(chan struct{}, 1)
	defer close(ready)

	go func() {
		ready <- struct{}{}
	}()

	for {
		select {
		case click := <-rx:
			direction := 0
			idx := 0
			switch click.Button {
			case i3.MiddleClick:
				idx = slices.IndexFunc(d.locations, func(loc *time.Location) bool {
					return loc.String() == d.currentLocation.String()
				})
				if idx < 0 || idx > len(d.locations)-1 {
					continue
				}
				d.verbose = !d.verbose
			case i3.LeftClick:
				direction = 1
			case i3.RightClick:
				direction = -1
			}

			if direction != 0 {
				newIdx := idx + direction
				if newIdx >= len(d.locations) {
					newIdx = 0
				} else if newIdx < 0 {
					newIdx = len(d.locations) - 1
				}
				d.currentLocation = d.locations[newIdx]
			}

			d.print(tx, time.Now(), c)
		case <-ready:
			d.print(tx, time.Now(), c)
			go func() {
				time.Sleep(1 * time.Second)
				ready <- struct{}{}
			}()
		}
	}
}
