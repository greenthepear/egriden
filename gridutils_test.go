package egriden

import (
	"testing"

	"github.com/greenthepear/imggg"
)

func TestGridUtilities(t *testing.T) {
	lv := NewBaseLevel("test")
	l1 := lv.CreateSimpleGridLayerOnTop("test", 20, 12, 6, Sparse, 0, 50)

	if l1.Anchor.Y != 50 {
		t.Errorf("Anchor is not 50, instead %v!", l1.Anchor.Y)
	}

	// Testing associated cells with the anchor
	shouldbe := [...]int{-3, -2, -2, -1, -1, 0, 0, 1, 1, 2}
	for i, s := range shouldbe {
		x, y := screenXYtoGrid(*l1, 0, float64(i*10))
		if y != s {
			t.Errorf(`Wrong screen XY to grid conversion with yoffset %v!
is:		%v %v
should be:	%v %v`,
				int(l1.Anchor.Y), x, y, x, s)
		}
		ax, ay := snapScreenXYtoCellAnchor(*l1, 0, float64(i*10))
		if ay != float64(s*l1.cellDimensions.Height) {
			t.Errorf(`Wrong anchor point calculation!
is:		%v %v
should be:	%v %v`,
				ax, ay, ax, float64(s*l1.cellDimensions.Height))
		}
	}

	l2 := lv.CreateGridLayerOnTop("test2", GridLayerParameters{
		GridDimensions: Dimensions{5, 5},
		CellDimensions: Dimensions{5, 5},
		PaddingVector:  imggg.Pt(2.0, 1.0),
		Anchor:         imggg.Pt(-5.0, -5.0),
	})

	type gapTest struct {
		screenpos  imggg.Point[float64]
		gridpos    imggg.Point[int]
		outsideGap bool
	}
	forTest := [...]gapTest{
		{imggg.Pt(0.0, 0.0), imggg.Pt(0, 0), true},
		{imggg.Pt(1.0, 0.0), imggg.Pt(0, 0), false},
		{imggg.Pt(2.0, 0.0), imggg.Pt(1, 0), true},
		{imggg.Pt(-1.0, -1.0), imggg.Pt(0, 0), true},
	}
	for _, e := range forTest {
		c, b := l2.CellAtScreenPosWithPadding(e.screenpos.X, e.screenpos.Y)
		cx, cy := c.XY()
		if c != l2.CellAtScreenPos(e.screenpos.X, e.screenpos.Y) ||
			cx != e.gridpos.X || cy != e.gridpos.Y || b != e.outsideGap {
			t.Errorf(`Wrong cell for padding layer! (%v)
returned:	%v, %v
target: 	%v, %v`, e.screenpos, c.Coordinate, b, e.gridpos, e.outsideGap)
		}
	}

}
