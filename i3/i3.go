// Package i3 provides data structures for interacting with the i3bar/swaybar
// protocol.
package i3

import "syscall"

const (
	// LeftClick represents a left mouse click.
	LeftClick = iota + 1
	// MiddleClick represents a left mouse click.
	MiddleClick
	// RightClick represents a left mouse click.
	RightClick
)

// Header is the first thing that i3bar will read to determine how this program
// will interact with it.
type Header struct {
	Version     int            `json:"version"`
	StopSignal  syscall.Signal `json:"stop_signal,omitempty"`
	ContSignal  syscall.Signal `json:"cont_signal,omitempty"`
	ClickEvents bool           `json:"click_events,omitempty"`
}

// Block is a single section of the i3bar.
type Block struct {
	FullText            string `json:"full_text"`
	ShortText           string `json:"short_text,omitempty"`
	Color               string `json:"color,omitempty"`
	Background          string `json:"background,omitempty"`
	Border              string `json:"border,omitempty"`
	BorderTop           int    `json:"border_top,omitempty"`
	BorderRight         int    `json:"border_right,omitempty"`
	BorderBottom        int    `json:"border_bottom,omitempty"`
	BorderLeft          int    `json:"border_left,omitempty"`
	MinWidth            int    `json:"min_width,omitempty"`
	Align               string `json:"align,omitempty"`
	Urgent              bool   `json:"urgent,omitempty"`
	Name                string `json:"name,omitempty"`
	Instance            string `json:"instance,omitempty"`
	Separator           bool   `json:"separator,omitempty"`
	SeparatorBlockWidth int    `json:"separator_block_width,omitempty"`
	Markup              string `json:"markup,omitempty"`
}

// ClickEvent is the data sent to this program via STDIN when a click is
// registered on the i3bar.
type ClickEvent struct {
	Name      string   `json:"name"`
	Instance  string   `json:"instance"`
	Button    int      `json:"button"`
	Modifiers []string `json:"modifiers"`
	X         int      `json:"x"`
	Y         int      `json:"y"`
	RelativeX int      `json:"relative_x"`
	RelativeY int      `json:"relative_y"`
	OutputX   int      `json:"output_x"`
	OutputY   int      `json:"output_y"`
	Width     int      `json:"width"`
	Height    int      `json:"height"`
}
