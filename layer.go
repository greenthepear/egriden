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

// Creates a grid layer at the lowest empty Z and returns a pointer to it.
//
// See drawMode constants for which one you can use,
// but for small grids Sparce/Dense doesn't make much of a difference.
func (g *EgridenAssets) CreateGridLayerOnTop(name string, squareLength int, width, height int, drawMode drawMode, XOffset, YOffset float64) *GridLayer {
	ln := len(g.gridLayers)
	g.gridLayers = append(g.gridLayers, newGridLayer(name, ln, squareLength, width, height, drawMode, XOffset, YOffset))
	return g.gridLayers[ln]
}

// False visibility disables drawing both the Sprites and custom draw scripts
// of all Gobjects.
func (l *GridLayer) SetVisibility(to bool) {
	l.Visible = to
}

// Returns Gobject at x y, nil if empty. Panics if out of bounds.
func (l GridLayer) GobjectAt(x, y int) Gobject {
	if x >= l.Width || y >= l.Height {
		panic("GobjectAt() panic! Out of bounds.")
	}
	if l.mode == Sparce {
		return l.mapMat[vec{x, y}]
	}
	return l.sliceMat[y][x]
}

func (l GridLayer) IsOccupiedAt(x, y int) bool {
	return l.GobjectAt(x, y) != nil
}

// Adds Gobject to the layer at x y. Will overwrite the any existing Gobject there.
func (l *GridLayer) AddObject(o Gobject, x, y int) {
	o.setXY(x, y)
	if l.mode == Sparce {
		if l.mapMat[vec{x, y}] != nil {
			l.mapMat[vec{x, y}].markForDeletion()
		}
		l.mapMat[vec{x, y}] = o
		return
	}
	if l.sliceMat[y][x] != nil {
		l.sliceMat[y][x].markForDeletion()
	}
	l.sliceMat[y][x] = o
}
