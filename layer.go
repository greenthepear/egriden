package egriden

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

// For optimization there are a couple of ways a grid layer can be draw depending if it
// changes frequently (Static or not) and if it has many (Dense) or few (Sparce)
// gobjects most of the time.
type DrawMode int

const (
	//Used for sparcely populated grids, ranges over a map for drawing
	Sparce DrawMode = iota
	//Used for thickly populated grids, ranges over a slice for drawing
	Dense
	//Used for layers that don't get updated often, creates ebiten.Image of the the entire layer.
	//Can be refreshed with GridLayer.RefreshImage().
	Static
)

// TODO: replace with image.Point maybe?
type vec struct {
	x, y int
}

type Dimensions struct {
	Width, Height int
}

// Width and height as ints
func (d Dimensions) WH() (int, int) {
	return d.Width, d.Height
}

type GridLayer struct {
	Name            string // Name of the layer, for convenience sake
	z               int
	cellDimensions  Dimensions
	layerDimensions Dimensions

	// Defines the "gaps" between cells:
	// point's X for horizontal gaps length and Y for vertical.
	Padding image.Point
	// If false no sprite will be drawn, nor layers' gobjects draw scripts
	// executed.
	Visible     bool
	mode        DrawMode
	mapMat      map[vec]Gobject
	sliceMat    [][]Gobject
	staticImage *ebiten.Image

	// Anchor is the top left point from which the layer is drawn,
	// default being (0,0). Can be anywhere, off screen or not.
	Anchor        image.Point
	numOfGobjects int

	level Level
}

func newGridLayer(
	name string, z int, cellDims Dimensions, gridDims Dimensions,
	drawMode DrawMode, anchor image.Point, padding image.Point) *GridLayer {

	var mapMat map[vec]Gobject = nil
	var sliceMat [][]Gobject = nil
	if drawMode == Sparce {
		mapMat = make(map[vec]Gobject, gridDims.Width*gridDims.Height)
	} else {
		sliceMat = make([][]Gobject, gridDims.Height)
		for i := range sliceMat {
			sliceMat[i] = make([]Gobject, gridDims.Width)
		}
	}
	return &GridLayer{
		Name:            name,
		z:               z,
		cellDimensions:  cellDims,
		layerDimensions: gridDims,
		Visible:         true,
		mode:            drawMode,
		mapMat:          mapMat,
		sliceMat:        sliceMat,
		staticImage:     nil,
		Anchor:          anchor,
		numOfGobjects:   0,
		Padding:         padding,
	}
}

func (le *BaseLevel) addGridLayer(l *GridLayer) *GridLayer {
	ln := len(le.gridLayers)
	l.z = ln
	le.gridLayers = append(le.gridLayers, l)
	l.level = le
	return le.gridLayers[ln]
}

// Creates a grid layer with square cells and no padding within the level.
// Also returns the pointer to it.
func (le *BaseLevel) CreateSimpleGridLayerOnTop(
	name string, squareLength int, width, height int,
	drawMode DrawMode, XOffset, YOffset int) *GridLayer {

	return le.addGridLayer(
		newGridLayer(
			name, 0,
			Dimensions{squareLength, squareLength},
			Dimensions{width, height},
			drawMode,
			image.Point{XOffset, YOffset},
			image.Point{0, 0}))
}

// Shorthand for [(*BaseLevel).CreateSimpleGridLayerOnTop]
// for the current level
func (g *EgridenAssets) CreateSimpleGridLayerOnTop(
	name string, squareLength int, width, height int,
	drawMode DrawMode, XOffset, YOffset int) *GridLayer {

	return g.Level().(*BaseLevel).CreateSimpleGridLayerOnTop(
		name, squareLength, width, height, drawMode, XOffset, YOffset)
}

type GridLayerParameters struct {
	// Width and height of the layer's grid
	GridDimensions Dimensions

	// Width and height of individual cells
	CellDimensions Dimensions
	// Defines the "gaps" between cells:
	// point's X for horizontal gaps length and Y for vertical.
	PaddingVector image.Point

	// Layer's [(GridLayer).Anchor]
	Anchor image.Point
	// Layer's [DrawMode]
	Mode DrawMode
}

// Creates a grid layer with custom parameters within the level and returns the pointer to it.
// If you want a simple square grid layer use [(*BaseLevel).CreateSimpleGridLayerOnTop].
func (le *BaseLevel) CreateGridLayerOnTop(name string, params GridLayerParameters) *GridLayer {
	return le.addGridLayer(
		newGridLayer(
			name, 0,
			params.CellDimensions,
			params.GridDimensions,
			params.Mode,
			params.Anchor,
			params.PaddingVector,
		))
}

// Shorthand for [(*BaseLevel).CreateGridLayerOnTop]
// for the current level
func (g *EgridenAssets) CreateGridLayerOnTop(name string, params GridLayerParameters) *GridLayer {
	return g.Level().(*BaseLevel).CreateGridLayerOnTop(name, params)
}

// False visibility disables drawing both the Sprites and custom draw scripts
// of all Gobjects.
func (l *GridLayer) SetVisibility(to bool) {
	l.Visible = to
}

// Returns the Z level
func (l *GridLayer) Z() int {
	return l.z
}

func (l *GridLayer) anchor() image.Point {
	return l.Anchor
}

// Layer's anchor point as two floats.
func (l *GridLayer) AnchorXYf() (float64, float64) {
	return float64(l.Anchor.X), float64(l.Anchor.Y)
}

// Width and height of the grid.
func (l *GridLayer) Dimensions() (int, int) {
	return l.layerDimensions.Width, l.layerDimensions.Height
}

// Returns a GridLayer at z in the current Level, returns nil if out of bounds.
func (g EgridenAssets) GridLayer(z int) *GridLayer {
	return g.Level().GridLayer(z)
}

func (g EgridenAssets) GridLayers() []*GridLayer {
	return g.Level().GridLayers()
}

// Draw all GridLayers of the current Level in their Z order.
func (g EgridenAssets) DrawAllGridLayers(on *ebiten.Image) {
	g.Level().(*BaseLevel).DrawAllGridLayers(on)
}

// Draw all free layers of the current Level in their Z order.
func (g EgridenAssets) DrawAllFreeLayers(on *ebiten.Image) {
	g.Level().(*BaseLevel).DrawAllFreeLayers(on)
}
