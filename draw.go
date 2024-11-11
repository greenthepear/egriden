package egriden

import (
	"github.com/greenthepear/imggg"
	"github.com/hajimehoshi/ebiten/v2"
)

type Layer interface {
	DrawSprite(o Gobject, on *ebiten.Image)
	Static() bool
	anchor() imggg.Point[float64]
}

// Returns ebiten.DrawImageOptions of a Gobject's SpritePack applied
// to the layer
func appliedDrawOptionsForPosition(o Gobject, layer Layer, x, y float64) *ebiten.DrawImageOptions {
	copy := *o.SpritePack().DrawOptions
	r := &copy
	drawX, drawY := o.ScreenPos(layer).XY()

	// In static layers, the layer anchor offsets are handled in the draw
	// function itself, so they need to be subtracted here
	if layer.Static() {
		xRealign, yRealign := layer.anchor().X, layer.anchor().Y
		r.GeoM.Translate(
			drawX-float64(xRealign), drawY-float64(yRealign))
	} else {
		r.GeoM.Translate(
			drawX, drawY)
	}

	return r

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

// Refresh image of a static grid layer
func (l *GridLayer) RefreshImage() {
	if l.mode != Static {
		return
	}
	img := ebiten.NewImage(
		l.layerDimensions.Width*l.cellDimensions.Width,
		l.layerDimensions.Height*l.cellDimensions.Height)
	l.drawFromSliceMat(img)
	l.staticImage = img
}

func (l GridLayer) DrawSprite(o Gobject, on *ebiten.Image) {
	x, y := o.GridPos().XY()
	on.DrawImage(o.Sprite(),
		appliedDrawOptionsForPosition(o, &l, float64(x), float64(y)))
}

func (fl FreeLayer) DrawSprite(o Gobject, on *ebiten.Image) {
	x, y := o.GridPos().XY()
	on.DrawImage(o.Sprite(),
		appliedDrawOptionsForPosition(o, &fl, float64(x), float64(y)))
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
	case Static:
		if l.staticImage == nil {
			l.RefreshImage()
		}
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(l.Anchor.X), float64(l.Anchor.Y))
		on.DrawImage(l.staticImage, op)
	}
}

func (fl FreeLayer) internalDraw(on *ebiten.Image) {
	for _, k := range fl.gobjects.keys {
		if k.OnDraw() != nil {
			k.OnDraw()(k, on, &fl)
			continue
		}
		k.DrawSprite(on, &fl)
	}
}

// Draw the layer
func (fl FreeLayer) Draw(on *ebiten.Image) {
	if !fl.Visible {
		return
	}
	if fl.static {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(fl.Anchor.X), float64(fl.Anchor.Y))
		on.DrawImage(fl.staticImage, op)
		return
	}
	fl.internalDraw(on)
}

// Refresh/create image of a static free layer
func (fl *FreeLayer) RefreshImage() {
	if !fl.static {
		return
	}
	fl.internalDraw(fl.staticImage)
}
