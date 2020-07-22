package primitive

import (
	"fmt"
	"math"
	"image/color"
	"strings"
)

type Color struct {
	R, G, B, A int
}

func RGBADiffColor(c0, c1 Color) float64 {
	var diff float64 = 0.0
	diff += math.Pow(float64(c0.R + c1.R), 2)
	diff += math.Pow(float64(c0.G + c1.G), 2)
	diff += math.Pow(float64(c0.B + c1.B), 2)
	diff = math.Sqrt(diff) / math.Sqrt(3*(math.Pow(256, 2)))
	return diff
}

func MakeColor(c color.Color) Color {
	r, g, b, a := c.RGBA()
	return Color{int(r / 257), int(g / 257), int(b / 257), int(a / 257)}
}

func MakeHexColor(x string) Color {
	x = strings.Trim(x, "#")
	var r, g, b, a int
	a = 255
	switch len(x) {
	case 3:
		fmt.Sscanf(x, "%1x%1x%1x", &r, &g, &b)
		r = (r << 4) | r
		g = (g << 4) | g
		b = (b << 4) | b
	case 4:
		fmt.Sscanf(x, "%1x%1x%1x%1x", &r, &g, &b, &a)
		r = (r << 4) | r
		g = (g << 4) | g
		b = (b << 4) | b
		a = (a << 4) | a
	case 6:
		fmt.Sscanf(x, "%02x%02x%02x", &r, &g, &b)
	case 8:
		fmt.Sscanf(x, "%02x%02x%02x%02x", &r, &g, &b, &a)
	}
	return Color{r, g, b, a}
}

func (c *Color) NRGBA() color.NRGBA {
	return color.NRGBA{uint8(c.R), uint8(c.G), uint8(c.B), uint8(c.A)}
}
