package module

import (
	"time"

	"github.com/jmbaur/gobar/i3"
)

type Datetime struct {
	Format string
}

func (d Datetime) Run(c chan Update, position int) error {
	for {
		c <- Update{
			Block: i3.Block{
				FullText: time.Now().Format(d.Format),
			},
			Position: position,
		}
		time.Sleep(1 * time.Second)
	}
}
