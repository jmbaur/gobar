package module

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/jmbaur/gobar/i3"
)

type Update struct {
	Block    i3.Block
	Position int
}

type Module interface {
	Run(c chan Update, position int)
}

func Run(modules ...Module) error {
	header := i3.Header{
		Version:     1,
		StopSignal:  syscall.SIGSTOP,
		ContSignal:  syscall.SIGCONT,
		ClickEvents: false, // TODO(jared): handle STDIN
	}
	if data, err := json.Marshal(header); err == nil {
		fmt.Printf("%s\n", data)
	}

	done := make(chan struct{}, 1)
	signals := make(chan os.Signal, 1)
	updates := make(chan Update)
	signal.Notify(signals)
	signal.Ignore(syscall.SIGURG)

	go func() {
		for {
			sig := <-signals
			switch sig {
			case syscall.SIGCONT:
			case syscall.SIGSTOP:
			case syscall.SIGKILL, syscall.SIGINT, syscall.SIGTERM:
				done <- struct{}{}
			}
		}
	}()

	for i, m := range modules {
		go m.Run(updates, i)
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
		if data, err := json.Marshal(blocks); err == nil {
			if isDone {
				fmt.Printf("%s\n]\n", data)
				break
			} else {
				fmt.Printf("%s,\n", data)
			}
		}
	}

	return nil
}
