package module

import (
	"log"
	"time"
)

type Datetime struct {
	Format string
	Log    *log.Logger
}

func (d Datetime) Interval() time.Duration {
	return 1 * time.Second
}

func (d Datetime) String() string {
	defer d.Log.Println("Updated datetime module")
	return time.Now().Format(d.Format)
}
