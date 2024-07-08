package egriden

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

// For optimization there are a couple of ways a grid layer can be draw depending if it
// changes frequently (Static or not) and if it has many (Dense) or few (Sparce)
// gobjects most of the time.
type drawMode int

const (
	//Used for sparcely populated grids, ranges over a map for drawing
	Sparce drawMode = iota
	//Used for thickly populated grids, ranges over a slice for drawing
	Dense
	//Used for layers that don't get updated often, creates ebiten.Image of the the entire layer.
	//Can be refreshed with GridLayer.RefreshImage().
	Static
)

type vec struct {
	x, y int
}

type Dimensions struct {
	width, height int
}

type GridLayer struct {
	Name            string // Name of the layer, for convenience sake
	z               int
	cellDimensions  Dimensions
	layerDimensions Dimensions
	Visible         bool
	mode            drawMode
	mapMat          map[vec]Gobject
	sliceMat        [][]Gobject
	staticImage     *ebiten.Image
	AnchorPt        image.Point
	numOfGobjects   int

	level Level
}

func newGridLayer(
	name string, z int, cellWidth int, cellHeight int, lwidth, lheight int, drawMode drawMode, XOffset, YOffset int) *GridLayer {

	var mapMat map[vec]Gobject = nil
	var sliceMat [][]Gobject = nil
	if drawMode == Sparce {
		mapMat = make(map[vec]Gobject, lwidth*lheight)
	} else {
		sliceMat = make([][]Gobject, lheight)
		for i := range sliceMat {
			sliceMat[i] = make([]Gobject, lwidth)
		}
	}
	return &GridLayer{
		Name:            name,
		z:               z,
		cellDimensions:  Dimensions{cellWidth, cellHeight},
		layerDimensions: Dimensions{lwidth, lheight},
		Visible:         true,
		mode:            drawMode,
		mapMat:          mapMat,
		sliceMat:        sliceMat,
		staticImage:     nil,
		AnchorPt:        image.Point{int(XOffset), int(YOffset)},
		numOfGobjects:   0,
	}
}

// Creates a grid layer at the lowest empty Z and returns a pointer to it.
//
// See drawMode constants for which one you can use,
// but for small grids Sparce/Dense doesn't make much of a difference.
func (le *BaseLevel) CreateSimpleGridLayerOnTop(name string, squareLength int, width, height int, drawMode drawMode, XOffset, YOffset int) *GridLayer {
	ln := len(le.gridLayers)
	newLayer := newGridLayer(
		name, ln, squareLength, squareLength,
		width, height, drawMode, XOffset, YOffset)
	le.gridLayers = append(le.gridLayers, newLayer)
	newLayer.level = le
	return le.gridLayers[ln]
}

// Short hand for BaseLevel.CreateSimpleGridLayerOnTop() for the current level
func (g *EgridenAssets) CreateSimpleGridLayerOnTop(name string, squareLength int, width, height int, drawMode drawMode, XOffset, YOffset int) *GridLayer {
	return g.Level().(*BaseLevel).CreateSimpleGridLayerOnTop(name, squareLength, width, height, drawMode, XOffset, YOffset)
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

func (l *GridLayer) Anchor() image.Point {
	return l.AnchorPt
}

func (l *GridLayer) AnchorXYf() (float64, float64) {
	return float64(l.AnchorPt.X), float64(l.AnchorPt.Y)
}

func (l *GridLayer) Dimensions() (int, int) {
	return l.layerDimensions.width, l.layerDimensions.height
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
