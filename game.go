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

func (g EgridenAssets) GridLayer(z int) (*GridLayer, error) {
	if z >= len(g.gridLayers) {
		return nil, fmt.Errorf("no grid layer %d (number of layers %d)", z, len(g.gridLayers))
	}
	return g.gridLayers[z], nil
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
	gobjects := slices.DeleteFunc( //Slow, we're ranging through the slice anyway
		g.gobjectsWithUpdateScripts,
		func(o Gobject) bool {
			return o.isMarkedForDeletion()
		})

	for _, o := range gobjects {
		o.OnUpdate()()
	}
}
