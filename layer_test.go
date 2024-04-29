package egriden

import (
	"testing"
)

func TestLayerCreation(t *testing.T) {
	g := &EgridenAssets{}
	g.InitEgridenComponents()

	l := g.CreateGridLayerOnTop("test", 16, 10, 12, Sparce, 0, 0)
	if l != g.GridLayer(0) {
		t.Errorf("returned layer not the same as retrieved (%p != %p)", l, g.GridLayer(0))
	}
	w, h := l.Width, l.Height

	if w != 10 || h != 12 {
		t.Errorf("wrong layer dimensions (%d x %d != %d x %d)", w, h, 10, 12)
	}

	testGobj := NewBaseGobject("tester", EmptySpritePack()).Build()
	l.AddObject(testGobj, 1, 1)

	if l.GobjectAt(1, 1) == nil {
		t.Errorf("goboject not present at added location")
	}
}
