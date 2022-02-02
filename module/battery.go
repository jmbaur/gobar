package module

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
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
		log.Println(err)
	}
	defer f.Close()
	data, err := ioutil.ReadAll(f)
	if err != nil {
		log.Println(err)
	}
	log.Println("Updated battery module")
	return fmt.Sprintf("%s: %s%%", b.Name, bytes.Trim(data, "\n"))
}
