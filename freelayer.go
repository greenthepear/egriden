package egriden

import (
	"iter"

	"container/list"

	"github.com/greenthepear/imggg"
)

// A free layer is a layer where the default drawing position of Gobjects is only
// determined by their XY coordinates and can be anywhere on the screen or outside of it.
type FreeLayer struct {
	Name string
	z    int

	gobjects *list.List

	Visible bool
	Anchor  imggg.Point[float64]

	thinkers *list.List
}

func newFreeLayer(
	name string, z int, visible bool,
	anchor imggg.Point[float64]) *FreeLayer {

	return &FreeLayer{
		Name:     name,
		z:        z,
		Visible:  visible,
		gobjects: list.New(),
		Anchor:   anchor,
		thinkers: list.New(),
	}
}

// Creates a new FreeLayer. Before adding to a level, it will have Z of 0.
func NewFreeLayer(name string, anchor imggg.Point[float64]) *FreeLayer {
	return newFreeLayer(name, 0, true, anchor)
}

// Add layer on top, return it's Z.
func (le *BaseLevel) AddFreeLayer(l *FreeLayer) int {
	ln := len(le.freeLayers)
	l.z = ln
	le.freeLayers = append(le.freeLayers, l)
	return ln
}

// Replaces a layer at specified Z. Panics if out of bounds.
func (le *BaseLevel) ReplaceFreeLayerAt(l *FreeLayer, z int) {
	if z >= len(le.freeLayers) {
		panic("z out of bounds")
	}
	l.z = z
	le.freeLayers[z] = l
}

// Deletes a layer at specified Z, updates Z of the rest of the layers.
// Panics if out of bounds.
func (le *BaseLevel) DeleteFreeLayerAt(z int) {
	if z >= len(le.freeLayers) {
		panic("z out of bounds")
	}
	le.freeLayers = append(le.freeLayers[:z], le.freeLayers[z+1:]...)
	for zi, l := range le.freeLayers {
		l.z = zi
	}
}

// Creates a new FreeLayer and returns a pointer to it.
func (le *BaseLevel) CreateFreeLayerOnTop(
	name string, anchor imggg.Point[float64]) *FreeLayer {

	z := len(le.freeLayers)
	newLayer := newFreeLayer(name, z, true, anchor)
	le.freeLayers = append(le.freeLayers, newLayer)
	return le.freeLayers[z]
}

// FreeLayer at given z layer, returns nil if out of bounds
func (le *BaseLevel) FreeLayer(z int) *FreeLayer {
	if z >= len(le.freeLayers) || z < 0 {
		return nil
	}
	return le.freeLayers[z]
}

func (le *FreeLayer) CurrentAnchor() imggg.Point[float64] {
	return le.Anchor
}

func (le *FreeLayer) SetVisibility(to bool) {
	le.Visible = to
}

func (fl *FreeLayer) AddGobject(o Gobject, x, y float64) {
	o.setScreenPos(x, y)
	if o.OnUpdate() != nil {
		o.setThinkerElement(fl.thinkers.PushBack(o))
	}
	o.setGobjectElement(fl.gobjects.PushBack(o))
}

func (fl *FreeLayer) MoveGobjectTo(o Gobject, x, y float64) {
	o.setScreenPos(x, y)
}

func (fl *FreeLayer) DeleteGobject(o Gobject) {
	if o.thinkerElement() != nil {
		fl.thinkers.Remove(o.thinkerElement())
		o.setThinkerElement(nil)
	}
	fl.gobjects.Remove(o.gobjectElement())
	o.setThinkerElement(nil)
}

func (fl *FreeLayer) Z() int {
	return fl.z
}

func (fl FreeLayer) gobjectRange() iter.Seq[Gobject] {
	return func(yield func(Gobject) bool) {
		for e := fl.gobjects.Front(); e != nil; e = e.Next() {
			o, ok := e.Value.(Gobject)
			if !ok {
				panic("list element isn't a Gobject")
			}
			if !yield(o) {
				return
			}
		}
	}
}

func (fl FreeLayer) AllGobjects() iter.Seq[Gobject] {
	return func(yield func(Gobject) bool) {
		for o := range fl.gobjectRange() {
			if !yield(o) {
				return
			}
		}
	}
}

// Delete all gobjects
func (fl *FreeLayer) Clear() {
	fl.thinkers.Init()
	fl.gobjects.Init()
}

func (fl *FreeLayer) RunThinkers() {
	for e := fl.thinkers.Front(); e != nil; e = e.Next() {
		o, ok := e.Value.(Gobject)
		if !ok {
			panic("non-gobject in thinker list")
		}
		o.OnUpdate()(o, fl)
	}
}

// Shortcut for g.Level().CreateFreeLayerOnTop().
// Level implementation must have BaseLevel component.
//
// Deprecated: Use method directly from Level
func (g *EgridenAssets) CreateFreeLayerOnTop(
	name string, xOffset, yOffset float64) *FreeLayer {

	return g.Level().CreateFreeLayerOnTop(name,
		imggg.Point[float64]{X: xOffset, Y: yOffset})
}
