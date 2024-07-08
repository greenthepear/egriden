package egriden

import "image"

// A cell is a place on a grid layer
type Cell struct {
	Coordinate image.Point
	layerPtr   *GridLayer
}

func (c Cell) XY() (int, int) {
	return c.Coordinate.X, c.Coordinate.Y
}

func (c Cell) XYf() (float64, float64) {
	return float64(c.Coordinate.X), float64(c.Coordinate.Y)
}

func (c Cell) Layer() *GridLayer {
	return c.layerPtr
}

// Returns true if cell coordinate is within width and height of the grid layer
func (l GridLayer) IsXYwithinBounds(x, y int) bool {
	return x >= 0 && x < l.layerDimensions.width &&
		y >= 0 && y < l.layerDimensions.height
}

// Returns the layer's grid position corresponding to the given XY on the screen.
//
// Will return negative or otherwise out of bounds positions if XY is not within the grid
// on the screen so either check the screen bounds beforehand or grid bounds afterhand
// before accessing grid coordinates derived from this function.
func screenXYtoGrid[T, R int | float64](l GridLayer, x, y T) (R, R) {
	offx := int(x - T(l.AnchorPt.X))
	if offx < 0 {
		offx -= l.cellDimensions.width - 1
	}
	offy := int(y - T(l.AnchorPt.Y))
	if offy < 0 {
		offy -= l.cellDimensions.height - 1
	}

	return R(offx / l.cellDimensions.width),
		R(offy / l.cellDimensions.height)
}

func (l *GridLayer) CellAtScreenPos(x, y int) (Cell, bool) {
	foundx, foundy := screenXYtoGrid[int, int](*l, x, y)
	return Cell{
			Coordinate: image.Pt(foundx, foundy),
			layerPtr:   l},
		l.IsXYwithinBounds(foundx, foundy)
}

// Checks if XY is within bounds on the screen, taking into account the layer offsets.
func (l GridLayer) IsScreenXYwithinBounds(x, y int) bool {
	return l.IsXYwithinBounds(screenXYtoGrid[int, int](l, x, y))
}

// Returns the anchor (top left point on the screen) of a cell in a grid layer according to screen's XY.
//
// Like [ScreenXYtoGrid] it can return positions out of grid bounds if XY is.
//
// To visualize:
// ┌─┬─┬─┐          ┌─┬─┬─┐
// │⠐│ │⠄│          │⠁│ │⠁│
// ├─┼─┼─┤ becomes: ├─┼─┼─┤
// │⠠│ │ │          │⠁│ │ │
// └─┴─┴─┘          └─┴─┴─┘
func SnapScreenXYtoCellAnchor[T, R int | float64](l GridLayer, x, y T) (R, R) {
	ax, ay := screenXYtoGrid[T, R](l, x, y)
	return ax * R(l.cellDimensions.width), ay * R(l.cellDimensions.height)
}
