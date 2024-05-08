package egriden

import "github.com/hajimehoshi/ebiten/v2"

func createDrawImageOptionsForXY(x, y float64) *ebiten.DrawImageOptions {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(x, y)
	return op
}

type Layer interface {
	DrawSprite(o Gobject, on *ebiten.Image)
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
		createDrawImageOptionsForXY(
			float64(x)*float64(l.SquareLength)+l.XOffset,
			float64(y)*float64(l.SquareLength)+l.YOffset))
}

func (fl FreeLayer) DrawSprite(o Gobject, on *ebiten.Image) {
	x, y := o.XY()
	on.DrawImage(o.Sprite(), createDrawImageOptionsForXY(
		float64(x)+fl.XOffset, float64(y)+fl.YOffset))
}

// Draw the layer
func (l GridLayer) Draw(screen *ebiten.Image) {
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
				o.OnDraw()(screen, &l)
				continue
			}

			o.DrawSprite(screen, &l)
		}
	case Dense:
		l.drawFromSliceMat(screen)
	case Static:
		if l.staticImage == nil {
			l.RefreshImage()
		}
		screen.DrawImage(l.staticImage,
			createDrawImageOptionsForXY(l.XOffset*float64(l.SquareLength), l.YOffset*float64(l.SquareLength)))
	}
}

// Draw the layer
func (fl FreeLayer) Draw(on *ebiten.Image) {
	if !fl.Visible {
		return
	}
	for _, k := range fl.gobjects.keys {
		if k.OnDraw() != nil {
			k.OnDraw()(on, &fl)
			continue
		}
		k.DrawSprite(on, &fl)
	}
}

// Refresh image of a static free layer
func (fl *FreeLayer) RefreshImage() {
	if !fl.static {
		return
	}
	fl.Draw(fl.staticImage)
}
