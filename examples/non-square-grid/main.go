package main

import (
	"fmt"
	"image/color"
	"log"

	"github.com/greenthepear/egriden"
	"github.com/greenthepear/imggg"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Game struct {
	egriden.EgridenAssets
}

func (g *Game) Draw(screen *ebiten.Image) {
	l0 := g.GridLayer(0)
	l0.DebugDrawCheckerBoard(
		color.White,
		color.RGBA{0x99, 0x99, 0x99, 0xff}, screen)

	cx, cy := ebiten.CursorPosition()
	c, b := l0.CellAtScreenPosWithPadding(float64(cx), float64(cy))
	s := fmt.Sprintf("Pointing at %v, is outside gap: %v", c.Coordinate, b)
	ebitenutil.DebugPrint(screen, s)
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
		PaddingVector: imggg.Pt(4.0, 4.0),
		Anchor:        imggg.Pt(60.0, 76.0),
	})

	ebiten.SetWindowSize(640, 640)
	//Run the game
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
