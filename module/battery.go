package module

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

type Battery struct {
	Name string
}

func (b Battery) String() string {
	f, err := os.Open(fmt.Sprintf("/sys/class/power_supply/%s/capacity", b.Name))
	if err != nil {
		log.Println(err)
		return "ERR"
	}
	defer f.Close()
	data, err := ioutil.ReadAll(f)
	if err != nil {
		log.Println(err)
		return "ERR"
	}
	return fmt.Sprintf("%s: %s%%", b.Name, bytes.Trim(data, "\n"))
}
