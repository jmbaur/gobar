// Package color provides opinionated colors.
package color

// Color is a way to 'theme' the returned hex values of colors.
type Color struct {
	Variant string // "dark" or "light"
}

// Green provides an opinionated 'green' hex value.
func (c Color) Green() string {
	return "#b5bd68"
}

// Red provides an opinionated 'red' hex value.
func (c Color) Red() string {
	return "#cc6666"
}

// Yellow provides an opinionated 'yellow' hex value.
func (c Color) Yellow() string {
	return "#f0c674"
}

// Normal provides an opinionated foreground hex value.
func (c Color) Normal() string {
	if c.Variant == "light" {
		return "#000000"
	}

	return "#ffffff"
}
