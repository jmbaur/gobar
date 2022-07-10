package module

import (
	"log"
	"time"

	col "github.com/jmbaur/gobar/color"
	"github.com/jmbaur/gobar/i3"
)

type locationInfo struct {
	loc     *time.Location
	verbose bool
}

// Datetime is a module for printing the date and time.
type Datetime struct {
	// For example: Local, UTC, Europe/Zurich, etc.
	Timezones []string `mapstructure:"timezones"`
	// Whether to show all timezones at once. If false, the timezones can be
	// toggled with a middle click.
	ShowAllTimezones bool `mapstructure:"show_all_timezones"`

	currentLocation int
	locations       []locationInfo
	shortFormat     string
	longFormat      string
}

func (d *Datetime) print(tx chan []i3.Block, t time.Time) {
	longFormat := d.longFormat

	blocks := []i3.Block{}

	if d.ShowAllTimezones {
		for _, locInfo := range d.locations {
			if !locInfo.verbose {
				longFormat = d.shortFormat
			}
			blocks = append(blocks, i3.Block{
				Name:      "datetime",
				Instance:  locInfo.loc.String(),
				FullText:  t.In(locInfo.loc).Format(longFormat),
				Color:     col.Normal,
				ShortText: t.In(locInfo.loc).Format(d.shortFormat),
				MinWidth:  len(d.shortFormat),
			})
		}
	} else {
		blocks = []i3.Block{{
			Name:      "datetime",
			Instance:  "datetime",
			FullText:  t.In(d.locations[d.currentLocation].loc).Format(longFormat),
			Color:     col.Normal,
			ShortText: t.In(d.locations[d.currentLocation].loc).Format(d.shortFormat),
			MinWidth:  len(d.shortFormat),
		}}
	}

	tx <- blocks
}

func (d *Datetime) Run(tx chan []i3.Block, rx chan i3.ClickEvent) {
	now := time.Now()

	d.shortFormat = "15:04:05 MST"
	d.longFormat = time.RFC1123

	if len(d.Timezones) == 0 {
		d.Timezones = append(d.Timezones, "Local")
		d.shortFormat = "15:04:05"
		d.longFormat = "Mon, 02 Jan 2006 15:04:05"
	}

	tzMap := map[string]struct{}{}
	for _, tz := range d.Timezones {
		loc, err := time.LoadLocation(tz)
		if err != nil {
			log.Printf("error parsing timezone: %v", err)
		}
		// avoid duplicate timezones
		name, _ := now.In(loc).Zone()
		if _, ok := tzMap[name]; ok {
			continue
		}
		tzMap[name] = struct{}{}
		d.locations = append(d.locations, locationInfo{
			loc:     loc,
			verbose: false,
		})
	}

	d.currentLocation = 0

	ready := make(chan struct{}, 1)
	defer close(ready)

	go func() {
		ready <- struct{}{}
	}()

	for {
		select {
		case click := <-rx:
			direction := 0
			switch click.Button {
			case i3.MiddleClick:
				d.locations[d.currentLocation].verbose = !d.locations[d.currentLocation].verbose
			case i3.LeftClick:
				direction = 1
			case i3.RightClick:
				direction = -1
			}

			if direction != 0 {
				idx := d.currentLocation + direction
				if idx >= len(d.locations) {
					idx = 0
				} else if idx < 0 {
					idx = len(d.locations) - 1
				}
				d.currentLocation = idx
			}

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
