package charpix

import "image"
import "image/color"
import "fmt"

const EmptyPixel = "  "

// FormatColor converts from an image.Color to an ANSI escape code setting
// the terminal to that color.
func FormatColor(bg bool, c color.Color) string {
	ground := 48
	if bg {
		ground = 38
	}

	red, green, blue, _ := c.RGBA()

	return fmt.Sprintf(
		"\x1b[%d;2;%d;%d;%dm",
		ground,
		red,
		green,
		blue,
	)
}

// ColoredRune represents a string of characters with a 24bit true-color
// foreground and background.
type ColoredRune struct {
	fg color.RGBA
	bg color.RGBA
	string
}

// RGBA converts a ColoredRune to an RGBA color.
func (c ColoredRune) RGBA() (r, g, b, a uint32) {
	r, g, b, a = c.bg.RGBA()
	r *= 4
	g *= 4
	b *= 4
	a *= 4
	return
}

func colorRuneModel(c color.Color) color.Color {
	red, green, blue, _ := c.RGBA()
	cr := ColoredRune{
		color.RGBA{},
		color.RGBA{uint8(red / 4), uint8(green / 4), uint8(blue / 4), ^uint8(0)},
		EmptyPixel,
	}
	return cr
}

// ColoredRuneModel is needed to construct images.
var ColoredRuneModel = color.ModelFunc(colorRuneModel)

// how many indecies wide an element is in our array.
// we'll start with in-place storage.
const indexWidth = 1

// Charpix is an Image with the ColoredRuneModel.
type Charpix struct {
	Pix    []ColoredRune
	Stride int
	Rect   image.Rectangle
}

// {{{ implement Image interface

// ColorModel returns the model for the Charpix, which is always
// ColoredRuneModel.
func (p *Charpix) ColorModel() color.Model {
	return ColoredRuneModel
}

// Bounds implements image.Image.Bounds by directly returning the Rect.
func (p *Charpix) Bounds() image.Rectangle {
	return p.Rect
}

// At implements image.Image.At by returning the color at the given location.
func (p *Charpix) At(x, y int) color.Color {
	if !(image.Point{x, y}.In(p.Rect)) {
		return ColoredRune{}
	}

	return p.Pix[p.PixOffset(x, y)]
}

// }}} end implement Image interface

// PixOffset returns the index of the first element of Pix that corresponds to
// the pixel at (x, y).
func (p *Charpix) PixOffset(x, y int) int {
	return (y-p.Rect.Min.Y)*p.Stride + (x-p.Rect.Min.X)*indexWidth
}

// SetRune sets a pixel.
func (p *Charpix) SetRune(x, y int, c ColoredRune) {
	if !(image.Point{x, y}.In(p.Rect)) {
		return
	}

	i := p.PixOffset(x, y)
	p.Pix[i] = c
}

// RowToLine converts a row of pixel runes into a styled string.
func (p *Charpix) RowToLine(y int) string {

}

// New creates a new instance of Charpix given a rectangle specifying dimensions
func New(rect image.Rectangle) *Charpix {
	stride := rect.Dx() * indexWidth
	pix := make([]ColoredRune, stride*rect.Dy())
	return &Charpix{pix, stride, rect}
}
