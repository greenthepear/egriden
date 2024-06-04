package egriden

import (
	"image"
	"image/png"
	"log"
	"os"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
)

func ebitenimageToFile(img *ebiten.Image, path string) {
	imageImage := image.NewRGBA(img.Bounds())
	dx, dy := img.Bounds().Dx(), img.Bounds().Dy()
	for x := range dx {
		for y := range dy {
			imageImage.Set(x, y, img.At(x, y))
		}
	}
	f, err := os.Create(path)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	if err = png.Encode(f, imageImage); err != nil {
		log.Fatal(err)
	}
}

func TestFreeLayers(t *testing.T) {
	g := EgridenAssets{}
	g.InitEgridenComponents()

	testoffx, testoffy := 10.0, 20.0
	fl1 := g.CreateFreeLayerOnTop("freelayer1", testoffx, testoffy)
	offx, offy := fl1.Offsets()
	if offx != testoffx || offy != testoffy {
		t.Errorf("Offsets didn't get applied! (%.0f, %.0f) != (%.0f, %.0f)",
			testoffx, testoffy, offx, offy)
	}
	if fl1.static {
		t.Errorf("Layer shouldn't be static!")
	}
	seq, err := CreateImageSequenceFromFolder(
		"unrev", "./examples/gridsweeper/Graphics/unrevealed")

	if err != nil {
		t.Error(err)
	}
	o := NewBaseGobject(
		"tester", NewSpritePackWithSequence(seq))
	fl1.AddGobject(
		o.Build(),
		5, 5)

	eimg := ebiten.NewImage(100, 100)
	fl1.Draw(eimg)
	cornerpixel := o.Sprite().originalImage
	if eimg.At(int(testoffx)+5, int(testoffy+5)) != cornerpixel {
		ebitenimageToFile(eimg, "_test_bad_image.png")
		t.Errorf("Image didn't draw properly, saved to `_test_bad_image.png`")
	}
}
