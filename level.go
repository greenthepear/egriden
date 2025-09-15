package egriden

import (
	"github.com/greenthepear/gunc"
	"github.com/hajimehoshi/ebiten/v2"
)

// Levels are essentially different collections of layers.
// They are often called scenes in game dev.
type Level interface {
	Name() string
	Index() int

	GridLayer(int) *GridLayer

	CreateGridLayerOnTop(name string, params GridLayerParameters) *GridLayer
	CreateSimpleGridLayerOnTop(
		name string, squareLength int, width, height int,
		drawMode DrawMode, XOffset, YOffset float64) *GridLayer

	ReplaceGridLayerAt(z int, name string, param GridLayerParameters) *GridLayer

	FreeLayer(int) *FreeLayer

	CreateFreeLayerOnTop(name string, xOffset, yOffset float64) *FreeLayer

	addGobjectWithOnUpdate(o Gobject, l Layer)
}

type gobjectWithLayer struct {
	o Gobject
	l Layer
}

type BaseLevel struct {
	name  string
	index int

	gridLayers                []*GridLayer
	gobjectsWithUpdateScripts []gobjectWithLayer

	freeLayers []*FreeLayer
}

// Initialize by creating slices for layers
func (le *BaseLevel) Init() {
	le.gridLayers = make([]*GridLayer, 0)
	le.gobjectsWithUpdateScripts = make([]gobjectWithLayer, 0)
}

func NewBaseLevel(name string) *BaseLevel {
	le := &BaseLevel{name: name}
	le.Init()
	return le
}

func (le BaseLevel) Name() string {
	return le.name
}

func (le BaseLevel) Index() int {
	return le.index
}

// Returns a GridLayer at z, returns nil if out of bounds
func (le BaseLevel) GridLayer(z int) *GridLayer {
	if z >= len(le.gridLayers) || z < 0 {
		return nil
	}
	return le.gridLayers[z]
}

// Draws all grid layers according to their Z order
func (le BaseLevel) DrawAllGridLayers(on *ebiten.Image) {
	for _, l := range le.gridLayers {
		l.Draw(on)
	}
}

// Draws all free layers according to their Z order
func (le BaseLevel) DrawAllFreeLayers(on *ebiten.Image) {
	for _, l := range le.freeLayers {
		l.Draw(on)
	}
}

func (le *BaseLevel) addGobjectWithOnUpdate(o Gobject, l Layer) {
	le.gobjectsWithUpdateScripts =
		append(le.gobjectsWithUpdateScripts, gobjectWithLayer{o, l})
}

// UNTESTED! Run all the onUpdate() functions of Gobjects that have them
func (le *BaseLevel) RunUpdateScripts() {
	marked := 0
	for _, elem := range le.gobjectsWithUpdateScripts {
		if elem.o.isMarkedForDeletion() {
			marked++
			continue
		}
		elem.o.OnUpdate()(elem.o, elem.l)
	}

	if marked > 0 {
		le.gobjectsWithUpdateScripts = gunc.Filter(le.gobjectsWithUpdateScripts,
			func(ol gobjectWithLayer) bool {
				return !ol.o.isMarkedForDeletion()
			})

	}
}
