package module

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

type Battery struct {
	Name string
}

func (b Battery) Interval() time.Duration {
	return 30 * time.Second
}

func (b Battery) String() string {
	f, err := os.Open(fmt.Sprintf("/sys/class/power_supply/%s/capacity", b.Name))
	if err != nil {
		return fmt.Sprintf("%s: %s", b.Name, err)
	}
	defer f.Close()
	data, err := ioutil.ReadAll(f)
	if err != nil {
		return fmt.Sprintf("%s: %s", b.Name, err)
	}
	return fmt.Sprintf("%s: %s%%", b.Name, bytes.Trim(data, "\n"))
}
