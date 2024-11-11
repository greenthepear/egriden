package main

import (
	"image/color"
	"log"
	"math/rand/v2"

	"github.com/greenthepear/egriden"
	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	egriden.EgridenAssets
}

var defGridLen = 16
var defGridHeight = 20
var defGridWidth = 20

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0xeb, 0xeb, 0xeb, 0xff})
	g.DrawAllGridLayers(screen)
	g.DrawAllFreeLayers(screen)
}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 340, 360
}

func main() {
	g := &Game{}
	g.InitEgridenAssets()
	offx, offy := 10, 30
	lbg := g.CreateSimpleGridLayerOnTop("Background",
		defGridLen, defGridWidth, defGridHeight, egriden.Static, offx, offy)
	lbo := g.CreateSimpleGridLayerOnTop("Bombs",
		defGridLen, defGridWidth, defGridHeight, egriden.Sparse, offx, offy)
	lre := g.CreateSimpleGridLayerOnTop("Reveal tiles",
		defGridLen, defGridWidth, defGridHeight, egriden.Dense, offx, offy)
	lui := g.CreateFreeLayerOnTop("UI", 0, 0)

	seq, err := egriden.CreateImageSequenceFromPaths("backtile", "Graphics/backtile.png")
	if err != nil {
		log.Fatal(err)
	}
	sprBackground := egriden.NewSpritePackWithSequence(seq)
	objBacktile := egriden.NewBaseGobject("backtile", sprBackground)

	seq, err = egriden.CreateImageSequenceFromPaths("backtile", "Graphics/explode.png")
	if err != nil {
		log.Fatal(err)
	}
	sprBomb := egriden.NewSpritePackWithSequence(seq)
	objBomb := egriden.NewBaseGobject("bomb", sprBomb)

	seq, err = egriden.CreateImageSequenceFromFolder("unrevealed", "Graphics/unrevealed/")
	if err != nil {
		log.Fatal(err)
	}
	sprReveal := egriden.NewSpritePackWithSequence(seq)
	objRevealTile := egriden.NewBaseGobject("reveal", sprReveal)

	seq, err = egriden.CreateImageSequenceFromPaths("smiley", "Graphics/smiley.png")
	if err != nil {
		log.Fatal(err)
	}
	sprSmiley := egriden.NewSpritePackWithSequence(seq)
	objSmiley := egriden.NewBaseGobject("smiley", sprSmiley)
	// TODO: remove magic numbers
	lui.AddGobject(objSmiley.Build(), 162, 5)

	for x := range defGridWidth {
		for y := range defGridHeight {
			lbg.AddGobject(objBacktile.Build(), x, y)
			backtileCopy := objRevealTile.Build()
			backtileCopy.SetFrame(rand.IntN(3))
			lre.AddGobject(backtileCopy, x, y)
		}
	}

	for range 10 {
		lbo.AddGobject(objBomb.Build(), rand.IntN(defGridWidth), rand.IntN(defGridHeight))
	}

	ebiten.SetWindowSize(680, 720)
	ebiten.SetWindowTitle("Gridsweeper")

	lre.Visible = false
	lbo.Visible = false

	if err = ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
