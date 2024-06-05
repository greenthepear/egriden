package main

import (
	"fmt"
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
	g.DrawAllGridLayers(screen)
	g.DrawAllFreeLayers(screen)
}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 320, 320
}

func main() {
	g := &Game{}
	g.InitEgridenAssets()
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

	//Testing free layers
	lfree := g.CreateFreeLayerOnTop("free test", 21, 21)
	lfree.AddGobject(objBomb.Build(), 80, 0)

	ebiten.SetWindowSize(640, 640)
	ebiten.SetWindowTitle("Gridsweeper")

	//lre.SetVisibility(false)

	for x := range lre.Width {
		for y := range lre.Height {
			fmt.Printf("%v\n\t L %p\n", lre.GobjectAt(x, y), lre.GobjectAt(x, y).SpritePack().DrawOptions)
		}
	}

	if err = ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
