package module

import (
	"time"
)

type Datetime struct {
	Format string
}

func (d Datetime) Interval() time.Duration {
	return 1 * time.Second
}

func (d Datetime) String() string {
	return time.Now().Format(d.Format)
}
