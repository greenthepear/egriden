package egriden

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
)

type drawMode int

const (
	Sparce drawMode = iota //Used for sparcely populated grids, ranges over a map for drawing
	Dense                  //Used for thickly populated grids, ranges over a slice for drawing
	Static                 //Used for layers that don't get updated often, creates ebiten.Image of the the entire layer
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

func (g *EgridenAssets) CreateGridLayerOnTop(name string, squareLength int, width, height int, drawMode drawMode, XOffset, YOffset float64) *GridLayer {
	ln := len(g.gridLayers)
	g.gridLayers = append(g.gridLayers, newGridLayer(name, ln, squareLength, width, height, drawMode, XOffset, YOffset))
	return g.gridLayers[ln]
}

func (l *GridLayer) SetVisibility(to bool) {
	l.Visible = to
}

func (l GridLayer) GobjectAt(x, y int) Gobject {
	if l.mode == Sparce {
		return l.mapMat[vec{x, y}]
	}
	return l.sliceMat[y][x]
}

func (l GridLayer) IsOccupiedAt(x, y int) bool {
	return l.GobjectAt(x, y) != nil
}

func (l *GridLayer) AddObject(o Gobject, x, y int) {
	if Warnings && l.IsOccupiedAt(x, y) {
		fmt.Printf(
			"Egriden WARNING: Gobject already exists at (%d,%d) in layer %s (%d). It will be overwritten.",
			x, y, l.Name, l.Z)
	}

	o.setXY(x, y)
	if l.mode == Sparce {
		l.mapMat[vec{x, y}] = o
		return
	}
	l.sliceMat[y][x] = o
}
