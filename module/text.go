package module

import (
	col "github.com/jmbaur/gobar/color"
	"github.com/jmbaur/gobar/i3"
)

// Text is a module that will just print static text content.
type Text struct {
	Content string
}

// Run implements Module.
func (t *Text) Run(tx chan []i3.Block, _ chan i3.ClickEvent, c col.Color) {
	tx <- []i3.Block{{
		Name:      "text",
		Instance:  t.Content,
		FullText:  t.Content,
		ShortText: t.Content,
		MinWidth:  len(t.Content),
		Color:     c.Normal(),
	}}
}
