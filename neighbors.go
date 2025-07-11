package egriden

import (
	"github.com/greenthepear/imggg"
)

type Steps []imggg.Point[int]

// Like [(Cell).GetNeighbors]. Get a slice of neighboring cells, dictated by the steps slice
// and whenever the cells return true for checkFunc.
func (c Cell) GetNeighborsFunc(
	steps Steps, includeSelf bool, inbounds bool,
	checkFunc func(Cell) bool) []Cell {

	l := c.layerPtr
	r := make([]Cell, 0)
	if includeSelf {
		r = append(r, c)
	}

	for _, step := range steps {
		nc := l.CellAt(c.Coordinate.Add(step).XY())
		if inbounds && !nc.IsWithinBounds() {
			continue
		}
		if !checkFunc(nc) {
			continue
		}
		r = append(r, nc)
	}
	return r
}

// Get a slice of neighboring cells dictated by steps vectors. A step to the left by one
// would be Pt(-1, 0), a step to the top right by two would be Pt(2, 2) etc.
//
// includeSelf will add the origin cell at the first position.
// If inbounds is true, the method discards any cell that would be out of bounds in the grid.
func (c Cell) GetNeighbors(
	steps []imggg.Point[int], includeSelf bool, inbounds bool) []Cell {

	return c.GetNeighborsFunc(
		steps, includeSelf, inbounds, func(c Cell) bool { return true })
}

// Like [(Cell).GetNeighborsFunc] but it returns a "set" - a map of empty structs.
func (c Cell) GetNeighborsSetFunc(
	steps []imggg.Point[int], includeSelf bool, inbounds bool,
	checkFunc func(Cell) bool) map[Cell]struct{} {

	l := c.layerPtr
	r := make(map[Cell]struct{})
	if includeSelf {
		r[c] = struct{}{}
	}

	for _, step := range steps {
		nc := l.CellAt(c.Coordinate.Add(step).XY())
		if inbounds && !nc.IsWithinBounds() {
			continue
		}
		if !checkFunc(nc) {
			continue
		}
		r[nc] = struct{}{}
	}
	return r
}

// Like [(Cell).GetNeighbors] but it returns a "set" - a map of empty structs.
func (c Cell) GetNeighborsSet(
	steps []imggg.Point[int], includeSelf bool, inbounds bool) map[Cell]struct{} {

	return c.GetNeighborsSetFunc(
		steps, includeSelf, inbounds, func(c Cell) bool { return true })
}

/// Pre-made step slices

// Moore neighborhood - the 8 cells around the cell
var King Steps = Steps{
	imggg.Pt(-1, -1), imggg.Pt(0, -1), imggg.Pt(1, -1),
	imggg.Pt(-1, 0), imggg.Pt(1, 0),
	imggg.Pt(-1, 1), imggg.Pt(0, 1), imggg.Pt(1, 1),
}

// Von Neumann neighborhood - cell up, right, down, left
var Rook Steps = Steps{
	imggg.Pt(0, -1), imggg.Pt(1, 0), imggg.Pt(0, 1), imggg.Pt(-1, 0),
}

// The 4 cells to the diagonals
var Bishop Steps = Steps{
	imggg.Pt(-1, -1), imggg.Pt(1, -1), imggg.Pt(1, 1), imggg.Pt(-1, 1),
}
