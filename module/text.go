package module

import (
	"fmt"

	col "github.com/jmbaur/gobar/color"
	"github.com/jmbaur/gobar/i3"
)

type Text struct {
	Content string
}

func (t *Text) Run(tx chan Update, rx chan i3.ClickEvent, position int) {
	fullText := t.Content
	for {
		tx <- Update{
			Block: i3.Block{
				Name:     "text",
				Instance: "text",
				FullText: fullText,
				Color:    col.Normal,
			},
			Position: position,
		}
		event := <-rx
		fullText = fmt.Sprintf("%+v", event)
	}
}
