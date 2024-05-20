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
	GridLayers() []*GridLayer

	FreeLayer(int) *FreeLayer
	FreeLayers() []*FreeLayer
}

type BaseLevel struct {
	name  string
	index int

	gridLayers                []*GridLayer
	gobjectsWithUpdateScripts []Gobject

	freeLayers []*FreeLayer
}

// Initialize by creating slices for layers
func (le *BaseLevel) Init() {
	le.gridLayers = make([]*GridLayer, 0)
	le.gobjectsWithUpdateScripts = make([]Gobject, 0)
}

func NewBaseLevel(name string) *BaseLevel {
	le := &BaseLevel{name: name}
	le.Init()
	return le
}

func (le *BaseLevel) Name() string {
	return le.name
}

func (le *BaseLevel) Index() int {
	return le.index
}

// Returns a GridLayer at z, returns nil if out of bounds
func (le *BaseLevel) GridLayer(z int) *GridLayer {
	if z >= len(le.gridLayers) || z < 0 {
		return nil
	}
	return le.gridLayers[z]
}

// Returns slice of all current grid layers
func (le *BaseLevel) GridLayers() []*GridLayer {
	return le.gridLayers
}

// Draws all grid layers according to their Z order
func (le *BaseLevel) DrawAllGridLayers(on *ebiten.Image) {
	for _, l := range le.gridLayers {
		l.Draw(on)
	}
}

// Draws all free layers according to their Z order
func (le *BaseLevel) DrawAllFreeLayers(on *ebiten.Image) {
	for _, l := range le.freeLayers {
		l.Draw(on)
	}
}

// Run all the onUpdate() functions of Gobjects that have them
func (le *BaseLevel) RunUpdateScripts() {
	marked := 0
	for _, o := range le.gobjectsWithUpdateScripts {
		if o.isMarkedForDeletion() {
			marked++
			continue
		}
		o.OnUpdate()()
	}

	if marked > 0 {
		le.gobjectsWithUpdateScripts = gunc.Filter(le.gobjectsWithUpdateScripts,
			func(o Gobject) bool {
				return o.isMarkedForDeletion()
			})

	}
}
