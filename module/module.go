package module

import (
	"encoding/json"
	"fmt"

	"github.com/jmbaur/gobar/i3"
)

type Update struct {
	Block    i3.Block
	Position int
}

type Module interface {
	Run(c chan Update, position int) error
}

func Run(modules ...Module) {
	header := i3.Header{
		Version: 1,
		// StopSignal:  10,
		// ContSignal:  12,
		// ClickEvents: true,
	}
	if data, err := json.Marshal(header); err == nil {
		fmt.Printf("%s\n", data)
	}

	updates := make(chan Update)
	blocks := make([]i3.Block, len(modules))

	for i, m := range modules {
		go m.Run(updates, i)
	}

	fmt.Printf("[\n")
	for {
		select {
		case u := <-updates:
			if u.Position > len(blocks)-1 {
				continue
			}
			blocks[u.Position] = u.Block
		}
		if data, err := json.Marshal(blocks); err == nil {
			fmt.Printf("%s\n", data)
		}
	}
}
