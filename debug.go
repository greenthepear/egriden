package egriden

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// Draws a checkerboard pattern with specified colors that represents
// the cells of the GridLayer
func (l *GridLayer) DebugDrawCheckerBoard(
	c1, c2 color.Color, on *ebiten.Image) {

	w, h := l.layerDimensions.WH()
	currentColor := c1
	flipper := true
	for x := range w {
		for y := range h {
			if flipper {
				currentColor = c1
			} else {
				currentColor = c2
			}
			rec := l.CellAt(x, y).BoundsRectangle()

			vector.DrawFilledRect(on,
				float32(rec.Min.X), float32(rec.Min.Y),
				float32(rec.Dx()), float32(rec.Dy()),
				currentColor, false)

			flipper = !flipper
		}
		if w%2 == 0 {
			flipper = !flipper
		}
	}
}
