package egriden

import "image"

// A cell is a place on a grid layer
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

// Returns true if cell coordinate is within width and height of the gridlayer
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
func screenXYtoGrid[T, R int | float64](l GridLayer, x, y T) (R, R) {
	offx := int(x - T(l.Anchor.X))
	if offx < 0 {
		offx -= l.cellDimensions.Width - 1
	}
	offy := int(y - T(l.Anchor.Y))
	if offy < 0 {
		offy -= l.cellDimensions.Height - 1
	}

	return R(offx / l.cellDimensions.Width),
		R(offy / l.cellDimensions.Height)
}

// Returns a cell at the given screen position, taking into accout the layer
// anchor. Also returns is it within bounds.
func (l *GridLayer) CellAtScreenPos(x, y int) (Cell, bool) {
	foundx, foundy := screenXYtoGrid[int, int](*l, x, y)
	return Cell{
			Coordinate: image.Pt(foundx, foundy),
			layerPtr:   l},
		l.IsXYwithinBounds(foundx, foundy)
}

// Checks if XY is within bounds on the screen, taking into account
// the layer anchor.
func (l GridLayer) IsScreenXYwithinBounds(x, y int) bool {
	return l.IsXYwithinBounds(screenXYtoGrid[int, int](l, x, y))
}

func snapScreenXYtoCellAnchor[T, R int | float64](l GridLayer, x, y T) (R, R) {
	ax, ay := screenXYtoGrid[T, R](l, x, y)
	return ax * R(l.cellDimensions.Width), ay * R(l.cellDimensions.Height)
}

// Returns the top left point of the cell on the screen, aka the anchor or
// the draw point if the cell has a object with a sprite
// (sprite draw offset is irrelevant here).
func (c Cell) Anchor() image.Point {
	return image.Point{
		c.Coordinate.X*c.layerPtr.cellDimensions.Width + c.layerPtr.Anchor.X,
		c.Coordinate.Y*c.layerPtr.cellDimensions.Height + c.layerPtr.Anchor.Y,
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
