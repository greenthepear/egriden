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

	FreeLayer(int) *FreeLayer
	AddFreeLayer(l *FreeLayer) int
	DeleteFreeLayerAt(z int)
	ReplaceFreeLayerAt(l *FreeLayer, z int)

	DrawAllGridLayers(*ebiten.Image)
	DrawAllFreeLayers(*ebiten.Image)

	AllGridLayers() iter.Seq2[int, *GridLayer]
	AllFreeLayers() iter.Seq2[int, *FreeLayer]
	AllLayers(gridLayersFirst bool) iter.Seq[Layer]

	setIndex(int)

	CreateGridLayerOnTop(name string, params GridLayerParameters) *GridLayer
	CreateSimpleGridLayerOnTop(
		name string, squareLength int, width, height int,
		drawMode DrawMode, XOffset, YOffset float64) *GridLayer

	CreateAndReplaceGridLayerAt(
		z int, name string, param GridLayerParameters) *GridLayer

	CreateFreeLayerOnTop(name string, anchor imggg.Point[float64]) *FreeLayer
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
	return &BaseLevel{name: name}
}

func (le BaseLevel) Name() string {
	return le.name
}

func (le BaseLevel) Index() int {
	return le.index
}

/// Get layers

// Returns a GridLayer at z, returns nil if out of bounds
func (le BaseLevel) GridLayer(z int) *GridLayer {
	if z >= len(le.gridLayers) || z < 0 {
		return nil
	}
	return le.gridLayers[z]
}

// FreeLayer at given z layer, returns nil if out of bounds
func (le *BaseLevel) FreeLayer(z int) *FreeLayer {
	if z >= len(le.freeLayers) || z < 0 {
		return nil
	}
	return le.freeLayers[z]
}

/// Iterators

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

/// Draw layers

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

/// GridLayer management

func (le *BaseLevel) addGridLayer(l *GridLayer) *GridLayer {
	ln := len(le.gridLayers)
	l.z = ln
	le.gridLayers = append(le.gridLayers, l)
	return le.gridLayers[ln]
}

// Add layer on top, return it's Z.
func (le *BaseLevel) AddGridLayer(l *GridLayer) int {
	le.addGridLayer(l)
	return l.z
}

// Replaces a layer at specified Z. Panics if out of bounds.
func (le *BaseLevel) ReplaceGridLayerAt(l *GridLayer, z int) {
	if z >= len(le.gridLayers) {
		panic("z out of bounds")
	}
	l.z = z
	le.gridLayers[z] = l
}

// Deletes a layer at specified Z, updates Z of the rest of the layers.
// Panics if out of bounds.
func (le *BaseLevel) DeleteGridLayerAt(z int) {
	if z >= len(le.gridLayers) {
		panic("z out of bounds")
	}
	le.gridLayers = append(le.gridLayers[:z], le.gridLayers[z+1:]...)
	for zi, l := range le.gridLayers {
		l.z = zi
	}
}

// Creates a grid layer with custom parameters within the level and returns the
// pointer to it.
func (le *BaseLevel) CreateGridLayerOnTop(
	name string, params GridLayerParameters) *GridLayer {

	return le.addGridLayer(NewGridLayer(name, params))
}

/// FreeLayer management

// Add layer on top, return it's Z.
func (le *BaseLevel) AddFreeLayer(l *FreeLayer) int {
	ln := len(le.freeLayers)
	l.z = ln
	le.freeLayers = append(le.freeLayers, l)
	return ln
}

// Replaces a layer at specified Z. Panics if out of bounds.
func (le *BaseLevel) ReplaceFreeLayerAt(l *FreeLayer, z int) {
	if z >= len(le.freeLayers) {
		panic("z out of bounds")
	}
	l.z = z
	le.freeLayers[z] = l
}

// Deletes a layer at specified Z, updates Z of the rest of the layers.
// Panics if out of bounds.
func (le *BaseLevel) DeleteFreeLayerAt(z int) {
	if z >= len(le.freeLayers) {
		panic("z out of bounds")
	}
	le.freeLayers = append(le.freeLayers[:z], le.freeLayers[z+1:]...)
	for zi, l := range le.freeLayers {
		l.z = zi
	}
}

// Creates a new FreeLayer and returns a pointer to it.
func (le *BaseLevel) CreateFreeLayerOnTop(
	name string, anchor imggg.Point[float64]) *FreeLayer {

	z := len(le.freeLayers)
	newLayer := newFreeLayer(name, z, true, anchor)
	le.freeLayers = append(le.freeLayers, newLayer)
	return le.freeLayers[z]
}

// Clears and creates a new gridlayer at specified index.
//
// Probably temporary. TODO: Remove
func (le *BaseLevel) CreateAndReplaceGridLayerAt(
	z int, name string, params GridLayerParameters) *GridLayer {

	old := le.GridLayer(z)
	if old == nil {
		return nil
	}
	old.Clear()

	le.gridLayers[z] = newGridLayer(
		name, z,
		params.CellDimensions,
		params.GridDimensions,
		params.Mode,
		params.Anchor,
		params.PaddingVector,
	)
	return le.gridLayers[z]
}

// Creates a grid layer with square cells and no padding within the level.
// Also returns the pointer to it.
//
// Deprecated: Just use CreateGridLayerOnTop.
func (le *BaseLevel) CreateSimpleGridLayerOnTop(
	name string, squareLength int, width, height int,
	drawMode DrawMode, XOffset, YOffset float64) *GridLayer {

	return le.addGridLayer(
		newGridLayer(
			name, 0,
			Dimensions{squareLength, squareLength},
			Dimensions{width, height},
			drawMode,
			imggg.Pt(XOffset, YOffset),
			imggg.Pt(0.0, 0.0),
		),
	)
}
