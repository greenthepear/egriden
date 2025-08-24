package egriden

import (
	"iter"

	"github.com/greenthepear/imggg"
)

func mod(a, b float64) float64 {
	ai, bi := int(a), int(b)
	m := ai % bi
	if ai < 0 && bi < 0 {
		m -= bi
	}
	if ai < 0 && bi > 0 {
		m += bi
	}
	return float64(m)
}

// A cell is a place on a GridLayer
type Cell struct {
	// Coordinate of the cell as a point. Can be out of layer bounds.
	Coordinate imggg.Point[int]
	layerPtr   *GridLayer
}

func (l *GridLayer) CellAt(x, y int) Cell {
	return Cell{
		imggg.Pt(x, y),
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

// Returns the layer's grid position corresponding to the given XY on the
// screen.
//
// Will return negative or otherwise out of bounds positions if XY is not within
// the grid on the screen so either check the screen bounds beforehand or grid
// bounds afterhand before accessing grid coordinates derived from this
// function.
func screenXYtoGrid(l GridLayer, x, y float64) (int, int) {
	offx := x - l.Anchor.X
	if offx < 0 {
		offx -= float64(l.cellDimensions.Width) - 1
	}
	offy := y - l.Anchor.Y
	if offy < 0 {
		offy -= float64(l.cellDimensions.Height) - 1
	}

	return int(offx) / (l.cellDimensions.Width + int(l.Padding.X)),
		int(offy) / (l.cellDimensions.Height + int(l.Padding.Y))
}

// Returns a cell at the given screen position, taking into account the layer
// anchor. Ignores padding, gaps are seen here as extensions to the right and
// bottom of the returned cell. To check for padding use
// [(*GridLayer).CellAtScreenPosWithPadding].
func (l *GridLayer) CellAtScreenPos(x, y float64) Cell {
	foundx, foundy := screenXYtoGrid(*l, x, y)
	return Cell{
		Coordinate: imggg.Pt(foundx, foundy),
		layerPtr:   l}
}

// Like [(*GridLayer).CellAtScreenPos] but also returns a bool as true if
// the point is NOT within a padding gap. Even if it is, the returned cell works
// like in [(*GridLayer).CellAtScreenPos].
func (l *GridLayer) CellAtScreenPosWithPadding(x, y float64) (Cell, bool) {
	withinx := mod(
		x-l.Anchor.X, float64(l.cellDimensions.Width)+l.Padding.X)
	withiny := mod(
		y-l.Anchor.Y, float64(l.cellDimensions.Height)+l.Padding.Y)
	return l.CellAtScreenPos(x, y),
		withinx <= float64(l.cellDimensions.Width) &&
			withiny <= float64(l.cellDimensions.Height)
}

// Checks if XY is within bounds on the screen, taking into account
// the layer anchor.
func (l GridLayer) IsScreenXYwithinBounds(x, y float64) bool {
	return l.IsXYwithinBounds(screenXYtoGrid(l, x, y))
}

//lint:ignore U1000 used for tests
func snapScreenXYtoCellAnchor(l GridLayer, x, y float64) (float64, float64) {
	ax, ay := screenXYtoGrid(l, x, y)
	return float64(ax*l.cellDimensions.Width) + l.Padding.X,
		float64(ay*l.cellDimensions.Height) + l.Padding.Y
}

// Returns the top left point of the cell on the screen, aka the anchor or
// the draw point if the cell has a object with a sprite
// (sprite draw offset is irrelevant here).
func (c Cell) Anchor() imggg.Point[float64] {
	return imggg.Pt[float64](
		float64(c.Coordinate.X)*
			(float64(c.layerPtr.cellDimensions.Width)+c.layerPtr.Padding.X)+
			c.layerPtr.Anchor.X,
		float64(c.Coordinate.Y)*
			(float64(c.layerPtr.cellDimensions.Height)+c.layerPtr.Padding.Y)+
			c.layerPtr.Anchor.Y,
	)
}

// Cell's bounds on the screen as a rectangle, without padding.
//
// Think of it as a rectangular area on the screen where any screen point
// put into [(*GridLayer).CellAtScreenPosWithPadding] would return this cell c
// and true, as the bounds rectangle doesn't feature padding gaps.
func (c Cell) BoundsRectangle() imggg.Rectangle[float64] {
	w, h := c.layerPtr.cellDimensions.WH()
	return imggg.Rectangle[float64]{
		Min: c.Anchor(),
		Max: c.Anchor().Add(imggg.Pt(float64(w), float64(h))),
	}
}

// Gobject within the cell, simply calls [GridLayer.GobjectAt],
// so it can be nil.
func (c Cell) Gobject() Gobject {
	return c.Layer().GobjectAt(c.XY())
}

// Returns whenever the cell is empty or not
func (c Cell) HasGobject() bool {
	return c.Gobject() != nil
}

// Iterator for all cells in a GridLayer, iterates a row at a time, top to
// bottom left to right.
func (l GridLayer) AllCells() iter.Seq[Cell] {
	w, h := l.Dimensions()
	return func(yield func(Cell) bool) {
		for hi := range h {
			for wi := range w {
				if !yield(l.CellAt(wi, hi)) {
					return
				}
			}
		}
	}
}
