package egriden

import (
	"fmt"
	"slices"

	"github.com/hajimehoshi/ebiten/v2"
)

type EgridenAssets struct {
	gridLayers                []*GridLayer
	gobjectsWithUpdateScripts []Gobject
}

// Run this while initalizing the game, before adding any layers
func (g *EgridenAssets) InitEgridenComponents() {
	g.gridLayers = make([]*GridLayer, 0)
	g.gobjectsWithUpdateScripts = make([]Gobject, 0)
}

// Returns a GridLayer at z, panics if out of bounds
func (g EgridenAssets) GridLayer(z int) *GridLayer {
	if z >= len(g.gridLayers) {
		panic(fmt.Sprintf("no grid layer %d (number of layers %d)", z, len(g.gridLayers)))
	}
	return g.gridLayers[z]
}

func (g EgridenAssets) GridLayers() []*GridLayer {
	return g.gridLayers
}

// Draw all GridLayers in their Z order. Use this in the Draw() function.
func (g EgridenAssets) DrawAllLayers(screen *ebiten.Image) {
	for _, l := range g.gridLayers {
		l.Draw(screen)
	}
}

// Run all the OnUpdate() scripts of Gobjects that have them.
// Use this in the Update() function.
func (g *EgridenAssets) RunUpdateScripts() {
	marked := 0
	for _, o := range g.gobjectsWithUpdateScripts {
		if o.isMarkedForDeletion() {
			marked++
			continue
		}
		o.OnUpdate()()
	}

	if marked > 0 {
		g.gobjectsWithUpdateScripts = slices.DeleteFunc( //Slow, we're ranging through the slice anyway
			g.gobjectsWithUpdateScripts,
			func(o Gobject) bool {
				return o.isMarkedForDeletion()
			})
	}
}
