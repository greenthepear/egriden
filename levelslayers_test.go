package egriden

import (
	"testing"
)

func TestLevelsLayersGobjects(t *testing.T) {
	g := &EgridenAssets{}
	g.InitEgridenAssets()

	l := g.CreateGridLayerOnTop(
		"test",
		GridLayerParameters{
			GridDimensions: Dimensions{10, 12},
			CellDimensions: Dimensions{16, 16},
			Mode:           Sparse,
		},
	)
	if l != g.GridLayer(0) {
		t.Errorf("returned layer not the same as retrieved (%p != %p)", l, g.GridLayer(0))
	}

	w, h := l.layerDimensions.WH()
	if w != 10 || h != 12 {
		t.Errorf("wrong layer dimensions (%d x %d != %d x %d)", w, h, 10, 12)
	}
	if l.IsXYwithinBounds(12, 12) {
		t.Errorf("should be within layer bound")
	}

	testGobj := NewBaseGobject("tester", EmptySpritePack())
	l.AddGobject(testGobj.Build(), 1, 1)
	if l.GobjectAt(1, 1) == nil {
		t.Errorf("goboject not present at added location")
		if gx, gy := l.GobjectAt(1, 1).GridPos().XY(); gx != 1 || gy != 1 {
			t.Errorf("gobject xy not applied (%d, %d) != (%d, %d)",
				gx, gy, 1, 1)
		}
	}

	l2 := g.AddLevel(NewBaseLevel("tester level"))
	if l2.Index() != 1 {
		t.Errorf("Added level of wrong index (is %v)", l2.Index())
	}

	g.NextLevel()
	if g.Level().Name() != "tester level" {
		t.Errorf("level assignment and iteration failed (name is `%s`)", g.Level().Name())
	}

	l2l := g.CreateGridLayerOnTop(
		"test",
		GridLayerParameters{
			GridDimensions: Dimensions{2, 2},
			CellDimensions: Dimensions{10, 10},
			Mode:           Dense,
		},
	)
	testGobjCopy := testGobj.Build()
	l2l.AddGobject(testGobjCopy, 0, 1)
	if !l2l.IsOccupiedAt(0, 1) {
		t.Errorf("gobject not present on level 2\nGobjectAt returns: %v\nfull map: \n%v",
			l2l.GobjectAt(0, 1), l2l.mapMat)
	}

	l2l.MoveGobjectTo(testGobjCopy, 0, 0)
	if found := l2l.GobjectAt(0, 0); l2l.IsOccupiedAt(0, 1) || found == nil || found.isMarkedForDeletion() {
		t.Errorf("MoveGobjectTo failed\ntarget space is: %v,\nstart space is: %v",
			l2l.GobjectAt(0, 0), l2l.GobjectAt(0, 1))
	}
	l2l.DeleteAt(0, 0)
	if l2l.IsOccupiedAt(0, 0) {
		t.Errorf("deletion failed: %v", l2l.GobjectAt(0, 0))
	}
	gobjectCount := 0
	for range l2l.AllGobjects() {
		gobjectCount++
	}
	if gobjectCount != 0 {
		t.Errorf("wrong number of gobjects from iterator: %v (should be %v)",
			gobjectCount, 0)
	}

	l2l.AddGobject(NewBaseGobject("swap 0 0", EmptySpritePack()).Build(), 0, 0)
	l2l.AddGobject(NewBaseGobject("swap 1 1", EmptySpritePack()).Build(), 1, 1)
	l2l.SwapGobjectsAt(0, 0, 1, 1)
	if l2l.GobjectAt(0, 0).Name() != "swap 1 1" || l2l.GobjectAt(1, 1).Name() != "swap 0 0" {
		t.Errorf("SwapGobjectsAt failed, map:\n%v", l2l.mapMat)
	}
	l2l.SwapGobjectsAt(0, 0, 0, 1)
	if l2l.GobjectAt(0, 0) != nil || l2l.GobjectAt(0, 1).Name() != "swap 1 1" {
		t.Errorf("SwapGobjectsAt failed, map:\n%v", l2l.mapMat)
	}
	gobjectCount = 0
	for range l2l.AllGobjects() {
		gobjectCount++
	}
	if gobjectCount != 2 {
		t.Errorf("wrong number of gobjects from iterator: %v (should be %v)",
			gobjectCount, 2)
	}

	l2l.Clear()
	gobjectCount = 0
	for range l2l.AllGobjects() {
		gobjectCount++
	}
	if gobjectCount != 0 {
		t.Errorf("wrong number of gobjects from iterator after clearing: %v (should be %v)",
			gobjectCount, 0)
	}

	g.NextLevel()
	if g.Level().Name() != "Default" {
		t.Errorf("g.NextLevel didn't wrap around (all levels: %v)", g.Levels)
	}
}
