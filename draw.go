package egriden

import "github.com/hajimehoshi/ebiten/v2"

func createDrawImageOptionsForXY(x, y float64) *ebiten.DrawImageOptions {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(x, y)
	return op
}

func (l GridLayer) drawFromSliceMat(on *ebiten.Image) {
	for y := range l.Height {
		for x := range l.Width {
			o := l.sliceMat[y][x]
			if o == nil || !o.IsVisible() {
				continue
			}

			if !o.DoesDrawScriptOverwriteSprite() {
				on.DrawImage(o.Sprite(),
					createDrawImageOptionsForXY(
						float64(x)*float64(l.SquareLength)+l.XOffset,
						float64(y)*float64(l.SquareLength)+l.YOffset))
			}

			if o.OnDraw() != nil {
				o.OnDraw()(on)
			}

		}
	}
}

func (l *GridLayer) RefreshImage() {
	img := ebiten.NewImage(
		l.Width*l.SquareLength, l.Height*l.SquareLength)
	l.drawFromSliceMat(img)
	l.staticImage = img
}

func (l GridLayer) Draw(screen *ebiten.Image) {
	if !l.Visible {
		return
	}

	switch l.mode {
	case Sparce:
		for vec, o := range l.mapMat {
			if !o.IsVisible() {
				continue
			}
			if !o.DoesDrawScriptOverwriteSprite() {
				screen.DrawImage(o.Sprite(),
					createDrawImageOptionsForXY(
						float64(vec.x)*float64(l.SquareLength)+l.XOffset,
						float64(vec.y)*float64(l.SquareLength)+l.YOffset))
			}
			if o.OnDraw() != nil {
				o.OnDraw()(screen)
			}
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
