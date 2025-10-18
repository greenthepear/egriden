package egriden

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// Returns ebiten.DrawImageOptions of a Gobject's SpritePack applied
// to the layer
func appliedDrawOptionsForPosition(
	o Gobject, layer Layer) *ebiten.DrawImageOptions {

	copy := *o.SpritePack().DrawOptions
	copy.GeoM.Translate(o.ScreenPos(layer).XY())
	return &copy

}

func (l GridLayer) drawFromSliceMat(on *ebiten.Image) {
	for y := range l.layerDimensions.Height {
		for x := range l.layerDimensions.Width {
			o := l.sliceMat[y][x]
			if o == nil || !o.IsVisible() {
				continue
			}

			if o.OnDraw() != nil {
				o.OnDraw()(o, on, &l)
				continue
			}
			o.DrawSprite(on, &l)

		}
	}
}

func (l GridLayer) DrawLikeSprite(
	img *ebiten.Image, o Gobject, on *ebiten.Image) {

	on.DrawImage(img,
		appliedDrawOptionsForPosition(o, &l))
}

func (l GridLayer) DrawSprite(o Gobject, on *ebiten.Image) {
	l.DrawLikeSprite(o.Sprite(), o, on)
}

func (fl FreeLayer) DrawLikeSprite(img *ebiten.Image, o Gobject, on *ebiten.Image) {
	on.DrawImage(img,
		appliedDrawOptionsForPosition(o, &fl))
}

func (fl FreeLayer) DrawSprite(o Gobject, on *ebiten.Image) {
	fl.DrawLikeSprite(o.Sprite(), o, on)
}

// Draw the layer
func (l GridLayer) Draw(on *ebiten.Image) {
	if !l.Visible {
		return
	}

	switch l.mode {
	case Sparse:
		for _, o := range l.mapMat {
			if !o.IsVisible() {
				continue
			}

			if o.OnDraw() != nil {
				o.OnDraw()(o, on, &l)
				continue
			}

			o.DrawSprite(on, &l)
		}
	case Dense:
		l.drawFromSliceMat(on)
	default:
		panic("unknown DrawMode")
	}
}

func (fl FreeLayer) internalDraw(on *ebiten.Image) {
	for e := fl.gobjects.Front(); e != nil; e = e.Next() {
		o, ok := e.Value.(Gobject)
		if !ok {
			panic("list element isn't a Gobject")
		}
		if o.OnDraw() != nil {
			o.OnDraw()(o, on, &fl)
			continue
		}
		o.DrawSprite(on, &fl)
	}
}

// Draw the layer
func (fl FreeLayer) Draw(on *ebiten.Image) {
	if !fl.Visible {
		return
	}
	fl.internalDraw(on)
}
