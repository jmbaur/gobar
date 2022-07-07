package module

import (
	col "github.com/jmbaur/gobar/color"
	"github.com/jmbaur/gobar/i3"
)

type Text struct {
	Content string
}

func (t *Text) Run(tx chan i3.Block, rx chan i3.ClickEvent) {
	for {
		tx <- i3.Block{
			Name:      "text",
			Instance:  t.Content,
			FullText:  t.Content,
			ShortText: t.Content,
			MinWidth:  len(t.Content),
			Color:     col.Normal,
		}
	}
}
