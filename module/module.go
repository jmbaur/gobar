package module

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/jmbaur/gobar/i3"
)

type Update struct {
	Block    i3.Block
	Position int
}

// Module is a thing that can print to a block on the i3bar.
type Module interface {
	Run(tx chan Update, rx chan i3.ClickEvent, position int)
}

func parseStdin(tx chan i3.ClickEvent) {
	r := bufio.NewReader(os.Stdin)
	if _, err := r.ReadBytes('['); err != nil {
		log.Println(err)
		return
	}

	for {
		data, err := r.ReadBytes('}')
		if err != nil {
			log.Println(err)
			continue
		}

		var event i3.ClickEvent
		if err := json.Unmarshal(data, &event); err != nil {
			log.Println(err)
			continue
		}
		tx <- event

		if _, err := r.ReadBytes(','); err != nil {
			log.Println(err)
			continue
		}
	}
}

func handleSignals(signals chan os.Signal, done chan struct{}) {
	for {
		sig := <-signals
		switch sig {
		case syscall.SIGSTOP:
			// TODO(jared): stop running modules
		case syscall.SIGCONT:
			// TODO(jared): continue stopped modules
		case syscall.SIGKILL, syscall.SIGINT, syscall.SIGTERM:
			done <- struct{}{}
		}
	}
}

// Run is the entrypoint to running a list of modules.
func Run(modules ...Module) error {
	header := i3.Header{
		Version:     1,
		StopSignal:  syscall.SIGSTOP,
		ContSignal:  syscall.SIGCONT,
		ClickEvents: true,
	}
	headerData, err := json.Marshal(header)
	if err != nil {
		return err
	}
	fmt.Printf("%s\n", headerData)

	done := make(chan struct{}, 1)
	events := make(chan i3.ClickEvent)
	updates := make(chan Update)
	defer func() {
		close(done)
		close(events)
		close(updates)
	}()
	go parseStdin(events)

	signals := make(chan os.Signal, 1)
	signal.Notify(signals)
	go handleSignals(signals, done)

	for i, m := range modules {
		go m.Run(updates, events, i)
	}

	blocks := make([]i3.Block, len(modules))

	isDone := false
	fmt.Printf("[\n")
	for {
		select {
		case u := <-updates:
			if u.Position > len(blocks)-1 {
				continue
			}
			blocks[u.Position] = u.Block
		case <-done:
			isDone = true
		}
		data, err := json.Marshal(blocks)
		if err != nil {
			log.Println(err)
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
