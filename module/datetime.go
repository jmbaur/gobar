package module

import (
	"log"
	"time"
)

type Datetime struct {
	Format string
}

func (d Datetime) Interval() time.Duration {
	return 1 * time.Second
}

func (d Datetime) String() string {
	log.Println("Updated datetime module")
	return time.Now().Format(d.Format)
}
