package egriden

import (
	"iter"

	"container/list"

	"github.com/greenthepear/imggg"
	"github.com/hajimehoshi/ebiten/v2"
)

// A free layer is a layer where the default drawing position of Gobjects is only
// determined by their XY coordinates and can be anywhere on the screen or outside of it.
type FreeLayer struct {
	Name string
	z    int

	gobjects *list.List

	Visible     bool
	static      bool
	staticImage *ebiten.Image
	Anchor      imggg.Point[float64]

	level Level
}

// Options for a static free layer
type staticFreeLayerOp struct {
	width, height int
}

func newFreeLayer(
	name string, z int, visible bool, staticOptions *staticFreeLayerOp,
	xOffset, yOffset float64) *FreeLayer {

	paramStatic := false
	var img *ebiten.Image
	if staticOptions != nil {
		paramStatic = true
		img = ebiten.NewImage(staticOptions.width, staticOptions.height)
	}

	return &FreeLayer{
		Name:        name,
		z:           z,
		Visible:     visible,
		static:      paramStatic,
		gobjects:    list.New(),
		staticImage: img,
		Anchor:      imggg.Pt(xOffset, yOffset),
	}
}

// Creates a new FreeLayer and returns a pointer to it.
func (le *BaseLevel) CreateFreeLayerOnTop(
	name string, xOffset, yOffset float64) *FreeLayer {

	z := len(le.freeLayers)
	newLayer := newFreeLayer(name, z, true, nil, xOffset, yOffset)
	le.freeLayers = append(le.freeLayers, newLayer)
	newLayer.level = le
	return le.freeLayers[z]
}

// Creates a free layer whose image needs to to be refreshed to be updated.
// Remember to call Refresh() at least once after populating,
// otherwise you'll just get an empty image.
func (le *BaseLevel) CreateStaticFreeLayerOnTop(
	name string, imgWidth, imgHeight int, xOffset, yOffset float64) *FreeLayer {
	z := len(le.freeLayers)
	le.freeLayers = append(le.freeLayers, newFreeLayer(name, z, true,
		&staticFreeLayerOp{imgWidth, imgHeight}, xOffset, yOffset))
	return le.freeLayers[z]
}

// FreeLayer at given z layer, returns nil if out of bounds
func (le *BaseLevel) FreeLayer(z int) *FreeLayer {
	if z >= len(le.freeLayers) || z < 0 {
		return nil
	}
	return le.freeLayers[z]
}

func (le *FreeLayer) anchor() imggg.Point[float64] {
	return le.Anchor
}

func (le *FreeLayer) SetVisibility(to bool) {
	le.Visible = to
}

func (fl *FreeLayer) AddGobject(o Gobject, x, y float64) {
	o.setScreenPos(x, y)
	if o.OnUpdate() != nil {
		fl.level.addGobjectWithOnUpdate(o, fl)
	}
	e := fl.gobjects.PushBack(o)
	o.setListElement(e)
}

func (fl *FreeLayer) MoveGobjectTo(o Gobject, x, y float64) {
	o.setScreenPos(x, y)
}

func (fl *FreeLayer) DeleteGobject(o Gobject) {
	o.setMarkForDeletion(true)
	fl.gobjects.Remove(o.listElement())
}

func (fl *FreeLayer) Z() int {
	return fl.z
}

func (fl *FreeLayer) Static() bool {
	return fl.static
}

func (fl FreeLayer) AllGobjects() iter.Seq[Gobject] {
	return func(yield func(Gobject) bool) {
		for e := fl.gobjects.Front(); e != nil; e = e.Next() {
			o, ok := e.Value.(Gobject)
			if !ok {
				panic("list element isn't a Gobject")
			}
			if o != nil && !o.isMarkedForDeletion() {
				if !yield(o) {
					return
				}
			}
		}
	}
}

// Delete all gobjects
func (fl *FreeLayer) Clear() {
	for o := range fl.AllGobjects() {
		fl.DeleteGobject(o)
	}
}

// Shortcut for g.Level().CreateFreeLayerOnTop().
// Level implementation must have BaseLevel component.
func (g *EgridenAssets) CreateFreeLayerOnTop(
	name string, xOffset, yOffset float64) *FreeLayer {

	return g.Level().CreateFreeLayerOnTop(name, xOffset, yOffset)
}

// Shortcut for g.Level().CreateStaticFreeLayerOnTop().
// Level implementation must have BaseLevel component.
func (g *EgridenAssets) CreateStaticFreeLayerOnTop(
	name string, imgWidth, imgHeight int, xOffset, yOffset float64) *FreeLayer {

	return g.Level().CreateStaticFreeLayerOnTop(
		name, imgWidth, imgHeight, xOffset, yOffset)
}
