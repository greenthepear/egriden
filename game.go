package egriden

import (
	"github.com/greenthepear/imggg"
	"github.com/hajimehoshi/ebiten/v2"
)

// Egriden components to be embedded in your Game{} struct.
type EgridenAssets struct {
	Levels            []Level
	CurrentLevelIndex int
}

// Get the current level, returns nil if there are no levels. Index starts at
// 0 so the first level added will be returned if
// (EgridenAssets).CurrentLevelIndex was not changed
func (g *EgridenAssets) Level() Level {
	if g.CurrentLevelIndex >= len(g.Levels) {
		return nil
	}
	return g.Levels[g.CurrentLevelIndex]
}

// Get a level by it's name. Returns nil if not found.
func (g *EgridenAssets) LevelByName(name string) Level {
	for _, le := range g.Levels {
		if le.Name() == name {
			return le
		}
	}
	return nil
}

// Append level to the end of the list and return it
func (g *EgridenAssets) AddLevel(le Level) Level {
	g.Levels = append(g.Levels, le)
	idx := len(g.Levels) - 1
	le.setIndex(idx)
	return g.Levels[idx]
}

// Sets the current game's level or adds it if it's not in the assets already
func (g *EgridenAssets) SetLevelTo(le Level) {
	for i, rangeLe := range g.Levels {
		if rangeLe == le {
			g.CurrentLevelIndex = i
			return
		}
	}

	//If not found, add that level I guess
	g.AddLevel(le)
	g.CurrentLevelIndex = len(g.Levels) - 1
}

// Set the next level by iterating the level index
func (g *EgridenAssets) NextLevel() {
	g.CurrentLevelIndex = (g.CurrentLevelIndex + 1) % len(g.Levels)
}

/// Deprecated

// Shorthand for [Level.CreateSimpleGridLayerOnTop]
// for the current level
//
// Deprecated: Just use CreateGridLayerOnTop.
func (g *EgridenAssets) CreateSimpleGridLayerOnTop(
	name string, squareLength int, width, height int,
	drawMode DrawMode, XOffset, YOffset float64) *GridLayer {

	return g.Level().CreateSimpleGridLayerOnTop(
		name, squareLength, width, height, drawMode, XOffset, YOffset)
}

// Shorthand for [Level.CreateGridLayerOnTop]
// for the current level
//
// Deprecated: Use method directly from Level
func (g *EgridenAssets) CreateGridLayerOnTop(
	name string, params GridLayerParameters) *GridLayer {

	return g.Level().CreateGridLayerOnTop(name, params)
}

// Run this while initializing the game, before adding any layers. Creates a
// level called `Default`
//
// Deprecated: just use g.AddLevel(NewBaseLevel("Default"))
func (g *EgridenAssets) InitEgridenAssets() {
	g.AddLevel(NewBaseLevel("Default"))
}

// Returns a GridLayer at z in the current Level, returns nil if out of bounds.
//
// Deprecated: Use method directly from Level
func (g EgridenAssets) GridLayer(z int) *GridLayer {
	return g.Level().GridLayer(z)
}

// Draw all GridLayers of the current Level in their Z order.
//
// Deprecated: Use method directly from Level
func (g EgridenAssets) DrawAllGridLayers(on *ebiten.Image) {
	g.Level().DrawAllGridLayers(on)
}

// Draw all free layers of the current Level in their Z order.
//
// Deprecated: Use method directly from Level
func (g EgridenAssets) DrawAllFreeLayers(on *ebiten.Image) {
	g.Level().DrawAllFreeLayers(on)
}

// Shortcut for g.Level().CreateFreeLayerOnTop().
// Level implementation must have BaseLevel component.
//
// Deprecated: Use method directly from Level
func (g *EgridenAssets) CreateFreeLayerOnTop(
	name string, xOffset, yOffset float64) *FreeLayer {

	return g.Level().CreateFreeLayerOnTop(name,
		imggg.Point[float64]{X: xOffset, Y: yOffset})
}
