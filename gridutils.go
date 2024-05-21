package egriden

// Returns true if cell coordinate is within width and height of the grid layer
func (l GridLayer) IsXYwithinBounds(x, y int) bool {
	return x >= 0 && x < l.Width && y >= 0 && y < l.Height
}

// Returns the layer's grid position corresponding to the given XY on the screen.
//
// Will return negative or otherwise out of bounds positions if XY is not within the grid
// on the screen so either check the screen bounds beforehand or grid bounds afterhand
// before accessing grid coordinates derived from this function.
func ScreenXYtoGrid[T, R int | float64](l GridLayer, x, y T) (R, R) {
	offx := int(x - T(l.XOffset))
	if offx < 0 {
		offx -= l.SquareLength - 1
	}
	offy := int(y - T(l.YOffset))
	if offy < 0 {
		offy -= l.SquareLength - 1
	}

	return R(offx / l.SquareLength),
		R(offy / l.SquareLength)
}

// Checks if XY is within bounds on the screen, taking into account the layer offsets.
func (l GridLayer) IsScreenXYwithinBounds(x, y int) bool {
	return l.IsXYwithinBounds(ScreenXYtoGrid[int, int](l, x, y))
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
	ax, ay := ScreenXYtoGrid[T, R](l, x, y)
	return ax * R(l.SquareLength), ay * R(l.SquareLength)
}
