package egriden

import (
	"testing"
)

func TestFreeLayers(t *testing.T) {
	g := EgridenAssets{}
	g.InitEgridenAssets()

	testoffx, testoffy := 10.0, 20.0
	fl1 := g.CreateFreeLayerOnTop("freelayer1", testoffx, testoffy)
	offx, offy := fl1.Anchor.X, fl1.Anchor.Y
	if offx != testoffx || offy != testoffy {
		t.Errorf("Offsets didn't get applied! (%v, %v) != (%v, %v)",
			testoffx, testoffy, offx, offy)
	}
	if fl1.static {
		t.Errorf("Layer shouldn't be static!")
	}
	seq, err := CreateImageSequenceFromFolder(
		"unrev", "./examples/gridsweeper/Graphics/unrevealed")

	if err != nil {
		t.Error(err)
	}
	o := NewBaseGobject(
		"tester", NewSpritePackWithSequence(seq))
	fl1.AddGobject(
		o.Build(),
		5, 5)
}
