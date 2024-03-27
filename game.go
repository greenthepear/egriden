package egriden

import (
	"fmt"
	"slices"

	"github.com/hajimehoshi/ebiten/v2"
)

var Warnings bool = true

type EgridenAssets struct {
	gridLayers                []*GridLayer
	gobjectsWithUpdateScripts []Gobject
}

func (g *EgridenAssets) InitEgridenComponents() {
	g.gridLayers = make([]*GridLayer, 0)
	g.gobjectsWithUpdateScripts = make([]Gobject, 0)
}

func (g EgridenAssets) GridLayer(z int) *GridLayer {
	if z >= len(g.gridLayers) {
		panic(fmt.Sprintf("no grid layer %d (number of layers %d)", z, len(g.gridLayers)))
	}
	return g.gridLayers[z]
}

func (g EgridenAssets) GridLayers() []*GridLayer {
	return g.gridLayers
}

func (g EgridenAssets) DrawAllLayers(screen *ebiten.Image) {
	for _, l := range g.gridLayers {
		l.Draw(screen)
	}
}

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
