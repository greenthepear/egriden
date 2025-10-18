package egriden

import (
	"testing"

	"github.com/greenthepear/imggg"
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
	if found := l2l.GobjectAt(0, 0); l2l.IsOccupiedAt(0, 1) || found == nil {
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

	l2lReplaced := g.Level().
		CreateAndReplaceGridLayerAt(0, "Replaced layer", GridLayerParameters{})

	if l2lReplaced == l2l || l2lReplaced.Name != "Replaced layer" {
		t.Errorf("layer not replaced properly")
	}

	g.NextLevel()
	if g.Level().Name() != "Default" {
		t.Errorf("g.NextLevel didn't wrap around (all levels: %v)", g.Levels)
	}
}

// Test 0.4.0 methods
func TestLevelsLayerNew(t *testing.T) {
	g := &EgridenAssets{}
	if g.Level() != nil {
		t.Errorf("non-nil where there should be no level")
	}

	level0 := NewBaseLevel("Level 0")
	level1 := NewBaseLevel("Level 1")

	g.AddLevel(level0)
	g.AddLevel(level1)
	if level0.Index() != 0 || level0.Name() != "Level 0" ||
		level1.Index() != 1 || level1.Name() != "Level 1" {

		t.Errorf("Levels not properly initialized:\n%v\n%v", level0, level1)
	}

	if g.Level() != level0 {
		t.Errorf("Wrong level returned: %v", g.Level())
	}

	l0 := NewGridLayer("Layer 0", GridLayerParameters{
		CellDimensions: Dimensions{10, 10},
		GridDimensions: Dimensions{3, 4},
		Mode:           Sparse,
	})
	l1 := NewGridLayer("Layer 1", GridLayerParameters{
		CellDimensions: Dimensions{10, 10},
		GridDimensions: Dimensions{4, 5},
		Mode:           Sparse,
	})
	l0z := level0.AddGridLayer(l0)
	level0.AddGridLayer(l1)
	if l0z != 0 {
		t.Errorf("Wrong z returned: %v", l0z)
	}
	if level0.GridLayer(l0z) != l0 {
		t.Errorf("Layers should be the same: %v != %v",
			l0, level0.GridLayer(l0z))
	}
	level0.DeleteGridLayerAt(l0z)
	if level0.GridLayer(l0z) == l0 ||
		level0.GridLayer(l0z) != l1 ||
		l1.Z() != l0z {

		t.Errorf("Wrong layer after deletion: %v", level0.GridLayer(l0z))
	}

	freeLayers := []*FreeLayer{
		NewFreeLayer("Free layer 0", imggg.Point[float64]{}),
		NewFreeLayer("Free layer 1", imggg.Point[float64]{}),
		NewFreeLayer("Free layer 2", imggg.Point[float64]{}),
		NewFreeLayer("Free layer 3", imggg.Point[float64]{}),
	}
	rfl0 := NewFreeLayer("Free layer replacement", imggg.Point[float64]{})
	for _, fl := range freeLayers {
		level1.AddFreeLayer(fl)
	}
	if len(level1.freeLayers) != 4 {
		t.Errorf("Wrong number of layers: %v", len(level1.freeLayers))
	}
	level1.ReplaceFreeLayerAt(rfl0, 1)
	if level1.FreeLayer(1) != rfl0 ||
		level1.FreeLayer(1).Name != "Free layer replacement" {

		t.Errorf("Free layer not replaced correctly: %v != %v",
			level1.FreeLayer(1), rfl0)
	}
	if len(level1.freeLayers) != 4 {
		t.Errorf("Wrong number of layers: %v", len(level1.freeLayers))
	}
	level1.DeleteFreeLayerAt(1)
	if len(level1.freeLayers) != 3 {
		t.Errorf("Wrong number of layers: %v", len(level1.freeLayers))
	}
	for z, l := range level1.AllFreeLayers() {
		if l.Z() != z {
			t.Errorf("Wrong z for layer: %v %v", z, l)
		}
	}
	if level1.FreeLayer(2).Name != "Free layer 3" {
		t.Errorf("Wrong layer name after deletion: %v", level1.FreeLayer(3))
	}
}
