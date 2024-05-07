package main

import (
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
	g.DrawAllLayers(screen)
}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 320, 320
}

func main() {
	g := &Game{}
	g.InitEgridenComponents()
	lbg := g.CreateGridLayerOnTop("Background",
		defGridLen, defGridWidth, defGridHeight, egriden.Static, 0, 0)
	lbo := g.CreateGridLayerOnTop("Bombs",
		defGridLen, defGridWidth, defGridHeight, egriden.Sparce, 0, 0)
	lre := g.CreateGridLayerOnTop("Reveal tiles",
		defGridLen, defGridWidth, defGridHeight, egriden.Dense, 0, 0)

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

	for x := range defGridWidth {
		for y := range defGridHeight {
			lbg.AddGobject(objBacktile.Build(), x, y)
			backtileCopy := objRevealTile.Build()
			backtileCopy.SetFrame(rand.IntN(3))
			lre.AddGobject(backtileCopy, x, y)
		}
	}

	for range 6 {
		lbo.AddGobject(objBomb.Build(), rand.IntN(defGridWidth), rand.IntN(defGridHeight))
	}

	ebiten.SetWindowSize(640, 640)
	ebiten.SetWindowTitle("Gridsweeper")

	//lre.SetVisibility(false)

	if err = ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
