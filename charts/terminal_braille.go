package charts

import (
	"math"
)

// Braille pattern constants
// Braille Unicode characters are composed of 8 dots arranged in a 2x4 grid:
// 1 4
// 2 5
// 3 6
// 7 8
const (
	brailleBase = 0x2800 // Unicode base for braille patterns
)

// Braille dot masks
const (
	brailleDot1 = 0x01
	brailleDot2 = 0x02
	brailleDot3 = 0x04
	brailleDot4 = 0x08
	brailleDot5 = 0x10
	brailleDot6 = 0x20
	brailleDot7 = 0x40
	brailleDot8 = 0x80
)

// BrailleCanvas represents a canvas for drawing with braille characters
type BrailleCanvas struct {
	Width  int
	Height int
	pixels [][]bool
}

// NewBrailleCanvas creates a new braille canvas
// Each braille character represents 2x4 pixels, so canvas dimensions are multiplied
func NewBrailleCanvas(width, height int) *BrailleCanvas {
	// Each braille char is 2 pixels wide and 4 pixels tall
	pixelWidth := width * 2
	pixelHeight := height * 4

	pixels := make([][]bool, pixelHeight)
	for i := range pixels {
		pixels[i] = make([]bool, pixelWidth)
	}

	return &BrailleCanvas{
		Width:  width,
		Height: height,
		pixels: pixels,
	}
}

// SetPixel sets a pixel at the given coordinates
func (c *BrailleCanvas) SetPixel(x, y int) {
	if x >= 0 && x < c.Width*2 && y >= 0 && y < c.Height*4 {
		c.pixels[y][x] = true
	}
}

// DrawLine draws a line between two points using Bresenham's algorithm
func (c *BrailleCanvas) DrawLine(x0, y0, x1, y1 int) {
	dx := abs(x1 - x0)
	dy := abs(y1 - y0)
	sx := -1
	if x0 < x1 {
		sx = 1
	}
	sy := -1
	if y0 < y1 {
		sy = 1
	}
	err := dx - dy

	x, y := x0, y0
	for {
		c.SetPixel(x, y)

		if x == x1 && y == y1 {
			break
		}

		e2 := 2 * err
		if e2 > -dy {
			err -= dy
			x += sx
		}
		if e2 < dx {
			err += dx
			y += sy
		}
	}
}

// DrawPoint draws a single point (useful for scatter plots)
func (c *BrailleCanvas) DrawPoint(x, y int) {
	c.SetPixel(x, y)
}

// DrawCurve draws a smooth curve through multiple points
func (c *BrailleCanvas) DrawCurve(points []Point) {
	if len(points) < 2 {
		return
	}

	for i := 0; i < len(points)-1; i++ {
		x0 := int(math.Round(points[i].X))
		y0 := int(math.Round(points[i].Y))
		x1 := int(math.Round(points[i+1].X))
		y1 := int(math.Round(points[i+1].Y))
		c.DrawLine(x0, y0, x1, y1)
	}
}

// Point represents a point in 2D space
type Point struct {
	X, Y float64
}

// Render converts the canvas to a string of braille characters
func (c *BrailleCanvas) Render() string {
	result := make([]rune, c.Width*c.Height)
	idx := 0

	for charY := 0; charY < c.Height; charY++ {
		for charX := 0; charX < c.Width; charX++ {
			// Calculate braille pattern for this 2x4 block
			pattern := brailleBase

			// Map pixels to braille dots
			pixelX := charX * 2
			pixelY := charY * 4

			// Top-left column
			if pixelY < len(c.pixels) && pixelX < len(c.pixels[pixelY]) && c.pixels[pixelY][pixelX] {
				pattern |= brailleDot1
			}
			if pixelY+1 < len(c.pixels) && pixelX < len(c.pixels[pixelY+1]) && c.pixels[pixelY+1][pixelX] {
				pattern |= brailleDot2
			}
			if pixelY+2 < len(c.pixels) && pixelX < len(c.pixels[pixelY+2]) && c.pixels[pixelY+2][pixelX] {
				pattern |= brailleDot3
			}
			if pixelY+3 < len(c.pixels) && pixelX < len(c.pixels[pixelY+3]) && c.pixels[pixelY+3][pixelX] {
				pattern |= brailleDot7
			}

			// Top-right column
			if pixelY < len(c.pixels) && pixelX+1 < len(c.pixels[pixelY]) && c.pixels[pixelY][pixelX+1] {
				pattern |= brailleDot4
			}
			if pixelY+1 < len(c.pixels) && pixelX+1 < len(c.pixels[pixelY+1]) && c.pixels[pixelY+1][pixelX+1] {
				pattern |= brailleDot5
			}
			if pixelY+2 < len(c.pixels) && pixelX+1 < len(c.pixels[pixelY+2]) && c.pixels[pixelY+2][pixelX+1] {
				pattern |= brailleDot6
			}
			if pixelY+3 < len(c.pixels) && pixelX+1 < len(c.pixels[pixelY+3]) && c.pixels[pixelY+3][pixelX+1] {
				pattern |= brailleDot8
			}

			result[idx] = rune(pattern)
			idx++
		}
	}

	// Convert to string with newlines
	output := ""
	for i := 0; i < len(result); i++ {
		output += string(result[i])
		if (i+1)%c.Width == 0 && i < len(result)-1 {
			output += "\n"
		}
	}

	return output
}

// Clear clears the canvas
func (c *BrailleCanvas) Clear() {
	for y := range c.pixels {
		for x := range c.pixels[y] {
			c.pixels[y][x] = false
		}
	}
}

// FillArea fills the area under a curve (for area charts)
func (c *BrailleCanvas) FillArea(points []Point, baselineY int) {
	if len(points) < 2 {
		return
	}

	for i := 0; i < len(points); i++ {
		x := int(math.Round(points[i].X))
		y := int(math.Round(points[i].Y))

		// Draw vertical line from baseline to point
		startY, endY := baselineY, y
		if startY > endY {
			startY, endY = endY, startY
		}

		for fillY := startY; fillY <= endY; fillY++ {
			c.SetPixel(x, fillY)
		}
	}
}

// GetBrailleCharacter returns a single braille character for a given pattern
func GetBrailleCharacter(dots [8]bool) rune {
	pattern := brailleBase
	if dots[0] {
		pattern |= brailleDot1
	}
	if dots[1] {
		pattern |= brailleDot2
	}
	if dots[2] {
		pattern |= brailleDot3
	}
	if dots[3] {
		pattern |= brailleDot4
	}
	if dots[4] {
		pattern |= brailleDot5
	}
	if dots[5] {
		pattern |= brailleDot6
	}
	if dots[6] {
		pattern |= brailleDot7
	}
	if dots[7] {
		pattern |= brailleDot8
	}
	return rune(pattern)
}
