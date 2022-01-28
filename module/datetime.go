package module

import "time"

type Datetime struct {
	Format string
}

func (d Datetime) String() string {
	t := time.Now()
	return t.Format(d.Format)
}
