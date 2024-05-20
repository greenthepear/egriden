package egriden

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// A free layer is a layer where the default drawing position of Gobjects is only
// determined by their XY coordinates and can be anywhere on the screen or outside of it.
type FreeLayer struct {
	Name string
	z    int

	gobjects gobjectSet

	Visible          bool
	static           bool
	staticImage      *ebiten.Image
	XOffset, YOffset float64
}

// Options for a static free layer
type StaticFreeLayerOp struct {
	width, height int
}

func newFreeLayer(name string, z int, visible bool, staticOptions *StaticFreeLayerOp, xOffset, yOffset float64) *FreeLayer {
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
		gobjects:    newGobjectSet(),
		staticImage: img,
		XOffset:     xOffset,
		YOffset:     yOffset,
	}
}

// Creates a new FreeLayer and returns a pointer to it.
func (le *BaseLevel) CreateFreeLayerOnTop(name string, xOffset, yOffset float64) *FreeLayer {
	z := len(le.freeLayers)
	le.freeLayers = append(le.freeLayers, newFreeLayer(name, z, true, nil, xOffset, yOffset))
	return le.freeLayers[z]
}

// Creates a free layer whose image needs to to be refreshed to be updated.
// Remember to call Refresh() at least once after populating,
// otherwise you'll just get an empty image.
func (le *BaseLevel) CreateStaticFreeLayerOnTop(
	name string, imgWidth, imgHeight int, xOffset, yOffset float64) *FreeLayer {
	z := len(le.freeLayers)
	le.freeLayers = append(le.freeLayers, newFreeLayer(name, z, true,
		&StaticFreeLayerOp{imgWidth, imgHeight}, xOffset, yOffset))
	return le.freeLayers[z]
}

// Retruns FreeLayer at given z layer, returns nil if out of bounds
func (le *BaseLevel) FreeLayer(z int) *FreeLayer {
	if z >= len(le.freeLayers) || z < 0 {
		return nil
	}
	return le.freeLayers[z]
}

func (le *BaseLevel) FreeLayers() []*FreeLayer {
	return le.freeLayers
}

func (le *FreeLayer) SetVisibility(to bool) {
	le.Visible = to
}

func (fl *FreeLayer) AddGobject(o Gobject, x, y int) {
	o.setXY(x, y)
	fl.gobjects.Add(o)
}

func (fl *FreeLayer) MoveGobjectTo(o Gobject, x, y int) {
	_, ok := fl.gobjects.m[o]
	if !ok {
		panic("Gobject does not exist in layer")
	}
	o.setXY(x, y)
}

func (fl *FreeLayer) DeleteGobject(o Gobject) {
	_, ok := fl.gobjects.m[o]
	if !ok {
		panic("Gobject does not exist in layer")
	}
	o.markForDeletion()
	fl.gobjects.Delete(o)
}

func (fl *FreeLayer) Z() int {
	return fl.z
}
