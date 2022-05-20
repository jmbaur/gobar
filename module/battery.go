package module

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/jmbaur/gobar/i3"
)

type Battery struct {
	Name string
}

func getFileContents(f *os.File) (string, error) {
	data, err := ioutil.ReadAll(f)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s", bytes.Trim(data, "\n")), nil
}

func (b Battery) Run(c chan Update, position int) error {
	capacityFile, err := os.Open(fmt.Sprintf("/sys/class/power_supply/%s/capacity", b.Name))
	if err != nil {
		return fmt.Errorf("failed to run battery module: %v", err)
	}
	defer capacityFile.Close()

	capacityLevelFile, err := os.Open(fmt.Sprintf("/sys/class/power_supply/%s/capacity_level", b.Name))
	if err != nil {
		return fmt.Errorf("failed to run battery module: %v", err)
	}
	defer capacityLevelFile.Close()

	for {
		var fullText string
		if capacity, err := getFileContents(capacityFile); err != nil {
			capacityLevel, capacityLevelErr := getFileContents(capacityLevelFile)
			if capacityLevelErr != nil {
				fullText = fmt.Sprintf("%s: n/a", b.Name)
			}
			fullText = fmt.Sprintf("%s: %s", b.Name, capacityLevel)
		} else {
			fullText = fmt.Sprintf("%s: %s%%", b.Name, capacity)
		}
		c <- Update{
			Block: i3.Block{
				FullText: fullText,
			},
			Position: position,
		}
		time.Sleep(30 * time.Second)
	}
}
