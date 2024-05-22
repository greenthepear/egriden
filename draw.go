package egriden

import "github.com/hajimehoshi/ebiten/v2"

type Layer interface {
	DrawSprite(o Gobject, on *ebiten.Image)
	Offsets() (float64, float64)
}

func appliedDrawOptionsForPosition(o Gobject, layer Layer, x, y float64) *ebiten.DrawImageOptions {
	copy := *o.SpritePack().DrawOptions
	r := &copy
	xoffset, yoffset := layer.Offsets()
	spriteXoffset, spriteYoffset := o.SpritePack().XOffset, o.SpritePack().YOffset
	switch l := layer.(type) {
	case *GridLayer:
		if l.mode == Static {
			xoffset, yoffset = 0, 0
		}
		r.GeoM.Translate(
			float64(x)*float64(l.SquareLength)+xoffset+spriteXoffset,
			float64(y)*float64(l.SquareLength)+yoffset+spriteYoffset)
	case *FreeLayer:
		if l.static {
			xoffset, yoffset = 0, 0
		}
		r.GeoM.Translate(
			float64(x)+l.XOffset+spriteXoffset,
			float64(y)+l.YOffset+spriteYoffset)
	}
	return r

}

func (l GridLayer) drawFromSliceMat(on *ebiten.Image) {
	for y := range l.Height {
		for x := range l.Width {
			o := l.sliceMat[y][x]
			if o == nil || !o.IsVisible() {
				continue
			}

			if o.OnDraw() != nil {
				o.OnDraw()(on, &l)
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
		l.Width*l.SquareLength, l.Height*l.SquareLength)
	l.drawFromSliceMat(img)
	l.staticImage = img
}

func (l GridLayer) DrawSprite(o Gobject, on *ebiten.Image) {
	x, y := o.XY()
	on.DrawImage(o.Sprite(),
		appliedDrawOptionsForPosition(o, &l, float64(x), float64(y)))
}

func (fl FreeLayer) DrawSprite(o Gobject, on *ebiten.Image) {
	on.DrawImage(o.Sprite(), o.SpritePack().DrawOptions)
}

// Draw the layer
func (l GridLayer) Draw(on *ebiten.Image) {
	if !l.Visible {
		return
	}

	switch l.mode {
	case Sparce:
		for _, o := range l.mapMat {
			if !o.IsVisible() {
				continue
			}

			if o.OnDraw() != nil {
				o.OnDraw()(on, &l)
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
		op.GeoM.Translate(l.XOffset, l.YOffset)
		on.DrawImage(l.staticImage, op)
	}
}

func (fl FreeLayer) internalDraw(on *ebiten.Image) {
	for _, k := range fl.gobjects.keys {
		if k.OnDraw() != nil {
			k.OnDraw()(on, &fl)
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
		op.GeoM.Translate(fl.XOffset, fl.YOffset)
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
