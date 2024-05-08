package egriden

import (
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

type GridLayer struct {
	Name             string
	Z                int
	SquareLength     int
	Width, Height    int
	Visible          bool
	mode             drawMode
	mapMat           map[vec]Gobject
	sliceMat         [][]Gobject
	staticImage      *ebiten.Image
	XOffset, YOffset float64
	NumOfGobjects    int
}

func newGridLayer(name string, z int, squareLength int, width, height int, drawMode drawMode, XOffset, YOffset float64) *GridLayer {
	var mapMat map[vec]Gobject = nil
	var sliceMat [][]Gobject = nil
	if drawMode == Sparce {
		mapMat = make(map[vec]Gobject, width*height)
	} else {
		sliceMat = make([][]Gobject, height)
		for i := range sliceMat {
			sliceMat[i] = make([]Gobject, width)
		}
	}
	return &GridLayer{
		Name:          name,
		Z:             z,
		SquareLength:  squareLength,
		Width:         width,
		Height:        height,
		Visible:       true,
		mode:          drawMode,
		mapMat:        mapMat,
		sliceMat:      sliceMat,
		staticImage:   nil,
		XOffset:       XOffset,
		YOffset:       YOffset,
		NumOfGobjects: 0,
	}
}

func (l GridLayer) IsXYwithinBounds(x, y int) bool {
	return x >= 0 && x < l.Width && y >= 0 && x < l.Height
}

// Creates a grid layer at the lowest empty Z and returns a pointer to it.
//
// See drawMode constants for which one you can use,
// but for small grids Sparce/Dense doesn't make much of a difference.
func (le *BaseLevel) CreateGridLayerOnTop(name string, squareLength int, width, height int, drawMode drawMode, XOffset, YOffset float64) *GridLayer {
	ln := len(le.gridLayers)
	le.gridLayers = append(le.gridLayers, newGridLayer(name, ln, squareLength, width, height, drawMode, XOffset, YOffset))
	return le.gridLayers[ln]
}

// Short hand for BaseLevel.CreateGridLayerOnTop() for the current level
func (g *EgridenAssets) CreateGridLayerOnTop(name string, squareLength int, width, height int, drawMode drawMode, XOffset, YOffset float64) *GridLayer {
	return g.Level().(*BaseLevel).CreateGridLayerOnTop(name, squareLength, width, height, drawMode, XOffset, YOffset)
}

// False visibility disables drawing both the Sprites and custom draw scripts
// of all Gobjects.
func (l *GridLayer) SetVisibility(to bool) {
	l.Visible = to
}
