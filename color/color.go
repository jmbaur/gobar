package color

type Color struct {
	Variant string // "dark" or "light"
}

func (c Color) Green() string {
	return "#11ab00"
}

func (c Color) Red() string {
	return "#ff4053"
}

func (c Color) Yellow() string {
	return "#bf8c00"
}

func (c Color) Normal() string {
	if c.Variant == "light" {
		return "#000000"
	}

	return "#ffffff"
}
