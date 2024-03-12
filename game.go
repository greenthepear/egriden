package egriden

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
)

var Warnings bool = true

type EgridenGame struct {
	gridLayers []*GridLayer
}

func NewEgridenGame() *EgridenGame {
	return &EgridenGame{
		gridLayers: make([]*GridLayer, 0),
	}
}

func (g EgridenGame) GridLayer(z int) (*GridLayer, error) {
	if z >= len(g.gridLayers) {
		return nil, fmt.Errorf("no grid layer %d (number of layers %d)", z, len(g.gridLayers))
	}
	return g.gridLayers[z], nil
}

func (g EgridenGame) GridLayers() []*GridLayer {
	return g.gridLayers
}

func (g EgridenGame) DrawAllLayers(screen *ebiten.Image) {
	for _, l := range g.gridLayers {
		l.draw(screen)
	}
}
