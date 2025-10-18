package egriden

import (
	"iter"

	"github.com/greenthepear/imggg"
	"github.com/hajimehoshi/ebiten/v2"
)

// Levels are essentially different collections of layers.
// They are often called scenes in game dev.
type Level interface {
	Name() string
	Index() int

	GridLayer(int) *GridLayer

	AddGridLayer(l *GridLayer) int
	DeleteGridLayerAt(z int)
	ReplaceGridLayerAt(l *GridLayer, z int)

	AddFreeLayer(l *FreeLayer) int
	DeleteFreeLayerAt(z int)
	ReplaceFreeLayerAt(l *FreeLayer, z int)

	CreateGridLayerOnTop(name string, params GridLayerParameters) *GridLayer
	CreateSimpleGridLayerOnTop(
		name string, squareLength int, width, height int,
		drawMode DrawMode, XOffset, YOffset float64) *GridLayer

	CreateAndReplaceGridLayerAt(
		z int, name string, param GridLayerParameters) *GridLayer

	FreeLayer(int) *FreeLayer

	CreateFreeLayerOnTop(name string, anchor imggg.Point[float64]) *FreeLayer

	DrawAllGridLayers(*ebiten.Image)
	DrawAllFreeLayers(*ebiten.Image)

	AllGridLayers() iter.Seq2[int, *GridLayer]
	AllFreeLayers() iter.Seq2[int, *FreeLayer]
	AllLayers(gridLayersFirst bool) iter.Seq[Layer]

	setIndex(int)
}

type BaseLevel struct {
	name  string
	index int

	gridLayers []*GridLayer
	freeLayers []*FreeLayer
}

func (le *BaseLevel) setIndex(i int) {
	le.index = i
}

func NewBaseLevel(name string) *BaseLevel {
	le := &BaseLevel{name: name}
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

func (le BaseLevel) AllGridLayers() iter.Seq2[int, *GridLayer] {
	return func(yield func(int, *GridLayer) bool) {
		for z, l := range le.gridLayers {
			if !yield(z, l) {
				return
			}
		}
	}
}

func (le BaseLevel) AllFreeLayers() iter.Seq2[int, *FreeLayer] {
	return func(yield func(int, *FreeLayer) bool) {
		for z, l := range le.freeLayers {
			if !yield(z, l) {
				return
			}
		}
	}
}

func (le BaseLevel) AllLayers(gridLayersFirst bool) iter.Seq[Layer] {
	return func(yield func(Layer) bool) {
		if gridLayersFirst {
			for _, gl := range le.AllGridLayers() {
				if !yield(gl) {
					return
				}
			}
			for _, fl := range le.AllFreeLayers() {
				if !yield(fl) {
					return
				}
			}
		} else {
			for _, fl := range le.AllFreeLayers() {
				if !yield(fl) {
					return
				}
			}
			for _, gl := range le.AllGridLayers() {
				if !yield(gl) {
					return
				}
			}
		}
	}
}

// Draws all grid layers according to their Z order
func (le BaseLevel) DrawAllGridLayers(on *ebiten.Image) {
	for _, l := range le.AllGridLayers() {
		l.Draw(on)
	}
}

// Draws all free layers according to their Z order
func (le BaseLevel) DrawAllFreeLayers(on *ebiten.Image) {
	for _, l := range le.AllFreeLayers() {
		l.Draw(on)
	}
}
