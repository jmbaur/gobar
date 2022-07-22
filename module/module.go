package module

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/jmbaur/gobar/config"
	"github.com/jmbaur/gobar/i3"
	"github.com/mitchellh/mapstructure"
	"golang.org/x/exp/slices"
)

// Module is a thing that can print to a block on the i3bar.
type Module interface {
	Run(tx chan []i3.Block, rx chan i3.ClickEvent)
}

var header = i3.Header{
	Version:     1,
	StopSignal:  syscall.SIGSTOP,
	ContSignal:  syscall.SIGCONT,
	ClickEvents: true,
}

func parseStdin(state []moduleState) {
	r := bufio.NewReader(os.Stdin)

	if _, err := r.ReadBytes('['); err != nil {
		log.Printf("error reading to opening bracket: %v", err)
		return
	}

	var parseComma bool
	for {
		if parseComma {
			if _, err := r.ReadBytes(','); err != nil {
				log.Printf("error reading to comma: %v", err)
				break
			}
		}
		b, err := r.ReadBytes('}')
		if err != nil {
			if err != io.EOF {
				log.Printf("error reading to closing brace: %v", err)
			}
			break
		}
		var event i3.ClickEvent
		if err := json.Unmarshal(b, &event); err != nil {
			log.Printf("error parsing click event: %v", err)
		}

		for i, modState := range state {
			if modState.name == event.Name {
				state[i].clickChan <- event
			}
		}
		parseComma = true
	}
}

func handleSignals(signals chan os.Signal, done chan struct{}) {
	for {
		sig := <-signals
		switch sig {
		case syscall.SIGSTOP:
		case syscall.SIGCONT:
		case syscall.SIGKILL, syscall.SIGINT, syscall.SIGTERM:
			done <- struct{}{}
		}
	}
}

type moduleState struct {
	name      string
	mod       Module
	clickChan chan i3.ClickEvent
	blocks    []i3.Block
	position  int
}

func decodeToState(cfg *config.Config) []moduleState {
	state := []moduleState{}

	for _, maybeModAny := range cfg.Modules {
		var mod Module
		maybeMod, ok := maybeModAny.(map[any]any)
		if !ok {
			continue
		}
		if maybeName, ok := maybeMod["module"]; ok {
			name, ok := maybeName.(string)
			if !ok {
				continue
			}
			switch name {
			case "battery":
				mod = &Battery{}
			case "datetime":
				mod = &Datetime{}
			case "memory":
				mod = &Memory{}
			case "network":
				mod = &Network{}
			case "text":
				mod = &Text{}
			default:
				log.Printf("module '%s' not found", maybeName)
				continue
			}
			if err := mapstructure.Decode(maybeMod, &mod); err != nil {
				log.Printf("failed to decode map structure: %v", err)
				continue
			}
			state = append(state, moduleState{
				name:      name,
				mod:       mod,
				clickChan: make(chan i3.ClickEvent),
				blocks:    []i3.Block{},
			})
		}
	}

	return state
}

// Run is the entrypoint to running a list of modules.
func Run(cfg *config.Config) error {
	state := decodeToState(cfg)

	headerData, err := json.Marshal(header)
	if err != nil {
		return err
	}
	fmt.Printf("%s\n", headerData)

	done := make(chan struct{}, 1)
	blocksChan := make(chan []i3.Block)
	defer func() {
		close(done)
		close(blocksChan)
		for _, v := range state {
			close(v.clickChan)
		}
	}()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals)
	go handleSignals(signals, done)

	for _, modState := range state {
		go modState.mod.Run(blocksChan, modState.clickChan)
	}

	go parseStdin(state)

	isDone := false
	fmt.Printf("[\n")
	for {
		select {
		case blocks := <-blocksChan:
			if len(blocks) == 0 {
				continue
			}

			pos := slices.IndexFunc(state, func(modState moduleState) bool {
				return modState.name == blocks[0].Name
			})
			if pos == -1 {
				continue
			}

			state[pos].blocks = blocks
		case <-done:
			isDone = true
		}

		blockSlice := []i3.Block{}
		for _, modState := range state {
			blockSlice = append(blockSlice, modState.blocks...)
		}
		data, err := json.Marshal(blockSlice)
		if err != nil {
			log.Printf("failed to marshal blocks to JSON: %v\n", err)
			continue
		}

		if isDone {
			fmt.Printf("%s\n]\n", data)
			break
		}

		fmt.Printf("%s,\n", data)
	}

	return nil
}
