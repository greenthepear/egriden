package egriden

import (
	"image"
)

func mod(a, b int) int {
	m := a % b
	if a < 0 && b < 0 {
		m -= b
	}
	if a < 0 && b > 0 {
		m += b
	}
	return m
}

// A cell is a place on a GridLayer
type Cell struct {
	// Coordinate of the cell as a point. Can be out of layer bounds.
	Coordinate image.Point
	layerPtr   *GridLayer
}

func (l *GridLayer) CellAt(x, y int) Cell {
	return Cell{
		image.Point{x, y},
		l,
	}
}

func (l *GridLayer) CellOfGobject(o Gobject) Cell {
	return l.CellAt(o.GridPos().XY())
}

// Cell.Coordinate's X and Y as ints
func (c Cell) XY() (int, int) {
	return c.Coordinate.X, c.Coordinate.Y
}

// Cell.Coordinate's X and Y as two floats
func (c Cell) XYf() (float64, float64) {
	return float64(c.Coordinate.X), float64(c.Coordinate.Y)
}

// Pointer to the grid layer which the cell is in
func (c Cell) Layer() *GridLayer {
	return c.layerPtr
}

// Returns true if cell coordinate is within width and height of the GridLayer
func (l GridLayer) IsXYwithinBounds(x, y int) bool {
	return x >= 0 && x < l.layerDimensions.Width &&
		y >= 0 && y < l.layerDimensions.Height
}

// A cell can have negative coordinates or one's beyond the width and height,
// this checks if that is not the case.
func (c Cell) IsWithinBounds() bool {
	return c.layerPtr.IsXYwithinBounds(c.XY())
}

// Returns the layer's grid position corresponding to the given XY on the screen.
//
// Will return negative or otherwise out of bounds positions if XY is not within the grid
// on the screen so either check the screen bounds beforehand or grid bounds afterhand
// before accessing grid coordinates derived from this function.
func screenXYtoGrid(l GridLayer, x, y int) (int, int) {
	offx := x - l.Anchor.X
	if offx < 0 {
		offx -= l.cellDimensions.Width - 1
	}
	offy := y - l.Anchor.Y
	if offy < 0 {
		offy -= l.cellDimensions.Height - 1
	}

	return offx / (l.cellDimensions.Width + l.Padding.X),
		offy / (l.cellDimensions.Height + l.Padding.Y)
}

// Returns a cell at the given screen position, taking into account the layer
// anchor. Ignores padding, gaps are seen here as extensions to the right and bottom of
// the returned cell. To check for padding use [(*GridLayer).CellAtScreenPosWithPadding].
func (l *GridLayer) CellAtScreenPos(x, y int) Cell {
	foundx, foundy := screenXYtoGrid(*l, x, y)
	return Cell{
		Coordinate: image.Pt(foundx, foundy),
		layerPtr:   l}
}

// Like [(*GridLayer).CellAtScreenPos] but also returns a bool as true if
// the point is NOT within a padding gap. Even if it is, the returned cell works like in
// [(*GridLayer).CellAtScreenPos].
func (l *GridLayer) CellAtScreenPosWithPadding(x, y int) (Cell, bool) {
	withinx := mod(x-l.Anchor.X, l.cellDimensions.Width+l.Padding.X)
	withiny := mod(y-l.Anchor.Y, l.cellDimensions.Height+l.Padding.Y)
	return l.CellAtScreenPos(x, y),
		withinx <= l.cellDimensions.Width &&
			withiny <= l.cellDimensions.Height
}

// Checks if XY is within bounds on the screen, taking into account
// the layer anchor.
func (l GridLayer) IsScreenXYwithinBounds(x, y int) bool {
	return l.IsXYwithinBounds(screenXYtoGrid(l, x, y))
}

func snapScreenXYtoCellAnchor(l GridLayer, x, y int) (int, int) {
	ax, ay := screenXYtoGrid(l, x, y)
	return ax * (l.cellDimensions.Width + l.Padding.X),
		ay * (l.cellDimensions.Height + l.Padding.Y)
}

// Returns the top left point of the cell on the screen, aka the anchor or
// the draw point if the cell has a object with a sprite
// (sprite draw offset is irrelevant here).
func (c Cell) Anchor() image.Point {
	return image.Point{
		c.Coordinate.X*
			(c.layerPtr.cellDimensions.Width+c.layerPtr.Padding.X) +
			c.layerPtr.Anchor.X,
		c.Coordinate.Y*
			(c.layerPtr.cellDimensions.Height+c.layerPtr.Padding.Y) +
			c.layerPtr.Anchor.Y,
	}
}

// Cell's bounds on the screen as a rectangle.
//
// Think of it as a rectangular area on the screen where any screen point
// put into [(*GridLayer).CellAtScreenPos] would return this cell c.
func (c Cell) BoundsRectangle() image.Rectangle {
	w, h := c.layerPtr.cellDimensions.WH()
	return image.Rectangle{
		c.Anchor(),
		c.Anchor().Add(image.Pt(w, h))}
}

// Gobject within the cell, simply calls [GridLayer.GobjectAt], so it can be nil.
func (c Cell) Gobject() Gobject {
	return c.Layer().GobjectAt(c.XY())
}
