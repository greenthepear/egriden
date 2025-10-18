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

func (lf *FreeLayer) CurrentAnchor() imggg.Point[float64] {
	return lf.Anchor
}

func (lf *FreeLayer) SetVisibility(to bool) {
	lf.Visible = to
}

func (fl *FreeLayer) AddGobject(o Gobject, x, y float64) {
	o.setScreenPos(x, y)
	if o.OnUpdate() != nil {
		o.setThinkerElement(fl.thinkers.PushBack(o))
	}
	o.setGobjectElement(fl.gobjects.PushBack(o))

	if o.OnAdd() != nil {
		o.OnAdd()(o, fl)
	}
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
	if o.OnDelete() != nil {
		o.OnDelete()(o, fl)
	}
}

func (fl *FreeLayer) Z() int {
	return fl.z
}

func (fl FreeLayer) gobjectRange() iter.Seq[Gobject] {
	return func(yield func(Gobject) bool) {
		var next *list.Element
		for e := fl.gobjects.Front(); e != nil; e = next {
			next = e.Next() // In case of deletion during iteration
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

// Iterates over all Gobjects in the layer. Deleting gobjects during iteration
// might cause headaches.
func (fl FreeLayer) AllGobjects() iter.Seq[Gobject] {
	return func(yield func(Gobject) bool) {
		for o := range fl.gobjectRange() {
			if !yield(o) {
				return
			}
		}
	}
}

// Delete all gobjects. This does not trigger OnDelete.
func (fl *FreeLayer) Clear() {
	fl.thinkers.Init()
	fl.gobjects.Init()
}

func (fl *FreeLayer) RunThinkers() {
	var next *list.Element
	for e := fl.thinkers.Front(); e != nil; e = next {
		next = e.Next() // In case of deletion during iteration
		o, ok := e.Value.(Gobject)
		if !ok {
			panic("non-gobject in thinker list")
		}
		o.OnUpdate()(o, fl)
	}
}
