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
	seq, err := CreateImageSequenceFromFolder(
		"unrev", "./examples/gridsweeper/Graphics/unrevealed")

	if err != nil {
		t.Error(err)
	}

	zero := 0
	updateFuncTestCounter := &zero
	o := NewBaseGobject(
		"tester", NewSpritePackWithSequence(seq))
	o.OnUpdateFunc = func(self Gobject, l Layer) {
		*updateFuncTestCounter += 1
	}
	fl1.AddGobject(
		o.Build(),
		5, 5)
	specificGobject := o.Build()
	fl1.AddGobject(
		specificGobject,
		12, 15)
	gobjectCount := 0

	for range fl1.AllGobjects() {
		gobjectCount++
	}

	if gobjectCount != 2 {
		t.Errorf("not enough gobjects")
	}

	fl1.RunThinkers()
	if *updateFuncTestCounter != 2 {
		t.Errorf(
			"gobject update scripts didn't update counter correctly, value: %d",
			*updateFuncTestCounter)
	}

	fl1.DeleteGobject(specificGobject)

	gobjectCount = 0
	for range fl1.AllGobjects() {
		gobjectCount++
	}
	if gobjectCount != 1 {
		t.Errorf("too many gobjects")
	}

	fl1.RunThinkers()
	if *updateFuncTestCounter != 3 {
		t.Errorf(
			"gobject update scripts didn't update counter correctly, value: %d",
			*updateFuncTestCounter)
	}
	fl1.Clear()
	fl1.RunThinkers()
	if *updateFuncTestCounter != 3 {
		t.Errorf(
			"gobject update scripts didn't update counter correctly, value: %d",
			*updateFuncTestCounter)
	}
}
