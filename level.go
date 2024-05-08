package egriden

import (
	"fmt"
	"slices"

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

// Returns a GridLayer at z, panics if out of bounds
func (le *BaseLevel) GridLayer(z int) *GridLayer {
	if z >= len(le.gridLayers) {
		panic(fmt.Sprintf("no grid layer %d (number of layers %d)", z, len(le.gridLayers)))
	}
	return le.gridLayers[z]
}

func (le *BaseLevel) GridLayers() []*GridLayer {
	return le.gridLayers
}

func (le *BaseLevel) DrawAllGridLayers(screen *ebiten.Image) {
	for _, l := range le.gridLayers {
		l.Draw(screen)
	}
}

func (le *BaseLevel) DrawAllFreeLayers(screen *ebiten.Image) {
	for _, l := range le.freeLayers {
		l.Draw(screen)
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
		le.gobjectsWithUpdateScripts = slices.DeleteFunc( //Slow, we're ranging through the slice anyway
			le.gobjectsWithUpdateScripts,
			func(o Gobject) bool {
				return o.isMarkedForDeletion()
			})
	}
}
