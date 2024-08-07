package main

import (
	"image"
	"image/color"
	"log"

	"github.com/greenthepear/egriden"
	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	egriden.EgridenAssets
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.GridLayer(0).DebugDrawCheckerBoard(color.Black, color.White, screen)
}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 320, 320
}

func main() {
	//Initialize
	g := &Game{}
	g.InitEgridenAssets()

	g.CreateGridLayerOnTop("Base", egriden.GridLayerParameters{
		GridDimensions: egriden.Dimensions{
			Width: 8, Height: 10},
		CellDimensions: egriden.Dimensions{
			Width: 20, Height: 12},
		PaddingVector: image.Point{0, 0},
	})

	//Run the game
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
