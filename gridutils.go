package egriden

// Returns the layer's grid position corresponding to the given XY on the screen.
//
// Will return negative or otherwise out of bounds positions if XY is not within the grid
// on the screen so either check the screen bounds beforehand or grid bounds afterhand
// before accessing grid coordinates derived from this function.
func ScreenXYtoGrid[T, R int | float64](l GridLayer, x, y T) (R, R) {
	return R(int(x-T(l.XOffset)) / l.SquareLength),
		R(int(y-T(l.YOffset)) / l.SquareLength)
}
