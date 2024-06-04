package egriden

import (
	"testing"
)

func TestFreeLayers(t *testing.T) {
	g := EgridenAssets{}
	g.InitEgridenComponents()

	testoffx, testoffy := 21.0, 1984.0
	fl1 := g.CreateFreeLayerOnTop("freelayer1", testoffx, testoffy)
	offx, offy := fl1.Offsets()
	if offx != testoffx || offy != testoffy {
		t.Errorf("Offsets didn't get applied! (%.0f, %.0f) != (%.0f, %.0f)",
			testoffx, testoffy, offx, offy)
	}

	if fl1.static {
		t.Errorf("Layer shouldn't be static!")
	}

	fl1.AddGobject(
		NewBaseGobject("tester", EmptySpritePack()).Build(), 1, 1)
}
