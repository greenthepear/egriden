package egriden

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
)

var Warnings bool = true

type EgridenAssets struct {
	gridLayers []*GridLayer
}

func (g *EgridenAssets) InitEgridenComponents() {
	g.gridLayers = make([]*GridLayer, 0)
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
