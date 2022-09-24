package color

type Color struct {
	Variant string // "dark" or "light"
}

func (c Color) Green() string {
	return "#b5bd68"
}

func (c Color) Red() string {
	return "#cc6666"
}

func (c Color) Yellow() string {
	return "#f0c674"
}

func (c Color) Normal() string {
	if c.Variant == "light" {
		return "#000000"
	}

	return "#ffffff"
}
