package egriden

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type EgridenAssets struct {
	Levels            []Level
	CurrentLevelIndex int
}

// Get the current level
func (g *EgridenAssets) Level() Level {
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
	le.(*BaseLevel).index = idx
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

// Run this while initalizing the game, before adding any layers. Creates a scene called `Default`
func (g *EgridenAssets) InitEgridenComponents() {
	g.AddLevel(NewBaseLevel("Default"))
}

// Returns a GridLayer at z in the current level, panics if out of bounds
func (g EgridenAssets) GridLayer(z int) *GridLayer {
	return g.Level().GridLayer(z)
}

func (g EgridenAssets) GridLayers() []*GridLayer {
	return g.Level().GridLayers()
}

// Draw all GridLayers of the current Level in their Z order. Use this in the Draw() function.
func (g EgridenAssets) DrawAllLayers(screen *ebiten.Image) {
	g.Level().(*BaseLevel).DrawAllLayers(screen)
}

func (g *EgridenAssets) RunUpdateScripts() {
	g.Level().(*BaseLevel).RunUpdateScripts()
}
