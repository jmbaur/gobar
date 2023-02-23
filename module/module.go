// Package module provides the logic for running each section of the i3/sway
// bar.
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

	col "github.com/jmbaur/gobar/color"
	"github.com/jmbaur/gobar/config"
	"github.com/jmbaur/gobar/i3"
	"github.com/mitchellh/mapstructure"
	"golang.org/x/exp/slices"
)

// Module is a thing that can print to a block on the i3bar.
type Module interface {
	Run(tx chan []i3.Block, rx chan i3.ClickEvent, c col.Color)
}

var header = i3.Header{
	Version:     1,
	StopSignal:  syscall.SIGUSR1,
	ContSignal:  syscall.SIGUSR2,
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

func handleSignals(signals chan os.Signal, done chan<- struct{}, pause chan<- bool) {
	for {
		switch <-signals {
		case syscall.SIGINT, syscall.SIGTERM:
			done <- struct{}{}
		case syscall.SIGUSR1:
			pause <- true
		case syscall.SIGUSR2:
			pause <- false
		}
	}
}

type moduleState struct {
	name      string
	mod       Module
	clickChan chan i3.ClickEvent
	blocks    []i3.Block
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

	pause := make(chan bool)

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR1, syscall.SIGUSR2)
	go handleSignals(signals, done, pause)

	c := col.Color{Variant: cfg.ColorVariant}

	for _, modState := range state {
		go modState.mod.Run(blocksChan, modState.clickChan, c)
	}

	go parseStdin(state)

	isDone := false
	isPaused := false
	fmt.Printf("[\n")
	for {
		// if we are currently paused, block until we are unpaused
		if isPaused {
			log.Println("paused")
			isPaused = <-pause
			log.Println("unpaused")
		}

		select {
		case paused := <-pause:
			{
				isPaused = paused
				continue
			}
		case blocks := <-blocksChan:
			{
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
			}
		case <-done:
			{
				isDone = true
			}
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

		fmt.Print(string(data))

		if isDone {
			fmt.Printf("\n]")
			break
		}

		fmt.Println(",")
	}

	return nil
}
