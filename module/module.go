package module

import (
	"fmt"
	"log"
	"time"
)

type Update struct {
	Content  string
	Position int
}

type Module interface {
	String() string
	Interval() time.Duration
}

type StatusLine struct {
	Separator     string
	moduleContent []string
}

func (s StatusLine) String() string {
	line := ""
	for i, content := range s.moduleContent {
		if i == 0 {
			line += fmt.Sprintf("%s", content)
		} else {
			line += fmt.Sprintf(" %s %s", s.Separator, content)
		}
	}
	return line
}

func Run(l *log.Logger, sep string, modules ...Module) {
	updater := make(chan Update)

	statusLine := StatusLine{Separator: "|"}

	for i, m := range modules {
		statusLine.moduleContent = append(statusLine.moduleContent, m.String())
		go func(m Module, position int) {
			t := time.NewTicker(m.Interval())
			for {
				_ = <-t.C
				updater <- Update{
					Content:  m.String(),
					Position: position,
				}
			}
		}(m, i)
	}

	for update := range updater {
		statusLine.moduleContent[update.Position] = update.Content
		fmt.Println(statusLine)
	}
}
