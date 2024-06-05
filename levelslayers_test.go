package egriden

import (
	"testing"
)

func TestLevelsLayersGobjects(t *testing.T) {
	g := &EgridenAssets{}
	g.InitEgridenAssets()

	l := g.CreateGridLayerOnTop("test", 16, 10, 12, Sparce, 0, 0)
	if l != g.GridLayer(0) {
		t.Errorf("returned layer not the same as retrieved (%p != %p)", l, g.GridLayer(0))
	}

	w, h := l.Width, l.Height
	if w != 10 || h != 12 {
		t.Errorf("wrong layer dimensions (%d x %d != %d x %d)", w, h, 10, 12)
	}
	if l.IsXYwithinBounds(12, 12) {
		t.Errorf("should be within layer bound")
	}

	testGobj := NewBaseGobject("tester", EmptySpritePack()).Build()
	l.AddGobject(testGobj, 1, 1)
	if l.GobjectAt(1, 1) == nil {
		t.Errorf("goboject not present at added location")
	}

	l2 := g.AddLevel(NewBaseLevel("tester level"))
	if l2.Index() != 1 {
		t.Errorf("Added level of wrong index (is %v)", l2.Index())
	}

	g.NextLevel()
	if g.Level().Name() != "tester level" {
		t.Errorf("level assignment and iteration failed (name is `%s`)", g.Level().Name())
	}

	l2l := g.CreateGridLayerOnTop("level 2 layer", 10, 2, 2, Dense, 0, 0)
	l2l.AddGobject(testGobj, 0, 1)
	if !l2l.IsOccupiedAt(0, 1) {
		t.Errorf("gobject not present on level 2\nGobjectAt returns: %v\nfull map: \n%v",
			l2l.GobjectAt(0, 1), l2l.mapMat)
	}

	l2l.MoveGobjectTo(testGobj, 0, 0)
	if found := l2l.GobjectAt(0, 0); l2l.IsOccupiedAt(0, 1) || found == nil || found.isMarkedForDeletion() {
		t.Errorf("MoveGobjectTo failed\ntarget space is: %v,\nstart space is: %v",
			l2l.GobjectAt(0, 0), l2l.GobjectAt(0, 1))
	}
	l2l.DeleteAt(0, 0)
	if l2l.IsOccupiedAt(0, 0) {
		t.Errorf("deletion failed: %v", l2l.GobjectAt(0, 0))
	}

	g.NextLevel()
	if g.Level().Name() != "Default" {
		t.Errorf("g.NextLevel didn't wrap around (all levels: %v)", g.Levels)
	}
}
