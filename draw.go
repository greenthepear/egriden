package egriden

import "github.com/hajimehoshi/ebiten/v2"

func createDrawImageOptionsForXY(x, y float64) *ebiten.DrawImageOptions {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(x, y)
	return op
}

func (l GridLayer) drawFromSliceMat(on *ebiten.Image) {
	for y := range l.height {
		for x := range l.width {
			o := l.sliceMat[y][x]
			if o == nil {
				continue
			}
			if o.OnDraw() != nil {
				o.OnDraw()(on)
				if o.DoesDrawScriptOverwriteSprite() {
					continue
				}
			}
			if !o.IsVisible() {
				continue
			}
			on.DrawImage(o.Sprite(),
				createDrawImageOptionsForXY(
					float64(x)*float64(l.squareLength)+l.xOffset,
					float64(y)*float64(l.squareLength)+l.yOffset))
		}
	}
}

func (l *GridLayer) RefreshImage() {
	img := ebiten.NewImage(
		l.width*l.squareLength, l.height*l.squareLength)
	l.drawFromSliceMat(img)
	l.staticImage = img
}

func (l GridLayer) Draw(screen *ebiten.Image) {
	if !l.visible {
		return
	}

	switch l.mode {
	case Sparce:
		for vec, o := range l.mapMat {
			if o.OnDraw() != nil {
				o.OnDraw()(screen)
				if o.DoesDrawScriptOverwriteSprite() {
					continue
				}
			}
			if !o.IsVisible() {
				continue
			}
			screen.DrawImage(o.Sprite(),
				createDrawImageOptionsForXY(
					float64(vec.x)*float64(l.squareLength)+l.xOffset,
					float64(vec.y)*float64(l.squareLength)+l.yOffset))
		}
	case Dense:
		l.drawFromSliceMat(screen)
	case Static:
		if l.staticImage == nil {
			l.RefreshImage()
		}
		screen.DrawImage(l.staticImage,
			createDrawImageOptionsForXY(l.xOffset*float64(l.squareLength), l.yOffset*float64(l.squareLength)))
	}
}
