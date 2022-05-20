package module

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"time"

	"github.com/jmbaur/gobar/color"
	"github.com/jmbaur/gobar/i3"
)

type Battery struct {
	Name string
}

func getFileContents(f *os.File) (string, error) {
	defer f.Seek(0, io.SeekStart)
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
		col := color.Normal

		if capacity, err := getFileContents(capacityFile); err != nil {
			capacityLevel, capacityLevelErr := getFileContents(capacityLevelFile)
			if capacityLevelErr != nil {
				fullText = fmt.Sprintf("%s: n/a", b.Name)
				col = color.Red
			} else {
				switch true {
				case capacityLevel == "Full":
					col = color.Green
				}
				fullText = fmt.Sprintf("%s: %s", b.Name, capacityLevel)
			}
		} else {
			if capInt, err := strconv.Atoi(capacity); err != nil {
			} else {
				switch true {
				case capInt > 80:
					col = color.Green
				case capInt < 20:
					col = color.Red
				}
				fullText = fmt.Sprintf("%s: %s%%", b.Name, capacity)
			}
		}
		c <- Update{
			Block: i3.Block{
				FullText: fullText,
				Color:    col,
			},
			Position: position,
		}
		time.Sleep(30 * time.Second)
	}
}
