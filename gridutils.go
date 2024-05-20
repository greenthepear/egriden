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
	return R(int(x-T(l.XOffset)) / l.SquareLength),
		R(int(y-T(l.YOffset)) / l.SquareLength)
}

// Checks if XY is within bounds on the screen, taking into account the layer offsets.
func (l GridLayer) IsScreenXYwithinBounds(x, y int) bool {
	return l.IsXYwithinBounds(ScreenXYtoGrid[int, int](l, x, y))
}
