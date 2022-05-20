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

func (b Battery) getFileContents(fileName string) (string, error) {
	f, err := os.Open(fmt.Sprintf("/sys/class/power_supply/%s/%s", b.Name, fileName))
	if err != nil {
		return "", err
	}
	defer f.Close()
	data, err := ioutil.ReadAll(f)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s", bytes.Trim(data, "\n")), nil
}

func (b Battery) Interval() time.Duration {
	return 30 * time.Second
}

func (b Battery) String() string {
	capacity, capacityErr := b.getFileContents("capacity")
	if capacityErr != nil {
		capacityLevel, capacityLevelErr := b.getFileContents("capacity_level")
		if capacityLevelErr != nil {
			return fmt.Sprintf("%s: n/a", b.Name)
		}
		return fmt.Sprintf("%s: %s", b.Name, capacityLevel)
	}
	return fmt.Sprintf("%s: %s%%", b.Name, capacity)
}
