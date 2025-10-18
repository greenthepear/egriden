package egriden

import (
	"testing"

	"github.com/greenthepear/imggg"
)

// Test OnAddFunc, OnUpdateFunc, OnDeleteFunc in both layer types
func TestGobjectCustomScripts(t *testing.T) {
	type triggerFlags struct {
		onAddTriggered    bool
		onUpdateTriggered bool
		onDeleteTriggered bool
	}
	flags := &triggerFlags{}
	triggerer := NewBaseGobject("Triggerer", EmptySpritePack())

	triggerer.OnAddFunc = func(self Gobject, l Layer) {
		flags.onAddTriggered = true
	}
	triggerer.OnUpdateFunc = func(self Gobject, l Layer) {
		flags.onUpdateTriggered = true
	}
	triggerer.OnDeleteFunc = func(self Gobject, l Layer) {
		flags.onDeleteTriggered = true
	}

	/// Free layer

	freeLayer := NewFreeLayer("Layer 0", imggg.Point[float64]{})
	builtO := triggerer.Build()
	freeLayer.AddGobject(builtO, 0, 0)
	if !flags.onAddTriggered {
		t.Errorf("OnAdd not triggered")
	}
	freeLayer.RunThinkers()
	if !flags.onUpdateTriggered {
		t.Errorf("OnUpdate not triggered")
	}
	freeLayer.DeleteGobject(builtO)
	if !flags.onDeleteTriggered {
		t.Errorf("OnDelete not triggered")
	}

	/// Grid layer

	flags = &triggerFlags{}
	gridLayer := NewGridLayer("Layer 0", GridLayerParameters{
		GridDimensions: Dimensions{1, 1},
		CellDimensions: Dimensions{1, 1},
	})
	built1 := triggerer.Build()
	gridLayer.AddGobject(built1, 0, 0)
	if !flags.onAddTriggered {
		t.Errorf("OnAdd not triggered")
	}
	gridLayer.RunThinkers()
	if !flags.onUpdateTriggered {
		t.Errorf("OnUpdate not triggered")
	}
	gridLayer.DeleteAt(builtO.GridPos().XY())
	if !flags.onDeleteTriggered {
		t.Errorf("OnDelete not triggered")
	}
}

// Test OnUpdate runs with RunThinkers, make sure layer keeps track of thinkers
// correctly after deletions and stuff.
func TestThinkers(t *testing.T) {
	fl := NewFreeLayer("Layer 0", imggg.Point[float64]{})

	type counter struct {
		updates int
	}

	state := &counter{updates: 0}

	thinker := NewBaseGobject("Thinker", EmptySpritePack())
	thinker.OnUpdateFunc = func(self Gobject, l Layer) {
		state.updates++
	}
	nonThinker := NewBaseGobject("Non-thinker", EmptySpritePack())
	for range 5 {
		fl.AddGobject(thinker.Build(), 0, 0)
		fl.AddGobject(thinker.Build(), 1, 1)
		fl.AddGobject(nonThinker.Build(), 2, 2)
		fl.AddGobject(nonThinker.Build(), 3, 3)
	}
	fl.RunThinkers()
	if state.updates != 10 {
		t.Errorf("Wrong number of state updates: %v", state.updates)
	}

	// Remove half of the thinkers and non-thinkers to make sure it doesn't
	// screw anything up.
	for o := range fl.AllGobjects() {
		if (o.Name() == "Thinker" &&
			o.ScreenPos(fl) == imggg.Pt(1.0, 1.0)) ||
			(o.Name() == "Non-thinker" &&
				o.ScreenPos(fl) == imggg.Pt(3.0, 3.0)) {

			fl.DeleteGobject(o)
		}
	}
	state.updates = 0
	fl.RunThinkers()
	if state.updates != 5 {
		t.Errorf("Wrong number of state updates: %v", state.updates)
	}
	fl.Clear()
	if state.updates != 5 {
		t.Errorf("Wrong number of state updates: %v", state.updates)
	}
}
