package egriden

import (
	"testing"
)

func TestGridUtilities(t *testing.T) {
	lv := NewBaseLevel("test")
	//l0 := lv.CreateGridLayerOnTop("test", 20, 12, 6, Sparce, 0, 0)
	l1 := lv.CreateSimpleGridLayerOnTop("test", 20, 12, 6, Sparce, 0, 50)

	if l1.Anchor.Y != 50 {
		t.Errorf("Anchor is not 50, instead %d!", l1.Anchor.Y)
	}

	shouldbe := [...]int{-3, -2, -2, -1, -1, 0, 0, 1, 1, 2}
	for i := range 10 {
		x, y := screenXYtoGrid[int, int](*l1, 0, i*10)
		if y != shouldbe[i] {
			t.Errorf(`Wrong screen XY to grid conversion with yoffset %d!
is:		%d %d
should be:	%d %d`,
				int(l1.Anchor.Y), x, y, x, shouldbe[i])
		}
		ax, ay := snapScreenXYtoCellAnchor[int, int](*l1, 0, i*10)
		if ay != shouldbe[i]*l1.cellDimensions.Height {
			t.Errorf(`Wrong anchor point calculation!
is:		%d %d
should be:	%d %d`,
				ax, ay, ax, shouldbe[i]*l1.cellDimensions.Height)
		}
	}
}
