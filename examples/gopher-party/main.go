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

var gopherPack egriden.SpritePack

func init() {
	seq, err := egriden.CreateImageSequenceFromFolder("gopher", "./Graphics/Gopher")
	if err != nil {
		log.Fatal(seq)
	}
	gopherPack = egriden.NewSpritePackWithSequence(seq)
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.DrawAllGridLayers(screen)
	g.DrawAllFreeLayers(screen)
}

func (g *Game) Update() error {
	g.RunUpdateScripts()
	return nil
}

const scrWidth = 30
const scrHeight = 30

func (g *Game) Layout(_, _ int) (int, int) {
	return scrWidth, scrHeight
}

type BouncyGopher struct {
	egriden.BaseGobject
	goingUp bool
}

func NewBouncyGopher() egriden.Gobject {
	o := &BouncyGopher{
		BaseGobject: egriden.NewBaseGobject("mr bounce", gopherPack),
		goingUp:     true,
	}

	//To be able to access BouncyGopher's fields we use the original
	//pointer `o` instead of `self`. We can do this since this function
	//returns the Gobject in its entirety, with a unique pointer being
	//created at the top. If you're using NewBaseGobject with .Build(),
	//use the `self` parameter instead.
	o.OnUpdateFunc = func(self egriden.Gobject, l egriden.Layer) {
		x, y := o.ScreenPos(l).XY()

		if y < 0 || y >= scrHeight-8 {
			o.NextFrame()
			o.goingUp = !o.goingUp
		}

		if o.goingUp {
			y--
		} else {
			y++
		}

		l.(*egriden.FreeLayer).MoveGobjectTo(o, x, y)
	}

	return o
}

func main() {
	g := &Game{}
	g.InitEgridenAssets()

	gopherGridL := g.CreateSimpleGridLayerOnTop(
		"gopher grid", 10, 3, 3, egriden.Dense, 0, 0)
	shakyGopher := egriden.NewBaseGobject("mr shake", gopherPack)
	shakyGopher.OnDrawFunc =
		func(self egriden.Gobject, i *ebiten.Image, l egriden.Layer) {
			self.SetDrawOffsets(float64(rand.IntN(2)-2), 0)
			self.DrawSprite(i, l)
		}

	spinnyGopher := egriden.NewBaseGobject("mr spin", gopherPack)
	spinnyGopher.OnDrawFunc =
		func(self egriden.Gobject, i *ebiten.Image, l egriden.Layer) {
			//Centered rotation
			s := self.Sprite().Bounds().Size()
			self.SpritePack().
				DrawOptions.GeoM.
				Translate(-float64(s.X)/2, -float64(s.Y)/2)
			self.SpritePack().
				DrawOptions.GeoM.Rotate(0.05)
			self.SpritePack().
				DrawOptions.GeoM.
				Translate(float64(s.X)/2, float64(s.Y)/2)
			self.DrawSprite(i, l)
		}

	gopherGridL.AddGobject(shakyGopher.Build(), 0, 0)
	gopherGridL.AddGobject(spinnyGopher.Build(), 0, 1)
	gopherGridL.AddGobject(shakyGopher.Build(), 2, 0)
	gopherGridL.AddGobject(spinnyGopher.Build(), 1, 0)
	gopherGridL.AddGobject(shakyGopher.Build(), 0, 2)
	gopherGridL.AddGobject(spinnyGopher.Build(), 1, 2)
	gopherGridL.AddGobject(shakyGopher.Build(), 2, 2)
	gopherGridL.AddGobject(spinnyGopher.Build(), 2, 1)

	bouncyL := g.CreateFreeLayerOnTop("bouncy layer", 0, 0)
	bouncyL.AddGobject(NewBouncyGopher(), 10, 12)

	ebiten.SetWindowSize(600, 600)

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
