package module

import (
	"time"

	col "github.com/jmbaur/gobar/color"
	"github.com/jmbaur/gobar/i3"
)

// Datetime is a module for printing the date and time.
type Datetime struct {
	// Format is the date format string, in Go form (see time.Layout).
	Format string
	// Interval determines how often to update the time in seconds.
	Interval time.Duration
}

func (d *Datetime) Run(tx chan Update, rx chan i3.ClickEvent, position int) {
	if d.Interval == 0 {
		d.Interval = 1
	}
	for {
		tx <- Update{
			Block: i3.Block{
				FullText: time.Now().Format(d.Format),
				Color:    col.Normal,
			},
			Position: position,
		}
		time.Sleep(d.Interval * time.Second)
	}
}
